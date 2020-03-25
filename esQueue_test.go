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
	wgPut.Add(goRoutineCnt)
	b.ResetTimer()

	put := func(start, end int) {
		for j := start; j < end; j++ {
			ok, _, _, _ := q.Put(j)
			for !ok {
				ok, _, _, _ = q.Put(j)
			}
		}
		wgPut.Done()
	}

	num := b.N / goRoutineCnt
	for i := 0; i < goRoutineCnt-1; i++ {
		go put(i*num, (i+1)*num)
	}
	go put((goRoutineCnt-1)*num, goRoutineCnt*num+(b.N%goRoutineCnt))

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
			// fmt.Printf("[%d]  xxxxxxx\n", i)
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
