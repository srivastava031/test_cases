package main

import (
	"crypto/tls"
	_"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"encoding/pem"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"software.sslmate.com/src/go-pkcs12"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {

	var b = make([]byte, 2024)
	req.Body.Read(b)
	defer req.Body.Close()

	fmt.Println("..", string(b))

	fmt.Println("recieved on handler")
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, string(b))
	fmt.Println("executed handler")

}
// func loadCertificate() (tls.Certificate, error) {
// 		cert := tls.Certificate{}
// 		b, _ := ioutil.ReadFile("/home/asus/RakeshMail/keystore.p12")
// 		// key, cer, err := pkcs12.Decode(b, "abcdefgh")
// 		// if err != nil {
// 		// 	return cert, err
// 		// }
// 		// cert.PrivateKey = key
// 		// cert.Certificate = [][]byte{cer.Raw}
// 		// return cert, nil
// }

func main() {
	fmt.Println("server running...")

	
	b, _ := ioutil.ReadFile("/home/asus/RakeshMail/keystore.p12")
	blocks, _ := pkcs12.ToPEM(b, "abcdefgh")

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, _ := tls.X509KeyPair(pemData, pemData)

	//cert, _ := tls.LoadX509KeyPair("/home/asus/RakeshMail/serv.cer", "/home/asus/RakeshMail/serv.key")

	// caCertPEM1, err := ioutil.ReadFile("/home/asus/Desktop/xx/clientTLSnonTLS/CA/cliCertPEM01.crt")
	// if err != nil {
	// 	fmt.Println("CA error")
	// }

	// roots := x509.NewCertPool()
	// ok := roots.AppendCertsFromPEM(caCertPEM1)
	// if !ok {
	// 	fmt.Println("x509.NewCertPool")
	// }

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloServer)

	cfg := &tls.Config{
		// ClientAuth: tls.RequireAndVerifyClientCert, //this will check the authentication of certificate, if it is signed by
		// ClientCAs:  roots,
		//InsecureSkipVerify: true,
		Certificates: []tls.Certificate{cert},
	}
	h2server := &http2.Server{
		MaxConcurrentStreams: 3,
	}
	srv := &http.Server{
		Addr:         "localhost:6514",
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		Handler:      h2c.NewHandler(mux, h2server), //here router is nothing but serveMux which itself is Handler,
		// as it is having method serveHttp attached to it, will server bot http1.1 and http2 without tls
		IdleTimeout: 100 * time.Second,
	}

	err := srv.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
