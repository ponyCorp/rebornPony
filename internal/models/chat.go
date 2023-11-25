package models

type Chat struct {
	ID             string
	ChildrenChats  []string
	ParentChat     string
	ChatID         string
	DisableRead    bool
	WelcomeMessage string
	RulesMessage   string
}
