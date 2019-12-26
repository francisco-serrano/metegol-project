package services

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/metegol-project/models"
	"github.com/metegol-project/viewmodels"
	"github.com/sirupsen/logrus"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
)

type Service struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewService(db *gorm.DB, logger *logrus.Logger) (*Service, error) {
	logger.Info("attempting to fetch users")

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		logger.WithError(err).Error("error while initializing service")
		return nil, err
	}

	logger.WithField("users", users).Info("users successfully fetched")

	var aux []string
	for _, u := range users {
		aux = append(aux, u.Name)
	}

	return &Service{
		db:     db,
		logger: logger,
	}, nil
}

func (s *Service) AddUsers(req viewmodels.AddUsersRequest) (*viewmodels.AddUsersResponse, error) {
	if err := s.validTournament(req.Tournament); err != nil {
		s.logger.WithError(err).WithField("tournament", req.Tournament).Error("error while validating tournament")
		return nil, err
	}

	s.logger.WithField("tournament", req.Tournament).Info("tournament validated")

	if err := s.validUsers(req); err != nil {
		s.logger.WithError(err).WithField("users", req.Users).Error("error while validating users")
		return nil, err
	}

	s.logger.WithField("users", req.Users).Info("users validated")

	teams := models.GenerateTeams(req.Users)
	matches := models.GenerateMatches(teams, 2)

	s.logger.WithField("matches", matches).Info("matches generated")

	var users []interface{}
	for _, u := range req.Users {
		user := models.User{
			Name:       u,
			Tournament: req.Tournament,
		}

		users = append(users, user)
	}

	s.logger.WithField("users", users).Info("attempting to insert new users into db")

	if err := gormbulk.BulkInsert(s.db, users, 2000); err != nil {
		s.logger.WithError(err).WithField("users", users).Error("error while inserting users")
		return nil, err
	}

	s.logger.Info("users successfully inserted into db")

	var matchesInsert []interface{}
	for _, m := range matches.Data {
		match := models.Match{
			LocalA:     m.Local.UserA,
			LocalB:     m.Local.UserB,
			VisitorA:   m.Visitor.UserA,
			VisitorB:   m.Visitor.UserB,
			Tournament: req.Tournament,
		}

		matchesInsert = append(matchesInsert, match)
	}

	s.logger.WithField("matches", matchesInsert).Info("attempting to insert new matches into db")

	if err := gormbulk.BulkInsert(s.db, matchesInsert, 2000); err != nil {
		s.logger.WithError(err).WithField("matches", matchesInsert).Error("error while inserting users")
		return nil, err
	}

	s.logger.Info("matches successfully inserted into db")

	return &viewmodels.AddUsersResponse{
		Message: "users added",
		Users:   req.Users,
	}, nil
}

func (s *Service) GetUsers() (*viewmodels.GetUsersResponse, error) {
	s.logger.Info("attempting to obtain users in db")

	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		s.logger.WithError(err).Error("error while obtaining users")
		return nil, err
	}

	s.logger.WithField("users", users).Info("users successfully found")

	return viewmodels.NewGetUsersResponse(users), nil
}

func (s *Service) GetMatches(tournament string) (*viewmodels.GetMatchesResponse, error) {
	var matches []models.Match

	s.logger.WithField("tournament", tournament).Info("attempting to obtain matches")

	if err := s.db.Where("tournament = ?", tournament).Find(&matches).Error; err != nil {
		s.logger.WithError(err).WithField("tournament", tournament).Error("error while obtaining matches")
		return nil, err
	}

	s.logger.Info("matches obtained successfully")

	return viewmodels.NewGetMatchesResponse(matches), nil
}

func (s *Service) PlayMatch(req viewmodels.PlayMatchRequest) (*viewmodels.PlayMatchResponse, error) {
	s.logger.WithField("matchID", req.MatchID).Info("attempting to find match in db")

	var match models.Match
	if err := s.db.First(&match, req.MatchID).Error; err != nil {
		s.logger.WithError(err).WithField("match", req.MatchID).Error("error while attempting to play match")
		return nil, err
	}

	s.logger.WithField("match", match).Info("match obtained")

	match.ScoreLocal = req.ScoreLocal
	match.ScoreVisitor = req.ScoreVisitor

	s.logger.WithField("match", match).Info("attempting to save updated match in db")

	if err := s.db.Save(&match).Error; err != nil {
		s.logger.WithError(err).WithField("match", match).Error("error while updating match")
		return nil, err
	}

	s.logger.Info("match successfully updated")

	return viewmodels.NewPlayMatchResponse(match), nil
}

