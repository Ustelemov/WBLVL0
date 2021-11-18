package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/ustelemov/WBLVL0/TestService/cache"
	"github.com/ustelemov/WBLVL0/TestService/event"
	"github.com/ustelemov/WBLVL0/TestService/handler"
	"github.com/ustelemov/WBLVL0/TestService/repository"
	"github.com/ustelemov/WBLVL0/TestService/server"
	"github.com/ustelemov/WBLVL0/TestService/service"
)

func initConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error while init config: %s", err.Error())
	}

	if err := godotenv.Load("../.env"); err != nil {
		logrus.Fatalf("Error while load .env: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DBName:   viper.GetString("db.dbname"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_Password"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("Error occured while creating postgres: %s", err.Error())
	}

	cache := cache.NewMapRepositoryCache()
	orderService := service.NewOrdersService(db, cache)

	err = orderService.GetAllOrdersInCache()

	if err != nil {
		logrus.Fatalf(err.Error())
	}

	nes, err := event.NewNatsJsonEventStore(viper.GetString("nats_url"), func(_ *nats.Conn, _ *nats.Subscription, _ error) {
		logrus.Errorf("error by nats-error-handler: %s", err)
	})

	if err != nil {
		logrus.Fatalf("Error occured while creating nats: %s", err.Error())
	}

	defer nes.Close()

	err = nes.SubscribeOnOrders(orderService.ProccessOrderMessage)

	if err != nil {
		logrus.Fatalf("error on subscribing: %s", err.Error())
	}

	handler := handler.NewHandler(orderService)
	server := new(server.Server)

	if err := server.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		logrus.Fatalf("Error occured while starting http-server: %s", err.Error())
	}

}
