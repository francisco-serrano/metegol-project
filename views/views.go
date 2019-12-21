package views

import (
	"encoding/json"
	"errors"
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

type AddUsersResponse struct {
	AddUsersRequest
}

type PlayMatchRequest struct {
	MatchID      uint `json:"match_id"`
	ScoreLocal   int  `json:"score_local"`
	ScoreVisitor int  `json:"score_visitor"`
}