func (s *Service) GetScores(tournament string) (*viewmodels.GetScoresResponse, error) {
	s.logger.WithField("tournament", tournament).Info("attempting to obtain users in db")

	users, err := s.getUsers(tournament)
	if err != nil {
		s.logger.WithError(err).WithField("tournament", tournament).Error("error while obtaining users")
		return nil, err
	}

	s.logger.WithField("users", users).Info("users successfully obtained from db")

	scores := make(map[string]int)
	for _, u := range *users {
		scores[u] = 0
	}

	s.logger.Info("attempting to find matched in db")

	var matches []models.Match
	if err := s.db.Where("tournament = ?", tournament).Find(&matches).Error; err != nil {
		s.logger.WithError(err).WithField("tournament", tournament).Error("error while obtaining matches")
		return nil, err
	}

	s.logger.Info("matches successfully found")

	for _, m := range matches {
		scores[m.LocalA] += m.ScoreLocal
		scores[m.LocalB] += m.ScoreLocal
		scores[m.VisitorA] += m.ScoreVisitor
		scores[m.VisitorB] += m.ScoreVisitor
	}

	return &viewmodels.GetScoresResponse{Scores: scores}, nil
}

func (s *Service) WipeData(tournament string) error {
	s.logger.WithField("tournament", tournament).Info("initializing transaction for data deletion")

	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			s.logger.WithField("recover", r).Info("transaction rollback")
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		s.logger.WithError(err).Error("error while beginning transaction")
		return err
	}

	s.logger.Info("attempting to delete matches")

	if err := s.db.Where("tournament = ?", tournament).Delete(&models.Match{}).Error; err != nil {
		s.logger.WithError(err).WithField("tournament", tournament).Error("error while deleting matches")
		tx.Rollback()
		return err
	}

	s.logger.Info("matches successfully deleted")
	s.logger.Info("attempting to delete tournaments")

	if err := s.db.Where("tournament = ?", tournament).Delete(&models.User{}).Error; err != nil {
		s.logger.WithError(err).WithField("tournament", tournament).Error("error while deleting users")
		tx.Rollback()
		return err
	}

	s.logger.Info("tournaments successfully deleted")

	if err := tx.Commit().Error; err != nil {
		s.logger.WithError(err).Error("error while committing transaction")
		return err
	}

	s.logger.Info("transaction successfully committed")

	return nil
}

func (s *Service) validTournament(tournament string) error {
	var user models.User

	s.logger.WithField("tournament", tournament).Info("attempting to validate tournament")

	var err error
	if err = s.db.Where("tournament = ?", tournament).First(&user).Error; err == nil {
		s.logger.WithError(err).WithField("tournament", tournament).Error("error while fetching user")
		return errors.New("cannot add users to an existing tournament")
	}

	if err == gorm.ErrRecordNotFound {
		s.logger.Info("tournament successfully validated")
		return nil
	}

	s.logger.WithError(err).WithField("tournament", tournament).Error("error while validating tournament")
	return err
}

func (s *Service) validUsers(req viewmodels.AddUsersRequest) error {
	s.logger.WithField("tournament", req.Tournament).Info("attempting to obtain users")

	existingUsers, err := s.getUsers(req.Tournament)
	if err != nil {
		s.logger.WithError(err).WithField("tournament", req.Tournament).Error("error while validating users")
		return err
	}

	for _, u := range req.Users {
		for _, u2 := range *existingUsers {
			if u == u2 {
				err := errors.New("users not valid")

				s.logger.WithError(err).
					WithField("user1", u).
					WithField("user2", u2).
					Error("error while validating user")

				return err
			}
		}
	}

	s.logger.Info("users successfully validated")

	return nil
}

func (s *Service) getUsers(tournament string) (*[]string, error) {
	s.logger.WithField("tournament", tournament).Info("attempting to fetch users from db")

	var users []models.User
	if err := s.db.Where("tournament = ?", tournament).Find(&users).Error; err != nil {
		s.logger.WithError(err).WithField("tournament", tournament).Error("error while obtaining users")
		return nil, err
	}

	s.logger.Info("users successfully obtained")

	var userNames []string
	for _, u := range users {
		userNames = append(userNames, u.Name)
	}

	return &userNames, nil
}
