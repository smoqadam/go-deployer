package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var logger *log.Logger

func main() {

	http.HandleFunc("/deploy", deployHandler)
	log.Fatalln(http.ListenAndServe(":3000", nil))
}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	if !isValidSignature(r, "MYSECRETKEY") {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("{\"result\":\"error\"}"))
	} else {
		deploy(w, r)
	}

}

func deploy(w http.ResponseWriter, r *http.Request) {
	// read config.yaml and parse it
	config, err := NewConfig("./config.yaml")
	if err != nil {
		panic(err)
	}

	// read request body
	p, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// create a payload struct based on request body
	payload, err := NewPayload(p)
	if err != nil {
		panic(err)
	}

	project := config.getProject(payload.Repo.Name)

	runner := NewRunner(project.Commands)

	if err := runner.run(); err != nil {
		runner.exec(project.OnFailure, project.Dir)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"result\":\"error\"}"))
	} else {
		w.Write([]byte("{\"result\":\"success\"}"))
	}
}

func isValidSignature(r *http.Request, key string) bool {
	gotHash := strings.SplitN(r.Header.Get("X-Hub-Signature"), "=", 2)
	if gotHash[0] != "sha1" {
		return false
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Cannot read the request body: %s\n", err)
		return false
	}

	hash := hmac.New(sha1.New, []byte(key))
	if _, err := hash.Write(b); err != nil {
		log.Printf("Cannot compute the HMAC for request: %s\n", err)
		return false
	}

	expectedHash := hex.EncodeToString(hash.Sum(nil))
	log.Println("EXPECTED HASH:", expectedHash, gotHash)
	return gotHash[1] == expectedHash
}
