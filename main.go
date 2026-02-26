package main

import (
	"app/pkg/router"
	"log"
	"os"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// r := gin.Default()
	r := gin.New()
	r.Use(gin.Recovery())

	router.SetupRoutes(r)

	if os.Getenv("PRODUCTION") == "true" {
		go func() {
			r.Run(":80")
		}()

		// log.Fatal(autotls.Run(r, "hainan888.top", "stockapp.sandhope.com"))
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("hainan888.top"),
			Cache:      autocert.DirCache("/app/cert"),
		}

		log.Fatal(autotls.RunWithManager(r, &m))
	} else {
		log.Println("Running in development mode - only HTTP server")
		log.Fatal(r.Run(":80"))
	}
}
