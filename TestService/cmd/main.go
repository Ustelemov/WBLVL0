package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"TestService/cache"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

	gin.SetMode(gin.ReleaseMode) //disable debug-mode messages

	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := initConfig(); err != nil {
		logrus.Fatalf("error while init config: %s", err.Error())
	}

	if err := godotenv.Load("../.env"); err != nil {
		logrus.Fatalf("error while load .env: %s", err.Error())
	}

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
		logrus.Fatalf(err.Error())
	}

	defer db.Close()

	cache := cache.NewMapRepositoryCache()
	orderService := service.NewOrdersService(db, cache)

	err = db.RunListenerDeamon("event_channel", orderService.UpdateCacheOnTringer)

	if err != nil {
		logrus.Fatalf(err.Error())
	}

	err = orderService.LoadAllOrdersInCache()

	if err != nil {
		logrus.Fatalf(err.Error())
	}

	nes, err := event.NewNatsJsonEventStore(event.Config{
		ConnecUrl:          viper.GetString("nats.url"),
		Max_reconnects:     viper.GetInt("nats.opts.max_reconnects"),
		Reconnect_wait_sec: viper.GetInt("nats.opts.reconnect_wait_sec"),
		Subject:            viper.GetString("nats.subject"),
		Subsc_queue:        viper.GetString("nats.subsc_queue"),
	})

	if err != nil {
		logrus.Fatalf(err.Error())
	}

	defer nes.Close()

	err = nes.QueueSubscribeOnOrders(orderService.ProccessOrderMessage)

	if err != nil {
		logrus.Fatalf("error occured while subscribing: %s", err.Error())
	}

	html_handler := http.FileServer(http.Dir("../web"))
	html_server := new(server.Server)

	go func() {
		if err := html_server.Run(viper.GetString("html_port"), html_handler); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error occured while starting html-server: %s", err.Error())
		}
	}()

	api_handler := handler.NewHandler(orderService)
	api_server := new(server.Server)

	go func() {
		if err := api_server.Run(viper.GetString("api_port"), api_handler.InitRoutes()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error occured while starting api-server: %s", err.Error())
		}
	}()

	logrus.Println("App started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	logrus.Println("App Shutting Down")

	if err := html_server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured while shutting down html-server")
	}
	if err := api_server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured while shutting down api-server")
	}
}
