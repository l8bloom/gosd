package gosd

import "testing"

func TestLoad(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatalf("Can't load stable-diffusion.cpp dynamic libraries.")
	}
}
