package main

import (
	"fmt"
	"math"
	"time"
)

type BucketLimit struct {
	rate       float64
	bucketSize float64
	unixNano   int64
	curWater   float64
}

func NewBucketLimit(rate float64, bucketSize int64) *BucketLimit {
	return &BucketLimit{
		bucketSize: float64(bucketSize),
		rate:       rate,
		unixNano:   time.Now().UnixNano(),
		curWater:   0,
	}
}

func (b *BucketLimit) reflesh() {
	now := time.Now().UnixNano()
	diffSec := float64(now-b.unixNano) / 1000 / 1000 / 1000
	b.curWater = math.Max(0, b.curWater-diffSec*b.rate)
	b.unixNano = now
	return
}

func (b *BucketLimit) Allow() bool {
	b.reflesh()

	if b.curWater < b.bucketSize {
		b.curWater = b.curWater + 1
		return true
	}
	return false
}

func main() {
	limit := NewBucketLimit(50, 100)
	m := make(map[int]bool)

	for i := 0; i < 1000; i++ {
		allow := limit.Allow()
		if allow {
			m[i] = true
			continue
		}
		m[i] = false
		time.Sleep(time.Millisecond * 10)
	}

	for i := 0; i < 1000; i++ {
		fmt.Println(" i ", i, " allow ", m[i])
	}
}
