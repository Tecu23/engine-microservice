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

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

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

	// Initialize a config struct
	cfg := config.LoadConfig(version)

	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	app := &application{
		config: cfg,
		logger: logger,
	}

	err := app.serve()
	if err != nil {
		logger.Error(err)
	}
}

func (app *application) serve() error {
	// Initialize authentication module, Maybe could be added to app struct
	auth.Initialize(app.config.AuthTokens)

	// Initliaze credentials for grpc, Maybe could be added to app struct
	creds, err := credentials.NewServerTLSFromFile(app.config.TLSCertFile, app.config.TLSKeyFile)
	if err != nil {
		return err
	}

	// Create new gRPC Server with the credentials and auth interceptor
	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor()),
	)

	// Register this server and maybe shpuld return the new server
	srv, err := server.RegisterServer(grpcServer, app.config)
	if err != nil {
		return err
	}

	// create a TCP connection on config port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", app.config.Port))
	if err != nil {
		return err
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
		app.logger.Info("shutting down server", map[string]string{
			"signal": s.String(),
		})

		// we now add a delay of 20 seconds for background jobs to complete
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// we then shutdown the server and send any error through the error channel
		err := server.Shutdown(ctx, grpcServer, srv)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks", map[string]string{
			"addr": lis.Addr().String(),
		})

		// now we wait for the all connections to be closed to showdown the server
		app.wg.Wait()
		shutdownError <- nil
	}()

	// Starting listening for requests
	app.logger.Infof("Chess Engine gRPC server is running on port %d...", app.config.Port)
	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	// checking shutdown errors
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", map[string]string{
		"addr": lis.Addr().String(),
	})

	return nil
}
