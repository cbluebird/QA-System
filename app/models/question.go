package models

type Question struct {
	ID           int    `json:"id"`
	SurveyID     int    `json:"survey_id"`     //问卷ID
	SerialNum    int    `json:"serial_num"`    //题目序号
	Suject       string `json:"subject"`       //题目
	Descrption   string `json:"description"`   //题目描述
	Required     bool   `json:"required"`      //是否必填
	QuestionType int    `json:"question_type"` //题目类型 1单选2多选3填空4简答5图片
	Reg          string `json:"reg"`           //正则表达式
}
