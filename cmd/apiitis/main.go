package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Asmitshukl/apiitis/internal/config"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//setup datbase

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /",func(w http.ResponseWriter , r *http.Request) {
		w.Write([]byte("Welcome to just api"))
	})
	//setup server
	fmt.Println(" server started")
	
	server :=http.Server{
		Addr: cfg.HttpServer.Addr,
		Handler: router,
	}

	done :=make(chan os.Signal ,1)
	signal.Notify(done,os.Interrupt, syscall.SIGINT , syscall.SIGTERM)

	go func(){
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("failed to start server")
	}
	}()

	<-done
	slog.Info("shutting down the server")

	ctx,cancel := context.WithTimeout(context.Background(),5 * time.Second)
	defer cancel()



	// err := server.Shutdown(ctx)

	if err := server.Shutdown(ctx) ; err != nil {
		slog.Error("failed to shutdown server" ,slog.String("error" , err.Error()))
	}

	slog.Info("server stopped successfully")
}
