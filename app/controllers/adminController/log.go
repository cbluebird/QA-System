package adminController

import (
	"QA-System/app/apiException"
	"QA-System/app/services/adminService"
	"QA-System/app/services/sessionService"
	"QA-System/app/utils"

	"github.com/gin-gonic/gin"
)

func GetLogMsg(c *gin.Context){
	//鉴权
	_, err := sessionService.GetUserSession(c)
	if err != nil {
		c.Error(err)
		utils.JsonErrorResponse(c, apiException.NotLogin)
		return
	}
	data, err := adminService.GetLastLinesFromLogFile("app.log", 5)
	if err != nil {
		c.Error(err)
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, data)
}