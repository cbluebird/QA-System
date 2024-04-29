package userController

import (
	"QA-System/app/apiException"
	"QA-System/app/services/adminService"
	"QA-System/app/services/userService"
	"QA-System/app/utils"

	"QA-System/config/config"
	"github.com/gabriel-vasile/mimetype"
	"image/jpeg"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type SubmitServeyData struct {
	ID            int                         `json:"id" binding:"required"`
	QuestionsList []userService.QuestionsList `json:"questions_list"`
}

func SubmitSurvey(c *gin.Context) {
	var data SubmitServeyData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	// 判断问卷问题和答卷问题数目是否一致
	survey, err := userService.GetSurveyByID(data.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	questions, err := userService.GetQuestionsBySurveyID(survey.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	if len(questions) != len(data.QuestionsList) {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 判断填写时间是否在问卷有效期内
	if !survey.Deadline.IsZero() && survey.Deadline.Before(time.Now()) {
		utils.JsonErrorResponse(c, apiException.TimeBeyondError)
		return
	}
	// 逐个判断问题答案
	for _, q := range data.QuestionsList {
		question, err := userService.GetQuestionByID(q.QuestionID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		// 判断必填字段是否为空
		if question.Required && q.Answer == "" {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		// 判断正则是否匹配
		if question.Reg != "" {
			match, err := regexp.MatchString(question.Reg, q.Answer)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			if !match {
				utils.JsonErrorResponse(c, apiException.RegError)
				return
			}
		}
		// 判断唯一字段是否唯一
		if question.Unique {
			unique, err := userService.CheckUnique(data.ID, q.QuestionID, question.SerialNum, q.Answer)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			if !unique {
				utils.JsonErrorResponse(c, apiException.UniqueError)
				return
			}

		}
	}
	// 提交问卷
	err = userService.SubmitSurvey(data.ID, data.QuestionsList)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

type GetSurveyData struct {
	ID int `form:"id" binding:"required"`
}

type SurveyData struct {
	ID        int                    `json:"id"`
	Time      string                 `json:"time"`
	Desc      string                 `json:"desc"`
	Img       string                 `json:"img"`
	Questions []userService.Question `json:"questions"`
}

// 用户获取问卷
func GetSurvey(c *gin.Context) {
	var data GetSurveyData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	// 获取问卷
	survey, err := adminService.GetSurveyByID(data.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 判断填写时间是否在问卷有效期内
	if !survey.Deadline.IsZero() && survey.Deadline.Before(time.Now()) {
		utils.JsonErrorResponse(c, apiException.TimeBeyondError)
		return
	}
	// 获取相应的问题
	questions, err := userService.GetQuestionsBySurveyID(survey.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 构建问卷响应
	questionsResponse := make([]map[string]interface{}, 0)
	for _, question := range questions {
		options, err := userService.GetOptionsByQuestionID(question.ID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		optionsResponse := make([]map[string]interface{}, 0)
		for _, option := range options {
			optionResponse := map[string]interface{}{
				"option_type": option.OptionType,
				"content":     option.Content,
				"serial_num":  option.SerialNum,
			}
			optionsResponse = append(optionsResponse, optionResponse)
		}
		questionMap := map[string]interface{}{
			"id":            question.ID,
			"serial_num":    question.SerialNum,
			"subject":       question.Subject,
			"describe":      question.Description,
			"required":      question.Required,
			"unique":        question.Unique,
			"img":           question.Img,
			"question_type": question.QuestionType,
			"reg":           question.Reg,
			"options":       optionsResponse,
		}
		questionsResponse = append(questionsResponse, questionMap)
	}
	response := map[string]interface{}{
		"id":        survey.ID,
		"time":      survey.Deadline.Format("2006-01-02 15:04:05"),
		"desc":      survey.Desc,
		"img":       survey.Img,
		"questions": questionsResponse,
	}

	utils.JsonSuccessResponse(c, response)
}

// 上传图片
func UploadImg(c *gin.Context) {
	// 保存图片文件
	file, err := c.FormFile("img")
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 检查文件类型是否为图像
	if !isImageFile(file) {
		utils.JsonErrorResponse(c, apiException.PictureError)
		return
	}
	// 检查文件大小是否超出限制
	if file.Size > 10<<20 { // 10MB，1MB = 1024 * 1024 bytes
		utils.JsonErrorResponse(c, apiException.PictureSizeError)
		return
	}
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "tempdir")
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}() // 在处理完之后删除临时目录及其中的文件
	// 在临时目录中创建临时文件
	tempFile := filepath.Join(tempDir, file.Filename)
	f, err := os.Create(tempFile)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}()
	// 将上传的文件保存到临时文件中
	src, err := file.Open()
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	defer func() {
		if err := src.Close(); err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}()

	_, err = io.Copy(f, src)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 判断文件的MIME类型是否为图片
	mime, err := mimetype.DetectFile(tempFile)
	if err != nil || !strings.HasPrefix(mime.String(), "image/") {
		utils.JsonErrorResponse(c, apiException.PictureError)
		return
	}
	// 保存原始图片
	filename := uuid.New().String() + ".jpg" // 修改扩展名为.jpg
	dst := "./static/" + filename
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	// 转换图像为JPG格式并压缩
	jpgFile := filepath.Join(tempDir, "compressed.jpg")
	err = convertAndCompressImage(dst, jpgFile)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	// 替换原始文件为压缩后的JPG文件
	err = os.Rename(jpgFile, dst)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	urlHost := config.Config.GetString("url.host")
	url := urlHost + "static/" + filename

	utils.JsonSuccessResponse(c, url)
}

// 仅支持常见的图像文件类型
func isImageFile(file *multipart.FileHeader) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	return allowedTypes[file.Header.Get("Content-Type")]
}

// 用于转换和压缩图像的函数
func convertAndCompressImage(srcPath, dstPath string) error {
	srcImg, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}

	// 调整图像大小（根据需要进行调整）
	resizedImg := resize.Resize(300, 0, srcImg, resize.Lanczos3)

	// 创建新的JPG文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 以JPG格式保存调整大小的图像，并设置压缩质量为90
	err = jpeg.Encode(dstFile, resizedImg, &jpeg.Options{Quality: 90})
	if err != nil {
		return err
	}

	return nil
}
