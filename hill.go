package main

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/BigJk/goexmars"
)

type hill struct {
	Config config
	Age    int
	Entrys []hillentry
	Killed []hillentry
	Queue  chan warrior `json:"-"`
}

func (h hill) Len() int {
	return len(h.Entrys)
}

func (h hill) Less(i, j int) bool {
	return h.Entrys[i].Score > h.Entrys[j].Score
}

func (h hill) Swap(i, j int) {
	h.Entrys[i], h.Entrys[j] = h.Entrys[j], h.Entrys[i]
}

func (h *hill) JSON() []byte {
	b, _ := json.Marshal(h)
	return b
}

func (h hill) Index() int {
	for i := 0; i < len(h.Entrys); i++ {
		if h.Entrys[i].Warrior == nil {
			return i
		}
	}
	return len(h.Entrys) - 1
}

func (h *hill) Fight() {
	for i := 0; i < len(h.Entrys); i++ {
		if h.Entrys[i].Warrior == nil {
			continue
		}
		h.Entrys[i].Reset()
	}

	for i := 0; i < len(h.Entrys); i++ {
		if h.Entrys[i].Warrior == nil {
			continue
		}
		for j := 0; j < len(h.Entrys); j++ {
			if i == j || h.Entrys[j].Warrior == nil {
				continue
			}

			w, l, t := goexmars.Fight2Warriors(h.Entrys[i].Warrior.Code, h.Entrys[j].Warrior.Code, h.Config.Coresize, h.Config.Cycles, h.Config.MaxProcess, h.Config.Rounds, h.Config.MaxWarriorLength, h.Config.MinSep, h.Config.FixedPos)

			h.Entrys[i].Wins += w
			h.Entrys[i].Loses += l
			h.Entrys[i].Ties += t

			h.Entrys[j].Wins += l
			h.Entrys[j].Loses += w
			h.Entrys[j].Ties += l
		}
	}

	for i := 0; i < len(h.Entrys); i++ {
		h.Entrys[i].CalculateScore()
		h.Entrys[i].Age++
	}

	h.Age++

	sort.Sort(h)
}

func (h *hill) Worker(name string) {
	for {
		next := <-h.Queue
		log.Println(h.Config.HillType, ": Got", next.Name, "From Queue")
		h.Append(next)
		h.Fight()
		log.Println(h.Config.HillType, ": Finished Fighting...")
		hills[name] = *h
	}
}

func (h *hill) Append(w warrior) {
	he := createHillEntry(&w)
	hIndex := h.Index()
	if h.Entrys[hIndex].Warrior != nil {
		h.PrependKilled(h.Entrys[hIndex])
	}
	h.Entrys[hIndex] = he
}

func (h *hill) PrependKilled(he hillentry) {
	h.Killed = append([]hillentry{he}, h.Killed...)
	h.Killed = h.Killed[:len(h.Killed)-1]
}

func (h *hill) Valid(w warrior) bool {
	return true //w.Type == h.Config.HillType //&& w.Lines <= h.Config.MaxWarriorLength
}

func (h *hill) HasWarrior(w warrior) bool {
	for i := 0; i < len(h.Entrys); i++ {
		if h.Entrys[i].Warrior == nil {
			continue
		}
		if h.Entrys[i].Warrior.Name == w.Name || h.Entrys[i].Warrior.Checksum == w.Checksum {
			return true
		}
	}
	return false
}
