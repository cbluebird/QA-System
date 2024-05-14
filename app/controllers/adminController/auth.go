package adminController

import (
	"QA-System/app/apiException"
	"QA-System/app/models"
	"QA-System/app/services/adminService"
	"QA-System/app/services/sessionService"
	"QA-System/app/utils"
	"QA-System/config/config"
	"errors"

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
		c.Error(err)
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//判断密码是否正确
	user, err := adminService.GetAdminByUsername(data.Username)
	if err != nil {
		c.Error(err)
		if err == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, apiException.UserNotFind)
			return
		} else {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}
	if user.Password != data.Password {
		c.Error(errors.New("密码错误"))
		utils.JsonErrorResponse(c, apiException.NoThatPasswordOrWrong)
		return
	}
	//设置session
	err = sessionService.SetUserSession(c, user)
	if err != nil {
		c.Error(err)
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type RegisterData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Key      int    `json:"key" binding:"required"`
}

// 注册
func Register(c *gin.Context) {
	var data RegisterData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err)
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//判断是否有权限
	adminKey := config.Config.GetInt("key")
	if adminKey != data.Key {
		c.Error(errors.New("没有权限"))
		utils.JsonErrorResponse(c, apiException.NotSuperAdmin)
		return
	}
	//判断用户是否存在
	err = adminService.IsAdminExist(data.Username)
	if err == nil {
		c.Error(err)
		utils.JsonErrorResponse(c, apiException.UserExist)
		return
	}
	//创建用户
	err = adminService.CreateAdmin(models.User{
		Username:  data.Username,
		Password:  data.Password,
		AdminType: 1,
	})
	if err != nil {
		c.Error(err)
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
