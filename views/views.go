package views

type AddUsersRequest struct {
	Users []string `json:"users"`
}

type PlayMatchRequest struct {
	MatchID      uint `json:"match_id"`
	ScoreLocal   int  `json:"score_local"`
	ScoreVisitor int  `json:"score_visitor"`
}
