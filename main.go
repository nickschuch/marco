package main

import (
    "math/rand"
    "net/url"
    "net/http"
    "net/http/httputil"
)

func handler(w http.ResponseWriter, r *http.Request) {
    answers := []string{
      "http://www.google.com",
      "http://172.17.0.2:2368",
    }

    remote, err := url.Parse(answers[rand.Intn(len(answers))])
    if err != nil {
      panic(err)
    } 

    proxy := httputil.NewSingleHostReverseProxy(remote)
    proxy.ServeHTTP(w, r)
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":80", nil)
}

