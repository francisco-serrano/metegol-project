package models

import "github.com/jinzhu/gorm"

type Match struct {
	gorm.Model
	LocalA       string
	LocalB       string
	VisitorA     string
	VisitorB     string
	ScoreLocal   int
	ScoreVisitor int
	Tournament   string
}

type User struct {
	gorm.Model
	Name       string
	Tournament string
}
