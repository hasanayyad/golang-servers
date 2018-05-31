package main

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!\n")
}

func world(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "World!\n")
}

func timeStamp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "\nServer time stamp: " + time.Now().Format("2006/01/02 3:04:05.00 PM"))
}

func log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		fmt.Println("Handler function called - " + name)
		h(w, r)
		timeStamp(w, r)
	}
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/hello", log(hello))
	http.HandleFunc("/world", log(world))
	server.ListenAndServe()
}
