package models

import "time"

type Survey struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`    //问卷标题
	Desc     string    `json:"desc"`     //问卷描述
	Img      string    `json:"img"`      //问卷图片
	Deadline time.Time `json:"deadline"` //截止时间
	Status   int       `json:"status"`   //问卷状态  1:未发布 2:已发布
}
