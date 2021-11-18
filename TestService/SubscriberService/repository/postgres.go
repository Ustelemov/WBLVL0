package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/ustelemov/WBLVL0/TestService/schema"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

const (
	ordersJsonTable = "ordersjson"
)

type PostgresDB struct {
	DB       *sqlx.DB
	listener *pq.Listener
}

func NewPostgresDB(cfg Config) (*PostgresDB, error) {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("host = %v port = %v user = %v dbname = %v password = %v sslmode = %v",
			cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, fmt.Errorf("error while connect to postgres: %s", err.Error())
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("error while ping to postgres: %s", err.Error())
	}

	listener := createListner(cfg)

	return &PostgresDB{DB: db, listener: listener}, nil

}

func createListner(config Config) *pq.Listener {
	conString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Username, config.Password, config.Host,
		config.Port, config.DBName)

	minRecInterval := 10 * time.Second
	maxRecInterval := time.Minute

	listenerEventCallback := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			logrus.Printf("error from db-listener: %s", err.Error())
		}
	}

	return pq.NewListener(conString, minRecInterval, maxRecInterval, listenerEventCallback)

}

func (postgresDB *PostgresDB) RunListenerDeamon(channel string, f func([]byte)) {

	if err := postgresDB.listener.Listen(channel); err != nil {
		logrus.Printf(err.Error())
	}

	go func() {
		for {
			select {
			case notification := <-postgresDB.listener.Notify:
				f([]byte(notification.Extra))
			case <-time.After(90 * time.Second):
				go func() {
					err := postgresDB.listener.Ping()
					if err != nil {
						logrus.Fatalf("error while connecting to listener: %s", err)
					}
				}()
			}
		}
	}()
}

func (postgresDB *PostgresDB) SaveOrderInRepository(ord *schema.Order) error {
	json, err := json.Marshal(ord)

	if err != nil {
		logrus.Error("error while marshaling Order to JSON")
	}

	uuid := ord.OrderUID

	isExists, err := postgresDB.CheckExists(ord)

	if err != nil {
		return err
	}

	if isExists {
		return fmt.Errorf("cannot save cause order already in repository")
	}

	query := fmt.Sprintf("INSERT INTO %s (uuid, order_data) VALUES ($1,$2)", ordersJsonTable)

	postgresDB.DB.QueryRow(query, uuid, json)

	return nil

}

func (postgresDB *PostgresDB) CheckExists(ord *schema.Order) (bool, error) {

	orderJSON := schema.OrderJSON{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE uuid = $1", ordersJsonTable)

	err := postgresDB.DB.Get(&orderJSON, query, ord.OrderUID)

	if err == nil {
		return true, nil
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return false, fmt.Errorf("error while checking exists: %s", err)
}

func (postgresDB *PostgresDB) GetAllOrders() ([]schema.OrderJSON, error) {
	ordersJsonArr := make([]schema.OrderJSON, 0)

	query := fmt.Sprintf("SELECT * FROM %s", ordersJsonTable)
	err := postgresDB.DB.Select(&ordersJsonArr, query)

	if err != nil {
		return nil, fmt.Errorf("error while getting all orders: %s", err.Error())
	}

	return ordersJsonArr, nil
}
