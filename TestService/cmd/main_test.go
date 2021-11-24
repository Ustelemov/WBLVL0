package main

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/ustelemov/WBLVL0/TestService/event"
	"github.com/ustelemov/WBLVL0/TestService/repository"
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

	t.Run("DBOK", func(t *testing.T) {

		db, err := repository.NewPostgresDB(repository.Config{
			Host:                           viper.GetString("db.host"),
			Port:                           viper.GetString("db.port"),
			DBName:                         viper.GetString("db.dbname"),
			Username:                       viper.GetString("db.username"),
			Password:                       os.Getenv("DB_Password"),
			SSLMode:                        viper.GetString("db.sslmode"),
			Listener_Max_Reconnect_Seconds: viper.GetInt("db.listener.max_reconnect_seconds"),
			Listener_Min_Reconnect_Seconds: viper.GetInt("db.listener.min_reconnect_seconds"),
			Listener_Ping_NoEvent_Seconds:  viper.GetInt("db.listener.ping_noevent_seconds"),
		})

		if err != nil {
			t.Fatalf(err.Error())
		}

		defer db.Close()
	})
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
