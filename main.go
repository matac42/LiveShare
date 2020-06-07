package main

import (
	"github.com/gin-gonic/gin"
	"github.com/matac42/LiveShare/router"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("html/*.html")
	router.Router(r)
}
