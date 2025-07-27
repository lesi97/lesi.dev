package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/api.lesi.dev/internal/app"
	"github.com/lesi97/api.lesi.dev/internal/router"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	application, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer application.DB.Close(context.Background())

	routes := router.SetupRoutes(application)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		IdleTimeout: time.Minute,
		Handler: routes,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	application.Logger.Printf("we are running on port %d\n", port)

	err = server.ListenAndServe() 
	if err != nil {
		application.Logger.Fatal(err)
	}
}