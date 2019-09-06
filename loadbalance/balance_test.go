package loadbalance

import (
	"context"
	"fmt"
	"micro/registry"
	"sync"
	"testing"

	"github.com/luci/go-render/render"
)

func TestRandomTest(t *testing.T) {
	balance := &RandomBalance{}

	var weights = [3]int{50, 100, 150}

	var nodes []*registry.Node

	for i := 0; i < 4; i++ {
		node := &registry.Node{
			IP:     fmt.Sprintf("127.0.0.%d", i),
			Port:   8080,
			Weight: weights[i%3],
		}

		fmt.Println("node ", render.Render(node))
		nodes = append(nodes, node)
	}

	cntStat := make(map[string]int)
	for i := 0; i < 10000; i++ {
		node, err := balance.Select(context.TODO(), nodes)
		if err != nil {
			t.Fatalf("select failed, %v", err)
			continue
		}

		cntStat[node.IP]++
	}

	for key, val := range cntStat {
		fmt.Println(" key ", key, " val ", val)
	}
}

func TestRoundTest(t *testing.T) {
	balance := &RoundRobinBalance{
		index:     0,
		gcd:       make(chan int, 1),
		maxWeight: make(chan int, 1),
		mu:        new(sync.Mutex),
		curWeight: -1,
	}

	var weights = [4]int{2, 4, 8, 10}

	var nodes []*registry.Node

	for i := 1; i < 5; i++ {
		node := &registry.Node{
			IP:     fmt.Sprintf("127.0.0.%d", i),
			Port:   8080,
			Weight: weights[i%4],
		}

		fmt.Println("node ", render.Render(node))
		nodes = append(nodes, node)
	}

	cntStat := make(map[string]int)
	for i := 0; i < 10000; i++ {
		node, err := balance.Select(context.TODO(), nodes)
		if err != nil {
			t.Fatalf("select failed, %v", err)
			continue
		}

		cntStat[node.IP]++
	}

	for key, val := range cntStat {
		fmt.Println(" key ", key, " val ", val)
	}
}
