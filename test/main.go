package main

import (
	"fmt"
	"time"

	"github.com/AlessandroRuggiero/go-socketcast"
)

func main() {
	pool := socketcast.CreatePool(&socketcast.Config{})
	fmt.Println(pool)
	pool.Start()
	time.Sleep(time.Second)
}
