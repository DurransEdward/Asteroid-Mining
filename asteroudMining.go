package main

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

const numberAsteroids int = 10

func main() {
	// /\/\/\/\/\/\/\/\/\/\ Creating arrays for the starting coordinates and randomly generating asteroid coordinates /\/\/\/\/\/\/\/\/\/\

	startingCoordinates := [3]float64{0, 0, 0}

	var asteroidCoordinates [numberAsteroids + 1][3]float64 // The first row is all zeros (to represent the starting location)

	/* This for loop will randomly generate the asteroid coordinates */

	for a := 1; a <= numberAsteroids; a++ {
		for i := 0; i < 3; i++ {
			asteroidCoordinates[a][i] = rand.Float64()*2000 - 1000
		}
	}

	// /\/\/\/\/\/\/\/\/\/\ Insert starting coordinates and asteroid coordinates into SQL database /\/\/\/\/\/\/\/\/\/\

	var db *sql.DB // Declare pointer to database object
	var err error

	db, err = sql.Open("postgres", "postgres://postgres:fvcszfs-22rcaXX@localhost:5432/test?sslmode=disable") // Connect db to the required database
	if err != nil {
		fmt.Println("Something is broken... 1")
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO startingCoordinatesAndAsteroidCoordinates (id, x, y, z) VALUES (0,0,0,0);")
	if err != nil {
		fmt.Println("Something is broken... 2")
	}

	var stmt *sql.Stmt

	stmt, err = db.Prepare("INSERT INTO startingCoordinatesAndAsteroidCoordinates(id, x, y, z) VALUES($1, $2, $3, $4)")
	if err != nil {
		fmt.Println("Something is broken... 3")
	}

	for a := 1; a <= 10; a++ {
		_, err = stmt.Exec(a, asteroidCoordinates[a][0], asteroidCoordinates[a][1], asteroidCoordinates[a][2])
		if err != nil {
			fmt.Println("Something is broken... 4", err)
		}
	}

	// /\/\/\/\/\/\/\/\/\/\ Insert starting coordinates and asteroid coordinates into Redis database /\/\/\/\/\/\/\/\/\/\

	var ctx = context.Background()
	var rdb *redis.Client

	var connectionOptions = redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	rdb = redis.NewClient(&connectionOptions)

	rdb.Set(ctx, "starting_x_coordinate", 0, 0) // Enter starting coordinates into Redis database
	rdb.Set(ctx, "starting_y_coordinate", 0, 0) // Enter starting coordinates into Redis database
	rdb.Set(ctx, "starting_z_coordinate", 0, 0) // Enter starting coordinates into Redis database

	var keyString string
	
	for a := 1; a <= numberAsteroids; a++ {
		for i := 0; i <= 2; i++{
			switch i {
				case 1: keyString = "asteroid_" + strconv.Itoa(a) + "_x_coordinate"
				
				case 2: keyString = "asteroid_" + strconv.Itoa(a) + "_y_coordinate"

				case 3: keyString = "asteroid_" + strconv.Itoa(a) + "_z_coordinate"
			}

			rdb.Set(ctx, keyString, asteroidCoordinates[a][i], 0)
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
asteroid visited by the spaceship, then pathThroughAsteroids[j] = i. Since the spaceship starts at location 0, we set pathThroughAsteroids[0] = 0 by default.
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
