package main

import (
	"fmt"
	"sync"
)

type BlockingQueue struct {
	m    sync.Mutex
	c    sync.Cond
	data []interface{}
	size int
}

func NewBlockingQueue(capacity int) *BlockingQueue {
	q := new(BlockingQueue)
	q.c = sync.Cond{L: &q.m}
	q.size = capacity
	return q
}

// methods
func (q *BlockingQueue) Put(item interface{}) {
	q.c.L.Lock()
	defer q.c.L.Unlock()

	for q.isFull() {
		q.c.Wait()
	}

	q.data = append(q.data, item)
	q.c.Broadcast()
}

func (q *BlockingQueue) Take() interface{} {
	q.c.L.Lock()
	defer q.c.L.Unlock()

	for q.isEmpty() {
		q.c.Wait()
	}

	result := q.data[0]
	q.data = q.data[1:len(q.data)]
	q.c.Broadcast()
	return result
}

func (q *BlockingQueue) isFull() bool {
	return q.size == len(q.data)
}

func (q *BlockingQueue) isEmpty() bool {
	return len(q.data) == 0
}

func main() {
	fmt.Println("This is the prototype of Blocking Bounded Queue using mutex and cond in Go")

	//using go routines

	bq := NewBlockingQueue(1)
	var wg sync.WaitGroup

	// Producers
	for i:=0; i<10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bq.Put("A")
		}()

		
	}

	// Consumers
	for i:=0; i<10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			val := bq.Take() // blocking call
			fmt.Printf("Got %v %v\n", val, i)
		}(i)
	}

	wg.Wait()
	
	
	/*
	bq := NewBlockingQueue(1)
	done := make(chan bool)

	// slow writer
	go func() {
		bq.Put("A")
		time.Sleep(1000 * time.Millisecond)
		bq.Put("B")
		time.Sleep(1000 * time.Millisecond)
		bq.Put("C")
	}()

	// reader will be blocked
	go func() {
		item := bq.Take()
		fmt.Printf("Got %v\n", item)
		item = bq.Take()
		fmt.Printf("Got %v\n", item)
		item = bq.Take()
		fmt.Printf("Got %v\n", item)
		done <- true
	}()

	// block while done
	<-done*/

}
/*
🔹 Key Takeaway

-> Signal() is fine if you can guarantee only one kind of waiter (e.g., one consumer thread waiting for items).

-> With multiple producers and consumers, you can’t guarantee the signal wakes the right type → use Broadcast().
*/