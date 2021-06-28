package main

import (
	"fmt"
	"math"
	"math/rand"
)

const numberAsteroids int = 10

func main() {
	startingCoordinates := [3]float64{0, 0, 0}	

	var asteroidCoordinates [numberAsteroids + 1][3]float64 // The first row is all zeros (to represent the starting location)
	
	/* This for loop will randomly generate the asteroid coordinates */
	
	for i := 1; i <= numberAsteroids; i++ {
		for j := 0; j < 3; j++ {
			asteroidCoordinates[i][j] = rand.Float64()*2000 - 1000
		}
	}

	fmt.Println("Asterod coordinates:", "\n\n", asteroidCoordinates, "\n\n")

	pathThroughAsteroids := pathFinderMethod1(startingCoordinates, asteroidCoordinates)

	fmt.Println("Path through asteroids:", pathThroughAsteroids)
}

/*
pathFinderMethod1:

Function Arguments:

1. The starting coordinates of the spaceship
2. The coordinates of the asteroids

Function Returns:

1. An array variable called pathThroughAsteroids which records the path the spaceship takes through the asteroids. For example, if asteroid i is the jth
asteroid visited by the spaceship, then pathThroughAsteroids[i] = j. Since the spaceship starts at location 0, we set pathThroughAsteroids[0] = 0 by default.
The method used to find this path is described in section 1 of the Self-Assigned Go/SQL/Github Asteroid Mining Task Description document.
*/
func pathFinderMethod1(startingCoordinates [3]float64, asteroidCoordinates [numberAsteroids + 1][3]float64) [numberAsteroids + 1]int {
	currentLocation := 0 // Start at starting location which has id = 0

	var visitedAsteroids [numberAsteroids + 1]bool /* The entries of this variable array will tell us if the corresponding asteroid has been visited yet.
	false => Asteroud not yet visited, true => asteroid previously visited. Note that, by default, Go will initialse this variable to be full of false */

	visitedAsteroids[0] = true // The starting location is 0 which the spaceship visits immediately

	var pathThroughAsteroids [numberAsteroids + 1]int // Variable explained in function description

	pathThroughAsteroids[0] = 0 // The spaceship's initial location is the starting coordinates, which have id = 0

	for a := 1; a <= numberAsteroids; a++ {
		currentLocation = goToClosestUnvisitedAsteroid(startingCoordinates, currentLocation, asteroidCoordinates, &visitedAsteroids)

		pathThroughAsteroids[a] = currentLocation

		visitedAsteroids[currentLocation] = true
	}

	return pathThroughAsteroids
}

/*
goToClosestUnvisitedAsteroid:

Function Arguments:

1. The starting coordinates of the spaceship
2. The coordinates of the asteroids
3. A pointer to an array recording which asteroids have already been visited

Function Returns:

1. The id of the asteroid the spaceship should travel to next
*/
func goToClosestUnvisitedAsteroid(startingCoordinates [3]float64, currentLocation int, asteroidCoordinates [numberAsteroids + 1][3]float64, visitedAsteroids *[numberAsteroids + 1]bool) int {
	currentBestNextAsteroid := -1          // Initialise to -1 (which isn't an asteroid) as an error detector
	currentBestDistance := math.MaxFloat64 // Ridicuously big initial current best distance which will be overridden by the first calculated distance
	var distance float64

	if currentLocation == 0 {
		for a := 1; a <= numberAsteroids; a++ {
			distance = pythagThm3D(startingCoordinates[0]-asteroidCoordinates[a][0], startingCoordinates[1]-asteroidCoordinates[a][1], startingCoordinates[2]-asteroidCoordinates[a][2])
			if distance < currentBestDistance {
				currentBestNextAsteroid = a
				currentBestDistance = distance
			}
		}
	} else {
		for a := 1; a <= numberAsteroids; a++ {
			if visitedAsteroids[a] {
			} else {
				distance = pythagThm3D(asteroidCoordinates[currentLocation][0]-asteroidCoordinates[a][0], asteroidCoordinates[currentLocation][1]-asteroidCoordinates[a][1], asteroidCoordinates[currentLocation][2]-asteroidCoordinates[a][2])
				if distance < currentBestDistance {
					currentBestNextAsteroid = a
					currentBestDistance = distance
				}
			}
		}
	}

	return currentBestNextAsteroid
}

/*
pythagThm3d

Function Arguments:

1. The coordinates of a point in three dimensional Euclidian space

Function Returns:

1. The distance from the origin to the point
*/
func pythagThm3D(x, y, z float64) float64 {
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2) + math.Pow(z, 2))
}
































/*
FUTURE DEVELOPMENT. DO NOT READ!

bruteForceShortestPathFinder:

Function Arguments:

1. The starting coordinates of the spaceship
2. The coordinates of the asteroids

Function Returns:

1. An array variable called pathThroughAsteroids which records the path the spaceship takes through the asteroids. For example, if asteroid i is the
jth asteroid visited by the spaceship, then pathThroughAsteroids[i] = j. Since the spaceship starts at location 0, we set pathThroughAsteroids[0] = 0 by default.
The method used to find this path is brute force. I.e. checking every possible path and finding the shortest.
*/

/*
func bruteForceShortestPathFinder(startingCoordinates [3]float64, asteroidCoordinates [numberAsteroids + 1][3]float64) [numberAsteroids + 1]int {



}
*/
