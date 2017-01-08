package main

import "time"

type hillentry struct {
	Warrior *warrior
	Entry   time.Time
	Age     int
	Wins    int
	Ties    int
	Loses   int
	Score   float32
}

func createHillEntry(w *warrior) hillentry {
	h := hillentry{}
	h.Warrior = w
	h.Entry = time.Now()
	return h
}

func (he *hillentry) CalculateScore() {
	he.Score = ((float32(he.Wins)*3.0 + float32(he.Ties)) / float32(he.Wins+he.Loses+he.Ties+1)) * 100.0
}

func (he *hillentry) Reset() {
	he.Loses = 0
	he.Ties = 0
	he.Wins = 0
	he.Score = 0
}
