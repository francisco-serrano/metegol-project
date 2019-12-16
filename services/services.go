package services

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/metegol-project/models"
	"github.com/metegol-project/views"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
)

type Service struct {
	Db      *gorm.DB
	users   []string
	teams   models.Teams
	matches models.Matches
}

func (s *Service) AddUsers(req views.AddUsersRequest) gin.H {
	if err := s.validUsers(req); err != nil {
		panic(err)
	}

	s.users = append(s.users, req.Users...)
	s.teams = models.GenerateTeams(s.users)
	s.matches = models.GenerateMatches(s.teams, 2)

	var matches []interface{}
	for _, m := range s.matches.Data {
		match := models.Match{
			LocalA:   m.Local.UserA,
			LocalB:   m.Local.UserB,
			VisitorA: m.Visitor.UserA,
			VisitorB: m.Visitor.UserB,
		}

		matches = append(matches, match)
	}

	if err := gormbulk.BulkInsert(s.Db, matches, 2000); err != nil {
		panic(err)
	}

	return gin.H{
		"message": "users added",
		"users":   s.users,
	}
}

func (s *Service) GetUsers() gin.H {
	return gin.H{
		"users": s.users,
	}
}

func (s *Service) GetMatches() gin.H {
	var matches []models.Match
	if err := s.Db.Find(&matches).Error; err != nil {
		panic(err)
	}

	return gin.H{
		"matches": matches,
	}
}

func (s *Service) PlayMatch(req views.PlayMatchRequest) gin.H {
	var match models.Match

	if err := s.Db.First(&match, req.MatchID).Error; err != nil {
		panic(err)
	}

	match.ScoreLocal = req.ScoreLocal
	match.ScoreVisitor = req.ScoreVisitor
	match.Played = true

	if err := s.Db.Save(&match).Error; err != nil {
		panic(err)
	}

	return gin.H{
		"new-match": match,
	}
}

func (s *Service) GetScores() gin.H {
	scores := make(map[string]int)
	for _, u := range s.users {
		scores[u] = 0
	}

	var matches []models.Match

	if err := s.Db.Find(&matches).Error; err != nil {
		panic(err)
	}

	for _, m := range matches {
		scores[m.LocalA] += m.ScoreLocal
		scores[m.LocalB] += m.ScoreLocal
		scores[m.VisitorA] += m.ScoreVisitor
		scores[m.VisitorB] += m.ScoreVisitor
	}

	fmt.Println(scores)

	return gin.H{
		"scores": scores,
	}
}

func (s *Service) validUsers(req views.AddUsersRequest) error {
	for _, u := range req.Users {
		for _, u2 := range s.users {
			if u == u2 {
				return errors.New("users not valid")
			}
		}
	}

	return nil
}
