package utils

import (
	"QA-System/app/apiException"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JsonResponse(c *gin.Context, httpStatusCode int, code int, msg string, data interface{}) {
	c.JSON(httpStatusCode, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func JsonSuccessResponse(c *gin.Context, data interface{}) {
	JsonResponse(c, http.StatusOK, 200, "OK", data)
}

func JsonErrorResponse(c *gin.Context, err *apiException.Error) {
	JsonResponse(c, http.StatusOK, err.Code, err.Msg, nil)
}
