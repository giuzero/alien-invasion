package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

//number of steps
var steps int = 10

//true: will "clean" the map just after landing, destroying cities invade by two or more
//alien. After t0 you'll have at most one alien per city. Will they still exist?
var doIwantFightsAtT0 bool = false

//true: enable kills at the end of the turn/step.
//false: when an alien arrives in a city, could find one or more aliens there.
//Will kill an destroy during the step
var killAndDestroyAtTheEndOfShift bool = false

//destroy and kill aliens if and only if the city is invaded by exactly 2 aliens
var strictTwoAliensDestroyPolicy bool = false

func main() {
	aliens := CheckingArgs(os.Args)
	if aliens <= 0 {
		return
	}
	//ASSUMPTION: map file is always well formed and there are no isolated cities (rows containg just the city name, with no available directions)
	//I'll use map datatype to keep reference on data and because in this use case order is not important.
	cityMapFile, err := ioutil.ReadFile("marvinsplan.txt")
	if err != nil {
		fmt.Println("Can't read the map file.", err)
		return
	}

	cityMap := CreateMap(cityMapFile)

	//I want 2 different behaviours based on the number of
	//the aliens N compared to the number of the city C:
	//1) N <= C
	//Since I don't want to lose cities at t0 in this case,
	//I assume that only one alien could land in each city
	//2) N > C
	//Aliens land wherever they want

	//these are mirrored maps:
	//alientstatus[who]=where
	//invaders[where][]{who}
	status := make(map[int]string)
	invaders := make(map[string][]int)

	numberOfCitiesBeforeSiege := len(cityMap)

	if aliens <= numberOfCitiesBeforeSiege {
		status, invaders = ExclusiveLanding(cityMap, aliens)
	} else {
		status, invaders = Landing(cityMap, aliens)
	}

	//summary of landing event
	fmt.Println(WhereAliens(status))

	if doIwantFightsAtT0 {
		//Destroying cities with more than one alien, killing aliens. Possibly ending the world and kill all the aliens at t0.
		fmt.Println("\t--- CLEANING ---")
		OutStepKillAndDestroy(status, invaders, cityMap, false)
		fmt.Println("\t--- CLEANING ---")
		if len(status) < 1 || (len(cityMap) < 1 && len(status) < 1) {
			fmt.Println("Your cleaning erased the Earth...")
			os.Exit(1)
		}
	}

	for i := 1; i <= steps; i++ {

		for k, v := range status {
			if len(cityMap[v]) == 0 {
				fmt.Printf("Alien %d is trapped in %s!\n", k, v)
				continue
			}
			currentCity := status[k]
			nextCity := NextDestination(currentCity, cityMap)

			status[k] = nextCity
			invaders[currentCity] = RemoveAlien(invaders[currentCity], k)
			invaders[nextCity] = append(invaders[nextCity], k)
			fmt.Printf("Alien %d is going to %s\n", k, status[k])

			if !killAndDestroyAtTheEndOfShift {
				//STRICT MODE: Just two alien and no more are needed to destroy the city
				if (len(invaders[nextCity]) > 1 && !strictTwoAliensDestroyPolicy) || (len(invaders[nextCity]) == 2 && strictTwoAliensDestroyPolicy) {
					InStepKillAndDestroy(nextCity, status, invaders, cityMap)
				}

			}

		}
		if killAndDestroyAtTheEndOfShift {
			OutStepKillAndDestroy(status, invaders, cityMap, strictTwoAliensDestroyPolicy)
		}

		fmt.Printf("\nSnapshot at the END of STEP %d: aliens alive %d of %d; standing cities: %d of %d.\n\n", i, len(status), aliens, len(cityMap), numberOfCitiesBeforeSiege)
		if len(cityMap) < 1 && len(status) >= 1 {
			fmt.Println("Earth died at step ", i)
			break
		}
		if len(cityMap) >= 1 && len(status) < 1 {
			fmt.Println("Aliens went extinct at step ", i)
			break
		}
		if len(cityMap) < 1 && len(status) < 1 {
			fmt.Println("Total destruction at step ", i)
			break
		}
	}

	PrintMap(cityMap)

}

func CheckingArgs(inputArgs []string) int {
	if len(inputArgs) <= 1 {
		fmt.Println("Wrong number of arguments, try again. I need one positive integer number as argument.\ne.g.: >> go run alienivasion 4")
		return 0
	}
	if numberOfAliens, err := strconv.Atoi(inputArgs[1]); err != nil {
		fmt.Printf("Can't comprehend: '%s'\nTry again. I need one positive integer number as argument.\ne.g.: >> go run alienivasion 4", inputArgs[1])
		return 0
	} else {
		if numberOfAliens == 0 {
			fmt.Println("No alien is attacking us. Earth is safe.")
		}
		if numberOfAliens < 0 {
			fmt.Println("Can't have a negative number of aliens.\nI need one positive integer number as argument.\ne.g.: >> go run alienivasion 4")
		}
		return numberOfAliens
	}
}

