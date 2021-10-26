package main

import (
	"github.com/pallat/todowasm/todo"
	dom "honnef.co/go/js/dom/v2"
)

var signal = make(chan int)

func keepAlive() {
	for {
		<-signal
	}
}

func main() {
	d := dom.GetWindow().Document()
	app := todo.New(d)
	_ = app

	keepAlive()
}
