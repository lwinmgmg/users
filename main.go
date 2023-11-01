package main

import (
	"fmt"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/lwinmgmg/user/controllers/v1"
	"github.com/lwinmgmg/user/env"
	"github.com/lwinmgmg/user/grpc/server"
)

var (
	Env = env.GetEnv()
)

func main() {
	var wg sync.WaitGroup
	app := gin.Default()
	app.Use(cors.Default())
	v1.DefineRoutes(app)
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Run(fmt.Sprintf("%v:%v", Env.Settings.Host, Env.Settings.Port))
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.StartServer()
	}()
	wg.Wait()
}
