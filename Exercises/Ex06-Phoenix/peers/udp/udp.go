package udp

import (
	"fmt"
	"net"
	"./conn"
)

func checkError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}

type Connection struct {
	PacketConn net.PacketConn
	Addr *net.UDPAddr
}

func Init(port int) Connection {
	packetConn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))
	return Connection{packetConn, addr};
}

func Broadcast(c Connection, msg string) {
	c.PacketConn.WriteTo([]byte(msg), c.Addr)
}

func Receive(c Connection) string {
	buf := make([]byte, 1024)
	n, _, err := c.PacketConn.ReadFrom(buf[0:])
	checkError(err)

	return string(buf[0:n])
}

