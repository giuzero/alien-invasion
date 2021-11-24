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

## Destruction behaviour modifiers
User can decide how destruction happens
```
//true: will "clean" the map just after landing, destroying cities invaded by two (strict = true)
//or more (strict = false) aliens. After t0 you'll have at most one alien per city. Will they still exist?
var doIwantFightsAtT0 bool = false

//true: enable kills at the end of the turn/step.
//false: when an alien arrives in a city, could find one or more aliens already there.
//Will kill an destroy during the step
var killAndDestroyAtTheEndOfShift bool = false

//destroy and kill aliens if and only if the city is invaded by exactly 2 aliens
var strictTwoAliensDestroyPolicy bool = false
```