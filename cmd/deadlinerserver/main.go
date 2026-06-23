package main

import (
	"log"

	"github.com/aritxonly/deadlinerserver/internal/app/bootstrap"
	"github.com/aritxonly/deadlinerserver/internal/config"
	"github.com/aritxonly/deadlinerserver/internal/utils/logutil"
)

func main() {
	logutil.ConfigureRuntime()

	cfg, err := config.Load(config.DefaultPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	log.Printf(
		"BOOT service=%s http=%s kitex=%s db=%s",
		cfg.Service.Name,
		cfg.HTTP.Address,
		cfg.Service.Address,
		cfg.Database.Driver,
	)
	if err := bootstrap.Run(cfg); err != nil {
		log.Printf("EXIT err=%q", err.Error())
	}
}
