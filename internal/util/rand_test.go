package util

import "testing"

func TestRandom(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default",
			args: args{
				length: 5,
			},
			want: "12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Random(tt.args.length); len(got) != len(tt.want) {
				t.Errorf("Random() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomCharset(t *testing.T) {
	type args struct {
		length  int
		charset string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default",
			args: args{
				length:  5,
				charset: "A",
			},
			want: "AAAAA",
		},
		{
			name: "other",
			args: args{
				length:  5,
				charset: "b",
			},
			want: "bbbbb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomCharset(tt.args.length, tt.args.charset); got != tt.want {
				t.Errorf("RandomCharset() = %v, want %v", got, tt.want)
			}
		})
	}
}
