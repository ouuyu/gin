package main

import (
	"log"
	"main/common"
	"main/model"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	common.SetupGinLog()
	common.SysLog("Gin Templates Version " + common.Version + " started.")

	err := model.InitDB()
	if err != nil {
		common.FatalLog("failed to init db: " + err.Error())
		return
	}
	defer func() {
		err := model.CloseDB()
		if err != nil {
			common.FatalLog("failed to close db: " + err.Error())
		}
	}()

	common.SysLog("db initialized")

	server := gin.Default()
	var port = os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(*common.Port)
	}
	err = server.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}
