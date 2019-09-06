package loadbalance

import (
	"context"
	"errors"
	"micro/registry"
	"sync"
)

type RoundRobinBalance struct {
	index     int
	mu        *sync.Mutex
	gcd       chan int
	curWeight int
	maxWeight chan int
}

func (r *RoundRobinBalance) Name() string {
	return "roundrobin"
}

func gcdNormal(x, y int) int {
	var n int
	if x > y {
		n = y
	} else {
		n = x
	}
	for i := n; i >= 1; i-- {
		if x%i == 0 && y%i == 0 {
			return i
		}
	}
	return 1
}

func (r *RoundRobinBalance) getGcd() int {
	select {
	case val := <-r.gcd:
		return val
	}
}

func (r *RoundRobinBalance) calGcd(nodes []*registry.Node) {
	//TODO未实现各个节点的权重的最大公约数
	r.gcd <- 2
}

func (r *RoundRobinBalance) getMaxWeight() int {
	select {
	case val := <-r.maxWeight:
		return val
	}
}

func (r *RoundRobinBalance) calMaxWeight(nodes []*registry.Node) {
	max := -1

	for _, v := range nodes {
		if v.Weight > max {
			max = v.Weight
		}
	}
	r.maxWeight <- max
}

func (r *RoundRobinBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
	if len(nodes) == 0 {
		err = errors.New("ErrNotHaveNodes")
		return
	}

	go r.calMaxWeight(nodes)
	go r.calGcd(nodes)
	r.mu.Lock()
	defer r.mu.Unlock()
	gcd := r.getGcd()
	maxWeight := r.getMaxWeight()
	for {
		r.index = (r.index + 1) % len(nodes)

		if r.index == 0 {
			r.curWeight = r.curWeight - gcd
			if r.curWeight <= 0 {
				r.curWeight = maxWeight
			}
		}
		node = nodes[r.index]
		if node.Weight >= r.curWeight {
			return
		}
	}
	return
}
