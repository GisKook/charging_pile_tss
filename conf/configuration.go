package conf

import (
	"encoding/json"
	"os"
)

type ProducerConf struct {
	Addr        string
	Count       int
	TopicWeChat string
}

type ConsumerConf struct {
	Addr     string
	Topic    string
	Channels []string
}

type NsqConfiguration struct {
	Producer *ProducerConf
	Consumer *ConsumerConf
}

type RedisConfigure struct {
	Addr              string
	MaxIdle           int
	MaxActive         int
	IdleTimeOut       int
	Passwd            string
	TranInterval      int
	OfflineThreshold  int
	ReOnlineThreshold int
}

type DBConfigure struct {
	Host             string
	Port             string
	User             string
	Passwd           string
	DbName           string
	ListenPriceTable string
	TransInterval    uint8
}

type Configuration struct {
	Uuid  string
	Nsq   *NsqConfiguration
	Redis *RedisConfigure
	DB    *DBConfigure
}

var G_conf *Configuration

func ReadConfig(confpath string) (*Configuration, error) {
	file, _ := os.Open(confpath)
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	G_conf = &config

	return &config, err
}

func GetConf() *Configuration {
	return G_conf
}
