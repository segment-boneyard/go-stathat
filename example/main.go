package main

import "github.com/tj/go-gracefully"
import "time"
import ".."
import "os"

func main() {
	key := os.Args[1]
	stats := stathat.New(key)
	quit := make(chan bool)
	stats.Verbose = true

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				stats.Count("something", 5)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				stats.Value("whatever", 123)
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	gracefully.Shutdown()

	close(quit)
	stats.Close()
}
