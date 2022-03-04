package main

import (
	"testing"
)

func TestSetup(t *testing.T) {
	svcDir = "testdata/services"

	_, err := Setup()
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}
}

func TestSetup_MissingSvcDir(t *testing.T) {
	svcDir = "/tmp/this/dir/hopefully/doesnt/exist"

	_, err := Setup()
	if err == nil {
		t.Errorf("expected error, received none")
	}
}

func TestSetup_MissingCerts(t *testing.T) {
	svcDir = "testdata/services"
	certDir = "/tmp/this/dir/hopefully/doesnt/exist"

	_, err := Setup()
	if err == nil {
		t.Errorf("expected error, received none")
	}
}
