package main

import (
	"net"
	"fmt"
	"os"
	"io/ioutil"
)

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}

func main() {
	 // init
	url := "cloudflare-workers.akhtariev.workers.dev:80"

	//  tcpAddr, err := net.ResolveTCPAddr("tcp4", url)
	//  checkError(err)
	 conn, err := net.Dial("tcp", url)
	 checkError(err)
  
	 // send message
	 _, err = conn.Write([]byte("GET / HTTP/1.1\r\n\r\n"))
	  checkError(err)
  
	 // receive message
	 result, err := ioutil.ReadAll(conn)
	 checkError(err)
	 fmt.Println(string(result))
     os.Exit(0)
}
