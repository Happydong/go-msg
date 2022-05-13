package model

import "time"

type SmsTemplate struct {
	Id              string     `gorm:"primary_key;column:id" json:"id"`
	CreateTime      *time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime      *time.Time `gorm:"column:update_time" json:"updateTime"`
	Code            string
	Name            string
	SignName        string
	TemplateCode    string
	AccessKeyId     string
	AccessKeySecret string
	CodeExpired     int
}

func (SmsTemplate) TableName() string {
	return "sms_template"
}
