package main

import "testing"

func TestSlugFromTitle(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{title: "我的第一篇 Go 文章", want: "wo-de-di-yi-pian-go-wen-zhang"},
		{title: "My First Go Article!", want: "my-first-go-article"},
	}

	for _, test := range tests {
		if got := slugFromTitle(test.title); got != test.want {
			t.Errorf("slugFromTitle(%q) = %q, want %q", test.title, got, test.want)
		}
	}
}
