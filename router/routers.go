package router

import (
	"github.com/matac42/LiveShare/oauth"

	"github.com/gin-gonic/gin"
)

// Router create gin router routine.
func Router(r *gin.Engine) {
	conf := oauth.CreateConf()

	v1 := r.Group("/")
	{
		v1.GET("oauthhome", oauth.HomeClient)
		v1.GET("google", conf.Google)
		v1.GET("callback", conf.CallBack)

	}

	r.Run(":8080")

}
