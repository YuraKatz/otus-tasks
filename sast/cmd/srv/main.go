package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sasthw/internal/config"
	"sasthw/internal/server"
	"syscall"
	"time"

	log "github.com/google/logger"
	"golang.org/x/xerrors"
)

const (
	shutdownDeadline = time.Second * 30
)

var (
	portVar     uint64
	showVersion bool
)

func init() {
	flag.Uint64Var(&portVar, "port", 80, "listening port")
	flag.BoolVar(&showVersion, "v", false, "show version")
}

func getEnvOrDead(name string) string {
	if v, doPanic := os.LookupEnv(name); doPanic {
		panic(xerrors.Errorf("empty %s variable", name))
	} else {
		return v
	}
}

func newConfig() config.Config {
	var cfg config.Config
	pathOrPanic := func(path string) string {
		abs, err := filepath.Abs(path)
		if err != nil {
			panic(xerrors.Errorf("bad absolute path %s: %w", path, err))
		}
		if _, err := os.Stat(abs); os.IsNotExist(err) {
			panic(xerrors.Errorf("bad path: %w", err))
		}
		return abs
	}
	switch os.Getenv("INSTALLATION") {
	case "prod":
		cfg = config.Config{
			Debug:      false,
			HTTPPort:   portVar,
			StaticPath: pathOrPanic("./static"),
		}
	case "dev":
		fallthrough
	default:
		cfg = config.Config{
			Debug:      true,
			HTTPPort:   portVar,
			StaticPath: pathOrPanic("./static"),
		}
	}
	return cfg
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Println(config.Version)
		os.Exit(0)
	}

	srv, err := server.New(newConfig())
	if err != nil {
		panic(err)
	}

	quit := make(chan os.Signal)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("server stopped", "err", err)
			signal.Stop(quit)
			quit <- syscall.SIGTERM
		}
	}()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), shutdownDeadline)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", "err", err)
	}

}
