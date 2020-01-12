package postgres

import (
	"database/sql"
	"fmt"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var pool *sql.DB

func GetDB() *sql.DB {
	// TODO: create a wrapper for transactions and bind it to context
	return pool
}

func InitDB() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetString("database.sslmode"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.WithError(err).Fatal("unable to open database")
	}

	db.SetMaxOpenConns(viper.GetInt("database.max_open_conns"))
	db.SetMaxIdleConns(viper.GetInt("database.max_idle_conns"))
	db.SetConnMaxLifetime(time.Minute * time.Duration(viper.GetInt("database.conn_max_lifetime")))

	if err := db.Ping(); err != nil {
		logrus.WithError(err).Fatal("unable to ping database")
	}

	pool = db

	runMigrations()
}

func runMigrations() {
	driver, err := postgres.WithInstance(pool, &postgres.Config{})
	if err != nil {
		logrus.WithError(err).Fatal("unable to migrate database")
	}

	m, err := migrate.NewWithDatabaseInstance("file:///opt/shorty-api/migrations", "postgres", driver)
	if err != nil {
		logrus.WithError(err).Fatal("unable to migrate database")
	}

	if err := m.Up(); err != nil {
		logrus.WithError(err).Fatal("unable to migrate database")
	}
}
