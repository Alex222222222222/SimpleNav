package main

import (
	"log"
	"net/http"
)

func InitServer() {

}

func Server() (err error) {
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		return err
	}

	return nil
}

// TODO: add handle function to return hidden_index.html by cookie
