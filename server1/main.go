package main

import (
    "fmt"
    "net/http"
    "log"
    "github.com/robjporter/CPVT-CloudNativeLab/lab"
)

var redisServer string
var serverCount string
var dbStartCount string

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    color := "red"
    content := "<html><head></head><body bgcolor='"+color+"'>"
    content += serverCount+" servers currently serving web content.<br>"
    content += dbStartCount+" server starts have been seen.<br>"
    content += lab.GetPageCount()+" page loads have happened.<br>"
    content += "</body></html>"
    fmt.Fprintf(w, content)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w,"OK\n")
}

func main() {
    http.HandleFunc("/", sayhelloName) // set router
    http.HandleFunc("/health", healthCheck) // set router

    urlsToRegister := []string{"/"}
    ipOfConsul := "172.17.0.2"
    portWeListenOn := "8080"

    result, err := lab.RegisterMe(ipOfConsul, urlsToRegister, portWeListenOn )
    redisServer = lab.GetServiceAddress("redis")
    serverCount = lab.GetServerCount("localhost-")
    dbStartCount = lab.GetDBStartCount()

    if result {
        fmt.Println("Server started and listening on port :"+portWeListenOn)
        err = http.ListenAndServe(":"+portWeListenOn, nil) // set listen port
        if err != nil {
            log.Fatal("ListenAndServe: ", err)
        }
    } else {
        fmt.Println(err)
    }
}

