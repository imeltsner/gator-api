package main

import (
	"fmt"

	"github.com/imeltsner/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		return
	}

	err = cfg.SetUser("isaac")
	if err != nil {
		return
	}

	cfg, err = config.Read()
	if err != nil {
		return
	}

	fmt.Printf("DB is: %v\nUser is %v\n", cfg.DBURL, cfg.CurrentUsername)
}
