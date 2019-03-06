package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/koeniglorenz/bwaas/pkg/buzzwords"
	"github.com/koeniglorenz/bwaas/pkg/store"
)

var logger *log.Logger

var s store.Store

var err error

func main() {
	logger = log.New(os.Stdout, "http: ", log.LstdFlags)

	if len(os.Args) < 2 {
		log.Fatal("Please provide the path to buzzwords.json as a argument to the programm call.")
	}

	p := os.Args[1]

	s, err = store.New(p)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = buzzwords.Init(p)
	if err != nil {
		log.Fatal(err.Error())
	}

	http.HandleFunc("/", handler)

	log.Println("Starting HTTP-Server at :8080...")
	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatalf("Error starting up HTTP-Server: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	logger.Println(r.Method, r.RequestURI, r.UserAgent(), r.RemoteAddr)

	t := r.Header.Get("accept")

//	bs := buzzwords.GetRandomWords()

	if t == "application/json" {
		b, err := s.GetJSON()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error formatting JSON: %v\n", err.Error())
			return
		}
		fmt.Fprintf(w, "%s", b)
		return
	} else {
		b := s.GetHTML()
		fmt.Fprintf(w, "%s", b)
		return
	}
}
