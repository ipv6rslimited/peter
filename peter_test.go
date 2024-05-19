/*
**
** peter_test
** Putting peter to the test
**
** Distributed under the COOL License.
**
** Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
** All Rights Reserved
**
*/

package peter

import (
  "bytes"
  "io"
  "net"
  "testing"
  "time"
)

type mockConn struct {
  io.Reader
  io.Writer
  closed bool
}

func (mc *mockConn) Close() error {
  mc.closed = true
  return nil
}

func (mc *mockConn) LocalAddr() net.Addr {
  return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080}
}

func (mc *mockConn) RemoteAddr() net.Addr {
  return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8081}
}

func (mc *mockConn) SetDeadline(t time.Time) error {
  return nil
}

func (mc *mockConn) SetReadDeadline(t time.Time) error {
  return nil
}

func (mc *mockConn) SetWriteDeadline(t time.Time) error {
  return nil
}

func TestPeterStart(t *testing.T) {
  clientConn := &mockConn{Reader: bytes.NewReader([]byte("request")), Writer: new(bytes.Buffer)}
  backendConn := &mockConn{Reader: bytes.NewReader([]byte("response")), Writer: new(bytes.Buffer)}

  peter := NewPeter(clientConn, backendConn)
  peter.Start()

  if !clientConn.closed || !backendConn.closed {
    t.Fatalf("Connections should be closed")
  }

  response := make([]byte, 8)
  clientConn.Writer.(*bytes.Buffer).Read(response)
  if string(response) != "response" {
    t.Fatalf("Expected 'response', got '%s'", string(response))
  }

  request := make([]byte, 7)
  backendConn.Writer.(*bytes.Buffer).Read(request)
  if string(request) != "request" {
    t.Fatalf("Expected 'request', got '%s'", string(request))
  }
}
