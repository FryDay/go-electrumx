package electrumx

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
)

const delim = byte('\n')

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrNodeConnected  = errors.New("node already connected")
	ErrNodeShutdown   = errors.New("node has shutdown")
)

type response struct {
	Id     uint64 `json:"id"`
	Method string `json:"method"`
	Error  any    `json:"error"`
}

type request struct {
	Id     uint64        `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type basicResp struct {
	Result string `json:"result"`
}

type container struct {
	content []byte
	err     error
}

type Node struct {
	transport *transport

	handlersLock sync.RWMutex
	handlers     map[uint64]chan *container

	pushHandlersLock sync.RWMutex
	pushHandlers     map[string][]chan *container

	errs chan error
	quit chan struct{}

	// nextId tags a request, and get the same id from server result.
	// Should be atomic operation for concurrence.
	// notice the max request limit, if reach to the max times,
	// 0 will be the next id. Assume the oldest has been deal completely.
	nextId uint64
}

// NewNode creates a new node.
func NewNode() *Node {
	n := &Node{
		handlers:     make(map[uint64]chan *container),
		pushHandlers: make(map[string][]chan *container),

		errs: make(chan error),
		quit: make(chan struct{}),
	}

	return n
}

// Errors returns any errors the node ran into while listening to messages.
func (n *Node) Errors() <-chan error {
	return n.errs
}

// Connect creates a new TCP connection to the specified address. If the TLS
// config is not nil, TLS is applied to the connection.
func (n *Node) Connect(ctx context.Context, addr string, config *tls.Config) error {
	if n.transport != nil {
		return ErrNodeConnected
	}

	transport, err := newTransport(ctx, addr, config)
	if err != nil {
		return err
	}
	n.transport = transport

	listenCtx, cancel := context.WithCancel(context.Background())
	go func() {
		n.transport.listen(listenCtx)
	}()

	// Quit the transport listening once the node shuts down
	go func() {
		<-n.quit
		cancel()
	}()

	go n.listen(listenCtx)

	return nil
}

// listen processes messages from the server.
func (n *Node) listen(ctx context.Context) {
	for {
		// Not exactly sure how this happened, but it must be a race condition
		// where disconnect and shutdown are called right after each other.
		//
		// Regardless, if there is no transport we should not be inside this
		// loop
		if n.transport == nil {
			log.Printf("Transport is nil inside Node.listen(), exiting loop")
			return
		}

		select {
		case <-ctx.Done():
			if DebugMode {
				log.Printf("node: listen: context finished, exiting loop")
			}
			return

		case err := <-n.transport.errors:
			n.errs <- fmt.Errorf("transport: %w", err)

		case bytes := <-n.transport.responses:
			result := &container{
				content: bytes,
			}

			msg := &response{}
			if err := json.Unmarshal(bytes, msg); err != nil {
				if DebugMode {
					log.Printf("unmarshal received message failed: %v", err)
				}

				result.err = fmt.Errorf("unmarshal received message failed: %v", err)
			} else if msg.Error != nil {
				result.err = errors.New(fmt.Sprint(msg.Error))
			}

			// subscribe message if returned message with 'method' field
			if len(msg.Method) > 0 {
				n.pushHandlersLock.RLock()
				handlers := n.pushHandlers[msg.Method]
				n.pushHandlersLock.RUnlock()

				for _, handler := range handlers {
					select {
					case handler <- result:
					default:
					}
				}
			}

			n.handlersLock.RLock()
			c, ok := n.handlers[msg.Id]
			n.handlersLock.RUnlock()

			if ok {
				c <- result
			}
		}
	}
}

// listenPush returns a channel of messages matching the method.
func (n *Node) listenPush(method string) <-chan *container {
	c := make(chan *container, 1)
	n.pushHandlersLock.Lock()
	n.pushHandlers[method] = append(n.pushHandlers[method], c)
	n.pushHandlersLock.Unlock()

	return c
}

// request makes a request to the server and unmarshals the response into v.
func (n *Node) request(ctx context.Context, method string, params []interface{}, v interface{}) error {
	select {
	case <-n.quit:
		return ErrNodeShutdown
	default:
	}

	msg := request{
		Id:     atomic.AddUint64(&n.nextId, 1),
		Method: method,
		Params: params,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	bytes = append(bytes, delim)
	if err := n.transport.SendMessage(ctx, bytes); err != nil {
		return err
	}

	c := make(chan *container, 1)

	n.handlersLock.Lock()
	n.handlers[msg.Id] = c
	n.handlersLock.Unlock()

	var resp *container
	select {
	case resp = <-c:
	case <-ctx.Done():
		return ctx.Err()
	}

	if resp.err != nil {
		return resp.err
	}

	n.handlersLock.Lock()
	delete(n.handlers, msg.Id)
	n.handlersLock.Unlock()

	if v != nil {
		err = json.Unmarshal(resp.content, v)
		if err != nil {
			return err
		}
	}

	return nil
}

// Disconnect shuts down the node. It is safe to call multiple times. Subsequent
// calls to the first one are no-ops.
// TODO: implement support for draining requests and waiting until everything is
// finished.
func (n *Node) Disconnect() {
	select {
	// Already called! Make it a no-op
	case <-n.quit:
		return

	default:
	}

	if n.transport == nil {
		log.Printf("WARNING: disconnecting node before transport is set up")
		return
	}

	if DebugMode {
		log.Printf("disconnecting node")
	}

	close(n.quit)

	n.transport.conn.Close()

	n.handlers = nil
	n.pushHandlers = nil
}
