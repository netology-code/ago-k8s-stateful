package main

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

const (
	defaultPort = "9999"
	defaultHost = "0.0.0.0"
)

const dataPath = "/data"

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(addr string) error {
	size := 16
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		return err
	}
	key := base64.RawStdEncoding.EncodeToString(buf)

	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("start handle %s by: %s", request.URL.Path, key)
		defer log.Printf("finish handle %s by: %s", request.URL.Path, key)
		switch request.Method {
		case http.MethodPost:
			{
				src, _, err := request.FormFile("file")
				if err != nil {
					log.Print(err)
					http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				name := uuid.New().String()
				dst, err := os.Create(filepath.Join(dataPath, name))
				if err != nil {
					log.Print(err)
					http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				_, err = io.Copy(dst, src)
				if err != nil {
					log.Print(err)
					http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				_, err = writer.Write([]byte(name))
				if err != nil {
					log.Print(err)
				}
			}
		default:
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	})

	server := &http.Server{Addr: addr, Handler: handler}

	log.Printf("server %s started at %s", key, addr)
	return server.ListenAndServe()
}
