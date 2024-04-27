package models

type Option struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"` //问题ID
	SerialNum  int    `json:"serial_num"`  //选项序号
	Context    string `json:"context"`     //选项内容
	OptionType int    `json:"option_type"` //选项类型 1文字2图片
}
