package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"eagain.net/go/bearproxy/flagx"
)

func readOneLine(p string) (string, error) {
	buf, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}
	idx := bytes.IndexByte(buf, '\n')
	if idx == -1 {
		return "", fmt.Errorf("file does not contain a complete line: %s", p)
	}
	if idx != len(buf)-1 {
		return "", errors.New("trailing junk after first line")
	}
	return string(buf[:len(buf)-1]), nil
}

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", prog)
	fmt.Fprintf(flag.CommandLine.Output(), "  %s [OPTS..]\n", prog)
	fmt.Fprintf(flag.CommandLine.Output(), "\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	var tokenPath string
	flag.StringVar(&tokenPath, "secret-token-file", "", "Path to `file` with secret token")
	var backendURL flagx.AbsURL
	flag.Var(&backendURL, "backend-url", "`URL` to proxy to")
	listenAddr := flagx.TCPAddr{
		Addr: &net.TCPAddr{
			Port: 8000,
		},
	}
	flag.Var(&listenAddr, "listen", "`host:port` to listen at")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
	}

	if tokenPath == "" {
		log.Fatal("flag -secret-token-file= is required")
	}
	if backendURL.URL == (url.URL{}) {
		log.Fatal("flag -backend-url= is required")
	}

	// TODO maybe support secret rotation, and thereby multiple tokens too
	token, err := readOneLine(tokenPath)
	if err != nil {
		log.Fatalf("cannot read secret token: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(&backendURL.URL)
	accessControl := &authenticatingProxy{
		token:   token,
		handler: proxy,
	}

	listener, err := net.ListenTCP("tcp", listenAddr.Addr)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer listener.Close()

	srv := &http.Server{
		Handler:           accessControl,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	defer srv.Close()
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("http listen error: %v", err)
	}
}
