package main

import (
	"QA-System/app/midwares"
	"QA-System/config/database"
	"QA-System/config/router"
	"QA-System/config/session"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.MysqlInit()
	database.MongodbInit()
	r := gin.Default()
	r.Use(midwares.ErrHandler())
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)
	r.Static("/static", "./static")
	r.Static("/xlsx", "./xlsx")
	session.Init(r)
	router.Init(r)
	err := r.Run()
	if err != nil {
		log.Fatal("ServerStartFailed", err)
	}

}
