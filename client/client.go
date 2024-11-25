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
    addr := fmt.Sprintf("%s:%d", ip, port)
    udpAddr, err := net.ResolveUDPAddr("udp", addr)
    if err != nil {
        return Client{}, err
    }

    conn, err := net.DialUDP("udp", nil, udpAddr)
    if err != nil {
        return Client{}, err
    }

    return Client{conn}, nil
}

func (c Client) Send(data any) error {
    bytes, err := json.Marshal(data)
    if err != nil {
        return err
    }

    _, err = c.conn.Write(bytes)
    if err != nil {
        return err
    }

    return nil
}

func (c Client) Close() {
    c.conn.Close()
}
