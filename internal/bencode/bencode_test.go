package bencode

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "basic string",
			input:    "5:hello",
			expected: []byte("hello"),
			wantErr:  false,
		},
		{
			name:     "longer string",
			input:    "10:helloworld",
			expected: []byte("helloworld"),
			wantErr:  false,
		},
		{
			name:     "single character",
			input:    "1:a",
			expected: []byte("a"),
			wantErr:  false,
		},
		{
			name:     "single number",
			input:    "i1e",
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "two digit number",
			input:    "i12e",
			expected: 12,
			wantErr:  false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "missing e",
			input:   "i12",
			wantErr: true,
		},
		{
			name:  "simple list",
			input: "l5:helloi1ee",
			expected: []interface{}{
				[]byte("hello"),
				1,
			},
		},
		{
			name:     "empty list",
			input:    "le",
			expected: []interface{}{},
			wantErr:  false,
		},
		{
			name:  "nested list with integer and string",
			input: "lli956e5:appleee",
			expected: []interface{}{
				[]interface{}{
					956,
					[]byte("apple"),
				},
			},
			wantErr: false,
		},
		{
			name:  "nested list",
			input: "ll5:helloi1eei2ee",
			expected: []interface{}{
				[]interface{}{
					[]byte("hello"),
					1,
				},
				2,
			},
			wantErr: false,
		},
		{
			name:  "list with multiple types",
			input: "l5:helloi1ei2ee",
			expected: []interface{}{
				[]byte("hello"),
				1,
				2,
			},
			wantErr: false,
		},
		{
			name:  "dictionary",
			input: "d3:foo3:bar5:helloi52ee",
			expected: map[string]interface{}{
				"foo":   []byte("bar"),
				"hello": 52,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := Decode([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeBencode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("decodeBencode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "simple string",
			args: args{
				value: []byte("hello"),
			},
			want: []byte("5:hello"),
		},
		{
			name: "longer string",
			args: args{
				value: []byte("helloworld"),
			},
			want: []byte("10:helloworld"),
		},
		{
			name: "integer",
			args: args{
				value: 123,
			},
			want: []byte("i123e"),
		},
		{
			name: "list",
			args: args{
				value: []interface{}{
					[]byte("hello"),
					123,
				},
			},
			want: []byte("l5:helloi123ee"),
		},
		{
			name: "dictionary",
			args: args{
				value: map[string]interface{}{
					"foo":   []byte("bar"),
					"hello": 123,
				},
			},
			want: []byte("d3:foo3:bar5:helloi123ee"),
		},
		{
			name: "nested dictionary",
			args: args{
				value: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": 123,
					},
				},
			},
			want: []byte("d3:food3:bari123eee"),
		},
		{
			name: "sorts keys",
			args: args{
				value: map[string]interface{}{
					"hello": 123,
					"foo":   []byte("bar"),
				},
			},
			want: []byte("d3:foo3:bar5:helloi123ee"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
