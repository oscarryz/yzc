package internal

import (
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		fileName string
		tokens []token
		want    *program
		wantErr bool
	}{
		{
			name: "Empty body",
			fileName: "empty.yz",
			tokens: []token{
				{pos(1, 5), EOF, "EOF"},
			},
			want: &program{
				&blockBody{
					expressions: []expression{
						&empty{},
					},
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
