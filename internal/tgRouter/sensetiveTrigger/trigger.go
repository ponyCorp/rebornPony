package sensetivetrigger

import (
	"fmt"
	"regexp"
	"strings"
)

type group struct {
	chatID int64
	words  []string
	regexp *regexp.Regexp
}
type sensetiveTrigger struct {
	groups map[int64]*group
}

func New() *sensetiveTrigger {
	return &sensetiveTrigger{
		groups: make(map[int64]*group),
	}
}

func join(words []string) string {
	if len(words) == 0 {
		return ""
	}
	var proccedWords []string

	for _, word := range words {
		start := ""
		end := ""
		w := word
		if !strings.HasPrefix(w, "*") {
			start += "\b"
		} else {
			w = strings.TrimPrefix(w, "*")
		}

		if !strings.HasSuffix(word, "*") {
			end += "\b"
		} else {
			w = strings.TrimSuffix(w, "*")
		}
		proccedWords = append(proccedWords, start+w+end)
	}
	return strings.Join(proccedWords, "|")
}
func (s *sensetiveTrigger) recompile(chatID int64) error {
	val, ok := s.groups[chatID]
	if !ok {
		return fmt.Errorf("group not found")
	}

	regxStr := join(val.words)
	exp, err := regexp.Compile(regxStr)
	if err != nil {
		return err
	}
	val.regexp = exp
	return nil
}

func (s *sensetiveTrigger) AddWords(chatID int64, words ...string) error {
	_, ok := s.groups[chatID]
	if !ok {
		s.groups[chatID] = &group{
			chatID: chatID,
			words:  words,
		}
	}

	s.groups[chatID].words = append(s.groups[chatID].words, words...)
	return s.recompile(chatID)
}

func (s *sensetiveTrigger) ChatIsSensetive(chatID int64) bool {
	_, ok := s.groups[chatID]
	return ok
}

func (s *sensetiveTrigger) MessageContainSensitiveWords(chatID int64, message string) bool {
	_, ok := s.groups[chatID]
	if !ok {
		return true
	}
	return s.groups[chatID].regexp.MatchString(message)
}
