package main

import (
	"fmt"
	"geoip2-http/geoip"
	"geoip2-http/internal/server"
)

func main() {
	err := geoip.ConnectDB()
	if err != nil {
		panic(fmt.Sprintf("cannot start database: %s", err))
	}
	println("connected to db")
	server := server.NewServer()
	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
