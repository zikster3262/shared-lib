package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFileName(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				url: "https://asura.gg/wp-content/uploads/2022/12/01-86.jpg",
			},
			want: "01-86.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetFileName(tt.args.url))
		})
	}
}

func TestGetIDFromChapterURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				url: "https://asura.gg/overpowered-sword-chapter-46/",
			},
			want: "46",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIDFromChapterURL(tt.args.url); got != tt.want {
				t.Errorf("GetIDFromChapterURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
