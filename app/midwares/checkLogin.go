package midwares

import (
	"QA-System/app/apiException"
	"QA-System/app/services/sessionService"
	"QA-System/app/utils"
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckLogin(c *gin.Context) {
	isLogin := sessionService.CheckUserSession(c)
	if !isLogin {
		c.Error(errors.New("未登录"))
		utils.JsonErrorResponse(c, apiException.NotLogin)
		c.Abort()
		return
	}
	c.Next()
}
