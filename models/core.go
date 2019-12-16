package models

import "fmt"

func GenerateTeams(users []string) Teams {
	teams := NewTeams()
	for _, userA := range users {
		for _, userB := range users {
			if userA != userB && !teams.Contains(userA, userB) {
				teams.Add(userA, userB)
			}
		}
	}
	return teams
}

func GenerateMatches(teams Teams, repetitions int) Matches {
	matches := NewMatches()

	for _, local := range teams.Data {
		for _, visitor := range teams.Data {
			if !teams.SameTeam(local, visitor) && !matches.Contains(local, visitor) {
				matches.Add(local, visitor)
			}
		}
	}
	matches.Expand(repetitions)

	return matches
}

type CoreMatch struct {
	Local   Team
	Visitor Team
}

type Matches struct {
	Data []CoreMatch
}

func NewMatches() Matches {
	var aux []CoreMatch
	return Matches{Data: aux}
}

func (m *Matches) Contains(local, visitor Team) bool {
	for _, v := range m.Data {
		if (v.Local == local && v.Visitor == visitor) || (v.Local == visitor && v.Visitor == local) {
			return true
		}
	}
	return false
}

func (m *Matches) Add(local, visitor Team) {
	m.Data = append(m.Data, CoreMatch{
		Local:   local,
		Visitor: visitor,
	})
}

func (m *Matches) Length() int {
	return len(m.Data)
}

func (m *Matches) Expand(n int) {

	aux := make([]CoreMatch, len(m.Data))

	copy(aux, m.Data)

	for i := 0; i < n; i++ {
		m.Data = append(m.Data, aux...)
	}
}

type Team struct {
	UserA string
	UserB string
}

type Teams struct {
	Data []Team
}

func NewTeams() Teams {
	var aux []Team
	return Teams{Data: aux}
}

func (m *Teams) Contains(userA, userB string) bool {
	for _, v := range m.Data {
		if (v.UserA == userA && v.UserB == userB) || (v.UserA == userB && v.UserB == userA) {
			return true
		}
	}
	return false
}

func (m *Teams) Add(userA, userB string) {
	m.Data = append(m.Data, Team{UserA: userA, UserB: userB})
}

func (m *Teams) Length() int {
	return len(m.Data)
}

func (m *Teams) SameTeam(local, visitor Team) bool {
	encountered := map[string]bool{}
	encountered[local.UserA] = true
	encountered[local.UserB] = true
	encountered[visitor.UserA] = true
	encountered[visitor.UserB] = true

	amountTeams := 0

	for _ = range encountered {
		amountTeams++
	}

	if amountTeams != 4 {
		return true
	}

	return (local.UserA == visitor.UserA && local.UserB == visitor.UserB) ||
		(local.UserA == visitor.UserB && local.UserB == visitor.UserA)
}

func Run() {
	users := []string{"A", "B", "C", "D"}
	teams := GenerateTeams(users)
	repetitions := 2

	fmt.Println("users", users)
	fmt.Println("teams", teams)

	matches := GenerateMatches(teams, repetitions)

	for _, m := range matches.Data {
		fmt.Println(m)
	}
}
