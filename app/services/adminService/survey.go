package adminService

import (
	"QA-System/app/models"
	"QA-System/app/services/mongodbService"
	"QA-System/config/config"
	"QA-System/config/database"
	"os"
	"sort"
	"strings"
	"time"
)

type Option struct {
	SerialNum  int    `json:"serial_num"`  //选项序号
	Content    string `json:"content"`     //选项内容
	OptionType int    `json:"option_type"` //选项类型 1文字2图片
}

type Question struct {
	ID           int      `json:"id"`
	SerialNum    int      `json:"serial_num"`    //题目序号
	Subject      string   `json:"subject"`       //问题
	Description  string   `json:"description"`   //问题描述
	Img          string   `json:"img"`           //图片
	Required     bool     `json:"required"`      //是否必填
	Unique       bool     `json:"unique"`        //是否唯一
	OtherOption  bool     `json:"other_option"`  //是否有其他选项
	QuestionType int      `json:"question_type"` //问题类型 1单选2多选3填空4简答5图片
	Reg          string   `json:"reg"`           //正则表达式
	Options      []Option `json:"options"`       //选项
}

func GetSurveyByID(id int) (models.Survey, error) {
	var survey models.Survey
	err := database.DB.Where("id = ?", id).First(&survey).Error
	return survey, err
}

func CreateSurvey(id int, title string, desc string, img string, questions []Question, status int, time time.Time) error {
	var survey models.Survey
	survey.UserID = id
	survey.Title = title
	survey.Desc = desc
	survey.Img = img
	survey.Status = status
	survey.Deadline = time
	err := database.DB.Create(&survey).Error
	if err != nil {
		return err
	}
	_,err = createQuestionsAndOptions(questions, survey.ID)
	return err
}

func UpdateSurveyStatus(id int, status int) error {
	var survey models.Survey
	err := database.DB.Model(&survey).Where("id = ?", id).Update("status", status).Error
	return err
}

func UpdateSurvey(id int, title string, desc string, img string, questions []Question, time time.Time) error {
	//遍历原有问题，删除对应选项
	var survey models.Survey
	var oldQuestions []models.Question
	var old_imgs []string
	new_imgs := make([]string, 0)
	//获取原有图片
	err := database.DB.Where("survey_id = ?", id).Find(&oldQuestions).Error
	if err != nil {
		return err
	}
	old_imgs, err = getOldImgs(id, oldQuestions)
	if err != nil {
		return err
	}
	//删除原有问题和选项
	for _, oldQuestion := range oldQuestions {
		err = database.DB.Where("question_id = ?", oldQuestion.ID).Delete(&models.Option{}).Error
		if err != nil {
			return err
		}
	}
	err = database.DB.Where("survey_id = ?", id).Delete(&models.Question{}).Error
	if err != nil {
		return err
	}
	//修改问卷信息
	err = database.DB.Model(&survey).Where("id = ?", id).Updates(map[string]interface{}{"title": title, "desc": desc, "img": img, "deadline": time}).Error
	if err != nil {
		return err
	}
	new_imgs = append(new_imgs, img)
	//重新添加问题和选项
	imgs,err := createQuestionsAndOptions(questions, id)
	if err != nil {
		return err
	}
	new_imgs = append(new_imgs, imgs...)
	urlHost := config.Config.GetString("url.host")
	//删除无用图片
	for _, old_img := range old_imgs {
		if !contains(new_imgs, old_img) {
			_ = os.Remove("./static/" + strings.TrimPrefix(old_img, urlHost+"/static/"))
		}
	}
	return nil
}

func UserInManage(uid int, sid int) bool {
	var survey models.Manage
	err := database.DB.Where("user_id = ? and survey_id = ?", uid, sid).First(&survey).Error
	return err == nil
}

func DeleteSurvey(id int) error {
	var survey models.Survey
	var questions []models.Question
	err := database.DB.Where("survey_id = ?", id).Find(&questions).Error
	if err != nil {
		return err
	}
	var answerSheets []mongodbService.AnswerSheet
	answerSheets,_, err = mongodbService.GetAnswerSheetBySurveyID(id,0,0)
	if err != nil {
		return err
	}
	//删除图片
	imgs, err := getDelImgs(id, questions, answerSheets)
	if err != nil {
		return err
	}
	urlHost := config.Config.GetString("url.host")
	for _, img := range imgs {
		_ = os.Remove("./static/" + strings.TrimPrefix(img, urlHost+"/static/"))
	}
	//删除答卷
	err = mongodbService.DeleteAnswerSheetBySurveyID(id)
	if err != nil {
		return err
	}
	//删除问题、选项、问卷、管理
	for _, question := range questions {
		err = database.DB.Where("question_id = ?", question.ID).Delete(&models.Option{}).Error
		if err != nil {
			return err
		}
	}
	err = database.DB.Where("survey_id = ?", id).Delete(&models.Question{}).Error
	if err != nil {
		return err
	}
	err = database.DB.Where("id = ?", id).Delete(&survey).Error
	if err != nil {
		return err
	}
	err = database.DB.Where("survey_id = ?",id).Delete(&models.Manage{}).Error
	return err
}

