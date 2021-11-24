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

func TestNextDestination(t *testing.T) {

	worldMap := make(map[string]map[string]string)
	worldMap["Rome"] = make(map[string]string)
	worldMap["Rome"]["Milan"] = "north"
	worldMap["Milan"] = make(map[string]string)
	worldMap["Milan"]["Rome"] = "south"
	worldMap["Milan"]["Venice"] = "east"
	worldMap["Venice"] = make(map[string]string)
	worldMap["Venice"]["Milan"] = "west"
	worldMap["Naples"] = make(map[string]string)

	expectedDirectionsFromRome := []string{"Milan"}
	expectedDirectionsFromMilan := []string{"Rome", "Venice"}
	expectedDirectionsFromNaples := []string{"Naples"}
	expectedDirectionsFromVenice := []string{"Milan"}

	tables := []struct {
		now   string
		world map[string]map[string]string
		next  []string
	}{
		{"Rome", worldMap, expectedDirectionsFromRome},
		{"Milan", worldMap, expectedDirectionsFromMilan},
		{"Naples", worldMap, expectedDirectionsFromNaples},
		{"Venice", worldMap, expectedDirectionsFromVenice},
	}

	for _, table := range tables {
		nextCity := NextDestination(table.now, table.world)
		if !stringInSlice(nextCity, table.next) {
			t.Errorf("Test Error")
		}
	}
}

func TestLanding(t *testing.T) {
	worldMap := make(map[string]map[string]string)
	worldMap["Rome"] = make(map[string]string)
	worldMap["Rome"]["Milan"] = "north"
	worldMap["Milan"] = make(map[string]string)
	worldMap["Milan"]["Rome"] = "south"
	worldMap["Milan"]["Venice"] = "east"
	worldMap["Venice"] = make(map[string]string)
	worldMap["Venice"]["Milan"] = "west"

	numberOfAliens1 := 8
	numberOfAliens2 := 3
	numberOfAliens3 := 2

	tables := []struct {
		world  map[string]map[string]string
		aliens int
	}{
		{worldMap, numberOfAliens1},
		{worldMap, numberOfAliens2},
		{worldMap, numberOfAliens3},
	}

	for _, table := range tables {
		result1, result2 := Landing(table.world, table.aliens)
		n := 0

		uniqueCities := make(map[string]int)
		for k := range result2 {
			n = n + len(result2[k])
		}
		for _, v := range result1 {
			if _, exist := uniqueCities[v]; !exist {
				uniqueCities[v] = 1
			} else {
				uniqueCities[v]++
			}
		}
		aliensByCity := 0
		for k := range uniqueCities {
			aliensByCity = aliensByCity + uniqueCities[k]
		}
		if !(n == table.aliens) ||
			!(aliensByCity == table.aliens) {
			t.Errorf("Test Error")
		}
	}

}

func TestExclusiveLanding(t *testing.T) {
	//will use this function only if aliens<=cities
	worldMap := make(map[string]map[string]string)
	worldMap["Rome"] = make(map[string]string)
	worldMap["Rome"]["Milan"] = "north"
	worldMap["Milan"] = make(map[string]string)
	worldMap["Milan"]["Rome"] = "south"
	worldMap["Milan"]["Venice"] = "east"
	worldMap["Venice"] = make(map[string]string)
	worldMap["Venice"]["Milan"] = "west"

	numberOfAliens2 := 3
	numberOfAliens3 := 2

	tables := []struct {
		world  map[string]map[string]string
		aliens int
	}{
		{worldMap, numberOfAliens2},
		{worldMap, numberOfAliens3},
	}

	for _, table := range tables {
		result1, result2 := ExclusiveLanding(table.world, table.aliens)
		n := 0

		uniqueCities := make(map[string]int)
		for k := range result2 {
			n = n + len(result2[k])
		}
		for _, v := range result1 {
			if _, exist := uniqueCities[v]; !exist {
				uniqueCities[v] = 1
			} else {
				uniqueCities[v]++
			}
		}
		aliensByCity := 0
		//checkMaxOneAlien := true
		for k := range uniqueCities {
			aliensByCity = aliensByCity + uniqueCities[k]
			if uniqueCities[k] > 1 {
				t.Errorf("Test Error: more than one alien in a city")
			}
		}
		if !(n == table.aliens) ||
			!(aliensByCity == table.aliens) {
			t.Errorf("Test Error")
		}
	}

}

func TestCheckingArgs(t *testing.T) {
	tables := []struct {
		args []string
		ret  int
	}{
		{[]string{}, 0},
		{[]string{"debugname"}, 0},
		{[]string{"debugname", "3,5"}, 0},
		{[]string{"debugname", "-4"}, -4},
		{[]string{"debugname", "4"}, 4},
		{[]string{"debugname", "0"}, 0},
		{[]string{"debugname", "test"}, 0},
	}

	for _, table := range tables {
		res := CheckingArgs(table.args)
		if res != table.ret {
			t.Errorf("Test failed")
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
