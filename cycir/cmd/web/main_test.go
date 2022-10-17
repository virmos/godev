package main

import "testing"

func TestRun(t *testing.T) {
	_, err := main()
	if err != nil {
		t.Error("failed main()")
	}
}
