package main

import (
	"fmt"
	"sync"
	"testing"

	"test/queue"
)

var goRoutineCnt = 1000

func checkPutAndGet(checkData []int) {
	l := len(checkData)
	for i := 0; i < l; i++ {
		if checkData[i] != i {
			fmt.Printf("\ncheckPutAndGet error!! lost data!!! i:%v checkData[i]:%v l:%v\n", i, checkData[i], l)
		}
	}
}
func BenchmarkEsQueueReadContention(b *testing.B) {
	fmt.Printf("enter : %d\n", b.N)
	var checkData = make([]int, b.N)
	q := queue.NewQueue(1024 * 1024)
	var wgGet sync.WaitGroup
	wgGet.Add(goRoutineCnt)
	var wgPut sync.WaitGroup
	wgPut.Add(1)
	b.ResetTimer()

	go func() {
		for i := 0; i < b.N; i++ {
			ok, n, p, g := q.Put(i)
			for !ok {
				fmt.Printf("---- %d  max: %d, i:%d, pusPos: %d, getPos: %d\n", n, b.N, i, p, g)
				ok, n, p, g = q.Put(i)
				fmt.Printf("---- %d  max: %d, i:%d, pusPos: %d, getPos: %d\n", n, b.N, i, p, g)
			}
		}
		wgPut.Done()
	}()

	for i := 0; i < goRoutineCnt; i++ {
		go func(i int) {
			for j := 0; j < b.N/goRoutineCnt; j++ {
				val, ok, _, _, _ := q.Get()
				for !ok {
					val, ok, _, _, _ = q.Get()
					// fmt.Printf("[%d]+++++ %d\n", i, f)
				}
				v := val.(int)
				checkData[v] = v
			}
			fmt.Printf("[%d]  xxxxxxx\n", i)
			wgGet.Done()
		}(i)
	}
	wgGet.Wait()
	fmt.Printf("~~~~~\n")
	wgPut.Wait()

	fmt.Printf("123123123123\n")
	for q.Quantity() > 0 {
		val, ok, _, _, _ := q.Get()
		for !ok {
			val, ok, _, _, _ = q.Get()
		}
		v := val.(int)
		checkData[v] = v
	}
	checkPutAndGet(checkData)
}
