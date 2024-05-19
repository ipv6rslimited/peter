# peter

**peter** is a piper, a simple bidirectional proxy that facilitates data transfer between two network connections.
It reads data from one connection and writes it to the other, and vice versa, until the connections are closed.

## Features
- Bidirectional data transfer between two network connections.
- Closes connections after data transfer is complete.

## Usage

To create a new proxy instance, use the NewPeter function by passing the client and backend connections as arguments.
```
clientConn, err := net.Dial("tcp", "client-address:port")
if err != nil {
  log.Fatal(err)
}

backendConn, err := net.Dial("tcp", "backend-address:port")
if err != nil {
  log.Fatal(err)
}

peter := NewPeter(clientConn, backendConn)
```
To start the data transfer, call the Start method on the Peter instance.
```
peter.Start()
```
This method blocks until both connections are closed, so run it in a go subroutine if you are doing multiple connections.

## License

Distributed under the COOLER License.

Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
All Rights Reserved

