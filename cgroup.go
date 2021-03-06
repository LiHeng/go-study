package main

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
)

const (
	MB = 1024 * 1024
)

func main() {
	// cobra_repo.Execute()
	// blocks := make([][MB]byte, 0)
	// fmt.Println("Child pid is", os.Getpid())
	//
	// for range time.Tick(time.Second) {
	// 	blocks = append(blocks, [MB]byte{})
	// 	printMemUsage()
	// }
	engine := gin.Default()
	_ = engine.Run()
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tSys = %v MiB \n", bToMb(m.Sys))
}

func bToMb(b uint64) uint64 {
	return b / MB
}
