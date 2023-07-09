package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	env "github.com/joho/godotenv"

	"github.com/ilivestrong/rules-engine/controllers"
	"github.com/ilivestrong/rules-engine/helpers"
	"github.com/ilivestrong/rules-engine/rules"
)

const envFile = ".env"

var loadEnv = env.Load

type service struct {
	Server *http.Server
	DBConn *pgx.Conn
}

func run() *service {
	err := loadEnv(envFile)
	if err != nil {
		fmt.Println("failed to load .env")
	}
	port, exist := os.LookupEnv("PORT")
	if !exist {
		fmt.Println("no port specified, defaulting to 4333")
		port = "4333"
	}
	port = fmt.Sprintf(":%s", port)
	fmt.Println("Listening at: ", port)

	fileManager := helpers.NewFileManager()

	rulesEngine, err := rules.NewRulesEngine(fileManager)
	if err != nil {
		fmt.Println(err)
	}

	dbCtx := context.Background()

	config := helpers.Config{
		User:   os.Getenv("DB_USER"),
		Pass:   os.Getenv("DB_PASSWORD"),
		Host:   os.Getenv("DB_HOST_NAME"),
		DBName: os.Getenv("DB_NAME"),
	}
	rulesDB, err := helpers.NewRulesEngineRepo(dbCtx, &config)
	if err != nil {
		fmt.Println(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/process", &controllers.CrediCardApprovalHandler{
		RulesEngine: rulesEngine,
		FileManager: fileManager,
		DBManager:   rulesDB,
	})

	s := &http.Server{
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error listening on port: %s\n", err)
		}
	}()

	return &service{
		Server: s,
		DBConn: rulesDB.Conn,
	}
}

func main() {
	svc := run()
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := svc.Server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shut down")
	}

	if err := svc.DBConn.Close(ctx); err != nil {
		log.Fatal("server forced to shut down")
	}
	log.Println("database connection closed")
	log.Println("server exiting")
}
