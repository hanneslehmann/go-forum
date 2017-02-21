package main

import (
	"crypto/tls"
	"net"
	"net/http"
)

type embeddedServer struct {

	/*
		  http://blog.davidvassallo.me/2015/06/17/practical-embedding-in-golang/
			custom struct that embeds golang's standard http.Server type
			another way of looking at this is that embeddedServer "inherits" from http.Server,
			though this is not strictly accurate. Have a look at note below for additional information
	*/
	http.Server
	webserverCertificate string
	webserverKey         string
}

func (srv *embeddedServer) ListenAndServeTLS(addr string, handler http.Handler) error {

	/*
		This is where we "hide" or "override" the default "ListenAndServeTLS" method so we modify it to accept
		hardcoded certificates and keys rather than the default filenames
		The default implementation of ListenAndServeTLS was obtained from:
		https://github.com/zenazn/goji/blob/master/graceful/server.go#L33
		and tls.X509KeyPair (http://golang.org/pkg/crypto/tls/#X509KeyPair) is used,
		rather than the default tls.LoadX509KeyPair
	*/

	config := &tls.Config{
		MinVersion: tls.VersionTLS10,
	}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.X509KeyPair([]byte(srv.webserverCertificate), []byte(srv.webserverKey))
	if err != nil {
		return err
	}

	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(conn, config)
	srv.Handler = handler
	return srv.Serve(tlsListener)
}
