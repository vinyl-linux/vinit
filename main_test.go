package main

import (
	"testing"
)

func TestSetup(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Errorf("unexpected panic: %#v", err)
		}
	}()

	svcDir = "testdata/services"

	Setup()
}
