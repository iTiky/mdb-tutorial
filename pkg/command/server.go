package command

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	v1 "github.com/itiky/mdb-tutorial/pkg/api/v1"
	"github.com/itiky/mdb-tutorial/pkg/common"
	"github.com/itiky/mdb-tutorial/pkg/mongodb"
	"github.com/itiky/mdb-tutorial/pkg/service"
	"github.com/itiky/mdb-tutorial/pkg/storage"
)

// serverCmd is a gRPC-server start command.
var serverCmd = &cobra.Command{
	Use:   "start",
	Short: "Start gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		logger := initLogger()

		// MongoDB connection
		mdbClient, err := mongodb.Connect(mongodb.Configuration{
			Url:  viper.GetString(common.MongoDBUrl),
			Port: viper.GetString(common.MongoDBPort),
		})
		if err != nil {
			logger.Fatalf("MongoDB connection failed: %v", err)
		}
		logger.Infof("MongoDB: connected")

		// Init dependencies
		storage, err := storage.NewStorage(
			storage.WithMongoDBClient(mdbClient),
			storage.WithDatabase(viper.GetString(common.MongoDBDatabase)),
			storage.WithLogger(logger),
		)
		if err != nil {
			logger.Fatalf("storage dep init: %v", err)
		}

		service, err := service.NewService(
			service.WithStorage(storage),
			service.WithLogger(logger),
		)
		if err != nil {
			logger.Fatalf("service dep init: %v", err)
		}

		// Start gRPC server
		server, err := v1.NewServer(
			v1.WithService(service),
			v1.WithLogger(logger),
			v1.WithCSVChunkSize(viper.GetInt(common.AppChunkSize)),
		)
		if err != nil {
			logger.Fatalf("server dep init: %v", err)
		}

		serverListener, err := net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString(common.ServerPort)))
		if err != nil {
			log.Fatalf("server listener init: %v", err)
		}

		go func() {
			if err := server.Serve(serverListener); err != nil {
				logger.Panicf("gRPC server crashed: %v", err)
			}
		}()
		logger.Infof("gRPC server: started at %s", serverListener.Addr().String())

		// Waiting for shutdown signals
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop

		// Shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		server.Stop()
		logger.Infof("gRPC server: stopped")
		_ = mdbClient.Disconnect(shutdownCtx)
		logger.Infof("MongoDB: disconnected")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
