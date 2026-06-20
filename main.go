package main

import (
	"log"

	"github.com/aritxonly/deadlinerserver/internal/config"
	deadlinerv1 "github.com/aritxonly/deadlinerserver/kitex_gen/deadliner/v1/deadlinerservice"
)

func main() {
	cfg, err := config.Load(config.DefaultPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	svr := deadlinerv1.NewServer(new(DeadlinerServiceImpl))

	log.Printf("starting %s on %s with %s database driver", cfg.Service.Name, cfg.Service.Address, cfg.Database.Driver)
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
