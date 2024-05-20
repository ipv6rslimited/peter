/*
**
** peter
** A piper
**
** Distributed under the COOLER License.
**
** Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
** All Rights Reserved
**
*/

package peter

import (
  "io"
  "net"
  "sync"
)

type Peter struct {
  client  net.Conn
  backend net.Conn
}

func NewPeter(client, backend net.Conn) *Peter {
  return &Peter{
    client:  client,
    backend: backend,
  }
}

func (p *Peter) Start() {
  var wg sync.WaitGroup

  copyConn := func(dst net.Conn, src net.Conn) {
    defer wg.Done()
    io.Copy(dst, src)
  }

  wg.Add(2)
  go func() {
    copyConn(p.backend, p.client)
  }()
  go func() {
    copyConn(p.client, p.backend)
  }()

  wg.Wait()
}
