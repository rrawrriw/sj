package main

import (
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

const (
	AppNamePrefix = "SJ"
)

func main() {
	host := os.Getenv("SJ_HOST")
	port := os.Getenv("SJ_PORT")
	srvRes := host + ":" + port

	publicDir := os.Getenv("SJ_PUBLIC_DIR")
	htmlDir := path.Join(publicDir, "html")

	srv := gin.Default()
	srv.Use(Serve("/", LocalFile(htmlDir, false)))
	srv.Static("/public", publicDir)

	srv.Run(srvRes)
}