type QuestionAnswers struct {
	Title    string `json:"title"`
	Answers []string `json:"answers"`
}

type AnswersResonse struct {
	QuestionAnswers []QuestionAnswers `json:"question_answers"`
	Time []string `json:"time"`
}



func GetSurveyAnswers(id int, num int, size int) (AnswersResonse, *int64, error) {
	var answerSheets []mongodbService.AnswerSheet
	var questions []models.Question
	data := make([]QuestionAnswers, 0)
	time := make([]string, 0)
	var total *int64
	//获取问题
	err := database.DB.Where("survey_id = ?", id).Find(&questions).Error
	if err != nil {
		return AnswersResonse{}, nil, err
	}
	//初始化data
	for _, question := range questions {
		var q QuestionAnswers
		q.Title = question.Subject
		q.Answers = make([]string, 0)
		data = append(data, q)
	}
	//获取答卷
	answerSheets,total, err = mongodbService.GetAnswerSheetBySurveyID(id,num,size)
	if err != nil {
		return AnswersResonse{}, nil, err
	}
	//填充data
	for _, answerSheet := range answerSheets {
		time = append(time, answerSheet.Time)
		for _, answer := range answerSheet.Answers {
			var question models.Question
			err = database.DB.Where("id = ?", answer.QuestionID).First(&question).Error
			if err != nil {
				return AnswersResonse{}, nil, err
			}
			for i, q := range data {
				if q.Title == question.Subject {
					data[i].Answers = append(data[i].Answers, answer.Content)
				}
			}
		}
	}
	return AnswersResonse{QuestionAnswers: data, Time: time}, total, nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getOldImgs(id int, questions []models.Question) ([]string, error) {
	var imgs []string
	var survey models.Survey
	err := database.DB.Where("id = ?", id).First(&survey).Error
	if err != nil {
		return nil, err
	}
	imgs = append(imgs, survey.Img)
	for _, question := range questions {
		imgs = append(imgs, question.Img)
		var options []models.Option
		err = database.DB.Where("question_id = ?", question.ID).Find(&options).Error
		if err != nil {
			return nil, err
		}
		for _, option := range options {
			if option.OptionType == 2 {
				imgs = append(imgs, option.Content)
			}
		}
	}
	return imgs, nil
}

func getDelImgs(id int, questions []models.Question, answerSheets []mongodbService.AnswerSheet) ([]string, error) {
	var imgs []string
	var survey models.Survey
	err := database.DB.Where("id = ?", id).First(&survey).Error
	if err != nil {
		return nil, err
	}
	imgs = append(imgs, survey.Img)
	for _, question := range questions {
		imgs = append(imgs, question.Img)
		var options []models.Option
		err = database.DB.Where("question_id = ?", question.ID).Find(&options).Error
		if err != nil {
			return nil, err
		}
		for _, option := range options {
			if option.OptionType == 2 {
				imgs = append(imgs, option.Content)
			}
		}
	}
	for _, answerSheet := range answerSheets {
		for _, answer := range answerSheet.Answers {
			var question models.Question
			err = database.DB.Where("id = ?", answer.QuestionID).First(&question).Error
			if err != nil {
				return nil, err
			}
			if question.QuestionType == 5 {
				imgs = append(imgs, answer.Content)
			}
		}

	}
	return imgs, nil
}

func createQuestionsAndOptions(questions []Question,sid int) ([]string,error) {
	var imgs []string
	for _, question := range questions {
		var q models.Question
		q.SerialNum = question.SerialNum
		q.SurveyID =  sid
		q.Subject = question.Subject
		q.Description = question.Description
		q.Img = question.Img
		q.Required = question.Required
		q.Unique = question.Unique
		q.OtherOption = question.OtherOption
		q.QuestionType = question.QuestionType
		imgs = append(imgs, question.Img)
		err := database.DB.Create(&q).Error
		if err != nil {
			return nil,err
		}
		for _, option := range question.Options {
			var o models.Option
			o.Content = option.Content
			o.QuestionID = q.ID
			o.SerialNum = option.SerialNum
			o.OptionType = option.OptionType
			if option.OptionType == 2 {
				imgs = append(imgs, option.Content)
			}
			err := database.DB.Create(&o).Error
			if err != nil {
				return nil,err
			}
		}
	}
	return imgs,nil
}

func GetAllSurveyByUserID(userId int) ([]interface{}, error) {
	var surveys []models.Survey
	err := database.DB.Model(models.Survey{}).Where("user_id = ?", userId).
		Order("CASE WHEN status = 2 THEN 0 ELSE 1 END, id DESC").Find(&surveys).Error
	response := getSurveyResponse(surveys)
	return response, err
}

func ProcessResponse(response []interface{}, pageNum, pageSize int, title string) ([]interface{}, *int64) {
	if title != "" {
		filteredResponse := make([]interface{}, 0)
		for _, item := range response {
			itemMap := item.(map[string]interface{})
			if strings.Contains(strings.ToLower(itemMap["title"].(string)), strings.ToLower(title)) {
				filteredResponse = append(filteredResponse, item)
			}
		}
		response = filteredResponse
	}
	num := int64(len(response))
	sort.Slice(response, func(i, j int) bool {
		return response[i].(map[string]interface{})["id"].(int) > response[j].(map[string]interface{})["id"].(int)
	})
	var sortedResponse []interface{}
	var status2Response, status1Response []interface{}
	for _, item := range response {
		itemMap := item.(map[string]interface{})
		if itemMap["status"].(int) == 2 {
			status2Response = append(status2Response, item)
		} else {
			status1Response = append(status1Response, item)
		}
	}
	sortedResponse = append(status2Response, status1Response...)

	startIdx := (pageNum - 1) * pageSize
	endIdx := startIdx + pageSize
	if endIdx > len(sortedResponse) {
		endIdx = len(sortedResponse)
	}
	pagedResponse := sortedResponse[startIdx:endIdx]

	return pagedResponse, &num
}

func GetAllSurvey(pageNum, pageSize int, title string) ([]interface{}, *int64) {
	var surveys []models.Survey
	var num int64
	query := database.DB.Model(models.Survey{}).
		Order("CASE WHEN status = 2 THEN 0 ELSE 1 END, id DESC")
	if title != "" {
		title := "%" + title + "%"
		query = query.Where("title LIKE ?", title)
		query.Find(&surveys)
		num = int64(len(surveys))
	} else {
		query.Find(&surveys)
		num = int64(len(surveys))
	}
	response := getSurveyResponse(surveys)

	startIdx := (pageNum - 1) * pageSize
	endIdx := startIdx + pageSize
	if endIdx > len(response) {
		endIdx = len(response)
	}
	pagedResponse := response[startIdx:endIdx]

	return pagedResponse, &num
}

func getSurveyResponse(surveys []models.Survey) []interface{} {
	response := make([]interface{}, 0)
	for _, survey := range surveys {
		surveyResponse := map[string]interface{}{
			"id":     survey.ID,
			"title":  survey.Title,
			"status": survey.Status,
			"num":    survey.Num,
		}
		response = append(response, surveyResponse)
	}
	return response
}

func GetManageredSurveyByUserID(userId int) ([]models.Manage, error) {
	var surveys []models.Manage
	err := database.DB.Model(models.Manage{}).Where("user_id = ?", userId).Order("id DESC").Find(&surveys).Error
	return surveys, err
}

func GetAllSurveyAnswers(id int) (AnswersResonse, error) {
	var data []QuestionAnswers
	var answerSheets []mongodbService.AnswerSheet
	var questions []models.Question
	var time []string
	err := database.DB.Where("survey_id = ?", id).Find(&questions).Error
	if err != nil {
		return AnswersResonse{}, err
	}
	for _, question := range questions {
		var q QuestionAnswers
		q.Title = question.Subject
		data = append(data, q)
	}
	answerSheets,_, err = mongodbService.GetAnswerSheetBySurveyID(id,0,0)
	if err != nil {
		return AnswersResonse{}, err
	}
	for _, answerSheet := range answerSheets {
		time = append(time, answerSheet.Time)
		for _, answer := range answerSheet.Answers {
			var question models.Question
			err = database.DB.Where("id = ?", answer.QuestionID).First(&question).Error
			if err != nil {
				return AnswersResonse{}, err
			}
			for i, q := range data {
				if q.Title == question.Subject {
					data[i].Answers = append(data[i].Answers, answer.Content)
				}
			}
		}
	}
	return AnswersResonse{QuestionAnswers: data, Time: time}, nil
}
