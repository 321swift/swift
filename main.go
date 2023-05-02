package main

import (
	"log"
	"os"
	"swift/node"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	node := node.NewNode(infoLog, errorLog)
	node.Start()
}
