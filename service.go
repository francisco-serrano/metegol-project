package main

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type Service struct {
	users   []string
	teams   Teams
	matches Matches
}

func (s *Service) AddUsers(req AddUsersRequest) gin.H {
	if err := s.validUsers(req); err != nil {
		panic(err)
	}

	s.users = append(s.users, req.Users...)
	s.teams = GenerateTeams(s.users)
	s.matches = GenerateMatches(s.teams, 2)

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
	return gin.H{
		"matches": s.matches,
	}
}

func (s *Service) validUsers(req AddUsersRequest) error {
	for _, u := range req.Users {
		for _, u2 := range s.users {
			if u == u2 {
				return errors.New("users not valid")
			}
		}
	}

	return nil
}
