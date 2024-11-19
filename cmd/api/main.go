// main entry of the microservice
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/Tecu23/engine-microservice/pkg/auth"
	"github.com/Tecu23/engine-microservice/pkg/config"
	"github.com/Tecu23/engine-microservice/pkg/server"
)

const version = "1.0.0"

type application struct {
	config *config.Config
	logger *log.Logger
	wg     sync.WaitGroup
}

func main() {
	// TODO: Should add config from cli with flags (port, log level, ...)

	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})

	// Initialize a config struct
	cfg, err := config.InitConfig()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	level, err := log.ParseLevel(cfg.Server.LogLevel)
	if err != nil {
		logger.WithError(err).Warn("Invalid log level, defaulting to 'info'")
		level = log.InfoLevel
	}
	logger.SetLevel(level)
	logger.WithFields(log.Fields{
		"port": cfg.Server.Port,
		"env":  os.Getenv("ENV"), // Asumming you set ENV=development|production|testing|staging
	}).Info("configuration loaded")

	app := &application{
		config: cfg,
		logger: logger,
	}

	err = app.serve()
	if err != nil {
		app.logger.WithError(err).Fatal("failed to serve")
	}
}

func (app *application) serve() error {
	// Initialize authentication module, Maybe could be added to app struct
	if err := auth.Initialize(&app.config.Auth); err != nil {
		return fmt.Errorf("failed to initialize authentication: %v", err)
	}

	// Initliaze credentials for grpc, Maybe could be added to app struct
	// creds, err := credentials.NewServerTLSFromFile(app.config.TLSCertFile, app.config.TLSKeyFile)
	// if err != nil {
	// 	return err
	// }

	// Create new gRPC Server with the credentials and auth interceptor
	grpcServer := grpc.NewServer(
		// grpc.Creds(creds),
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor()),
	)

	// Register this server and maybe shpuld return the new server
	srv, err := server.RegisterServer(grpcServer, app.config)
	if err != nil {
		return fmt.Errorf("failed to register server: %v", err)
	}

	// create a TCP connection on config port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", app.config.Server.Port))
	if err != nil {
		return fmt.Errorf("failed to start tcp server: %v", err)
	}

	// channel that receives possible errors regarding the graceful shutdown
	shutdownError := make(chan error)

	// graceful shutdown
	go func() {
		// creating a quit sygnal that will notify when SIGINT or SIGTERM are called
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// we sent the read the signal received and log it to the console
		s := <-quit
		app.logger.WithFields(log.Fields{"signal": s.String()}).Info("shutting down server")

		// we now add a delay of 20 seconds for background jobs to complete
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// we then shutdown the server and send any error through the error channel
		err := server.Shutdown(ctx, grpcServer, srv)
		if err != nil {
			shutdownError <- err
		}

		app.logger.WithFields(log.Fields{"addr": lis.Addr().String()}).
			Info("completing background tasks")

		// now we wait for the all connections to be closed to showdown the server
		app.wg.Wait()
		shutdownError <- nil
	}()

	// Starting listening for requests
	app.logger.Infof("Chess Engine gRPC server is running on port %d...", app.config.Server.Port)
	err = grpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to start grpc server: %v", err)
	}

	// checking shutdown errors
	err = <-shutdownError
	if err != nil {
		return fmt.Errorf("failed to shutdown the project: %v", err)
	}

	app.logger.WithFields(log.Fields{"signal": lis.Addr().String()}).Info("stopped server")
	return nil
}
