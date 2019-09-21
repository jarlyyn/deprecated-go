package httpserver

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

//Config http server config.
type Config struct {
	//Net net interface,"tcp" for example.
	Net string
	//Addr network addr.
	Addr string
	//BaseURL http scheme and host."http://127.0.0.1:8000" for example.
	BaseURL string
	//ReadTimeoutInSecond http conn read time out.
	ReadTimeoutInSecond int64
	//ReadTimeoutInSecond http conn read Header time out.
	ReadHeaderTimeoutInSecond int64
	//WriteTimeoutInSecond http conn write time out.
	WriteTimeoutInSecond int64
	//IdleTimeoutInSecond conn idle time out.
	IdleTimeoutInSecond int64
	//MaxHeaderBytes max header length in bytes.
	MaxHeaderBytes int
}

//Server create http server with config.
func (c *Config) Server() *http.Server {
	server := &http.Server{
		ReadTimeout:       time.Duration(c.ReadTimeoutInSecond) * time.Second,
		ReadHeaderTimeout: time.Duration(c.ReadHeaderTimeoutInSecond) * time.Second,
		WriteTimeout:      time.Duration(c.WriteTimeoutInSecond) * time.Second,
		IdleTimeout:       time.Duration(c.IdleTimeoutInSecond) * time.Second,
		MaxHeaderBytes:    c.MaxHeaderBytes,
	}
	server.ErrorLog = log.New(ioutil.Discard, "", 0)
	return server
}

//Listen listen net and addr in config.
//Return net listener and any error if raised.
func (c *Config) Listen() (net.Listener, error) {
	return net.Listen(c.Net, c.Addr)
}

//MustListen listen net and addr in config.
//Return net listener.
//Panic if any error raised.
func (c *Config) MustListen() net.Listener {
	l, err := net.Listen(c.Net, c.Addr)
	if err != nil {
		panic(err)
	}
	return l
}
