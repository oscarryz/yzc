package internal

import (
	"reflect"
	"testing"
)

func Test_parser_program(t *testing.T) {
	type fields struct {
		fileName string
		tokens   []token
	}
	EOF := token{pos(0, 0), EOF, "EOF"}
	tests := []struct {
		name    string
		fields  fields
		want    *program
		wantErr bool
	}{
		// TODO: Add test cases.

		{"Empty file",
			fields{
				"a.yz",
				[]token{EOF},
			},
			&program{&blockBody{
				[]expression{},
				[]statement{},
			},
			}, false,
		},
		{"Test literals",
			fields{
				"a.yz",
				[]token{token{pos(1, 1), INTEGER, "1"}, EOF},
			},
			&program{&blockBody{
				[]expression{
					&BasicLit{pos(1, 1), INTEGER, "1"},
				},
				[]statement{},
			}}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newParser(tt.fields.fileName, tt.fields.tokens)
			got, err := p.program()
			if (err != nil) != tt.wantErr {
				t.Errorf("program() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("program() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}
