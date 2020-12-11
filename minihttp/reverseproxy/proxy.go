package reverseproxy

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/url"
	"path"

	"github.com/tutamuniz/yasrp/minihttp"
	"github.com/tutamuniz/yasrp/minihttp/reverseproxy/configtypes"
	"github.com/tutamuniz/yasrp/miniutils/config"
)

type LocationMap map[string]configtypes.Location

type ReverseProxy struct {
	BindIP    string
	BindPort  uint16
	Locations LocationMap
}

func NewReverseProxy(BindIP string, BindPort uint16) (*ReverseProxy, error) {
	if net.ParseIP(BindIP) == nil {
		return nil, fmt.Errorf("Invalid IP Address:%s", BindIP)
	}
	revp := &ReverseProxy{BindIP: BindIP, BindPort: BindPort, Locations: make(LocationMap)}

	return revp, nil
}

func NewReverseProxyFromConfig(cfg config.Config) (*ReverseProxy, error) {
	rp, err := NewReverseProxy(cfg.BindIP, cfg.BindPort)
	if err != nil {
		return nil, err
	}
	for _, loc := range cfg.Locations {
		if err = rp.AddLocation(loc); err != nil {
			return nil, err
		}
	}

	return rp, nil
}

func (rp *ReverseProxy) Listen() {

	bindAddr := fmt.Sprintf("%s:%d", rp.BindIP, rp.BindPort)

	addr, err := net.ResolveTCPAddr("tcp", bindAddr)
	if err != nil {
		log.Fatalf("Listen(): Error resolving address %s", err.Error())
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("Listen(): %s", err.Error())
	}

	log.Printf("Starting YASRP...\n")
	log.Printf("Listening on %s\n", bindAddr)
	log.Print("List of Locations:\n")
	for _, loc := range rp.Locations {
		log.Printf("\t%s => %s \n", loc.Path, loc.Target)
	}

	log.Printf("Accepting connections.\n")
	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			log.Printf("Listen(): Error on Accept() %s", err.Error())
		}

		go rp.ConnectionHandler(conn)
	}

}

func (rp *ReverseProxy) AddLocation(loc configtypes.Location) error {
	path := removeLastSlash(loc.Path)
	if _, ok := rp.Locations[path]; ok {
		return fmt.Errorf("Location already exists : %s", loc.Path)
	}

	// create some better Target validation?!
	u, err := url.Parse(loc.Target)

	if err != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return fmt.Errorf("Invalid Target: %s", loc.Target)
	}

	rp.Locations[loc.Path] = loc

	return nil
}

func removeLastSlash(path string) string { // TODO - improve this function
	lastIdx := len(path) - 1
	if len(path) > 1 {
		for i := lastIdx; i > 1; i-- {
			if path[lastIdx] == '/' {
				lastIdx--
			} else {
				break
			}
		}
	}
	return path[:lastIdx+1]
}

func (rp *ReverseProxy) ConnectionHandler(conn *net.TCPConn) {
	defer conn.Close()

	log.Printf("Connection from: %s\n", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)

	req, err := minihttp.ParseRequest(reader)

	if err != nil {
		log.Printf("%s\n", err.Error())
		return
	}
	dir, _ := path.Split(req.URI)

	//rw := bufio.NewReadWriter(reader, writer)
	if loc, ok := rp.Locations[dir]; ok {

		writer := bufio.NewWriter(conn)
		resp, err := processLocation(req, loc)
		if err != nil {
			log.Printf("Error processing location %s", dir)
			return
		}
		_, _ = writer.Write(resp)
		writer.Flush()
	}

}

func processLocation(r *minihttp.Request, loc configtypes.Location) ([]byte, error) {

	log.Printf("Processing location %s\n", loc.Path)

	u, _ := url.Parse(loc.Target)

	hostname := u.Hostname()
	port := "80" // Remove magic numbers

	if u.Port() != "" {
		port = u.Port()
	}

	log.Printf("Connection on target %s:%s\n", hostname, port)

	client, err := net.Dial("tcp", hostname+":"+port)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	r.Headers["Host"] = u.Hostname()
	r.URI = u.Path

	_, err = client.Write(r.ToBytes())

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	reader := bufio.NewReader(client)
	resp, err := minihttp.ParseResponse(reader)

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	return resp.ToBytes(), nil
}
