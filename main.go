package main

import (
	"github.com/White-AK111/REST/config"
	"github.com/White-AK111/REST/fasth"
	ginGonic "github.com/White-AK111/REST/gin-gonic"
	"github.com/White-AK111/REST/gorilla"
	stdlibHttp "github.com/White-AK111/REST/stdlib-http"
	"log"
)

func main() {
	// init configuration
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("error on load configration file: %s", err)
	}

	switch cfg.Server.TypeOfServer {
	case "stdlib":
		stdlibHttp.Init(cfg)
	case "gorilla":
		gorilla.Init(cfg)
	case "gin":
		ginGonic.Init(cfg)
	case "fasthttp":
		fasth.Init(cfg)
	default:
		log.Fatal("Unknown service type.")
	}

}
