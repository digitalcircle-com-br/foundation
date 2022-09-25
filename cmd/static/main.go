package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	rootRouter := mux.NewRouter()
	withGz := gziphandler.GzipHandler(http.FileServer(http.Dir(".")))
	rootRouter.PathPrefix("/").Handler(withGz)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM)
	go func() {
		asig := <-sig
		log.Printf("Signal catch: %s", asig.String())
		os.Exit(0)
	}()

	err := http.ListenAndServe(":8080", rootRouter)
	if err != nil {
		panic(err)
	}
}
