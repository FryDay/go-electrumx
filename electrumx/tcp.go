package electrumx

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var DebugMode bool

type transport struct {
	conn      net.Conn
	responses chan []byte
	errors    chan error
}

func newConn(ctx context.Context, addr string, tlsConfig *tls.Config) (net.Conn, error) {
	if tlsConfig != nil {
		var d tls.Dialer
		d.Config = tlsConfig
		return d.DialContext(ctx, "tcp", addr)
	}
	var d net.Dialer
	return d.DialContext(ctx, "tcp", addr)
}

func newTransport(ctx context.Context, addr string, sslConfig *tls.Config) (*transport, error) {
	conn, err := newConn(ctx, addr, sslConfig)
	if err != nil {
		return nil, err
	}

	t := &transport{
		conn:      conn,
		responses: make(chan []byte),
		errors:    make(chan error),
	}

	return t, nil
}

func (t *transport) SendMessage(ctx context.Context, body []byte) error {
	if DebugMode {
		log.Printf("[debug] %s <- %s", t.conn.RemoteAddr(), body)
	}

	done := make(chan struct{})
	errs := make(chan error)
	go func() {
		if _, err := t.conn.Write(body); err != nil {
			errs <- err
			return
		}
		close(done)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("send message: %w", ctx.Err())

	case err := <-errs:
		return fmt.Errorf("send message: %w", err)

	case <-done:
		return nil
	}
}

func (t *transport) listen(ctx context.Context) {
	reader := bufio.NewReader(t.conn)
	responses := make(chan []byte)
	errs := make(chan error)
	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			// The Node should send server.ping request continuously with a
			// reasonable break in order to keep connection alive. If not
			// client will receive a disconnection error and encounter
			// io.EOF
			line, err := reader.ReadBytes(delim)
			if err != nil {
				// Check if we're done. Going to get all sorts of
				// errors once that happens.
				select {
				case <-done:
					return

				default:
				}

				// block until start handle error
				if DebugMode {
					log.Printf("transport encountered error: %s", err)
				}

				// The server closed the connection... This happens if we
				// send unsupported requests to it.
				if err == io.EOF {
					err = errors.New("server closed connection (potentially because we sent an unsupported request)")
				}
				errs <- err
				break
			}
			if DebugMode {
				log.Printf("[debug] %s -> %s", t.conn.RemoteAddr(), line)
			}

			responses <- line
		}
	}()

	for {
		select {
		case <-ctx.Done():
			if DebugMode {
				log.Printf("transport: listen: context finished, exiting loop")
			}
			return

		case err := <-errs:
			t.errors <- err
		case res := <-responses:
			t.responses <- res
		}
	}
}
