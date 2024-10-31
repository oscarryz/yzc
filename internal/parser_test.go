package internal

import (
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		tokens   []token
		want     *boc
		wantErr  bool
	}{
		{
			name:     "Nested directory",
			fileName: "parent/simple.yz",
			tokens: []token{
				{pos(0, 0), EOF, "EOF"},
			},
			want: &boc{
				Name:    "parent",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{
						&boc{
							Name: "simple",
							blockBody: &blockBody{
								expressions: []expression{},
								statements:  []statement{},
							},
						},
					},
					statements: []statement{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.fileName, tt.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// compare go recursively
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v\n want %v", got, tt.want)
			}
		})
	}
}
