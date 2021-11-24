package main

import (
	"reflect"
	"testing"
)

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
	if !(result == "Alien 1 is in Rome.\nAlien 7 is in Milan.\n" || result == "Alien 7 is in Milan.\nAlien 1 is in Rome.\n") {
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

func TestOutStepTestKillAndDestroy(t *testing.T) {

	worldMap2 := make(map[string]map[string]string)
	worldMap2["Rome"] = make(map[string]string)
	worldMap2["Rome"]["Milan"] = "north"
	worldMap2["Rome"]["Venice"] = "north"
	worldMap2["Milan"] = make(map[string]string)
	worldMap2["Milan"]["Rome"] = "south"
	worldMap2["Milan"]["Venice"] = "east"
	worldMap2["Venice"] = make(map[string]string)
	worldMap2["Venice"]["Milan"] = "west"
	worldMap2["Venice"]["Rome"] = "south"

	t0WorldMap2 := make(map[string]map[string]string)
	for k, v := range worldMap2 {
		t0WorldMap2[k] = v
	}

	worldMap1 := make(map[string]map[string]string)
	worldMap1["Lanciano"] = make(map[string]string)
	worldMap1["Lanciano"]["Turin"] = "north"
	worldMap1["Turin"] = make(map[string]string)
	worldMap1["Turin"]["Lanciano"] = "south"
	worldMap1["Turin"]["Gdansk"] = "east"
	worldMap1["Gdansk"] = make(map[string]string)
	worldMap1["Gdansk"]["Turin"] = "west"

	t0WorldMap1 := make(map[string]map[string]string)
	for k, v := range worldMap1 {
		t0WorldMap1[k] = v
	}

	alienStatus2 := make(map[int]string)
	alienStatus2[1] = "Rome"
	alienStatus2[2] = "Rome"
	alienStatus2[3] = "Rome"
	alienStatus2[4] = "Rome"
	alienStatus2[5] = "Milan"
	alienStatus2[9] = "Milan"
	alienStatus2[14] = "Venice"

	t0AlienStatus2 := make(map[int]string)
	for k, v := range alienStatus2 {
		t0AlienStatus2[k] = v
	}

	alienStatus1 := make(map[int]string)
	alienStatus1[1] = "Lanciano"
	alienStatus1[2] = "Turin"
	alienStatus1[3] = "Gdansk"

	t0AlienStatus1 := make(map[int]string)
	for k, v := range alienStatus1 {
		t0AlienStatus1[k] = v
	}

	invaders2 := make(map[string][]int)
	invaders2["Rome"] = []int{1, 2, 3, 4}
	invaders2["Milan"] = []int{5, 9}
	invaders2["Venice"] = []int{14}

	t0Invaders2 := make(map[string][]int)
	for k, v := range invaders2 {
		t0Invaders2[k] = v
	}

	invaders1 := make(map[string][]int)
	invaders1["Lanciano"] = []int{1}
	invaders1["Turin"] = []int{2}
	invaders1["Gdansk"] = []int{3}

	t0Invaders1 := make(map[string][]int)
	for k, v := range invaders1 {
		t0Invaders1[k] = v
	}

	//t1 Expectations: No changes
	OutStepKillAndDestroy(alienStatus1, invaders1, worldMap1, false)

	if !(reflect.DeepEqual(t0AlienStatus1, alienStatus1)) {
		t.Errorf("Test failed - t0Alienstatus1")
	}

	if !(reflect.DeepEqual(t0Invaders1, invaders1)) {
		t.Errorf("Test failed - t0Invaders1")
	}

	if !(reflect.DeepEqual(t0WorldMap1, worldMap1)) {
		t.Errorf("Test failed - t0Alienstatus1")
	}

	//t2 Milan destroyed, 5 and 9 killed

	//---
	t2WorldMap2 := make(map[string]map[string]string)
	t2WorldMap2["Rome"] = make(map[string]string)
	t2WorldMap2["Rome"]["Venice"] = "north"
	t2WorldMap2["Venice"] = make(map[string]string)
	t2WorldMap2["Venice"]["Rome"] = "south"

	t2AlienStatus2 := make(map[int]string)
	t2AlienStatus2[1] = "Rome"
	t2AlienStatus2[2] = "Rome"
	t2AlienStatus2[3] = "Rome"
	t2AlienStatus2[4] = "Rome"
	t2AlienStatus2[14] = "Venice"

	t2Invaders2 := make(map[string][]int)
	t2Invaders2["Rome"] = []int{1, 2, 3, 4}
	t2Invaders2["Venice"] = []int{14}
	//---

	OutStepKillAndDestroy(alienStatus2, invaders2, worldMap2, true)

	if !(reflect.DeepEqual(t2AlienStatus2, alienStatus2)) {
		t.Errorf("Test failed - t1Alienstatus1")
	}

	if !(reflect.DeepEqual(t2Invaders2, invaders2)) {
		t.Errorf("Test failed - t1Invaders2")
	}

	if !(reflect.DeepEqual(t2WorldMap2, worldMap2)) {
		t.Errorf("Test failed - t1WorldMap2")
	}

	//t3 no changes from t2
	OutStepKillAndDestroy(alienStatus2, invaders2, worldMap2, true)
	if !(reflect.DeepEqual(t2AlienStatus2, alienStatus2)) {
		t.Errorf("Test failed - t3Alienstatus1")
	}

	if !(reflect.DeepEqual(t2Invaders2, invaders2)) {
		t.Errorf("Test failed - t3Invaders2")
	}

	if !(reflect.DeepEqual(t2WorldMap2, worldMap2)) {
		t.Errorf("Test failed - t3WorldMap2")
	}

	//t4 Rome destroyed, 1, 2, 3, and 4 killed
	OutStepKillAndDestroy(alienStatus2, invaders2, worldMap2, false)

	//t4 expected alienStatus2, invaders2, worldMap2

	t4WorldMap2 := make(map[string]map[string]string)
	t4WorldMap2["Venice"] = make(map[string]string)

	t4Invaders2 := make(map[string][]int)
	t4Invaders2["Venice"] = []int{14}

	t4AlienStatus2 := make(map[int]string)
	t4AlienStatus2[14] = "Venice"

	if !(reflect.DeepEqual(t4AlienStatus2, alienStatus2)) {
		t.Errorf("Test failed - t4Alienstatus1")
	}

	if !(reflect.DeepEqual(t4Invaders2, invaders2)) {
		t.Errorf("Test failed - t4Invaders2")
	}

	if !(reflect.DeepEqual(t4WorldMap2, worldMap2)) {
		t.Errorf("Test failed - t4WorldMap2")
	}

	//t5 from the t0 copies kill and destroy in nonstrict mode -> expected same as t4
	OutStepKillAndDestroy(t0AlienStatus2, t0Invaders2, t0WorldMap2, false)

	if !(reflect.DeepEqual(t4AlienStatus2, t0AlienStatus2)) {
		t.Errorf("Test failed - t5Alienstatus1")
	}

	if !(reflect.DeepEqual(t4Invaders2, t0Invaders2)) {
		t.Errorf("Test failed - t5Invaders2")
	}

	if !(reflect.DeepEqual(t4WorldMap2, t0WorldMap2)) {
		t.Errorf("Test failed - t5WorldMap2")
	}
}

func TestInStepTestKillAndDestroy(t *testing.T) {

	worldMap := make(map[string]map[string]string)
	worldMap["Rome"] = make(map[string]string)
	worldMap["Rome"]["Milan"] = "north"
	worldMap["Rome"]["Venice"] = "north"
	worldMap["Milan"] = make(map[string]string)
	worldMap["Milan"]["Rome"] = "south"
	worldMap["Milan"]["Venice"] = "east"
	worldMap["Venice"] = make(map[string]string)
	worldMap["Venice"]["Milan"] = "west"
	worldMap["Venice"]["Rome"] = "south"
	worldMap["cityWithNoAliens"] = make(map[string]string)

	t0WorldMap := make(map[string]map[string]string)
	for k, v := range worldMap {
		t0WorldMap[k] = v
	}

	alienStatus := make(map[int]string)
	alienStatus[1] = "Rome"
	alienStatus[2] = "Rome"
	alienStatus[3] = "Rome"
	alienStatus[4] = "Rome"
	alienStatus[5] = "Milan"
	alienStatus[9] = "Milan"
	alienStatus[14] = "Venice"

	t0AlienStatus := make(map[int]string)
	for k, v := range alienStatus {
		t0AlienStatus[k] = v
	}

	invaders := make(map[string][]int)
	invaders["Rome"] = []int{1, 2, 3, 4}
	invaders["Milan"] = []int{5, 9}
	invaders["Venice"] = []int{14}

	t0Invaders := make(map[string][]int)
	for k, v := range invaders {
		t0Invaders[k] = v
	}

	cityNotInTheMap := "Pescara"

	//t1 Expectations: No changes
	InStepKillAndDestroy(cityNotInTheMap, alienStatus, invaders, worldMap)
	InStepKillAndDestroy("cityWithNoAliens", alienStatus, invaders, worldMap)

	if !(reflect.DeepEqual(t0AlienStatus, alienStatus)) {
		t.Errorf("Test failed - t0Alienstatus")
	}

	if !(reflect.DeepEqual(t0Invaders, invaders)) {
		t.Errorf("Test failed - t0Invaders")
	}

	if !(reflect.DeepEqual(t0WorldMap, worldMap)) {
		t.Errorf("Test failed - t0Alienstatus1")
	}

	//t2 Rome destroyed, 1, 2, 3, 4 killed

	//---
	t2WorldMap := make(map[string]map[string]string)
	t2WorldMap["Milan"] = make(map[string]string)
	t2WorldMap["Milan"]["Venice"] = "east"
	t2WorldMap["Venice"] = make(map[string]string)
	t2WorldMap["Venice"]["Milan"] = "west"
	t2WorldMap["cityWithNoAliens"] = make(map[string]string)

	t2AlienStatus := make(map[int]string)
	t2AlienStatus[14] = "Venice"
	t2AlienStatus[5] = "Milan"
	t2AlienStatus[9] = "Milan"

	t2Invaders := make(map[string][]int)
	t2Invaders["Venice"] = []int{14}
	t2Invaders["Milan"] = []int{5, 9}
	//---

	InStepKillAndDestroy("Rome", alienStatus, invaders, worldMap)

	if !(reflect.DeepEqual(t2AlienStatus, alienStatus)) {
		t.Errorf("Test failed - t1Alienstatus")
	}

	if !(reflect.DeepEqual(t2Invaders, invaders)) {
		t.Errorf("Test failed - t2Invaders")
	}

	if !(reflect.DeepEqual(t2WorldMap, worldMap)) {
		t.Errorf("Test failed - t2WorldMap")
	}

	//t3 will destroy Milan and 14

	//alienStatus and invaders will be the same, will remove Milan from map, as a city and as direction

	t3WorldMap := make(map[string]map[string]string)
	t3WorldMap["Venice"] = make(map[string]string)
	t3WorldMap["cityWithNoAliens"] = make(map[string]string)

	t3AlienStatus := make(map[int]string)
	t3AlienStatus[14] = "Venice"

	t3Invaders := make(map[string][]int)
	t3Invaders["Venice"] = []int{14}

	InStepKillAndDestroy("Milan", alienStatus, invaders, worldMap)
	if !(reflect.DeepEqual(t3AlienStatus, alienStatus)) {
		t.Errorf("Test failed - t3Alienstatus")
	}

	if !(reflect.DeepEqual(t3Invaders, invaders)) {
		t.Errorf("Test failed - t3Invaders")
	}

	if !(reflect.DeepEqual(t3WorldMap, worldMap)) {
		t.Errorf("Test failed - t3WorldMap")
	}

	newWorldMap := make(map[string]map[string]string)
	newWorldMap["Berlin"] = make(map[string]string)

	newAlienStatus := make(map[int]string)
	newAlienStatus[6] = "Berlin"

	newInvaders := make(map[string][]int)
	newInvaders["Berlin"] = []int{6}

	InStepKillAndDestroy("Berlin", newAlienStatus, newInvaders, newWorldMap)

	if !(len(newAlienStatus) == 0 && len(newInvaders) == 0 && len(newWorldMap) == 0) {
		t.Errorf("Test failed - Total destruction failed")
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
