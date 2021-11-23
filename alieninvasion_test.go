package main

import "testing"

func TestIndexOfAlien(t *testing.T) {

	aliens := []int{1, 2, 6, 8, 10, 67}
	alienName := 8
	index := IndexOfAlien(aliens, alienName)
	if index != 3 {
		t.Errorf("Test failed, got: %d, want: %d.", index, 3)
	}
}

func TestRemoveAlien(t *testing.T) {
	aliens := []int{1, 2, 6, 8, 10, 67}
	alienName := 8
	result := RemoveAlien(aliens, alienName)
	if IndexOfAlien(result, 1) == -1 ||
		IndexOfAlien(result, 2) == -1 ||
		IndexOfAlien(result, 6) == -1 ||
		IndexOfAlien(result, 10) == -1 ||
		IndexOfAlien(result, 67) == -1 ||
		IndexOfAlien(result, 8) != -1 ||
		len(result) != len(aliens)-1 {
		t.Errorf("Test failed")
	}

}

func TestWhereAliens(t *testing.T) {
	landedInvaders := map[int]string{
		1: "Rome",
		7: "Milan",
	}

	result := WhereAliens(landedInvaders)
	if result != "Alien 1 is in Rome.\nAlien 7 is in Milan.\n" {
		t.Errorf("Test failed\nExpecetd:\nAlien 1 is in Rome\nAlien 7 is in Milan\nGot:\n%s", result)
	}
}
