package main

import "github.com/thschue/environment-readiness-service/pkg/server"

func main() {
	serv := server.Config{
		Port: "8080",
	}
	err := serv.Serve()
	if err != nil {
		panic(err)
	}
}
