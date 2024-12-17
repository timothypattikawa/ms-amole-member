package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"time"
)

type DatabaseConnectPool struct {
	minConn           int
	maxConn           int
	maxLifeTime       time.Duration
	keepAliveInterval time.Duration
}

type DatabaseConnection struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
	pool     *DatabaseConnectPool
}

// postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
func (dc *DatabaseConnection) getConnectionUrlSource() string {
	u := url.URL{
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%d", dc.host, dc.port),
		User:   url.UserPassword(dc.user, dc.password),
		Path:   dc.dbname,
	}

	values := url.Values{}
	values.Add("sslmode", "disable")

	u.RawQuery = values.Encode()
	return u.String()
}

func (dc *DatabaseConnection) GetConnectionPgx() *pgxpool.Pool {
	source := dc.getConnectionUrlSource()
	log.Printf("Url Database: %s", source)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(source)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	config.MinConns = int32(dc.pool.minConn)
	config.MaxConns = int32(dc.pool.maxConn)
	config.MaxConnLifetime = dc.pool.maxLifeTime
	config.HealthCheckPeriod = dc.pool.keepAliveInterval

	poolWithConfig, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return poolWithConfig
}

func newDatabaseConnect(v *viper.Viper, dbname string) *DatabaseConnection {
	return &DatabaseConnection{
		host:     v.GetString(fmt.Sprintf("db.%s.host", dbname)),
		port:     v.GetInt(fmt.Sprintf("db.%s.port", dbname)),
		user:     v.GetString(fmt.Sprintf("db.%s.user", dbname)),
		password: v.GetString(fmt.Sprintf("db.%s.password", dbname)),
		dbname:   v.GetString(fmt.Sprintf("db.%s.schema", dbname)),
		pool: &DatabaseConnectPool{
			minConn:           v.GetInt(fmt.Sprintf("db.%s.min-conn", dbname)),
			maxConn:           v.GetInt(fmt.Sprintf("db.%s.max-conn", dbname)),
			maxLifeTime:       v.GetDuration(fmt.Sprintf("db.%s.max-life-time", dbname)),
			keepAliveInterval: v.GetDuration(fmt.Sprintf("db.%s.keep-alive-interval", dbname)),
		},
	}
}
