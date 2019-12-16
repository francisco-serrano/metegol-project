package views

type AddUsersRequest struct {
	Users      []string `json:"users"`
	Tournament string   `json:"tournament"`
}

type PlayMatchRequest struct {
	MatchID      uint `json:"match_id"`
	ScoreLocal   int  `json:"score_local"`
	ScoreVisitor int  `json:"score_visitor"`
}
