package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/pat"

	m "github.com/4sh/k8s-toleration-injection-webhook/pkg/mutate"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello %q", html.EscapeString(r.URL.Path))
}

func handleMutate(w http.ResponseWriter, r *http.Request) {
	// read the body / request
	injectionId := r.URL.Query().Get(":id")
	log.Printf("injecting toleration for %s", injectionId)

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		sendError(err, w)
		return
	}

	// mutate the request
	mutated, err := m.Mutate(body, injectionId, true)
	if err != nil {
		sendError(err, w)
		return
	}

	// and write it back
	w.WriteHeader(http.StatusOK)
	w.Write(mutated)
}

func sendError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%s", err)
}

func main() {
	log.Println("Starting server ...")

	var mux *pat.Router = pat.New()

	mux.Get("/", handleRoot)
	mux.Post("/mutate/{id}", handleMutate)

	s := &http.Server{
		Addr:           ":8443",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}

	log.Fatal(s.ListenAndServeTLS("./ssl/toleration-injection.pem", "./ssl/toleration-injection.key"))
}
