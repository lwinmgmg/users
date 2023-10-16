package main

import (
	"fmt"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/user/controllers"
	"github.com/lwinmgmg/user/grpc/server"
	"github.com/lwinmgmg/user/models"
	"gorm.io/gorm"
)

var (
	wg sync.WaitGroup
)

func printSomething(partner *models.Partner, db *gorm.DB) {
	fmt.Println(partner.NextCode(db))
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	app := gin.Default()
	app.Use(cors.Default())
	controllers.DefineRoutes(app)
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Run("0.0.0.0:8888")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.StartServer()
	}()
	wg.Wait()
}
