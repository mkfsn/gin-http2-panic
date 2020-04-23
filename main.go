package main

// $ openssl version
// LibreSSL 2.8.3
//go:generate openssl genrsa -out cert/private.key 4096
//go:generate openssl req -new -x509 -sha256 -days 1825 -key cert/private.key -subj "/C=TW/ST=Taiwan/CN=localhost:8443" -out cert/public.crt

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

const (
	imageFile = "./images/gopher.png"
	publicCert = "./cert/public.crt"
	privateKey = "./cert/private.key"
	address = "localhost:8334"
)

func main() {
	image, err := ioutil.ReadFile(imageFile)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "image/png", image)
	})
	go r.RunTLS(address, publicCert, privateKey)

	// $ curl --version
	// curl 7.64.1 (x86_64-apple-darwin19.0) libcurl/7.64.1 (SecureTransport) LibreSSL/2.8.3 zlib/1.2.11 nghttp2/1.39.2
	cmd := exec.Command("curl", fmt.Sprintf("https://%s", address), "-k", "--http2")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	cmd.Run()
}

