package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"time"
)

type config struct {
	HillType         string
	HillSize         int
	MaxWarriorLength int
	Coresize         int
	Cycles           int
	MaxProcess       int
	MinSep           int
	FixedPos         int
	Rounds           int
}

var hills map[string]hill

func main() {
	b, _ := ioutil.ReadFile("config.json")
	var configs map[string]config
	json.Unmarshal(b, &configs)

	if loadHills() != nil {
		hills = make(map[string]hill)

		for key, val := range configs {
			hills[key] = hill{val, 0, make([]hillentry, val.HillSize+1), make([]hillentry, 5), make(chan warrior, 25)}
			h := hills[key]
			go h.Worker(key)
		}
	} else {
		for key, val := range hills {
			val.Queue = make(chan warrior, 25)
			go val.Worker(key)
			hills[key] = val
		}
	}

	go saveWorker()

	initRest()
}

func loadHills() error {
	var buffer bytes.Buffer
	b, err := ioutil.ReadFile("hill.data")
	if err != nil {
		return err
	}
	buffer.Write(b)
	dec := gob.NewDecoder(&buffer)
	dec.Decode(&hills)
	return nil
}

func saveHills() {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(&hills)
	ioutil.WriteFile("hill.data", buffer.Bytes(), 0777)
}

func saveWorker() {
	for {
		time.Sleep(time.Second * 30)
		saveHills()
	}
}
