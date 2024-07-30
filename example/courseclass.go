package model

import "fmt"

// CourseClassType 课时分类
type CourseClassType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var COURSE_CLASS_TYPE = []CourseClassType{
	{ID: 1, Name: "必看系列"},
	{ID: 2, Name: "小学系列"},
	{ID: 3, Name: "初中系列"},
	{ID: 4, Name: "高中系列"},
	{ID: 5, Name: "推荐系列"},
}

// CourseClass 课时列表
type CourseClass struct {
	Id          int64  `json:"id" gorm:"column:id"`                       // 主键id
	CourseId    int64  `json:"course_id" gorm:"column:course_id"`         // 课时对应的课程id
	ClassTypeId int64  `json:"class_type_id" gorm:"column:class_type_id"` // 课时分类id
	Name        string `json:"name" gorm:"column:name"`                   // 课时名称
	Video       string `json:"video" gorm:"column:video"`                 // 课时视频
	WxState     int8   `json:"wx_state" gorm:"column:wx_state"`           // 微信上下架_1_上架_2_下架
	KsState     int8   `json:"ks_state" gorm:"column:ks_state"`           // 快手上下架_1_上架_2_下架
	DyState     int8   `json:"dy_state" gorm:"column:dy_state"`           // 抖音上下架_1_上架_2_下架
}

func (s *CourseClass) TableName() string {
	fmt.Println("ceshi")
	return "sys_course_class"
}
func (s CourseClass) TableName1() string {
	return "sys_course_class"
}
