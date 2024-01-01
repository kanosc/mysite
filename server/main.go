// server project main.go
package main

import (
	"fmt"

	"flag"

	"github.com/kanosc/mysite/server/router"
)

var mode = flag.String("mode", "production", "usage: -mode debug|production, production as default")

func main() {
	flag.Parse()
	router.InitRouter(*mode)
	router.Start(":5173")

	fmt.Println("Hello World!")
}
