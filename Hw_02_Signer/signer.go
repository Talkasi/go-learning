package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	//"runtime"
)

type Chunk struct {
	N   int
	Val string
}

// Конвейерная обработка задач
func ExecutePipeline(jobs ...job) {
	//runtime.GOMAXPROCS(0)
	channels := make([]chan interface{}, len(jobs)+1)
	wg := &sync.WaitGroup{}

	for i, curJob := range jobs {
		channels[i+1] = make(chan interface{}, MaxInputDataLen)
		wg.Add(1)
		go func(in, out chan interface{}, waiter *sync.WaitGroup) {
			defer waiter.Done()
			curJob(in, out)
			close(out)
		}(channels[i], channels[i+1], wg)
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	var wg = &sync.WaitGroup{}
	var mu = &sync.Mutex{}

	for val := range in {
		intVal := strconv.Itoa(val.(int))
		wg.Add(1)

		go func(out chan interface{}) {
			defer wg.Done()
			mu.Lock()
			md5 := DataSignerMd5(intVal)
			mu.Unlock()

			innerWg := &sync.WaitGroup{}
			innerC := make(chan Chunk, 100)

			innerWg.Add(2)
			go func(out chan Chunk) {
				defer innerWg.Done()
				innerC <- Chunk{
					N: 0,
					Val: DataSignerCrc32(intVal),
				}
			} (innerC)

			go func(out chan Chunk) {
				defer innerWg.Done()
				innerC <- Chunk{
					N: 1,
					Val: DataSignerCrc32(md5),
				}
			} (innerC)

			innerWg.Wait()

			a := <- innerC
			b := <- innerC

			close(innerC)

			if a.N > b.N {
				a, b = b, a
			}

			out <- a.Val + "~" + b.Val
		}(out)
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	const nTh = 6

	mainWg := &sync.WaitGroup{}
	for val := range in {
		mainWg.Add(1)
		go func(out chan interface{}) {
			defer mainWg.Done()

			var wg = &sync.WaitGroup{}
			tmpOut := make(chan interface{}, nTh)
			for th := 0; th < nTh; th++ {
				wg.Add(1)
				go func(out chan interface{}) {
					//fmt.Println("Start", th, "on", time.Now())
					defer wg.Done()
					res := Chunk{
						N: th,
						Val: DataSignerCrc32(strconv.Itoa(th) + val.(string)),
					}
					out <- res
					//fmt.Println("End", th, "on", time.Now())
				}(tmpOut)
			}
			wg.Wait()
			
			//fmt.Println("Wait ended on", time.Now())
			chanks := make([]string, nTh) 
			for th := 0; th < nTh; th++ {
				c := (<-tmpOut).(Chunk)
				chanks[c.N] = c.Val
			}
			close(tmpOut)
			res := strings.Join(chanks, "")
	
			out <- res
		} (out)
	}

	mainWg.Wait()
}

func CombineResults(in, out chan interface{}) {
	elems := make([]string, MaxInputDataLen)
	n := 0
	for el := range in {
		elems[n] = el.(string)
		n++
	}

	elems = elems[:n]
	sort.Strings(elems)

	var res string
	for i := 0; i < n-1; i++ {
		res += elems[i] + "_"
	}
	res += elems[n-1]

	out <- res
}

func main() {
	var recieved uint32
	freeFlowJobs := []job{
		job(func(in, out chan interface{}) {
			out <- uint32(1)
			fmt.Println("1 sent by 0 job")
			out <- uint32(3)
			fmt.Println("3 sent by 0 job")
			out <- uint32(4)
			fmt.Println("4 sent by 0 job")
		}),
		job(func(in, out chan interface{}) {
			for val := range in {
				fmt.Printf("%v recieved by 1 job\n", val.(uint32))
				out <- val.(uint32) * 3
				fmt.Printf("%v sent by 1 job\n", val.(uint32)*3)
				time.Sleep(time.Millisecond * 100)
			}
		}),
		job(func(in, out chan interface{}) {
			for val := range in {
				fmt.Println("collected", val)
				atomic.AddUint32(&recieved, val.(uint32))
			}
		}),
	}

	ExecutePipeline(freeFlowJobs...)
}
