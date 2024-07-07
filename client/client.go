package client

import (
    "fmt"
    "net"
    "encoding/json"
)


type Client struct {
    conn net.Conn
}

func New(ip string, port int) (Client, error) {
    // resolve UDP address
    addr := fmt.Sprintf("%s:%d", ip, port)
    udpAddr, err := net.ResolveUDPAddr("udp", addr)
    if err != nil {
        return Client{}, err
    }

    // dial UDP connection
    conn, err := net.DialUDP("udp", nil, udpAddr)
    if err != nil {
        return Client{}, err
    }

    return Client{conn}, nil
}

func (c Client) Send(data any) error {
    // encode data
    bytes, err := json.Marshal(data)
    if err != nil {
        return err
    }

    // send data
    _, err = c.conn.Write(bytes)
    if err != nil {
        return err
    }

    return nil
}

func (c Client) Close() {
    c.conn.Close()
}
