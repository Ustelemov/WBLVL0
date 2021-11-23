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
	Host                           string
	Port                           string
	Username                       string
	Password                       string
	DBName                         string
	SSLMode                        string
	Listener_Max_Reconnect_Seconds int
	Listener_Min_Reconnect_Seconds int
	Listener_Ping_NoEvent_Seconds  int
}

const (
	ordersJsonTable = "ordersjson"
)

type PostgresDB struct {
	DB       *sqlx.DB
	listener *pq.Listener
	config   Config
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

	return &PostgresDB{DB: db, listener: listener, config: cfg}, nil

}

func createListner(config Config) *pq.Listener {
	conString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Username, config.Password, config.Host,
		config.Port, config.DBName)

	minRecInterval := time.Duration(config.Listener_Min_Reconnect_Seconds) * time.Second
	maxRecInterval := time.Duration(config.Listener_Max_Reconnect_Seconds) * time.Second

	listenerEventCallback := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			logrus.Printf("error from db-listener: %s", err.Error())
		}
	}

	return pq.NewListener(conString, minRecInterval, maxRecInterval, listenerEventCallback)

}

func (postgresDB *PostgresDB) Close() {
	if postgresDB.listener != nil {
		postgresDB.listener.Close()
	}
}

func (postgresDB *PostgresDB) RunListenerDeamon(channel string, f func([]byte)) error {

	if err := postgresDB.listener.Listen(channel); err != nil {
		return fmt.Errorf("error while start listening to %s channel", channel)
	}

	go func() {
		for {
			select {
			case notification := <-postgresDB.listener.Notify:
				f([]byte(notification.Extra))
			case <-time.After(time.Duration(postgresDB.config.Listener_Ping_NoEvent_Seconds) * time.Second):
				go func() {
					err := postgresDB.listener.Ping()
					if err != nil {
						logrus.Fatalf("error while connecting to listener: %s", err)
					}
				}()
			}
		}
	}()

	return nil
}

func (postgresDB *PostgresDB) SaveOrderInRepository(ord *schema.Order) error {
	json, err := json.Marshal(ord)

	if err != nil {
		return fmt.Errorf("error while saving order in repository: cannot marshal Order to JSON")
	}

	uuid := ord.OrderUID

	isExists, err := postgresDB.CheckOrderExists(ord)

	if err != nil {
		return err
	}

	if isExists {
		return fmt.Errorf("error while saving order in repository: cannot save cause order already in repository")
	}

	query := fmt.Sprintf("INSERT INTO %s (uuid, order_data) VALUES ($1,$2)", ordersJsonTable)

	_, err = postgresDB.DB.Exec(query, uuid, json)

	if err != nil {
		return fmt.Errorf("error while saving order in repository: %s", err)
	}

	return nil

}

func (postgresDB *PostgresDB) CheckOrderExists(ord *schema.Order) (bool, error) {

	orderJSON := schema.OrderJSON{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE uuid = $1", ordersJsonTable)

	err := postgresDB.DB.Get(&orderJSON, query, ord.OrderUID)

	if err == nil {
		return true, nil
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return false, fmt.Errorf("error while checking order-exists in database: %s", err)
}

func (postgresDB *PostgresDB) GetAllOrders() ([]schema.OrderJSON, error) {
	ordersJsonArr := make([]schema.OrderJSON, 0)

	query := fmt.Sprintf("SELECT * FROM %s", ordersJsonTable)
	err := postgresDB.DB.Select(&ordersJsonArr, query)

	if err != nil {
		return nil, fmt.Errorf("error while getting all orders from database: %s", err.Error())
	}

	return ordersJsonArr, nil
}
