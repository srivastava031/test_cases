package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"

	"time"

	"golang.org/x/net/http2"
)

var ClientArr []http.Client

func CreateClientTLSConfig(caC, caK, ServerName string, SkipVerify, enableTls bool, minV, maxV uint16) (*tls.Config, error) {
	var roots *x509.CertPool
	if enableTls == false {
		return nil, nil
	}
	if (caC == "") || (caK == "") {
		roots = nil
	} else {
		caCertPEM, err := ioutil.ReadFile(caC)
		if err != nil {
			return nil, err
		}
		roots = x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(caCertPEM)
		if !ok {
			panic("failed to parse root certificate")
		}
	}
	clcrt, clkey := makeCertificate(caC, caK)

	cert, err := tls.X509KeyPair(clcrt.Bytes(), clkey.Bytes()) //if only CA and selfsigned Cert present
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		InsecureSkipVerify: SkipVerify,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            roots, //if CA certificate is present, else nil
		ServerName:         "localhost",
		MaxVersion:         maxV,
		MinVersion:         minV,
	}, nil
}

func CreateH2Client(allowHttp, tlsEnable bool, cfg *tls.Config) (http.Client, error) {
	// todo take values form json config file

	localAddr, err := net.ResolveIPAddr("ip", "127.0.0.6")
	if err != nil {
		fmt.Println(err)
	}
	localTCPAddr := net.TCPAddr{
		IP: localAddr.IP,
		//Port:""
	}
	dialer := &net.Dialer{
		LocalAddr: &localTCPAddr,
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	client := http.Client{ // http2 will use its parameter for client configuration but will not use it as transport
		Transport: &http.Transport{
			Dial: dialer.Dial,
			//TLSClientConfig: tlsConfig,
		},
	}
	if tlsEnable == false {
		client.Transport = &http2.Transport{
			//TLSClientConfig: cfg,
			AllowHTTP: allowHttp, //allow using plain http1.1/ without tls,but does not enable h2c// if tls is enabled then it cant be https/2
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) { // will pretend we are dialing tls endpoint,
				//so will give h2c feature->no tls req.
				return net.Dial(network, addr) //this is http2 without tls,, if this field is not present then it would be http2 with tls.
				// 	//TLSClientConfig:tlsConfig,

			},
		}
	} else {
		client.Transport = &http2.Transport{
			TLSClientConfig: cfg,
			AllowHTTP:       allowHttp, //allow using plain http1.1/ without tls,but does not enable h2c// if tls is enabled then it cant be https/2
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) { // will pretend we are dialing tls endpoint,so will give h2c feature->no tls req.
				return tls.Dial("tcp", "localhost:8080", cfg) //this is http2 without tls,, if this field is not present then it would be http2 with tls.
				// 	//TLSClientConfig:tlsConfig,

			},
		}
	}
	return client, nil
}

//==========================================================================================================================================================================
//===========================================================================================================================================================================
//Client contains Client instance
// var Client []http.Client

func getFromJSON(value string) interface{} {
	configfile := "/home/asus/go/src/clientInstances/config.JSON"
	var m map[string]interface{}
	content, err := ioutil.ReadFile(configfile)
	if err != nil {
		fmt.Sprintln(err)
	}
	json.Unmarshal(content, &m)
	val := m[value]
	return val
}

func CreateConnInstances(testConfigFile string) error {
	var enableTls = []bool{}
	var insecureSkipVerify = []bool{}
	//	var tlsClientCer = []string{}

	ConnCount := int((getFromJSON("ConnCount")).(float64))
	allowHttp := (getFromJSON("AllowHTTP")).(bool)
	tlsEnable := (getFromJSON("tlsEnable")).([]interface{})
	for k, _ := range tlsEnable {
		i := 0
		ii := tlsEnable[k].(bool)
		enableTls = append(enableTls, ii)
		i++
	}
	InsecureSkipVerify := (getFromJSON("InsecureSkipVerify")).([]interface{})
	for k, _ := range InsecureSkipVerify {
		i := 0
		ii := InsecureSkipVerify[k].(bool)
		insecureSkipVerify = append(insecureSkipVerify, ii)
		i++
	}
	// tlsClientCertificate := (getFromJSON("tlsClientCertificate")).(string)
	// tlsClientKey := (getFromJSON("tlsClientKey")).(string)
	serverName := (getFromJSON("servername")).(string)
	caC := (getFromJSON("CAcertificate")).(string)
	caK := (getFromJSON("CAkey")).(string)
	minV := uint16((getFromJSON("MinVersion")).(float64))
	maxV := uint16((getFromJSON("MaxVersion")).(float64))

	for count := 0; count < ConnCount; count++ {
		tlsConfig, _ := CreateClientTLSConfig(caC, caK, serverName, insecureSkipVerify[count], enableTls[count], minV, maxV)

		client, err := CreateH2Client(allowHttp, enableTls[count], tlsConfig)
		if err != nil {
			return err
		}
		ClientArr = append(ClientArr, client)
	}
	return nil
}

//SendRequest sends the request to the URL given
func SendRequest(client http.Client, method string, url string, header map[string]string, body []byte) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {

		fmt.Println(err)
	}

	if header != nil {
		for key, value := range header {
			req.Header[key] = []string{value}
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*1000*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	response, err := client.Do(req)
	if err != nil {

		fmt.Println(err)
	}

	fmt.Println("====>>>", response)
	//fmt.Println("====>>>", response.StatusCode)

}
func main() {
	CreateConnInstances("/home/asus/go/src/clientInstances/config.JSON")

	for _, v := range ClientArr {
		SendRequest(v, "GET", "https://localhost:8080/hello", nil, nil)
	}

}

//===================================================================================================
//>>>>>>>>>>>>>>>>>>><<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>><<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>><<

func makeCertificate(caCer, caKey string) (bytes.Buffer, bytes.Buffer) {
	var clcrt bytes.Buffer
	var clkey bytes.Buffer
	var ca *x509.Certificate
	var catls tls.Certificate
	var err error

	if (caCer != "") || (caKey != "") {

		catls, err = tls.LoadX509KeyPair(caCer, caKey)
		if err != nil {
			panic(err)
		}
		ca, err = x509.ParseCertificate(catls.Certificate[0])
		if err != nil {
			panic(err)
		}
	}

	// Prepare certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"TV"},
			Country:       []string{"IN"},
			Province:      []string{"KT"},
			Locality:      []string{"BE"},
			StreetAddress: []string{"ADDRESS"},
			PostalCode:    []string{"POSTAL_CODE"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey

	// Sign the certificate
	if (caCer == "") || (caKey == "") {

		certb, _ := x509.CreateCertificate(rand.Reader, cert, cert, pub, priv)
		pem.Encode(&clcrt, &pem.Block{Type: "CERTIFICATE", Bytes: certb})
		pem.Encode(&clkey, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

		return clcrt, clkey

	}

	certb, err := x509.CreateCertificate(rand.Reader, cert, ca, pub, catls.PrivateKey)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	pem.Encode(&clcrt, &pem.Block{Type: "CERTIFICATE", Bytes: certb})
	pem.Encode(&clkey, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	return clcrt, clkey

}
