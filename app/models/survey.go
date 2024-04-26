package models

import "time"

type Survey struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`      //问卷标题
	Deadline   time.Time `json:"deadline"`   //截止时间
	Stutus     int       `json:"status"`     //问卷状态  1:未发布 2:已发布
}
