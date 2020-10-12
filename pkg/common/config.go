package common

import (
	"github.com/spf13/viper"
)

const (
	// Application
	AppLogLevel    = "app.logLevel"
	AppChunkSize   = "app.chunkSize"
	AppTLSCertPath = "app.tls.certPath"
	AppTLSKeyPath  = "app.tls.keyPath"
	// Server
	ServerHost = "server.host"
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
	viper.SetDefault(AppTLSCertPath, "")
	viper.SetDefault(AppTLSKeyPath, "")
	// Server
	viper.SetDefault(ServerHost, "127.0.0.1")
	// MongoDB
	viper.SetDefault(MongoDBUrl, "localhost")
	viper.SetDefault(ServerPort, "27017")
	viper.SetDefault(MongoDBDatabase, "db")
}
