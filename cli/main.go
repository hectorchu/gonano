package main

import (
	"fmt"
	"log"

	"gonano/rpc"
)

func main() {
	rpc := rpc.Client{URL: "https://mynano.ninja/api/node"}
	_, count, _, err := rpc.BlockCount()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(count)
}
