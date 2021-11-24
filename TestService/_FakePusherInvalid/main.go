package main

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type InvalidMessage struct {
	Message string
}

func initConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {

	if err := initConfig(); err != nil {
		logrus.Fatalf("error while init config: %s", err.Error())
	}

	nc, err := nats.Connect(viper.GetString("nats.url"))

	if err != nil {
		logrus.Fatalf("error while creating nats-connection: %s", err.Error())
	}

	//enc, err := nats.NewEncodedConn(nc, "json")

	// if err != nil {
	// 	logrus.Fatalf("error while creating encoding nats-connection: %s", err.Error())
	// }

	defer nc.Close()

	i := 1

	for {
		invalidMessage := InvalidMessage{Message: gofakeit.Color()}

		err = nc.Publish(viper.GetString("nats.subject"), []byte(invalidMessage.Message))

		if err != nil {
			logrus.Fatal("error while publish invalid fake message: %s", err)
		}

		fmt.Printf("added %d invalid message from: %s\n", i, time.Now().Format(time.RFC1123))
		i += 1

		time.Sleep(20 * time.Second)

	}

}
