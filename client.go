package main

import (
	"github.com/riyadennis/chatterbox/internal"
	"github.com/rsms/gotalk"
	"log"
)

func main(){
	client()
}

func client() {
	s, err := gotalk.Connect("tcp", "localhost:1234")
	if err != nil {
		log.Fatalln(err)
	}
	greeting := &internal.GreetOut{}
	if err := s.Request("greet", internal.GreetIn{"Rasmus"}, greeting); err != nil {
		log.Fatalln(err)
	}
	log.Printf("greeting: %+v\n", greeting)
	s.Close()
}