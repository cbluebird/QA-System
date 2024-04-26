package adminController

import (
	"QA-System/app/apiException"
	"QA-System/app/services/adminService"
	"QA-System/app/services/sessionService"
	"QA-System/app/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录
func Login(c *gin.Context) {
	var data LoginData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//判断密码是否正确
	user, err := adminService.GetAdminByUsername(data.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, apiException.UserNotFind)
			return
		} else {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}
	if user.Password != data.Password {
		utils.JsonErrorResponse(c, apiException.NoThatPasswordOrWrong)
		return
	}
	//设置session
	err = sessionService.SetUserSession(c, user)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
