package server

import (
	"log"
	"sync"

	"github.com/bitly/go-nsq"
	"github.com/giskook/charging_pile_tss/conf"
)

type NsqConsumer struct {
	config    *conf.ConsumerConf
	channel   string
	waitGroup *sync.WaitGroup
	consumer  *nsq.Consumer
}

func NewNsqConsumer(conf *conf.ConsumerConf, ch string) *NsqConsumer {
	return &NsqConsumer{
		config:    conf,
		channel:   ch,
		waitGroup: &sync.WaitGroup{},
	}
}

func (s *NsqConsumer) recvNsq() {
	s.consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		data := message.Body
		GetServer().Redis.RecvNsqChargingPile(data)
		return nil
	}))
}

func (s *NsqConsumer) Start() {
	s.waitGroup.Add(1)
	defer func() {
		s.waitGroup.Done()
		err := recover()
		if err != nil {
			log.Println("err found")
			s.Stop()
		}

	}()

	config := nsq.NewConfig()

	var errmsg error
	s.consumer, errmsg = nsq.NewConsumer(s.config.Topic, s.channel, config)

	if errmsg != nil {
		//	panic("create consumer error -> " + errmsg.Error())
		log.Println("create consumer error -> " + errmsg.Error())
	}
	s.recvNsq()

	err := s.consumer.ConnectToNSQD(s.config.Addr)
	if err != nil {
		panic("Counld not connect to nsq -> " + err.Error())
	}
}

func (s *NsqConsumer) Stop() {
	s.waitGroup.Wait()

	errmsg := s.consumer.DisconnectFromNSQD(s.config.Addr)

	if errmsg != nil {
		log.Printf("stop consumer error ", errmsg.Error())
	}

	s.consumer.Stop()
}
