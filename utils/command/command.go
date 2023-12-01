package command

import (
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandParser struct {
	botName      string
	reSpecSymbol *regexp.Regexp
}
type command struct {
	Cmd string
	Arg string
}

func fidMinMax(x, y int) (int, int) {
	if x < y {
		return x, y
	}
	return y, x
}
func positiveOrNull(x int) int {
	if x < 0 {
		return 0
	}
	return x
}
func findMinMaxPositiveOrNull(num1, num2 int) int {
	n, m := fidMinMax(positiveOrNull(num1), positiveOrNull(num2))
	if n != 0 {
		return n
	}
	return m
}
func NewCommandParser(botName string) *CommandParser {

	return &CommandParser{
		botName:      botName,
		reSpecSymbol: regexp.MustCompile(`^([!/])[A-zА-я0-9]+$`),
	}
}
func (h *CommandParser) IsCommand(upd *tgbotapi.Update) bool {
	if upd == nil || upd.Message == nil {
		return false
	}
	if upd.Message.IsCommand() {
		return true
	}
	isCmd, _ := h.ParseCommand(upd.Message.Text)
	return isCmd
}
func (h *CommandParser) IsCommandWithArgs(upd *tgbotapi.Update) (bool, command) {
	if upd == nil || upd.Message == nil {
		return false, command{}
	}

	isCmd, cmd := h.ParseCommand(upd.Message.Text)
	return isCmd, cmd
}

// ParseCommand(text string)(isCommand bool, command string, args string)
func (c *CommandParser) ParseCommand(msg string) (bool, command) {
	cmdWithArg := command{}
	if len(msg) < 2 || (msg[0] != '/' && msg[0] != '!') {
		return false, command{}
	}

	index := findMinMaxPositiveOrNull(strings.Index(msg, " "), strings.Index(msg, "\n"))
	if index <= 0 {
		cmdWithArg.Cmd = msg
	} else {
		cmdWithArg.Cmd = msg[:index]
		cmdWithArg.Arg = msg[index+1:]
	}
	if indexIserName := strings.Index(cmdWithArg.Cmd, "@"); indexIserName >= 1 {
		user := cmdWithArg.Cmd[indexIserName+1:]
		if user != c.botName {
			return false, command{}
		}
		cmdWithArg.Cmd = cmdWithArg.Cmd[:indexIserName]
	}
	if len(cmdWithArg.Cmd) <= 1 || !c.reSpecSymbol.MatchString(cmdWithArg.Cmd) {

		return false, command{}
	}

	cmdWithArg.Cmd = cmdWithArg.Cmd[1:]
	return true, cmdWithArg

}
