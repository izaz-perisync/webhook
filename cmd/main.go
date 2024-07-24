package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webhook/api/handler"
	"webhook/api/router"
	"webhook/service"
	s3 "webhook/utils/s3"

	global "webhook"

	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"
	"github.com/urfave/negroni"
)

func main() {

	l := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	l.Println("connecting to db...")
	ctx := context.Background()
	location, _ := time.LoadLocation("Asia/Kolkata")

	config := global.GlobalConfig()
	// connStr := "host=localhost port=5431 user=postgres password=Pass@1234 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		log.Println("connect2 err", err)
		return
	}

	s3Client, err := s3.New("", config.AccessKey, config.SecretAccessKey, config.Region, config.Bucket, "")
	if err != nil {
		log.Println("s3 bucket err", err)
		return
	}

	s := service.New(config.JWT.Key, config.Discord.WebhookUrl, db, location, l, s3Client)

	h := handler.New(s)

	r := router.RouteBuilder(h)
	n := negroni.Classic()
	n.UseHandler(r)

	server := http.Server{
		Addr: fmt.Sprintf(":%d", 3005),
		Handler: handlers.CORS(
			handlers.ExposedHeaders([]string{"at", "At"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE", "PATCH"}),
			handlers.AllowedOrigins([]string{"*"}),
		)(n),
	}

	go func() {
		// Start listening.
		log.Println("service listening at", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Print("server err", err)
		}
	}()

	lock := make(chan os.Signal, 1)
	signal.Notify(lock, os.Interrupt, syscall.SIGTERM)
	<-lock

	server.Shutdown(ctx)

}
