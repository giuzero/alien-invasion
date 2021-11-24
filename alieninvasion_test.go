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

func TestInStepKillAndDestroy(t *testing.T) {

	worldMap2 := make(map[string]map[string]string)
	worldMap2["Rome"] = make(map[string]string)
	worldMap2["Rome"]["Milan"] = "north"
	worldMap2["Milan"] = make(map[string]string)
	worldMap2["Milan"]["Rome"] = "south"
	worldMap2["Milan"]["Venice"] = "east"
	worldMap2["Venice"] = make(map[string]string)
	worldMap2["Venice"]["Milan"] = "west"

	worldMap1 := make(map[string]map[string]string)
	worldMap1["Pescara"] = make(map[string]string)
	worldMap1["Pescara"]["Turin"] = "north"
	worldMap1["Turin"] = make(map[string]string)
	worldMap1["Turin"]["Pescara"] = "south"
	worldMap1["Turin"]["Gdansk"] = "east"
	worldMap1["Gdansk"] = make(map[string]string)
	worldMap1["Gdansk"]["Turin"] = "west"

	alienStatus2 := make(map[int]string)
	alienStatus2[1] = "Rome"
	alienStatus2[2] = "Rome"
	alienStatus2[3] = "Rome"
	alienStatus2[4] = "Rome"
	alienStatus2[5] = "Milan"
	alienStatus2[9] = "Milan"

	alienStatus1 := make(map[int]string)
	alienStatus1[1] = "Pescara"
	alienStatus1[2] = "Turin"
	alienStatus1[3] = "Gdansk"

	invaders2 := make(map[string][]int)
	invaders2["Rome"] = []int{1, 2, 3, 4}
	invaders2["Milan"] = []int{5, 9}

	invaders1 := make(map[string][]int)
	invaders1["Pescara"] = []int{1}
	invaders1["Turin"] = []int{2}
	invaders1["Gdansk"] = []int{3}

	tmp1, tmp2, tmp3 := alienStatus1, invaders1, worldMap1
	tmp4, tmp5, tmp6 := alienStatus2, invaders2, worldMap2
	ret1, ret2, ret3 := OutStepKillAndDestroy(tmp1, tmp2, tmp3, false)
	ret4, ret5, ret6 := OutStepKillAndDestroy(tmp4, tmp5, tmp6, true)

	if !(len(tmp1) == len(ret1) &&
		len(tmp2) == len(ret2) &&
		len(tmp3) == len(ret3)) {
		t.Errorf("Test failed")
	}

	if !(len(tmp4) == len(ret4)-1 &&
		len(tmp5) == len(ret5)-1 &&
		len(tmp6) == len(ret6)-1) {
		t.Errorf("Test failed")
	}

	for k := range alienStatus1 {
		if !(alienStatus1[k] == tmp1[k]) {
			t.Errorf("Test failed")
		}
	}
	for k := range invaders1 {
		if !(len(invaders1[k]) != len(tmp2[k])) {
			t.Errorf("Test failed")
		}
	}
	for k := range worldMap1 {
		if !(len(worldMap1[k]) != len(tmp3[k])) {
			t.Errorf("Test failed")
		}
	}
	for k := range alienStatus2 {
		if !(alienStatus2[k] == "Rome") {
			t.Errorf("Test failed")
		}
	}
	for k := range invaders2 {
		if !(len(invaders2[k]) != len(tmp2[k])) {
			t.Errorf("Test failed")
		}
	}
	for k := range worldMap1 {
		if !(len(worldMap1[k]) != len(tmp3[k])) {
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
