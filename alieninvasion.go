package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var steps int = 10000                          //number of steps
var doIwantFightsAtT0 bool = false             //setting it true, will "clean" the map, after t0 you'll have at most one alien per city. Will they still exist?
var killAndDestroyAtTheEndOfShift bool = false //true: enable multiple kills/
//false: when an alien arrives in a city, could find just another alien there. let's call it "precise mode"

func main() {
	aliens := checkingArgs(os.Args)
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

	cityMap := createMap(cityMapFile)

	//There is no specification about the landing so I want 2 different behaviours based on the number
	//the aliens N compared to the number of the city C:
	//1) N <= C
	//Since I don't want to lose cities at t0 in this case, I assume that only one alien could land in each city
	//(I would use always this kind of behaviour, if the fight consisted of only 2 aliens slap each other was a hard requirement)
	//2) N > C
	//Aliens land wherever they want

	//these are mirrored maps:
	//alientstatus[who]=where
	//invaders[where][]{who}
	status := make(map[int]string)
	invaders := make(map[string][]int)

	numberOfCitiesBeforeSiege := len(cityMap)

	if aliens <= numberOfCitiesBeforeSiege {
		status, invaders = exclusiveLanding(cityMap, aliens)
	} else {
		status, invaders = landing(cityMap, aliens)
	}

	fmt.Println(whereAliens(status))

	if doIwantFightsAtT0 {

		//Destroying cities with more than one alien, killing aliens. Possibly ending the world and kill all the aliens at t0.
		multipleKillAndDestroy(status, invaders, cityMap)
	}

	for i := 1; i <= steps; i++ {

		for k, v := range status {
			if len(cityMap[v]) == 0 {
				fmt.Printf("Alien %d is trapped in %s!\n", k, v)
				continue
			}
			currentCity := status[k]
			nextCity := nextDestination(currentCity, cityMap)

			status[k] = nextCity
			invaders[currentCity] = remove(invaders[currentCity], k)
			invaders[nextCity] = append(invaders[nextCity], k)
			fmt.Printf("Alien %d is going to %s\n", k, status[k])

			if !killAndDestroyAtTheEndOfShift {
				//Just two alien and no more are needed to destroy the city
				if len(invaders[nextCity]) > 1 {
					if len(invaders[nextCity]) > 2 {
						errors.New("this should be just 2")
					}
					preciseKillAndDestroy(nextCity, status, invaders, cityMap)
				}
			}

		}
		if killAndDestroyAtTheEndOfShift {
			multipleKillAndDestroy(status, invaders, cityMap)
		}

		fmt.Printf("END STEP %d: aliens alive %d of %d; standing cities: %d of %d.\n\n", i, len(status), aliens, len(cityMap), numberOfCitiesBeforeSiege)
		if len(status) < 1 || len(cityMap) < 1 {
			break
		}
	}
	fmt.Print(cityMap)
	printMap(cityMap)

}

func checkingArgs(inputArgs []string) int {
	if len(inputArgs) <= 1 {
		fmt.Println("Wrong number of arguments, try again. I need one integer number as argument.\ne.g.: >> go run alienivasion 4")
		return 0
	}
	if numberOfAliens, err := strconv.Atoi(inputArgs[1]); err != nil {
		fmt.Printf("Can't comprehend: '%s'\nTry again. I need one integer number as argument.\ne.g.: >> go run alienivasion 4", inputArgs[1])
		return 0
	} else {
		if numberOfAliens == 0 {
			fmt.Println("No alien is attacking us. Earth is safe.")
		}
		if numberOfAliens < 0 {
			fmt.Println("Can't have a negative number of aliens.\nI need one integer number as argument.\ne.g.: >> go run alienivasion 4")
		}
		return numberOfAliens
	}
}

//Output map: is a map of maps k1=city v1=map[k2:v2], k2=road to city, v2=direction
func createMap(inputFileMap []uint8) map[string]map[string]string {
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
func printMap(cityMap map[string]map[string]string) {
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
func exclusiveLanding(worldMap map[string]map[string]string, howManyAliens int) (map[int]string, map[string][]int) {
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
func landing(worldMap map[string]map[string]string, howManyAliens int) (map[int]string, map[string][]int) {

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
func nextDestination(city string, currentMapState map[string]map[string]string) string {
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
func whereAliens(alienStatus map[int]string) string {
	//I could have a lot of alien
	b := new(bytes.Buffer)
	for key, value := range alienStatus {
		fmt.Fprintf(b, "Alien %d is in %s.\n", key, value)
	}
	return b.String()
}

//given a city, it will kill its invaders, it will destroy the city and remove it from the map
//to be used in "precise mode"
func preciseKillAndDestroy(city string, alienStatus map[int]string, invaders map[string][]int, cityMap map[string]map[string]string) {
	fighters := invaders[city]
	if len(fighters) > 2 {
		fmt.Println("Multiple kills disabled. Set global variables\n\tdoIwantFightsAtT0\n and/or \n\tkillAndDestroyAtTheEndOfShift \nto true.")
		os.Exit(2)
	}
	fmt.Printf("%s has been destroyed by Alien %d and Alien %d!", city, fighters[0], fighters[1])
	delete(invaders, city)
	delete(alienStatus, fighters[0])
	delete(alienStatus, fighters[1])
	delete(cityMap, city)
	for k := range cityMap {
		//if it exists will be deleted
		delete(cityMap[k], city)
	}

}

//It look for cities with more then one invader, then destroy these cities and their invaders
func multipleKillAndDestroy(alienStatus map[int]string, invaders map[string][]int, cityMap map[string]map[string]string) {

	for city, aliens := range invaders {
		if len(aliens) > 1 {
			multipleKillMessage(city, aliens)
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
}

func remove(aliens []int, alien int) []int {
	aliens[indexOfAlien(aliens, alien)] = aliens[len(aliens)-1]
	return aliens[:len(aliens)-1]
}

func indexOfAlien(aliens []int, alien int) int {
	for k, v := range aliens {
		if alien == v {
			return k
		}
	}
	return -1 //not found.
}

func multipleKillMessage(city string, fighters []int) {
	fmt.Printf("DEBUG: ", fighters)
	fmt.Printf("Alien %d", fighters[0])
	for alien := range fighters[1 : len(fighters)-1] {
		fmt.Printf(", Alien %d ", fighters[1 : len(fighters)-1][alien])
	}
	fmt.Printf("and Alien %d destroyed %s", fighters[len(fighters)-1], city)
}
