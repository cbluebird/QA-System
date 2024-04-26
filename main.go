package main

import (
	"QA-System/app/midwares"
	"QA-System/config/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	r:=gin.Default()
	r.Use(midwares.ErrHandler())
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)

	err:=r.Run()
	if err != nil {
		log.Fatal("ServerStartFailed", err)
	}

}