package gosd

import "testing"

func TestLoad(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatal(err.Error())
	}
}
