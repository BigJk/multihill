package main

import (
	"crypto/sha1"
	"fmt"
	"regexp"
	"strings"
)

type warrior struct {
	Type     string `json:"-"`
	Name     string
	Author   string
	Strategy string
	Code     string `json:"-"`
	Lines    int    `json:"-"`
	Checksum string `json:"-"`
}

func (w *warrior) Print() {
	fmt.Println(w.Name)
	fmt.Println("Type:", w.Type)
	fmt.Println("Author:", w.Author)
	fmt.Println("Strategy:")
	fmt.Println(w.Strategy)
}

func (w *warrior) Valid() bool {
	return w.Type != "" && w.Name != "" && w.Author != "" && len(w.UniqueCode()) > 0
}

func (w *warrior) UniqueCode() string {
	return strings.ToLower(strings.Replace(strings.Replace(strings.Replace(w.Code, "\n", "", -1), " ", "", -1), "   ", "", -1))
}

var commentRegex = regexp.MustCompile(";.*")

func parseWarrior(s string) warrior {
	w := warrior{}
	code := strings.Replace(s, "\r", "", -1)
	lines := strings.Split(code, "\n")
	w.Code = ";assert 1\n"
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(strings.ToLower(lines[i]), ";redcode-") {
			w.Type = lines[i][9:]
		} else if strings.HasPrefix(strings.ToLower(lines[i]), ";name ") {
			w.Name = lines[i][6:]
		} else if strings.HasPrefix(strings.ToLower(lines[i]), ";author ") {
			w.Author = lines[i][8:]
		} else if strings.HasPrefix(strings.ToLower(lines[i]), ";strategy ") {
			if w.Strategy == "" {
				w.Strategy = lines[i][10:]
			} else {
				w.Strategy += "\n" + lines[i][10:]
			}
		} else if lines[i] != "" {
			w.Code += lines[i] + "\n"
			w.Lines++
		}
	}
	w.Checksum = fmt.Sprintf("%x", sha1.Sum([]byte(w.UniqueCode())))
	return w
}
