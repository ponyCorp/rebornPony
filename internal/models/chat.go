package models

type Chat struct {
	ID string
	//chat info
	ChatName       string
	ChatID         int64
	WelcomeMessage string
	RulesMessage   string
	KnownUsers     []string
	// service
	DeleteServiceMessage        bool
	DeleteServiceMessageTimeout int64
	IgnoreEvents                bool
	StormMode                   bool
	//forgiveness
	ForgivenessAlert   bool
	ForgivinessTimeout int64
	//tree
	ParentChat    string
	ChildrenChats []string
	EventsChannel int64
}
