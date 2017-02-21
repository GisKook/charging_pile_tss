package db

import (
	"database/sql"
	"fmt"
	"github.com/giskook/charging_pile_tss/base"
	"github.com/giskook/charging_pile_tss/conf"
	"github.com/lib/pq"
	"log"
	"strings"
	"sync"
	"time"
)

type DbSocket struct {
	Db             *sql.DB
	ChargingPrices map[uint64][]*base.ChargingPrice

	Listener *pq.Listener
}

var G_DBSocket_Mutex sync.Mutex
var G_DBSocket *DbSocket

func NewDbSocket(db_config *conf.DBConfigure) (*DbSocket, error) {
	defer G_DBSocket_Mutex.Unlock()
	G_DBSocket_Mutex.Lock()
	conn_string := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", db_config.User, db_config.Passwd, db_config.Host, db_config.Port, db_config.DbName)

	log.Println(conn_string)
	db, err := sql.Open(db_config.User, conn_string)

	if err != nil {
		return nil, err
	}
	log.Println("db open success")

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println(err.Error())
		}
	}

	G_DBSocket = &DbSocket{
		Db:             db,
		Listener:       pq.NewListener(conn_string, 60*time.Second, time.Minute, reportProblem),
		ChargingPrices: make(map[uint64][]*base.ChargingPrice),
	}

	return G_DBSocket, nil
}

func (db_socket *DbSocket) Listen(table string) error {
	return db_socket.Listener.Listen(table)
}

func (db_socket *DbSocket) ProcessNotify(_notify string) {
	log.Println(_notify)
	if strings.Contains(_notify, conf.GetConf().DB.ListenPriceTable) {
		notify := strings.TrimPrefix(_notify, conf.GetConf().DB.ListenPriceTable+"^")
		log.Println(notify)
		db_socket.parse_payload_price(notify)
	}

}

func (db_socket *DbSocket) WaitForNotification() {
	for {
		select {
		case notify := <-db_socket.Listener.Notify:
			db_socket.ProcessNotify(notify.Extra)
			break
		case <-time.After(90 * time.Second):
			go func() {
				db_socket.Listener.Ping()
			}()
			// Check if there's more work available, just in case it takes
			// a while for the Listener to notice connection loss and
			// reconnect.
			log.Println("received no work for 90 seconds, checking for new work")
			break
		}
	}
}

func (db_socket *DbSocket) Close() {
	db_socket.Listener.Close()
	db_socket.Db.Close()
}

func GetDBSocket() *DbSocket {
	defer G_DBSocket_Mutex.Unlock()
	G_DBSocket_Mutex.Lock()

	return G_DBSocket
}
