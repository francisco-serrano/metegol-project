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
	Played       bool
}
