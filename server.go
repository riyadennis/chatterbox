package main

import (
	"github.com/riyadennis/chatterbox/internal"
	"github.com/rsms/gotalk"
	"log"
)

func main(){
	server()
}

func server() {
	gotalk.Handle("greet", func(in internal.GreetIn) (internal.GreetOut, error) {
		return internal.GreetOut{"Hello " + in.Name}, nil
	})
	if err := gotalk.Serve("tcp", "localhost:1234", nil ); err != nil {
		log.Fatalln(err)
	}
}




