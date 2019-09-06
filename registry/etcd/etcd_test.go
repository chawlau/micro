package etcd

import (
	"context"
	"fmt"
	"micro/registry"
	"testing"
	"time"

	"github.com/luci/go-render/render"
)

func TestRegister(t *testing.T) {

	registryInst, err := registry.InitRegistry(context.TODO(), "etcd",
		registry.WithAddrs([]string{"127.0.0.1:2379"}),
		registry.WithTimeout(time.Second),
		registry.WithRegistryPath("/ibinarytree/koala/"),
		registry.WithHeartBeat(5),
	)
	if err != nil {
		t.Errorf("init registry failed, err:%v", err)
		return
	}

	service := &registry.Service{
		Name: "comment_service",
	}

	service.Nodes = append(service.Nodes, &registry.Node{
		IP:   "127.0.0.1",
		Port: 8801,
	},
		&registry.Node{
			IP:   "127.0.0.2",
			Port: 8802,
		},
	)
	registryInst.Register(context.TODO(), service)
	go func() {
		time.Sleep(time.Second * 10)
		service.Nodes = append(service.Nodes, &registry.Node{
			IP:   "127.0.0.3",
			Port: 8803,
		},
			&registry.Node{
				IP:   "127.0.0.4",
				Port: 8804,
			},
		)
		registryInst.Register(context.TODO(), service)
	}()
	for {
		service, err := registryInst.GetService(context.TODO(), "comment_service")
		if err != nil {
			t.Errorf("get service failed, err:%v", err)
			return
		}

		for _, node := range service.Nodes {
			fmt.Println("service nodes info ", service.Name, render.Render(node))
		}
		time.Sleep(time.Second * 5)
	}
}
