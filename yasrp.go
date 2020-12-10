package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/tutamuniz/yasrp/minihttp"
)

func main() {
	fmt.Println("Starting...")
	Listener("127.0.0.1:8080")
}

func Listener(fulladdr string) {
	laddr, err := net.ResolveTCPAddr("tcp", fulladdr)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		go ConnHandler(conn)
	}
}

// ADD X-FORWARD-FOR
func ConnHandler(conn *net.TCPConn) {
	defer conn.Close()
	//  Work with conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	fmt.Printf("Client: %s \n", conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)

	req, _ := minihttp.ParseRequest(reader)

	fmt.Println(req)

	res, _ := getRemote(req, "www.tjrn.jus.br:80")
	wr := bufio.NewWriter(conn)
	_, _ = wr.Write(res)
	wr.Flush()
}

func getRemote(r *minihttp.Request, addr string) ([]byte, error) {

	client, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	r.Headers["Host"] = "www.tjrn.jus.br"

	client.Write(r.ToBytes())

	reader := bufio.NewReader(client)
	resp, _ := minihttp.ParseResponse(reader)
	return resp.ToBytes(), nil
}
