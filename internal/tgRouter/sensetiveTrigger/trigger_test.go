package sensetivetrigger

import "testing"

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
				regxStr: "\b123\b|\b321\b",
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
				regxStr: "\b123\b",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sensetiveTrigger{
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
