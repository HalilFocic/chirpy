package main

import "testing"

// Function that replaces all bad words in a message with "****". It should be case sensitive. But it should look
func TestBadWorkReplacement(t *testing.T) {
	cases := []struct {
		message string
		want    string
	}{
		{"This is a kerfuffle", "This is a ****"},
		{"This is a sharbert", "This is a ****"},
		{"This is a fornax", "This is a ****"},
		{"This is a Fornax", "This is a ****"},
		{"This is a Fornax!", "This is a Fornax!"},
	}
	for _, c := range cases {
		actual := BadWorkReplacement(c.message)
		if actual != c.want {
			t.Errorf(" %q, want %q", actual, c.want)
		}
	}

}
