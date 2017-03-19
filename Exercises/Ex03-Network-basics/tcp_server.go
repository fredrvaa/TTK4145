package main
 
import (
    "fmt"
    "net"
    "os"
)

func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

func main() {
    ServerConn, err := net.Listen("tcp", ":20009")
    CheckError(err)
    defer ServerConn.Close()
 
    for {
	// Accept new connection
	c,err := ServerConn.Accept()
	CheckError(err)

	go handleConnection(c)
    }
}

func handleConnection(c net.Conn) {
    
    sendMessage(c, "Welcome!")
    for {
	buf := make([]byte, 1024)
	n,err := c.Read(buf)
	if err != nil {
	    break
	}

	fmt.Println("Received ",string(buf[0:n]), " from ",c.LocalAddr())

	if string(buf[0:n]) == "hei" {
	    sendMessage(c, "Hei igjen!")
	} else {
	    sendMessage(c, "Unknown command")
	}
    }
}

func sendMessage(c net.Conn, msg string) {
    buf := []byte(msg)
    _,err := c.Write(buf)
    CheckError(err)
}

