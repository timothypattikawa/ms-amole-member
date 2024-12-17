package config

import "github.com/spf13/viper"

type Config struct {
	DbPostgres *DatabaseConnection
	Svc        *ServerConfiguration
}

type ServerConfiguration struct {
	Port     string
	GrpcPort string
}

func NewConfig(v *viper.Viper) Config {
	return Config{
		DbPostgres: newDatabaseConnect(v, "postgres"),
		Svc:        newServerConfiguration(v),
	}
}

func newServerConfiguration(v *viper.Viper) *ServerConfiguration {
	return &ServerConfiguration{
		Port:     v.GetString("server.port"),
		GrpcPort: v.GetString("server.grpc_port"),
	}
}