//Output map: is a map of maps k1=city v1=map[k2:v2], k2=road to city, v2=direction
func CreateMap(inputFileMap []uint8) map[string]map[string]string {
	cityMapStringArray := strings.Split(string(inputFileMap), "\r\n")
	var cityMap = map[string]map[string]string{}
	for i := 0; i < len(cityMapStringArray); i++ {
		cityNameAndBoarders := strings.Split(cityMapStringArray[i], " ")
		cityMap[cityNameAndBoarders[0]] = map[string]string{}
		for j := 1; j < len(cityNameAndBoarders); j++ {
			splitterIndex := strings.Index(cityNameAndBoarders[j], "=")
			cityMap[cityNameAndBoarders[0]][cityNameAndBoarders[j][splitterIndex+1:]] = cityNameAndBoarders[j][:splitterIndex]
		}
	}
	return cityMap
}

//to print out the citymap in the same format of the file
func PrintMap(cityMap map[string]map[string]string) {
	fmt.Print("\n\tMAP:\n")
	for k, v := range cityMap {
		fmt.Printf(k)
		for city, direction := range v {
			fmt.Printf(" %s=%s", direction, city)
		}
		fmt.Printf("\n")
	}
}

//to land max one alien per city
func ExclusiveLanding(worldMap map[string]map[string]string, howManyAliens int) (map[int]string, map[string][]int) {
	tmp := make(map[string]map[string]string)
	for k, v := range worldMap {
		tmp[k] = v
	}

	alienStatus := make(map[int]string)
	invaders := make(map[string][]int)

	//"The iteration order over maps is not specified and is not guaranteed to be the same from one iteration to the next."
	//It is not important for my purpose since I need a "pseudo-random" item
	for i := 1; i <= howManyAliens; i++ {
		howManyCities := len(tmp)
		seed := rand.Intn(howManyCities)
		j := 0
		for k := range tmp {
			if j == seed {
				alienStatus[i] = k
				invaders[k] = append(invaders[k], i)
				delete(tmp, k)
				break
			}
			j++
		}
	}

	//since is not efficient accessing a key using the value i'm "mirroring"
	//these 2 maps:
	//alientstatus[who]=where
	//invaders[where][]{who}
	return alienStatus, invaders

}

//to land aliem randomly in the mapp
func Landing(worldMap map[string]map[string]string, howManyAliens int) (map[int]string, map[string][]int) {

	alienStatus := make(map[int]string)
	invaders := make(map[string][]int)

	for i := 1; i <= howManyAliens; i++ {
		seed := rand.Intn(len(worldMap))
		j := 0
		for k := range worldMap {
			if j == seed {
				alienStatus[i] = k
				invaders[k] = append(invaders[k], i)
				break
			}
			j++
		}
	}

	return alienStatus, invaders

}

//to choose the next alien destination
func NextDestination(city string, currentMapState map[string]map[string]string) string {
	availableDirections := currentMapState[city]
	var nextCity string = city
	//choosing currentMapState[city][0] everytime would have been more efficient, but I don't like this behaviour
	//([]->trapped)
	if len(availableDirections) == 0 {
		return nextCity
	}
	seed := rand.Intn(len(availableDirections))
	j := 0
	for k := range availableDirections {
		if j == seed {
			nextCity = k
		}
		j++
	}
	return nextCity
}

//print alien location
func WhereAliens(alienStatus map[int]string) string {
	//I could have a lot of alien
	b := new(bytes.Buffer)
	for key, value := range alienStatus {
		fmt.Fprintf(b, "Alien %d is in %s.\n", key, value)
	}
	return b.String()
}

//given a city, it will kill its invaders, it will destroy the city and remove it from the map
func InStepKillAndDestroy(city string, alienStatus map[int]string, invaders map[string][]int, cityMap map[string]map[string]string) (map[int]string, map[string][]int, map[string]map[string]string) {
	fighters := invaders[city]
	FightMessagePrinter(city, fighters)

	delete(invaders, city)
	for alienIndex := range fighters {
		delete(alienStatus, fighters[alienIndex])
	}
	delete(cityMap, city)
	for k := range cityMap {
		//if it exists will be deleted
		delete(cityMap[k], city)
	}
	return alienStatus, invaders, cityMap
}

//It looks for cities with more then one invader, then destroy these cities and their invaders
func OutStepKillAndDestroy(alienStatus map[int]string, invaders map[string][]int, cityMap map[string]map[string]string, strictMode bool) (map[int]string, map[string][]int, map[string]map[string]string) {

	for city, aliens := range invaders {
		if (len(aliens) > 1 && !strictMode) || (len(aliens) == 2 && strictMode) {
			FightMessagePrinter(city, aliens)
			for alienIndex := range aliens {
				delete(alienStatus, aliens[alienIndex])
			}
			delete(cityMap, city)
			for k := range cityMap {
				//if it exists will be deleted
				delete(cityMap[k], city)
			}
			delete(invaders, city)
		}
	}
	return alienStatus, invaders, cityMap
}

//remove one alien from the aliens array
func RemoveAlien(aliens []int, alien int) []int {
	aliens[IndexOfAlien(aliens, alien)] = aliens[len(aliens)-1]
	return aliens[:len(aliens)-1]
}

//find the index to remove
func IndexOfAlien(aliens []int, alien int) int {
	for k, v := range aliens {
		if alien == v {
			return k
		}
	}
	return -1 //not found.
}

//to print summary of detruction event
func FightMessagePrinter(city string, fighters []int) {

	fmt.Printf("%s has been destroyed by Alien %d", city, fighters[0])
	for alien := range fighters[1 : len(fighters)-1] {
		fmt.Printf(", Alien %d", fighters[1 : len(fighters)-1][alien])
	}
	fmt.Printf(" and Alien %d\n", fighters[len(fighters)-1])
}
