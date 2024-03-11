package models

import "gorm.io/gorm"

type Command struct {
	gorm.Model
	Script   string `json:"script"`
	Result   string `json:"result"`
	Executed bool   `json:"executed"`
}
