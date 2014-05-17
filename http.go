package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func initHttpServer(port int) {
	http.HandleFunc("/record.do", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, "{}")
		log.Printf("Got record from %s.", r.RemoteAddr)
	})
	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil {
			log.Fatalf("Exit %v", err)
		}
	}()
}
