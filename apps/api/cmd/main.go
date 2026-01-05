package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/api/internal/app"
	"github.com/lesi97/lesi.dev/api/internal/router"
	"github.com/lesi97/lesi.dev/api/internal/utils"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	application, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer application.DB.Close()

	routes := router.SetupRoutes(application)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		IdleTimeout: time.Minute,
		Handler: routes,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	utils.Startup(application.Logger, fmt.Sprintf(":%d", port))

	err = server.ListenAndServe() 
	if err != nil {
		application.Logger.Fatal(err)
	}
}