package infra

import (
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

type NatsMQ struct {
	nc *nats.Conn
}

func NewNatsMQ(vc *viper.Viper) *NatsMQ {
	host := vc.GetString("data.nats.host")
	port := vc.GetString("data.nats.port")
	nc, err := nats.Connect(host + ":" + port)
	if err != nil {
		panic(err)
	}
	nc.Subscribe("test", func(msg *nats.Msg) {})
	return &NatsMQ{
		nc: nc,
	}
}

func (n *NatsMQ) GetConn() *nats.Conn {
	return n.nc
}
