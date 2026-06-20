package main

import (
	"log"

	"github.com/aritxonly/deadlinerserver/internal/app/bootstrap"
	"github.com/aritxonly/deadlinerserver/internal/config"
)

func main() {
	cfg, err := config.Load(config.DefaultPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	svr, err := bootstrap.NewKitexServer(cfg)
	if err != nil {
		log.Fatalf("build server failed: %v", err)
	}

	log.Printf("starting %s on %s with %s database driver", cfg.Service.Name, cfg.Service.Address, cfg.Database.Driver)
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
