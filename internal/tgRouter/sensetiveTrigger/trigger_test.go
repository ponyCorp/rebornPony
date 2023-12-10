package sensetivetrigger

import (
	"reflect"
	"testing"
)

func Test_sensetiveTrigger_recompile(t *testing.T) {
	type fields struct {
		groups map[int64]*group
	}
	type args struct {
		chatID  int64
		words   []string
		regxStr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{

		{
			name: "Test_sensetiveTrigger_recompile",
			fields: fields{
				groups: make(map[int64]*group),
			},
			args: args{
				chatID:  1,
				words:   []string{"123", "321"},
				regxStr: "\\b123\\b|\\b321\\b",
			},
			wantErr: false,
		},
		//one word
		{
			name: "one word",
			fields: fields{
				groups: make(map[int64]*group),
			},
			args: args{
				chatID:  1,
				words:   []string{"123"},
				regxStr: "\\b123\\b",
			},
			wantErr: false,
		},
		//zero word
		{
			name: "empty word",
			fields: fields{
				groups: make(map[int64]*group),
			},
			args: args{
				chatID:  1,
				words:   []string{},
				regxStr: "",
			},
			wantErr: false,
		},

		//word  with *
		{
			name: "word  with *",
			fields: fields{
				groups: make(map[int64]*group),
			},
			args: args{
				chatID:  1,
				words:   []string{"*123*"},
				regxStr: "123",
			},
			wantErr: false,
		},
		//toster toster
		{
			name: "toster toster",
			fields: fields{
				groups: make(map[int64]*group),
			},
			args: args{
				chatID:  1,
				words:   []string{"toster", "toster"},
				regxStr: "\\btoster\\b|\\btoster\\b",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SensetiveTrigger{
				groups: tt.fields.groups,
			}
			s.groups[tt.args.chatID] = &group{
				chatID: tt.args.chatID,
				words:  tt.args.words,
			}

			if err := s.recompile(tt.args.chatID); (err != nil) != tt.wantErr {
				t.Errorf("sensetiveTrigger.recompile() error = %v, wantErr %v", err, tt.wantErr)
			}
			regStr := s.groups[tt.args.chatID].regexp.String()
			//	t.Log(regStr)
			if regStr != tt.args.regxStr {
				t.Errorf("sensetiveTrigger.recompile() error = %v, wantErr %v", regStr, tt.args.regxStr)
			}
		})
	}
}

func TestSensetiveTrigger_MessageContainSensitiveWords(t *testing.T) {
	type fields struct {
		groups map[int64]*group
	}
	type args struct {
		chatID    int64
		message   string
		words     []string
		regexpStr string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{

		{
			name: "TestSensetiveTrigger_MessageContainSensitiveWords",
			fields: fields{
				groups: make(map[int64]*group),
			},
			args: args{
				words:     []string{"toster", "toster"},
				regexpStr: "\\btoster\\b|\\btoster\\b",
				chatID:    1,
				message:   "toster toster",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			for _, w := range tt.args.words {
				err := s.AddWords(tt.args.chatID, w)
				if err != nil {
					t.Errorf("SensetiveTrigger.AddWords() error = %v", err)
					return
				}
			}
			regx, err := s.GetRegexpStr(tt.args.chatID)
			if err != nil {
				t.Errorf("SensetiveTrigger.GetRegexp() error = %v", err)
				return
			}

			if regx != tt.args.regexpStr {
				t.Errorf("sensetiveTrigger.recompile() error = %v, wantErr %v", regx, tt.args.regexpStr)
				return
			}
			got, err := s.MessageContainSensitiveWords(tt.args.chatID, tt.args.message)
			if err != nil {
				t.Errorf("SensetiveTrigger.MessageContainSensitiveWords() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("SensetiveTrigger.MessageContainSensitiveWords() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestSensetiveTrigger_AddWords(t *testing.T) {
	type fields struct {
		groups map[int64]*group
	}
	type args struct {
		chatID int64
		words  []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSensetiveTrigger_AddWords",
			fields: fields{
				groups: make(map[int64]*group),
			},
			args: args{
				chatID: 1,
				words:  []string{"toster", "toster"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SensetiveTrigger{
				groups: tt.fields.groups,
			}
			for _, w := range tt.args.words {

				err := s.AddWords(tt.args.chatID, w)
				if err != nil {
					t.Errorf("SensetiveTrigger.AddWords() error = %v", err)
				}
			}
			added, err := s.GetWords(tt.args.chatID)
			if err != nil {
				t.Errorf("SensetiveTrigger.GetWords() error = %v", err)
				return
			}
			if !reflect.DeepEqual(tt.args.words, added) {
				t.Errorf("SensetiveTrigger.GetWords() = %v, want %v", added, tt.args.words)

			}
		})
	}
}

func Test_join(t *testing.T) {
	type args struct {
		words  []string
		regStr string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_join",
			args: args{
				words:  []string{"toster", "toster"},
				regStr: "\\btoster\\b|\\btoster\\b",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if join(tt.args.words) != tt.args.regStr {
				t.Errorf("join() = %v, want %v", tt.args.regStr, tt.args.regStr)
			}
		})
	}
}

func Test_regexpMatch(t *testing.T) {
	type args struct {
		msg  string
		regx string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Test_regexpMatch",
			args: args{
				msg:  "toster toster",
				regx: "\\btoster\\b|\\btoster\\b",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test_regexpMatch",
			args: args{
				msg:  "tosters toster",
				regx: "toster",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test_regexpMatch boarder case",
			args: args{
				msg:  "tosters toster",
				regx: "\\btoster\\b",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := regexpMatch(tt.args.msg, tt.args.regx)
			if err != nil {
				t.Errorf("regexpMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("regexpMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
