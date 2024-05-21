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
  t.Run("BasicFunctionality", func(t *testing.T) {
    clientConn := &mockConn{Reader: bytes.NewReader([]byte("request")), Writer: new(bytes.Buffer)}
    backendConn := &mockConn{Reader: bytes.NewReader([]byte("response")), Writer: new(bytes.Buffer)}

    peter := NewPeter(clientConn, backendConn)
    go peter.Start()

    time.Sleep(50 * time.Millisecond)

    response := make([]byte, 8)
    n, err := clientConn.Writer.(*bytes.Buffer).Read(response)
    if err != nil && err != io.EOF {
      t.Fatalf("Failed to read response: %v", err)
    }
    if string(response[:n]) != "response" {
      t.Fatalf("Expected 'response', got '%s'", string(response[:n]))
    }

    request := make([]byte, 7)
    n, err = backendConn.Writer.(*bytes.Buffer).Read(request)
    if err != nil && err != io.EOF {
      t.Fatalf("Failed to read request: %v", err)
    }
    if string(request[:n]) != "request" {
      t.Fatalf("Expected 'request', got '%s'", string(request[:n]))
    }

    clientConn.Close()
    backendConn.Close()
  })

  t.Run("LargeDataTransfer", func(t *testing.T) {
    largeData := make([]byte, 1024*1024)
    for i := range largeData {
      largeData[i] = byte(i % 256)
    }

    clientConn := &mockConn{Reader: bytes.NewReader(largeData), Writer: new(bytes.Buffer)}
    backendConn := &mockConn{Reader: bytes.NewReader(largeData), Writer: new(bytes.Buffer)}

    peter := NewPeter(clientConn, backendConn)
    go peter.Start()

    time.Sleep(100 * time.Millisecond)

    clientBuffer := clientConn.Writer.(*bytes.Buffer)
    if clientBuffer.Len() != len(largeData) {
      t.Fatalf("Expected client buffer length %d, got %d", len(largeData), clientBuffer.Len())
    }

    backendBuffer := backendConn.Writer.(*bytes.Buffer)
    if backendBuffer.Len() != len(largeData) {
      t.Fatalf("Expected backend buffer length %d, got %d", len(largeData), backendBuffer.Len())
    }

    clientConn.Close()
    backendConn.Close()
  })

  t.Run("ConnectionInterruption", func(t *testing.T) {
    interruptingReader := &io.LimitedReader{R: bytes.NewReader([]byte("request")), N: 4}
    clientConn := &mockConn{Reader: interruptingReader, Writer: new(bytes.Buffer)}
    backendConn := &mockConn{Reader: bytes.NewReader([]byte("response")), Writer: new(bytes.Buffer)}

    peter := NewPeter(clientConn, backendConn)
    go peter.Start()

    time.Sleep(100 * time.Millisecond)

    backendBuffer := backendConn.Writer.(*bytes.Buffer)
    if backendBuffer.Len() != 4 {
      t.Fatalf("Expected backend buffer length %d, got %d", 4, backendBuffer.Len())
    }

    clientConn.Close()
    backendConn.Close()
  })

  t.Run("PartialReadWrite", func(t *testing.T) {
    clientConn := &mockConn{Reader: bytes.NewReader([]byte("request")), Writer: new(bytes.Buffer)}
    backendConn := &mockConn{
      Reader: io.MultiReader(bytes.NewReader([]byte("resp")), bytes.NewReader([]byte("onse"))),
      Writer: new(bytes.Buffer),
    }

    peter := NewPeter(clientConn, backendConn)
    go peter.Start()

    time.Sleep(100 * time.Millisecond)

    response := make([]byte, 8)
    n, err := clientConn.Writer.(*bytes.Buffer).Read(response)
    if err != nil && err != io.EOF {
      t.Fatalf("Failed to read response: %v", err)
    }
    if string(response[:n]) != "response" {
      t.Fatalf("Expected 'response', got '%s'", string(response[:n]))
    }

    request := make([]byte, 7)
    n, err = backendConn.Writer.(*bytes.Buffer).Read(request)
    if err != nil && err != io.EOF {
      t.Fatalf("Failed to read request: %v", err)
    }
    if string(request[:n]) != "request" {
      t.Fatalf("Expected 'request', got '%s'", string(request[:n]))
    }

    clientConn.Close()
    backendConn.Close()
  })

  t.Run("TimeoutHandling", func(t *testing.T) {
    clientConn := &mockConn{Reader: bytes.NewReader([]byte("request")), Writer: new(bytes.Buffer)}
    backendConn := &mockConn{Reader: bytes.NewReader([]byte("response")), Writer: new(bytes.Buffer)}

    time.AfterFunc(10*time.Millisecond, func() {
      clientConn.Close()
      backendConn.Close()
    })

    peter := NewPeter(clientConn, backendConn)
    go peter.Start()

    time.Sleep(50 * time.Millisecond)

    if !clientConn.closed {
      t.Fatalf("Expected client connection to be closed")
    }
    if !backendConn.closed {
      t.Fatalf("Expected backend connection to be closed")
    }
  })
}
