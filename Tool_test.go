package tool

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func Test_tool(t *testing.T) {
	std := log.New(os.Stderr, "", log.LstdFlags)
	fmt.Println("we will print the panic info")
	go func() {
		defer CutPanic(std, nil)
		panic("panic func1...")
	}()

	go func() {
		defer CutPanic(std, func() {
			fmt.Println("func2's after panic catch call...")
		})
		panic("panic func2...")
	}()
	time.Sleep(time.Second)
}

func Test_counter(t *testing.T) {
	cnt := NewCounter()
	fmt.Println(cnt.ReadV())

	cnt.Incr(2)
	fmt.Println(cnt.ReadV())

	cnt.Decr(1)
	fmt.Println(cnt.ReadV())

	fmt.Println(cnt.DecrToEqual(0))
	fmt.Println(cnt.ReadV())

	fmt.Println(cnt.EqualAndIncr(2))
	fmt.Println(cnt.ReadV())

	fmt.Println(cnt.EqualAndIncr(1))
	fmt.Println(cnt.ReadV())
}

func Test_Multi_counter(t *testing.T) {
	cnt := NewCounter()
	w := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func() {
			cnt.Incr(1)
			w.Done()
		}()
	}
	w.Wait()
	fmt.Println(cnt.ReadV())

	fmt.Println()
	for i := 0; i < 500; i++ {
		w.Add(1)
		go func() {
			cnt.Decr(1)
			w.Done()
		}()
	}
	w.Wait()
	fmt.Println(cnt.ReadV())
	fmt.Println()

	for i := 0; i < 500; i++ {
		w.Add(1)
		go func() {
			if cnt.DecrToEqual(1) {
				fmt.Println("get 0")
			}
			w.Done()
		}()
	}
	w.Wait()
	fmt.Println()

	for i := 0; i < 10; i++ {
		w.Add(1)
		go func() {
			cnt.EqualAndIncr(0)
			w.Done()
		}()
	}
	w.Wait()
	fmt.Println(cnt.ReadV())
}
