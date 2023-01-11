package main

import (
	"fmt"
	"github.com/OlegKapat/Rest-api-mongo/internal/config"
	"github.com/OlegKapat/Rest-api-mongo/pkg/logging"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/OlegKapat/Rest-api-mongo/internal/user"
	"github.com/julienschmidt/httprouter"
)

func main() {
	//logging.Init()
	logger := logging.GetLogger()
	logger.Info("Logger work")
	router := httprouter.New()

	cfg := config.GetConfig()
	logger.Info("Register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)
	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("Start application")
	var listenerErr error
	var listener net.Listener
	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketpath := path.Join(appDir, "app.sock")
		logger.Debugf("socket path %s", socketpath)

		logger.Info("listen unix socket")
		listener, listenerErr = net.Listen("unix", socketpath)
		logger.Infof("server is listenning unix socket %s", socketpath)

	} else {
		logger.Info("listen tcp port")
		listener, listenerErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("Server is listering port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}
	if listenerErr != nil {
		logger.Fatal(listenerErr)
	}

	server := http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
