package main

import (
	"fmt"

	"github.com/melonmanchan/dippa-writer/src/config"
	_ "github.com/melonmanchan/dippa-writer/src/models"
)

func main() {
	config := config.ParseConfig()
	fmt.Printf("%v", config)
}
