package main

import "sheinko.tk/copy_project/controllers"

var handler = controllers.Handler{}

func main() {
	handler.Initialize()
	handler.Run(":8888")
}
