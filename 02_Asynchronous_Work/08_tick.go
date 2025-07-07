package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	i := 0
	for tickTime := range ticker.C {
		i++
		fmt.Println("step", i, "time", tickTime)
		if i >= 5 {
			// Надо останавливать, иначе потечет
			ticker.Stop()
			break
		}
	}
	fmt.Println("total", i)

	//return

	// Не может быть остановлен и собран сборщиком мусора
	// использовать, если должен работать вечно
	c := time.Tick(time.Second)
	i = 0
	for tickTime := range c {
		i++
		fmt.Println("step", i, "time", tickTime)
		if i >= 5 {
			break
		}
	}
}
