package main

import (
	"crypto/tls"
	"flag"
	"log"
	"os"
	"net/http"
	"net/http/httputil"
	"net/url"
	"golang.org/x/crypto/acme/autocert"
)

func main () {

	email := flag.String("email", os.Getenv("PROXY_EMAIL"), "email address for let's encrypt account")
	listen := flag.String("listen", os.Getenv("PROXY_LISTEN"), "bind address e.g. 0.0.0.0 defaults to ::")
	destination := flag.String("destination", os.Getenv("PROXY_DESTINATION"), "proxy destination; defaults to http://localhost:8080")
	nohttp := flag.Bool("nohttp", false, "do not listen on http port 80; defaults to on")
	certdir := flag.String("certdir", "/etc/certs", "directory to store certificates; defaults to /etc/certs")

	flag.Parse()

	if *listen == "" {
		*listen = ":"
	}

	if *destination == "" {
		*destination = "http://localhost:8080"
	}

	log.Printf("Certificate directory is %+v ", *certdir)

	domains := flag.Args()
	if len(domains) == 0 {
		log.Fatalf("fatal; specify domains for tls certificate e.g. www.domain.com test.domain.com domain.com")
	}

	log.Printf("Serving http/https for domains %+v", domains)

	u, err := url.Parse(*destination)
	if err != nil {
		log.Fatal("error parsing destination")
	}

	log.Printf("Proxy destination %+v", *destination)

	director := func(req *http.Request) {
		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}
	http.Handle("/", proxy)

	certManager := autocert.Manager{
		Prompt:	autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domains...),
		Cache: autocert.DirCache(*certdir),
		Email: *email,
	}

	server := &http.Server {
		Addr: *listen + "https",
		TLSConfig: &tls.Config {
			GetCertificate: certManager.GetCertificate,
		},
	}

	if !*nohttp {
		log.Printf("Listening on %+v", *listen + ":http")
		go func() {
			h := certManager.HTTPHandler(nil)
			log.Fatal(http.ListenAndServe(*listen + "http", h))
		}()
	}

	log.Printf("Listening on %+v", *listen + ":https")
	log.Fatal(server.ListenAndServeTLS("", ""))
}
