package main
 
import (
    "fmt"
    "net"
    "time"
)
 
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}
 
func main() {
    Conn, err := net.Dial("udp", ":20009")
    CheckError(err)
 
    defer Conn.Close()
    for {
        msg := "hei"
        buf := []byte(msg)
        _,err := Conn.Write(buf)
        if err != nil {
            fmt.Println(msg, err)
        }
        time.Sleep(time.Second)
    }
}
