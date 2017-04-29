package main

import (
	"log"
	"net/http"
	"time"
)

// dd if=/dev/random bs=1 count=22| base64
var (
	HashKey  = []byte("xTDlKZ72kHxRRDb59lKcBppFLQVi6w==")
	BlockKey = []byte("Fc1kL2hL0KL2Pj49lJjgm+Hldulbmw==")
)

func main() {
	log.Print("starting dogf")

	r, err := NewRouter(HashKey, BlockKey, "postgres://postgres@database/dogf?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	s := &http.Server{
		Addr:           ":8000",
		Handler:        r,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
