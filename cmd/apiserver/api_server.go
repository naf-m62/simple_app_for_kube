package apiserver

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

const serverTimeout = 15 * time.Second

func Start(config *Config, db *sql.DB) {
	var (
		r = mux.NewRouter()
		h = NewHandler(db)
	)
	h.configureRouter(r)

	var srv = http.Server{
		Addr:              config.Host + ":" + config.Port,
		Handler:           r,
		ReadTimeout:       serverTimeout,
		ReadHeaderTimeout: serverTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("server error happened", err)
			return
		}
	}()

	var stop = make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	var ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("shutdown error")
	}
	log.Println("server stopped")
}
