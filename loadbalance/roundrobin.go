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

func gcd1(m int, n int) int {
	var r int
	for n > 0 {
		r = m % n
		m = n
		n = r
	}
	return m
}

func gcdNormal(m int, n int) int {
	if n > 0 {
		return gcd1(n, m%n)
	} else {
		return m
	}
}

func (r *RoundRobinBalance) getGcd() int {
	select {
	case val := <-r.gcd:
		return val
	}
}

func (r *RoundRobinBalance) calGcd(nodes []*registry.Node) {
	if len(nodes) < 2 {
		r.gcd <- nodes[0].Weight
		return
	}

	x := nodes[0].Weight
	for _, val := range nodes {
		if val.Weight > 0 && x > 0 {
			x = gcdNormal(x, val.Weight)
		} else {
			r.gcd <- 0
			return
		}
	}

	r.gcd <- x
	return
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
