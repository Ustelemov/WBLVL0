package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/ustelemov/WBLVL0/TestService/event"
)

func setup(t *testing.T) {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		t.Fatalf("error while init config: %s", err.Error())
	}

	if err = godotenv.Load("../.env"); err != nil {
		t.Fatalf("error while load .env: %s", err.Error())
	}
}

func TestMain_DBConnect(t *testing.T) {
	setup(t)

	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("host = %v port = %v user = %v dbname = %v password = %v sslmode = %v",
			viper.GetString("db.host"), viper.GetString("db.port"), viper.GetString("db.username"),
			viper.GetString("db.dbname"), os.Getenv("DB_Password"), viper.GetString("db.sslmode")))

	if err != nil {
		t.Fatalf("error while connect to postgres: %s", err.Error())
	}

	err = db.Ping()

	if err != nil {
		t.Fatalf("error while ping to postgres: %s", err.Error())
	}
}

func TestMain_NatsConnect(t *testing.T) {
	setup(t)

	t.Run("NatsOK", func(t *testing.T) {

		nes, err := event.NewNatsJsonEventStore(event.Config{
			ConnecUrl:          viper.GetString("nats.url"),
			Max_reconnects:     viper.GetInt("nats.opts.max_reconnects"),
			Reconnect_wait_sec: viper.GetInt("nats.opts.reconnect_wait_sec"),
			Subject:            viper.GetString("nats.subject"),
			Subsc_queue:        viper.GetString("nats.subsc_queue"),
		})

		if err != nil {
			t.Fatalf(err.Error())
		}

		defer nes.Close()

	})
}
