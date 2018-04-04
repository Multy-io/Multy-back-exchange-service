package main

import (
	core "Multy-back-exchange-service/core"
)

var manager = core.NewManager()


func main() {

	manager.StartListen()

}

