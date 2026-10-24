
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	emailVerifier "github.com/AfterShip/email-verifier"
)

func GetEmailVerification(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	verifier := emailVerifier.NewVerifier()

	ret, err := verifier.Verify(ps.ByName("email"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ret.Syntax.Valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, `{"error":"email address syntax is invalid"}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	bytes, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(bytes)
}

func main() {
	router := httprouter.New()

	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"status":"ok"}`)
	})

	router.GET("/v1/:email/verification", GetEmailVerification)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Email verification API listening on port %s", port)

	log.Fatal(server.ListenAndServe())
}
