package command

import (
	"reflect"
	"regexp"
	"testing"
)

func TestNewCommandParser(t *testing.T) {
	type args struct {
		botName string
	}
	tests := []struct {
		name string
		args args
		want *CommandParser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommandParser(tt.args.botName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommandParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandParser_ParseCommand(t *testing.T) {
	type fields struct {
		botName string
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  command
	}{
		{
			name:   "Command without bot name prefix",
			fields: fields{botName: "toster"},
			args:   args{msg: "/start"},
			want:   true,
			want1:  command{Cmd: "start", Arg: ""},
		},
		{
			name:   "Command with bot name prefix",
			fields: fields{botName: "toster"},
			args:   args{msg: "/help@toster"},
			want:   true,
			want1:  command{Cmd: "help", Arg: ""},
		},
		{
			name:   "Command with bot name prefix and arguments",
			fields: fields{botName: "toster"},
			args:   args{msg: "/search@toster keyword"},
			want:   true,
			want1:  command{Cmd: "search", Arg: "keyword"},
		},
		{
			name:   "Command without bot name prefix but with arguments",
			fields: fields{botName: "toster"},
			args:   args{msg: "/search keyword"},
			want:   true,
			want1:  command{Cmd: "search", Arg: "keyword"},
		},
		{
			name:   "Command with bot name prefix and invalid command",
			fields: fields{botName: "toster"},
			args:   args{msg: "/!#invalid@toster"},
			want:   false,
			want1:  command{Cmd: "", Arg: ""},
		},

		{
			name:   "Command without bot name prefix",
			fields: fields{botName: "toster"},
			args:   args{msg: "!start"},
			want:   true,
			want1:  command{Cmd: "start", Arg: ""},
		},
		{
			name:   "Command with bot name prefix",
			fields: fields{botName: "toster"},
			args:   args{msg: "!help@toster"},
			want:   true,
			want1:  command{Cmd: "help", Arg: ""},
		},
		{
			name:   "Command with bot name prefix and arguments",
			fields: fields{botName: "toster"},
			args:   args{msg: "!search@toster keyword"},
			want:   true,
			want1:  command{Cmd: "search", Arg: "keyword"},
		},
		{
			name:   "Command without bot name prefix but with arguments",
			fields: fields{botName: "toster"},
			args:   args{msg: "!search keyword"},
			want:   true,
			want1:  command{Cmd: "search", Arg: "keyword"},
		},
		{
			name:   "Command with bot name prefix and invalid command separator !",
			fields: fields{botName: "toster"},
			args:   args{msg: "!@invalid@toster"},
			want:   false,
			want1:  command{Cmd: "", Arg: ""},
		},

		{
			name:   "Command with different bot name prefix",
			fields: fields{botName: "toster"},
			args:   args{msg: "/search@otherbot"},
			want:   false,
			want1:  command{Cmd: "", Arg: ""},
		},
		{
			name:   "Command with different bot name prefix and arguments",
			fields: fields{botName: "toster"},
			args:   args{msg: "/search@otherbot keyword"},
			want:   false,
			want1:  command{Cmd: "", Arg: ""},
		},

		{
			name:   "Command with different bot name prefix  separator !",
			fields: fields{botName: "toster"},
			args:   args{msg: "!search@otherbot"},
			want:   false,
			want1:  command{Cmd: "", Arg: ""},
		},
		{
			name:   "Command with different bot name prefix and arguments  separator !",
			fields: fields{botName: "toster"},
			args:   args{msg: "!ыsearch@otherbot keyword"},
			want:   false,
			want1:  command{Cmd: "", Arg: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandParser{
				botName:      tt.fields.botName,
				reSpecSymbol: regexp.MustCompile(`^([!/])[A-zА-я0-9]+$`),
			}
			got, got1 := c.ParseCommand(tt.args.msg)
			if got != tt.want {
				t.Errorf("CommandParser.ParseCommand() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CommandParser.ParseCommand() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
