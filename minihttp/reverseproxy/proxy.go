package reverseproxy

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/tutamuniz/yasrp/cacheengine"
	"github.com/tutamuniz/yasrp/minihttp"
	"github.com/tutamuniz/yasrp/minihttp/cache"
	"github.com/tutamuniz/yasrp/minihttp/reverseproxy/configtypes"
	"github.com/tutamuniz/yasrp/miniutils/config"
)

type LocationMap map[string]configtypes.Location

type ReverseProxy struct {
	BindIP      string
	BindPort    uint16
	EnableCache bool
	CacheEngine cache.Cache
	Locations   LocationMap
}

func NewReverseProxy(BindIP string, BindPort uint16) (*ReverseProxy, error) {
	if net.ParseIP(BindIP) == nil {
		return nil, fmt.Errorf("Invalid IP Address:%s", BindIP)
	}
	revp := &ReverseProxy{BindIP: BindIP, BindPort: BindPort, Locations: make(LocationMap), EnableCache: false}

	return revp, nil
}

func NewReverseProxyFromConfig(cfg config.Config) (*ReverseProxy, error) {
	rp, err := NewReverseProxy(cfg.BindIP, cfg.BindPort)

	if err != nil {
		return nil, err
	}

	err = rp.SetCacheEngine(cfg.CacheEngine)

	if err != nil {
		return nil, err
	}

	rp.SetCacheStatus(cfg.EnableCache)

	for _, loc := range cfg.Locations {
		if err = rp.AddLocation(loc); err != nil {
			return nil, err
		}
	}

	return rp, nil
}

func (rp *ReverseProxy) SetCacheStatus(enable bool) error {
	if rp.CacheEngine == nil {
		return fmt.Errorf("CacheEngine not defined")
	}

	rp.EnableCache = enable

	return nil
}

func (rp *ReverseProxy) SetCacheEngine(engine string) error {
	eng, err := cacheengine.NewCacheEngine(engine)
	if err != nil {
		return err
	}

	rp.CacheEngine = *eng
	return nil
}

func (rp *ReverseProxy) Listen() {

	log.Printf("Starting YASRP...\n")

	bindAddr := fmt.Sprintf("%s:%d", rp.BindIP, rp.BindPort)

	addr, err := net.ResolveTCPAddr("tcp", bindAddr)
	if err != nil {
		log.Fatalf("Listen(): Error resolving address %s", err.Error())
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("Listen(): %s", err.Error())
	}
	log.Printf("Listening on %s\n", bindAddr)

	if rp.EnableCache {
		go rp.CacheEngine.StartEngine()
	}

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

	log.Printf("Request %s from %s\n", req.URI, conn.RemoteAddr().String())

	if err != nil {
		log.Printf("%s\n", err.Error())
		return
	}

	//	dir, _ := path.Split(req.URI)

	//rw := bufio.NewReadWriter(reader, writer)

	writer := bufio.NewWriter(conn)

	resp, err := rp.processLocation(req)
	if err != nil {
		log.Printf("Error processing location %s: %s", req.URI, err.Error())
	}
	_, _ = writer.Write(resp)
	writer.Flush()

}

func (rp *ReverseProxy) matchLocation(reqPath string) (string, *configtypes.Location) {
	for p, loc := range rp.Locations {
		if strings.HasPrefix(reqPath, p) {
			target, _ := url.Parse(loc.Target)
			newpath := strings.Replace(reqPath, loc.Path, target.Path, 1)
			return newpath, &loc
		}
	}
	return "", nil
}

// Process proxy request based on location mapping
func (rp *ReverseProxy) processLocation(r *minihttp.Request) ([]byte, error) {
	var hostname, port, scheme string

	urlPath := r.URI

	log.Printf("Processing location %s\n", urlPath)

	if rp.EnableCache {
		if rp.CacheEngine.InCache(urlPath) {
			e, err := rp.CacheEngine.Get(urlPath)
			if err != nil {
				log.Printf("Cache Error: %s\n", err.Error())
			} else {
				log.Printf("Cache HIT for %s\n", urlPath)
				return e.Resp.ToBytes(), nil
			}
		}
		log.Printf("Cache MISS for %s\n", urlPath)
	}

	newpath, loc := rp.matchLocation(r.URI)

	if newpath == "" {
		message := fmt.Sprintf("Resource <strong>%s</strong> not found.", r.URI)
		notfound, _ := minihttp.NewResponse(404, r.Headers, []byte(message))

		return notfound.ToBytes(), fmt.Errorf("Match Not found")
	}

	// overwrite the requested path with target's path.
	r.URI = newpath

	u, _ := url.Parse(loc.Target)

	hostname = u.Hostname()
	port = minihttp.DefaultHTTPPort

	if u.Scheme != "" {
		scheme = u.Scheme
	}

	if u.Port() != "" {
		port = u.Port()
	}

	log.Printf("Connecting target %s:%s\n", hostname, port)

	client, err := connectBackend(scheme, hostname, port)

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	// Set the target hostname
	r.Headers["Host"] = u.Hostname()

	_, err = client.Write(r.ToBytes())

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	reader := bufio.NewReader(client)
	resp, err := minihttp.ParseResponse(reader)

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	// Little bit confuse code with urlPath and r.URI. Solve it later

	if rp.EnableCache {
		if !rp.CacheEngine.InCache(urlPath) {
			log.Printf("Cache creating entry for %s\n", urlPath)
			rp.CacheEngine.Put(urlPath, cache.MakeCacheEntry(r, resp))
		}
	}
	return resp.ToBytes(), nil
}

// abstraction for backend connections
func connectBackend(scheme, host, port string) (net.Conn, error) {
	var conn net.Conn
	var err error

	if strings.ToLower(scheme) == "https" {
		config := tls.Config{InsecureSkipVerify: true}
		conn, err = tls.Dial("tcp", host+":"+port, &config)
	} else {
		conn, err = net.Dial("tcp", host+":"+port)
	}

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	return conn, nil

}
