package main

import (
	"example/todo/internal/config"
	"fmt"
)

func main() {
	cfg := config.LoadConfig()

	fmt.Println(cfg)

}