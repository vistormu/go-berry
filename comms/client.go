package comms

import (
    "fmt"
    "net"
    "encoding/json"

    "github.com/vistormu/go-berry/errors"
)


type Client struct {
    conn net.Conn
}

func NewUdpClient(ip string, port int) (Client, error) {
    addr := fmt.Sprintf("%s:%d", ip, port)
    udpAddr, err := net.ResolveUDPAddr("udp", addr)
    if err != nil {
        return Client{}, errors.New(errors.CONNECTION, err.Error())
    }

    conn, err := net.DialUDP("udp", nil, udpAddr)
    if err != nil {
        return Client{}, errors.New(errors.CONNECTION, err.Error())
    }

    return Client{conn}, nil
}

func (c Client) Send(data any) error {
    bytes, err := json.Marshal(data)
    if err != nil {
        return errors.New(errors.CLIENT_JSON, err.Error())
    }

    _, err = c.conn.Write(bytes)
    if err != nil {
        return errors.New(errors.CLIENT_SEND, err.Error())
    }

    return nil
}

func (c Client) Close() error {
    err := c.conn.Close()
    if err != nil {
        return errors.New(errors.CLIENT_CLOSE, err.Error())
    }

    return nil
}
