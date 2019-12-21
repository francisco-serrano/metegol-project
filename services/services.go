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
	db      *gorm.DB
	teams   models.Teams
	matches models.Matches
}

func NewService(db *gorm.DB) *Service {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		panic(err)
	}

	var aux []string
	for _, u := range users {
		aux = append(aux, u.Name)
	}

	return &Service{
		db: db,
	}
}

func (s *Service) AddUsers(req views.AddUsersRequest) (*gin.H, error) {
	if err := s.validTournament(req.Tournament); err != nil {
		return nil, err
	}

	if err := s.validUsers(req); err != nil {
		return nil, err
	}

	s.teams = models.GenerateTeams(req.Users)
	s.matches = models.GenerateMatches(s.teams, 2)

	var users []interface{}
	for _, u := range req.Users {
		user := models.User{
			Name:       u,
			Tournament: req.Tournament,
		}

		users = append(users, user)
	}

	if err := gormbulk.BulkInsert(s.db, users, 2000); err != nil {
		return nil, err
	}

	var matches []interface{}
	for _, m := range s.matches.Data {
		match := models.Match{
			LocalA:     m.Local.UserA,
			LocalB:     m.Local.UserB,
			VisitorA:   m.Visitor.UserA,
			VisitorB:   m.Visitor.UserB,
			Tournament: req.Tournament,
		}

		matches = append(matches, match)
	}

	if err := gormbulk.BulkInsert(s.db, matches, 2000); err != nil {
		return nil, err
	}

	return &gin.H{
		"message": "users added",
		"users":   users,
	}, nil
}

func (s *Service) GetUsers() (*gin.H, error) {
	var users []models.User

	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return &gin.H{
		"users": users,
	}, nil
}

func (s *Service) GetMatches(tournament string) (*gin.H, error) {
	var matches []models.Match
	if err := s.db.Where("tournament = ?", tournament).Find(&matches).Error; err != nil {
		return nil, err
	}

	return &gin.H{
		"matches": matches,
	}, nil
}

func (s *Service) PlayMatch(req views.PlayMatchRequest) (*gin.H, error) {
	var match models.Match

	if err := s.db.First(&match, req.MatchID).Error; err != nil {
		return nil, err
	}

	match.ScoreLocal = req.ScoreLocal
	match.ScoreVisitor = req.ScoreVisitor

	if err := s.db.Save(&match).Error; err != nil {
		return nil, err
	}

	return &gin.H{
		"new-match": match,
	}, nil
}

func (s *Service) GetScores(tournament string) gin.H {
	scores := make(map[string]int)
	for _, u := range s.getUsers(tournament) {
		scores[u] = 0
	}

	var matches []models.Match
	if err := s.db.Where("tournament = ?", tournament).Find(&matches).Error; err != nil {
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

func (s *Service) WipeData(tournament string) error {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := s.db.Where("tournament = ?", tournament).Delete(&models.Match{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := s.db.Where("tournament = ?", tournament).Delete(&models.User{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *Service) validTournament(tournament string) error {
	var user models.User

	var err error
	if err = s.db.Where("tournament = ?", tournament).First(&user).Error; err == nil {
		return errors.New("cannot add users to an existing tournament")
	}

	if err == gorm.ErrRecordNotFound {
		return nil
	}

	return err
}

func (s *Service) validUsers(req views.AddUsersRequest) error {
	existingUsers := s.getUsers(req.Tournament)

	for _, u := range req.Users {
		for _, u2 := range existingUsers {
			if u == u2 {
				return errors.New("users not valid")
			}
		}
	}

	return nil
}

func (s *Service) getUsers(tournament string) []string {
	var users []models.User
	if err := s.db.Where("tournament = ?", tournament).Find(&users).Error; err != nil {
		panic(err)
	}

	var userNames []string
	for _, u := range users {
		userNames = append(userNames, u.Name)
	}

	return userNames
}
