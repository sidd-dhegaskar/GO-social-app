package main

import (
	"fmt"
	"log"
)

func main(){
cfg:=config{
	addr:":8080",
}

	app:=&application{
		config:cfg,
	}
	mux:=app.mount()
	log.Fatal(app.run(mux))
	fmt.Println("server has start")
}
