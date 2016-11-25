package server

import (
	"github.com/giskook/charging_pile_tss/redis_socket"
)

type Server struct {
	MQ    *NsqSocket
	Redis *redis_socket.RedisSocket
}

var g_server *Server

func NewServer(nsq_socket *NsqSocket, redis_socket *redis_socket.RedisSocket) *Server {
	g_server = &Server{
		MQ:    nsq_socket,
		Redis: redis_socket,
	}

	return g_server
}

func (server *Server) Start() {
	server.MQ.Start()
	go server.Redis.DoWork()
}

func (server *Server) Stop() {
	server.MQ.Stop()
}

func GetServer() *Server {
	return g_server
}
