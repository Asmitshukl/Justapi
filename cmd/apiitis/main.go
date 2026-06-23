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
	"github.com/Asmitshukl/apiitis/internal/http/handlers/student"
	"github.com/Asmitshukl/apiitis/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//setup datbase

	_, err := sqlite.New(cfg)

	if err != nil {
		log.Fatalf("failed to setup database %s", err.Error())
	}
	slog.Info("database setup successfully",slog.String("env",cfg.Env))
	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students",student.New())
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
