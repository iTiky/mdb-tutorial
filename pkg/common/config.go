package common

import (
	"github.com/spf13/viper"
)

const (
	// Application
	AppLogLevel  = "app.logLevel"
	AppChunkSize = "app.chunkSize"
	// Server
	ServerPort = "server.port"
	// MongoDB
	MongoDBUrl      = "mdb.url"
	MongoDBPort     = "mdb.port"
	MongoDBDatabase = "mdb.database"
)

func init() {
	// Application
	viper.SetDefault(AppLogLevel, "info")
	viper.SetDefault(AppChunkSize, "3")
	// Server
	viper.SetDefault(MongoDBUrl, "2412")
	// MongoDB
	viper.SetDefault(MongoDBUrl, "localhost")
	viper.SetDefault(ServerPort, "27017")
	viper.SetDefault(MongoDBDatabase, "db")
}
