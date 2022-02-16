package store

import (
	"context"
	"datasync/services/events/data/source"
	"log"
	"strconv"
)

type StandingsService dbservice

//Standings contains league table stangings
type Standings struct {
	LeagueID    int    `firestore:"leagueId,omitempty"`
	Rank        int    `firestore:"rank,omitempty"`
	TeamID      int    `firestore:"teamId,omitempty"`
	TeamName    string `firestore:"teamName,omitempty"`
	Logo        string `firestore:"logo,omitempty"`
	Group       string `firestore:"group,omitempty"`
	Status      string `firestore:"status,omitempty"`
	Form        string `firestore:"forme,omitempty"`
	Description string `firestore:"description,omitempty"`
	AllStat     Stat   `firestore:"all,omitempty"`
	HomeStat    Stat   `firestore:"home,omitempty"`
	AwayStat    Stat   `firestore:"away,omitempty"`
	GoalsDiff   int    `firestore:"goalsDiff,omitempty"`
	Points      int    `firestore:"points,omitempty"`
	LastUpdated string `firestore:"lastUpdate,omitempty"`
}

//Stat contains  team statistics
type Stat struct {
	MatchsPlayed int `firestore:"matchesPlayed,omitempty"`
	Win          int `firestore:"win,omitempty"`
	Draw         int `firestore:"draw,omitempty"`
	Lose         int `firestore:"lose,omitempty"`
	GoalsFor     int `firestore:"goalsFor,omitempty"`
	GoalsAgainst int `firestore:"goalsAgainst,omitempty"`
}

func (service *StandingsService) Set(ctx context.Context, leagueName string, entity *source.StandingEntity) error {
	batch := service.client.fs.Batch()
	
	for _, standings := range entity.API.Standings[0] {
		s := Standings{}
		s.LeagueID = entity.API.LeagueID
		s.TeamID = standings.TeamID
		s.Rank = standings.Rank
		s.TeamName = standings.TeamName
		s.Logo  = standings.Logo
		s.Group = standings.Group
		s.Description = standings.Description
		s.Status = standings.Status
		s.Form = standings.Form
		s.AllStat.Draw = standings.AllStat.Draw
		s.AllStat.GoalsAgainst = standings.AllStat.GoalsAgainst
		s.AllStat.GoalsFor = standings.AllStat.GoalsFor
		s.AllStat.Lose = standings.AllStat.Lose
		s.AllStat.MatchsPlayed = standings.AllStat.MatchsPlayed
		s.AllStat.Win = standings.AllStat.Win
		s.HomeStat = Stat(standings.HomeStat)
		s.AwayStat = Stat(standings.AwayStat)
		s.GoalsDiff = standings.GoalsDiff
		s.Points = standings.Points
		s.LastUpdated = standings.LastUpdated

		leagueRef := service.client.fs.Collection("football").Doc(leagueName)
		docRef := leagueRef.
			Collection("leagues").
			Doc("leagueId_" + strconv.Itoa(entity.API.LeagueID)).
			Collection("standings").
			Doc(DocWithIDAndName(s.TeamID, s.TeamName))
		log.Println("Document Ref: " + docRef.Path)
		batch.Set(docRef, s)
	}

	_, err := batch.Commit(ctx)
	
	if err != nil {
		return err
	}

	return nil
}
