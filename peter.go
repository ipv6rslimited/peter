/*
**
** peter
** A piper
**
** Distributed under the COOL License.
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
  wg.Add(2)

  go func() {
    defer wg.Done()
    io.Copy(p.client, p.backend)
  }()

  go func() {
    defer wg.Done()
    io.Copy(p.backend, p.client)
  }()

  wg.Wait()
}
