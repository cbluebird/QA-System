package midwares

import (
	"QA-System/app/apiException"
	"QA-System/app/services/sessionService"
	"QA-System/app/utils"
	"github.com/gin-gonic/gin"
)

func CheckLogin(c *gin.Context) {
	isLogin := sessionService.CheckUserSession(c)
	if !isLogin {
		utils.JsonErrorResponse(c, apiException.NotLogin)
		c.Abort()
		return
	}
	c.Next()
}
