package redis_socket

import (
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/charging_pile_tss/base"
	"github.com/giskook/charging_pile_tss/conf"
	"github.com/giskook/charging_pile_tss/pb"
	"github.com/golang/protobuf/proto"
	"log"
	"sync"
	"time"
)

type RedisSocket struct {
	conf                *conf.RedisConfigure
	Pool                *redis.Pool
	ChargingPiles       []*Report.ChargingPileStatus
	Mutex_ChargingPiles sync.Mutex

	ChargingPilesChan      chan *Report.ChargingPileStatus
	ChargingCost           chan *base.ChargingCost
	StopChargingNotifyChan chan *base.StopChargingNotify

	ticker *time.Ticker
}

var G_RedisSocket *RedisSocket = nil

func GetRedisSocket() *RedisSocket {
	return G_RedisSocket
}
func NewRedisSocket(config *conf.RedisConfigure) (*RedisSocket, error) {
	if G_RedisSocket == nil {
		G_RedisSocket = &RedisSocket{
			conf: config,
			Pool: &redis.Pool{
				MaxIdle:     config.MaxIdle,
				MaxActive:   config.MaxActive,
				IdleTimeout: time.Duration(config.IdleTimeOut) * time.Second,
				Dial: func() (redis.Conn, error) {
					c, err := redis.Dial("tcp", config.Addr)
					if err != nil {
						log.Println(err.Error())
						return nil, err
					}

					if len(config.Passwd) > 0 {
						if _, err := c.Do("AUTH", config.Passwd); err != nil {
							log.Println(err.Error())
							c.Close()
							return nil, err
						}
					}

					return c, err
				},
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					if time.Since(t) < time.Minute {
						return nil
					}

					c.Do("PING")

					return nil
				},
			},
			ChargingPiles:     make([]*Report.ChargingPileStatus, 0),
			ChargingPilesChan: make(chan *Report.ChargingPileStatus),
			ticker:            time.NewTicker(time.Duration(config.TranInterval) * time.Second),
		}
	}

	return G_RedisSocket, nil
}

func (socket *RedisSocket) DoWork() {
	defer func() {
		socket.Close()
	}()

	for {
		select {
		case <-socket.ticker.C:
			go socket.ProcessChargingPile()
			go GetStatusChecker().Check()
		case p := <-socket.ChargingPilesChan:
			socket.Mutex_ChargingPiles.Lock()
			socket.ChargingPiles = append(socket.ChargingPiles, p)
			socket.Mutex_ChargingPiles.Unlock()
		}
	}
}

func (socket *RedisSocket) GetConn() redis.Conn {
	return socket.Pool.Get()
}

func (socket *RedisSocket) SetFeedbackChan(ch chan *base.ChargingCost) {
	socket.ChargingCost = ch
}

func (socket *RedisSocket) SetStopChargingNotifyChan(ch chan *base.StopChargingNotify) {
	socket.StopChargingNotifyChan = ch
}

func (socket *RedisSocket) Close() {
	close(socket.ChargingPilesChan)
	socket.ticker.Stop()
}

func (socket *RedisSocket) RecvNsqChargingPile(message []byte) {
	charging_pile_status := &Report.ChargingPileStatus{}
	err := proto.Unmarshal(message, charging_pile_status)
	if err != nil {
		log.Println("recv nsq unmarshal error")
	} else {
		log.Printf("<IN NSQ> %s %d \n", charging_pile_status.DasUuid, charging_pile_status.Cpid)
		socket.ChargingPilesChan <- charging_pile_status
		GetStatusChecker().Insert(charging_pile_status.Cpid, charging_pile_status.Timestamp, time.Now().Unix(), charging_pile_status.Id, charging_pile_status.StationId)
	}
}
