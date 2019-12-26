package viewmodels

import (
	"encoding/json"
	"errors"
	"github.com/metegol-project/models"
	"strings"
)

type AddUsersRequest struct {
	Users      []string `json:"users"`
	Tournament string   `json:"tournament"`
}

func (req *AddUsersRequest) UnmarshalJSON(data []byte) error {
	type Alias AddUsersRequest

	var req2 Alias
	if err := json.Unmarshal(data, &req2); err != nil {
		return err
	}

	if strings.Contains(req2.Tournament, " ") {
		return errors.New("tournament name cannot contain spaces")
	}

	*req = AddUsersRequest(req2)

	return nil
}

type PlayMatchRequest struct {
	MatchID      uint `json:"match_id"`
	ScoreLocal   int  `json:"score_local"`
	ScoreVisitor int  `json:"score_visitor"`
}

type BaseResponse struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

type AddUsersResponse struct {
	Message string   `json:"message"`
	Users   []string `json:"users"`
}

type userResponse struct {
	Name       string `json:"name"`
	Tournament string `json:"tournament"`
}

type GetUsersResponse struct {
	Users []userResponse `json:"users"`
}

func NewGetUsersResponse(users []models.User) *GetUsersResponse {
	var usersReturn []userResponse

	for _, u := range users {
		user := userResponse{
			Name:       u.Name,
			Tournament: u.Tournament,
		}

		usersReturn = append(usersReturn, user)
	}

	return &GetUsersResponse{Users: usersReturn}
}

type matchResponse struct {
	ID           uint   `json:"id"`
	LocalA       string `json:"local_a"`
	LocalB       string `json:"local_b"`
	VisitorA     string `json:"visitor_a"`
	VisitorB     string `json:"visitor_b"`
	ScoreLocal   int    `json:"score_local"`
	ScoreVisitor int    `json:"score_visitor"`
	Tournament   string `json:"tournament"`
}

type GetMatchesResponse struct {
	Matches []matchResponse `json:"matches"`
}

func NewGetMatchesResponse(matches []models.Match) *GetMatchesResponse {
	var matchesReturn []matchResponse

	for _, m := range matches {
		match := matchResponse{
			ID:           m.ID,
			LocalA:       m.LocalA,
			LocalB:       m.LocalB,
			VisitorA:     m.VisitorA,
			VisitorB:     m.VisitorB,
			ScoreLocal:   m.ScoreLocal,
			ScoreVisitor: m.ScoreVisitor,
			Tournament:   m.Tournament,
		}

		matchesReturn = append(matchesReturn, match)
	}

	return &GetMatchesResponse{Matches: matchesReturn}
}

type PlayMatchResponse struct {
	PlayedMatch matchResponse `json:"played_match"`
}

func NewPlayMatchResponse(match models.Match) *PlayMatchResponse {
	return &PlayMatchResponse{PlayedMatch: matchResponse{
		ID:           match.ID,
		LocalA:       match.LocalA,
		LocalB:       match.LocalB,
		VisitorA:     match.VisitorA,
		VisitorB:     match.VisitorB,
		ScoreLocal:   match.ScoreLocal,
		ScoreVisitor: match.ScoreVisitor,
		Tournament:   match.Tournament,
	}}
}

type GetScoresResponse struct {
	Scores map[string]int `json:"scores"`
}
