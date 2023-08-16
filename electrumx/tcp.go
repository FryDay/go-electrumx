package electrumx

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
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
	go t.listen()

	return t, nil
}

func (t *transport) SendMessage(ctx context.Context, body []byte) error {
	if DebugMode {
		log.Printf("%s [debug] %s <- %s", time.Now().Format("2006-01-02 15:04:05"), t.conn.RemoteAddr(), body)
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

func (t *transport) listen() {
	defer t.conn.Close()
	reader := bufio.NewReader(t.conn)
	for {
		// The Node should send server.ping request continuously with a
		// reasonable break in order to keep connection alive. If not
		// client will receive a disconnection error and encounter
		// io.EOF with following os.Exit(1).
		line, err := reader.ReadBytes(delim)
		if err != nil {
			// block until start handle error
			t.errors <- err
			break
		}
		if DebugMode {
			log.Printf("%s [debug] %s -> %s", time.Now().Format("2006-01-02 15:04:05"), t.conn.RemoteAddr(), line)
		}

		t.responses <- line
	}
}
