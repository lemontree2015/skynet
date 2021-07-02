package skynet

import (
	"testing"
)

func TestCriteriaBasic(t *testing.T) {
	criteria1 := NewCriteria()
	si1 := NewServiceInfo("fake-name1", "1.0.0")
	si2 := NewServiceInfo("fake-name1", "1.0.1")
	si3 := NewServiceInfo("fake-name2", "1.1.0")
	si4 := NewServiceInfo("fake-name3", "1.1.1")

	if !criteria1.Matches(si1) {
		t.Fatal("TestCriteriaBasic Error")
	}

	if !criteria1.Matches(si2) {
		t.Fatal("TestCriteriaBasic Error")
	}

	if !criteria1.Matches(si3) {
		t.Fatal("TestCriteriaBasic Error")
	}

	if !criteria1.Matches(si4) {
		t.Fatal("TestCriteriaBasic Error")
	}

	criteria1.AddService("fake-name1", "1.0.0")
	if !criteria1.Matches(si1) {
		t.Fatal("TestCriteriaBasic Error")
	}

	if criteria1.Matches(si2) {
		t.Fatal("TestCriteriaBasic Error")
	}

	if criteria1.Matches(si3) {
		t.Fatal("TestCriteriaBasic Error")
	}

	if criteria1.Matches(si4) {
		t.Fatal("TestCriteriaBasic Error")
	}

	criteria1.AddService("fake-name1", "1.0.0")
	criteria1.AddService("fake-name1", "")
	if !criteria1.Matches(si1) {
		t.Fatal("TestCriteriaBasic Error")
	}

	if !criteria1.Matches(si2) {
		t.Fatal("TestCriteriaBasic Error")
	}

	if criteria1.Matches(si3) {
		t.Fatal("TestCriteriaBasic Error")
	}

	if criteria1.Matches(si4) {
		t.Fatal("TestCriteriaBasic Error")
	}
}
