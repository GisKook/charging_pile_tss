package mq

import (
	"github.com/bitly/go-nsq"
	"github.com/giskook/charging_pile_tss/conf"
	"github.com/giskook/charging_pile_tss/redis_socket"
	"log"
)

type NsqSocket struct {
	conf      *conf.NsqConfiguration
	Consumers []*NsqConsumer
}

var G_NsqSocket *NsqSocket = nil

func GetNsqSocket(config *conf.NsqConfiguration) *NsqSocket {
	if G_NsqSocket == nil {
		G_NsqSocket = &NsqSocket{
			conf:      config,
			Consumers: make([]*NsqConsumer, 0),
		}
	}

	return G_NsqSocket
}

func (socket *NsqSocket) Start() {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()

	socket.consumerStart()
}

func (socket *NsqSocket) consumerStart() {
	for _, ch := range socket.conf.Consumer.Channels {
		consumer_conf := &NsqConsumerConf{
			Addr:    socket.conf.Consumer.Addr,
			Topic:   socket.conf.Consumer.Topic,
			Channel: ch,
			Handler: nsq.HandlerFunc(func(message *nsq.Message) error {
				data := message.Body
				redis_socket.GetRedisSocket().RecvNsqChargingPile(data)

				return nil

			}),
		}
		consumer := NewNsqConsumer(consumer_conf)
		consumer.Start()
		socket.Consumers = append(socket.Consumers, consumer)
	}
}

func (socket *NsqSocket) Stop() {
	for _, consumer := range socket.Consumers {
		consumer.Stop()
	}
}
