package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aritxonly/deadlinerserver/internal/config"
	hertzserver "github.com/cloudwego/hertz/pkg/app/server"
	kitexserver "github.com/cloudwego/kitex/server"
)

func Run(cfg config.Config) error {
	runtime, err := NewRuntime(cfg)
	if err != nil {
		return err
	}

	kitexServer, err := NewKitexServer(cfg, runtime)
	if err != nil {
		return err
	}

	httpServer, err := NewHertzServer(cfg, runtime)
	if err != nil {
		return err
	}

	errCh := make(chan error, 2)
	go func() {
		errCh <- fmt.Errorf("kitex server stopped: %w", kitexServer.Run())
	}()
	go func() {
		errCh <- fmt.Errorf("http server stopped: %w", httpServer.Run())
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signalCh)

	select {
	case sig := <-signalCh:
		log.Printf("STOP signal=%s", sig.String())
		if err := shutdownServers(httpServer, kitexServer); err != nil {
			return errors.Join(fmt.Errorf("received shutdown signal: %s", sig), err)
		}
		return nil
	case err := <-errCh:
		log.Printf("STOP cause=%q", err.Error())
		shutdownErr := shutdownServers(httpServer, kitexServer)
		if shutdownErr != nil {
			return errors.Join(err, shutdownErr)
		}
		return err
	}
}

func shutdownServers(httpServer *hertzserver.Hertz, kitexServer kitexserver.Server) error {
	var shutdownErr error

	if httpServer != nil {
		if err := httpServer.Shutdown(context.Background()); err != nil {
			shutdownErr = errors.Join(shutdownErr, fmt.Errorf("shutdown http server: %w", err))
		}
	}

	if kitexServer != nil {
		if err := kitexServer.Stop(); err != nil {
			shutdownErr = errors.Join(shutdownErr, fmt.Errorf("stop kitex server: %w", err))
		}
	}

	return shutdownErr
}
