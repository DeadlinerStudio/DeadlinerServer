package main

import (
	v1 "github.com/aritxonly/deadlinerserver/kitex_gen/deadliner/v1/deadlinerservice"
	"log"
)

func main() {
	svr := v1.NewServer(new(DeadlinerServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
