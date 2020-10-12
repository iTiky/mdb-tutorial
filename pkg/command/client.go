package command

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	v1 "github.com/itiky/mdb-tutorial/pkg/api/v1"
	"github.com/itiky/mdb-tutorial/pkg/common"
)

const (
	flagPageSkip        = "skip"
	flagPageLimit       = "limit"
	flagSortByName      = "sort-by-name"
	flagSortByPrice     = "sort-by-price"
	flagSortByTimestamp = "sort-by-timestamp"
)

// clientCmd is a gRPC-client debug root command.
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Debug gRPC client",
}

// GetClientListCmd returns a gRPC-client command for List() request.
func GetClientListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List price entries",
		Run: func(cmd *cobra.Command, args []string) {
			logger := initLogger()

			// parse inputs
			pageSkip, pageLimit := parseIntFlag(logger, flagPageSkip, cmd.Flags()), parseIntFlag(logger, flagPageLimit, cmd.Flags())
			sortByName, sortByPrice, sortByTimestamp := parseSortFlag(flagSortByName, cmd.Flags()), parseSortFlag(flagSortByPrice, cmd.Flags()), parseSortFlag(flagSortByTimestamp, cmd.Flags())

			// create gRPC client
			conn, err := createGRPCConnection()
			if err != nil {
				logger.Fatalf(err.Error())
			}
			defer conn.Close()

			client := v1.NewPriceEntryReaderClient(conn)

			// request
			requestCtx, requestCancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer requestCancel()
			resp, err := client.List(requestCtx, &v1.ListRequest{
				Pagination: &v1.PaginationParams{
					Skip:  uint32(pageSkip),
					Limit: uint32(pageLimit),
				},
				SortByName:      sortByName,
				SortByPrice:     sortByPrice,
				SortByTimestamp: sortByTimestamp,
			})
			if err != nil {
				logger.Fatalf("request failed: %v", err)
			}

			// print result
			if len(resp.Entries) == 0 {
				logger.Infof("no entries found")
				return
			}

			for _, entry := range resp.Entries {
				logger.Infof("%s\t->\t%s\t->\t%s",
					entry.ProductName,
					strconv.FormatInt(int64(entry.Price), 10),
					time.Unix(entry.Timestamp, 0).Format(time.RFC3339),
				)
			}
		},
	}
	cmd.Flags().Int(flagPageSkip, 0, "(optional) pagination param: skip")
	cmd.Flags().Int(flagPageLimit, 50, "(optional) pagination param: limit")
	cmd.Flags().String(flagSortByName, "", "(optional) sort param: by product name (ASC/DESC)")
	cmd.Flags().String(flagSortByPrice, "", "(optional) sort param: by price (ASC/DESC)")
	cmd.Flags().String(flagSortByTimestamp, "", "(optional) sort param: by timestamp (ASC/DESC)")

	return cmd
}

// GetClientFetchCmd returns a gRPC-client command for Fetch() request.
func GetClientFetchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fetch",
		Short:   "Fetch price entries CSV-file for specified URL arg",
		Example: "fetch {url_to_file}",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger := initLogger()

			// create gRPC client
			conn, err := createGRPCConnection()
			if err != nil {
				logger.Fatalf(err.Error())
			}
			defer conn.Close()

			client := v1.NewCSVFetcherClient(conn)

			// request
			requestCtx, requestCancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer requestCancel()
			_, err = client.Fetch(requestCtx, &v1.CSVFetchRequest{
				Url: args[0],
			})
			if err != nil {
				logger.Fatalf("request failed: %v", err)
			}

			logger.Infof("request: ok")
		},
	}

	return cmd
}

// GetClientFileServerCmd returns a file server command which provides CSV-files.
func GetClientFileServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "file-server",
		Short:   "Start a mock fileServer mirroring fileSystem",
		Example: "file-server {path_to_directory} {port}",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			logger := initLogger()

			// start file server
			fs := http.FileServer(http.Dir(args[0]))
			server := &http.Server{
				Addr:    fmt.Sprintf(":%s", args[1]),
				Handler: fs,
			}

			go func() {
				if err := server.ListenAndServe(); err != nil {
					logger.Fatalf("HTTP server: crashed: %v", err)
				}
			}()
			logger.Infof("HTTP server: started: %s", server.Addr)

			// Waiting for shutdown signals
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
			<-stop

			// Shutdown
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer shutdownCancel()
			err := server.Shutdown(shutdownCtx)
			logger.Infof("HTTP server: stopped: %v", err)
		},
	}

	return cmd
}

// parseIntFlag parses int cmd flag (crashes on failure).
func parseIntFlag(logger *logrus.Logger, flagName string, flags *pflag.FlagSet) int {
	v, err := flags.GetInt(flagName)
	if err != nil {
		logger.Fatalf("parsing %s flag: %v", flagName, err)
	}

	return v
}

// parseSortFlag converts cmd flag to gRPC SortOrder.
func parseSortFlag(flagName string, flags *pflag.FlagSet) v1.SortOrder {
	value := ""
	if v, err := flags.GetString(flagName); err == nil {
		value = strings.ToLower(v)
	}

	switch value {
	case "asc":
		return v1.SortOrder_Asc
	case "desc":
		return v1.SortOrder_Desc
	default:
		return v1.SortOrder_Undefined
	}
}

// createGRPCConnection creates a new gRPC client connection.
func createGRPCConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%s", viper.GetString(common.ServerPort)),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating client connection failed: %v", err)
	}

	return conn, nil
}

func init() {
	clientCmd.AddCommand(GetClientListCmd())
	clientCmd.AddCommand(GetClientFetchCmd())
	clientCmd.AddCommand(GetClientFileServerCmd())
	rootCmd.AddCommand(clientCmd)
}
