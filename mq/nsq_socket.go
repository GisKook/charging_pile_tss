package mq

import (
	"github.com/bitly/go-nsq"
	"github.com/giskook/charging_pile_tss/base"
	"github.com/giskook/charging_pile_tss/conf"
	"github.com/giskook/charging_pile_tss/db"
	"github.com/giskook/charging_pile_tss/pb"
	"github.com/giskook/charging_pile_tss/redis_socket"
	"github.com/golang/protobuf/proto"
	"log"
)

type NsqSocket struct {
	conf                   *conf.NsqConfiguration
	Consumers              []*NsqConsumer
	Producers              []*NsqProducer
	TransactionChan        chan string
	ChargingCostChan       chan *base.ChargingCost
	StopChargingNotifyChan chan *base.StopChargingNotify
}

var G_NsqSocket *NsqSocket = nil

func GetNsqSocket(config *conf.NsqConfiguration) *NsqSocket {
	nsq_producer_count := config.Producer.Count
	if G_NsqSocket == nil {
		G_NsqSocket = &NsqSocket{
			conf:                   config,
			Consumers:              make([]*NsqConsumer, 0),
			Producers:              make([]*NsqProducer, nsq_producer_count),
			TransactionChan:        make(chan string),
			ChargingCostChan:       make(chan *base.ChargingCost, 1024),
			StopChargingNotifyChan: make(chan *base.StopChargingNotify, 1024),
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

	producer := NewNsqProducer(socket.conf.Producer)
	producer.Start()
	socket.Producers[0] = producer
	socket.consumerStart()
	db.GetDBSocket().SetNotifyChan(socket.TransactionChan)
	redis_socket.GetRedisSocket().SetFeedbackChan(socket.ChargingCostChan)
	redis_socket.GetRedisSocket().SetStopChargingNotifyChan(socket.StopChargingNotifyChan)
	go socket.Notify()
	go socket.ProccessCost()
	go socket.NotifyStopCharging()
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

func (socket *NsqSocket) Send(topic string, value []byte) error {
	err := socket.Producers[0].Send(topic, value)

	return err
}

func (socket *NsqSocket) Notify() {
	for {
		select {
		case notify := <-socket.TransactionChan:
			socket.SendNotify(notify)
		}

	}
}

func (socket *NsqSocket) SendNotify(transcation_ids string) {
	paras := []*Report.Param{
		&Report.Param{
			Type:    Report.Param_STRING,
			Strpara: transcation_ids,
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_NOTIFY_TRANSCATION,
		Paras: paras,
	}

	data, _ := proto.Marshal(command)

	socket.Send(conf.GetConf().Nsq.Producer.TopicWeChat, data)
}

func (socket *NsqSocket) ProccessCost() {
	for {
		select {
		case feedback := <-socket.ChargingCostChan:
			socket.SendCost(feedback)
		}
	}
}

func (socket *NsqSocket) SendCost(cost *base.ChargingCost) {
	paras := []*Report.Param{
		&Report.Param{
			Type:  Report.Param_UINT32,
			Npara: uint64(cost.Cost),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REP_CHARGING_COST,
		Tid:   cost.Tid,
		Paras: paras,
	}

	data, _ := proto.Marshal(command)

	socket.Send(cost.Uuid, data)
}

func (socket *NsqSocket) NotifyStopCharging() {
	for {
		select {
		case notify := <-socket.StopChargingNotifyChan:
			socket.SendNotifyStopCharging(notify)
		}
	}
}

func (socket *NsqSocket) SendNotifyStopCharging(stop_charging_notify *base.StopChargingNotify) {
	command := &Report.Command{
		Type: Report.Command_CMT_REQ_STOP_CHARGING,
		Tid:  stop_charging_notify.Tid,
	}

	data, _ := proto.Marshal(command)

	socket.Send(stop_charging_notify.Uuid, data)
}
