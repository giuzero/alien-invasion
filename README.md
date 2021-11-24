# Alien Invasion
## Start the application
To start the application just insert the number of aliens that will invade the earth.
```
go run .\alieninvasion.go 4
```
## Landing behaviour
Being N the number of aliens and C the number of cities on the map, there are 2 type of landing:

>N<=C

In each city, at most one alien will land

>N>C

Random landing

## Map file example

In the same folder of ```alieninvasion.go``` there is the map file ```map_file.txt```.

This is a map example:

```
Seul south=Busan
Busan north=Seul south=Hiroshima east=Kyoto
Hiroshima north=Busan east=Osaka
Osaka north=Kyoto south=Kobe west=Hiroshima east=Nagoya
Kyoto south=Osaka west=Busan east=Kanazawa
Kobe north=Osaka
Kanazawa south=Nagoya west=Kyoto east=Tokyo
Nagoya north=Kanazawa west=Osaka east=Yokohama
Tokyo north=Nigata south=Yokohama west=Kanazawa
Nigata south=Tokyo
Yokohama north=Tokyo west=Nagoya
```

## Destruction behaviour modifiers

User can decide how destruction happens, setting global variables
```
//Destroy policy: destroy and kill aliens if and only if the city is invaded by exactly 2 aliens
var strictTwoAliensDestroyPolicy bool = false

//true: will "clean" the map just after landing, destroying cities invaded by two (strict = true)
//or more (strict = false) aliens. After t0 you'll have at most one alien per city. Will they still exist?
//Does not care about Destroy policy
var doIwantFightsAtT0 bool = false

//true: enable kills and destruction only at the end of the turn/step so after all the aliens that 
//survived the previous step move, if not trapped alone and if killing policy is respected.
//false: when an alien arrives in a city, could find one or more aliens already there.
//Would kill and destroy during the step, as soon as the alien arrives in the city. 
//If aliens don't move, don't kill. Must respect the destroy policy
var killAndDestroyAtTheEndOfStep bool = false
```

## TO DO
Autokill if trapped more than X steps
