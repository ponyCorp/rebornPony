package trie

import "testing"

func TestFilter_ContainsWord(t *testing.T) {

	type args struct {
		message string
		fWords  []string
	}
	tests := []struct {
		name   string
		filter *Filter
		args   args
		want   bool
	}{
		{
			name: "test",
			filter: &Filter{
				root: &TrieNode{
					children: make(map[rune]*TrieNode),
					isWord:   false,
				},
			},
			args: args{
				message: "testf",
				fWords:  []string{"test"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, word := range tt.args.fWords {
				tt.filter.AddWord(word)
			}

			if got := tt.filter.ContainsWord(tt.args.message); got != tt.want {
				t.Errorf("Filter.ContainsWord() = %v, want %v", got, tt.want)
			}
		})
	}
}
