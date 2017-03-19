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
    Conn, err := net.Dial("tcp", "127.0.0.1:20013")
    CheckError(err)
 
    defer Conn.Close()
    for {
	buf := make([]byte, 1024)
	n,err := Conn.Read(buf)
	if err != nil {
	    break
	}

	fmt.Println(string(buf[0:n]))

	fmt.Print("Enter message: ")
	var input string
	fmt.Scanln(&input)

        sendMessage(Conn, input)
    }
}

func sendMessage(c net.Conn, msg string) {
    buf := []byte(msg)
    _,err := c.Write(buf)
    CheckError(err)
}
