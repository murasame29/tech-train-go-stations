package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/env"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}

}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	notifyCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	defer stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	env, err := env.GetEnv()
	if err != nil {
		return err
	}

	srv := http.Server{
		Addr:    port,
		Handler: router.NewRouter(todoDB, env).Mux,
	}

	// Graceful Shutdown https://pkg.go.dev/net/http#Server.Shutdown
	idleConnsClosed := make(chan struct{})
	go func() error {
		<-notifyCtx.Done()
		log.Println("shutdown server")

		// タイムアウト設定
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()

		if err := srv.Shutdown(context.Background()); err != nil {
			close(idleConnsClosed)
			return err
		}
		// シャットダウンされたときにチャネルを閉じる
		close(idleConnsClosed)
		return nil
	}()

	// サーバを起動する
	go func() error {
		log.Println("starting server ...")
		if err := srv.ListenAndServe(); err != nil {
			return err
		}
		return nil
	}()

	// チャネルが閉じられるまで待つ
	<-idleConnsClosed

	return nil
}
