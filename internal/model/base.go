package model

import (
	"time"

	"gorm.io/gorm"
)

// 重新定义 BaseModel，主要是增加 json 标签
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createAt"`
	UpdatedAt time.Time `json:"updateAt"`
	// 通常在 JSON 输出中忽略 ("-")
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
