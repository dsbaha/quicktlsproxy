package main

import (
	"crypto/tls"
	"flag"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
)

var (
	email       = flag.String("email", os.Getenv("PROXY_EMAIL"), "email address for let's encrypt account")
	listen      = flag.String("listen", os.Getenv("PROXY_LISTEN"), "bind address e.g. 0.0.0.0 defaults to ::")
	destination = flag.String("destination", os.Getenv("PROXY_DESTINATION"), "proxy destination; defaults to http://localhost:8080")
	nohttp      = flag.Bool("nohttp", false, "do not listen on http port 80; defaults to on")
	certdir     = flag.String("certdir", "/etc/certs", "directory to store certificates; defaults to /etc/certs")
	quiet       = flag.Bool("quiet", false, "disable logging to console.")
	debug       = flag.Bool("debug", false, "console log send/receive messages.")
	domains     []string
)

func checkParams() {
	if *listen == "" {
		*listen = ":"
	}

	if *destination == "" {
		*destination = "http://localhost:8080"
	}

	domains = flag.Args()
	if len(domains) == 0 {
		log.Println("specify domains for tls certificate e.g. www.domain.com test.domain.com domain.com")
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Println("Using certificates stored in", *certdir)
	log.Println("Serving http/https for domains", domains)
	log.Println("Proxy destination", *destination)
}

func main() {
	log.Println("Starting quicktlsproxy")
	flag.Parse()
	checkParams()

	u, err := url.Parse(*destination)
	if err != nil {
		log.Fatalln("error parsing destination", err)
	}

	director := func(req *http.Request) {
		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}
	http.Handle("/", proxy)

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domains...),
		Cache:      autocert.DirCache(*certdir),
		Email:      *email,
	}

	server := &http.Server{
		Addr: *listen + "https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	var wg sync.WaitGroup

	switch *nohttp {
	case false:
		wg.Add(1)
		go startServer("http", &wg, &certManager, nil)
		fallthrough
	default:
		wg.Add(1)
		go startServer("https", &wg, &certManager, server)
	}

	wg.Wait()
	log.Println("Exiting ...")
}

func startServer(proto string, wg *sync.WaitGroup, certManager *autocert.Manager, server *http.Server) {

	defer wg.Done()

	log.Println("Listening on", *listen+proto)

	switch proto {
	case "http":
		h := certManager.HTTPHandler(nil)
		log.Println(http.ListenAndServe(*listen+proto, h))
	case "https":
		log.Println(server.ListenAndServeTLS("", ""))
	default:
		log.Fatalln("bad protocol type")
	}
}
