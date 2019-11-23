package main

import (
	"log"

	"go.i3wm.org/i3/v4"
)

func main() {
	recv := i3.Subscribe(i3.WindowEventType)
	for recv.Next() {
		ev := recv.Event().(*i3.WindowEvent)

		log.Printf("change: %s, %s", ev.Change, ev.Container.WindowProperties.Class)
	}
	log.Fatal(recv.Close())
}
