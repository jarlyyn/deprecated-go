package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

//MustListenAndServeHTTP listen and serve http server with given server,config and handler.
//Panic if any error raised.
func MustListenAndServeHTTP(server *http.Server, config Config, app http.Handler) {
	go func() {
		l := config.MustListen()
		defer l.Close()
		fmt.Println("Listening " + l.Addr().String())
		server.Handler = app
		err := server.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

//MustServeHTTP  serve http server with given server,listener  and handler.
//Panic if any error raised.
func MustServeHTTP(server *http.Server, l net.Listener, app http.Handler) {
	go func() {
		fmt.Println("Listening " + l.Addr().String())
		server.Handler = app
		err := server.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

}

// ShutdownHTTP  shutdown  http server.
func ShutdownHTTP(Server *http.Server) {
	WithContextShutdown(context.Background(), Server)
}

//ShutdownHTTPWithTimeout shutdown  http server ith given timeout.
func ShutdownHTTPWithTimeout(Server *http.Server, Timeout time.Duration) {
	ctx, fn := context.WithTimeout(context.Background(), Timeout)
	fn()
	WithContextShutdown(ctx, Server)

}

//WithContextShutdown shutdown  http server ith given context.
func WithContextShutdown(ctx context.Context, Server *http.Server) {
	fmt.Println("Http server shuting down...")
	Server.Shutdown(ctx)
	fmt.Println("Http server Stoped.")
}
