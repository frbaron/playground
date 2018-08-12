package main

import (
	"testing"
)

func TestWeeder(t *testing.T) {
	cases := []struct {
		in, weed, want string
	}{
		{"S01E01.mp4",		"WEB-HD", "S01E01.mp4"},
		{"S01E01 WEB-HD.mp4",	"WEB-HD", "S01E01.mp4"},
		{"S01E01  xVid.mp4",	"WEB-HD", "S01E01 xVid.mp4"},
		{"S01E01  xVid.mp4",	"xVid", "S01E01.mp4"},
	}

	for _, c := range cases {
		got := Weeder(c.in, c.weed)

		if got != c.want {
			t.Errorf("Weeder(%q, %q) == %q, want %q", c.in, c.weed, got, c.want)
		}
	}
}
