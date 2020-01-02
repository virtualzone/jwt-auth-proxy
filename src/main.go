package main

import (
	"log"
	"os"
)

func main() {
	log.Println("Starting server...")
	a := GetApp()
	GetDatatabase().connectMongoDb(GetConfig().MongoDbURL, GetConfig().MongoDbName)
	a.InitializePublicRouter()
	a.InitializeBackendRouter()
	a.InitializeTimers()
	readMailTemplatesFromFile()
	a.Run(GetConfig().PublicListenAddr, GetConfig().BackendListenAddr)
	GetDatatabase().disconnect()
	os.Exit(0)
}
