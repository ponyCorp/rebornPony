package sensetivetrigger

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type group struct {
	chatID int64
	words  []string
	regexp *regexp.Regexp
	regStr string
}
type SensetiveTrigger struct {
	groups map[int64]*group
	m      sync.RWMutex
}

func New() *SensetiveTrigger {
	return &SensetiveTrigger{
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
			start += "\\b"
		} else {
			w = strings.TrimPrefix(w, "*")
		}

		if !strings.HasSuffix(w, "*") {
			end += "\\b"
		} else {
			w = strings.TrimSuffix(w, "*")
		}
		proccedWords = append(proccedWords, start+w+end)
	}
	return strings.Join(proccedWords, "|")
}
func (s *SensetiveTrigger) recompile(chatID int64) error {
	val, ok := s.groups[chatID]
	if !ok {
		return fmt.Errorf("group not found")
	}
	//fmc.Printfln("words: %+v", val.words)
	regxStr := join(val.words)
	//	fmc.Printfln("regexpStr: %v", []byte(regxStr))
	exp, err := regexp.Compile(regxStr)
	if err != nil {
		return err
	}
	val.regexp = exp
	val.regStr = regxStr
	return nil
}

func (s *SensetiveTrigger) AddWords(chatID int64, words ...string) error {
	s.m.Lock()
	_, ok := s.groups[chatID]
	if !ok {
		s.groups[chatID] = &group{
			chatID: chatID,
		}
	}

	s.groups[chatID].words = append(s.groups[chatID].words, words...)
	//fmt.Printf("words in group: %+v\n", s.groups[chatID].words)
	s.m.Unlock()
	err := s.recompile(chatID)
	if err != nil {
		return err
	}

	//	fmc.Printfln("regexp: %s", s.groups[chatID].regexp.String())
	return err
}
func index(s []string, e string) int {
	for i, v := range s {
		if v == e {
			return i
		}
	}
	return -1
}

// DeleteWords
func (s *SensetiveTrigger) DeleteWords(chatID int64, words ...string) error {
	defer s.m.Unlock()
	s.m.Lock()
	sensWords, ok := s.groups[chatID]
	if !ok {
		return fmt.Errorf("group not found")
	}
	var newSensWords []string

	for _, word := range sensWords.words {
		if index(words, word) < 0 {
			newSensWords = append(newSensWords, word)
		}
	}
	s.groups[chatID].words = newSensWords
	return s.recompile(chatID)
}
func (s *SensetiveTrigger) ChatIsSensetive(chatID int64) bool {
	_, ok := s.groups[chatID]
	return ok
}
func regexpMatch(msg string, regx string) (bool, error) {
	reg, err := regexp.Compile(regx)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return reg.MatchString(msg), nil
}
func (s *SensetiveTrigger) MessageContainSensitiveWords(chatID int64, message string) (bool, error) {
	_, ok := s.groups[chatID]
	if !ok {
		return false, fmt.Errorf("group not found")
	}
	rRes := s.groups[chatID].regexp.MatchString(message)
	//fmt.Printf("SensetiveTrigger: [%v], message: %s, regexp: %s\n", rRes, message, s.groups[chatID].regexp.String())
	return rRes, nil
}

// func GetRegexp(chatID int64) (*regexp.Regexp, error)
func (s *SensetiveTrigger) GetRegexp(chatID int64) (*regexp.Regexp, error) {
	_, ok := s.groups[chatID]
	if !ok {
		return nil, fmt.Errorf("group not found")
	}
	return s.groups[chatID].regexp, nil
}

// func GetRegexpStr(chatID int64) (string, error)
func (s *SensetiveTrigger) GetRegexpStr(chatID int64) (string, error) {
	_, ok := s.groups[chatID]
	if !ok {
		return "", fmt.Errorf("group not found")
	}
	return s.groups[chatID].regStr, nil
}

// func GetWords(chatID int64) ([]string, error)
func (s *SensetiveTrigger) GetWords(chatID int64) ([]string, error) {
	_, ok := s.groups[chatID]
	if !ok {
		return nil, fmt.Errorf("group not found")
	}
	return s.groups[chatID].words, nil
}
