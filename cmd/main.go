package main

import (
	"fmt"
	"math/rand/v2"
	"slices"
)

type Player struct {
	// The health of the player.
	Health int

	// The probability that the player would miss their shot.
	MissChance int
}

type BeeType string

const (
	QueenBee  BeeType = "Queen Bee"
	WorkerBee BeeType = "Worker Bee"
	DroneBee  BeeType = "Drone Bee"
)

type Bee struct {
	// The type of bee.
	Type BeeType

	// The health of the bee.
	Health int

	// The probability that the bee would miss their shot.
	MissChance int
}

func main() {
	rounds := 0
	hits := 0
	stings := 0
	won := false

	player := Player{
		Health:     100,
		MissChance: 10,
	}

	bees := []Bee{
		// 1 queen bee.
		{QueenBee, 100, 10},

		// 5 worker bees.
		{WorkerBee, 75, 15},
		{WorkerBee, 75, 15},
		{WorkerBee, 75, 15},
		{WorkerBee, 75, 15},
		{WorkerBee, 75, 15},

		// 25 drone bees.
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
		{DroneBee, 60, 20},
	}

	fmt.Println("Destroy all bees...")

GameLoop:
	for {
		rounds += 1

		// Player's turn. Get a random bee from the list of bees.
		beeIndex := rand.IntN(len(bees))

		// Check if player actually manages to hit the bee.
		if rand.IntN(101) > player.MissChance {
			hits += 1

			// Determine the kind of bee the player hit and take the appropriate steps.
			switch bees[beeIndex].Type {
			case QueenBee:
				bees[beeIndex].Health -= 10
				if bees[beeIndex].Health <= 0 {
					msg := "You killed the Queen bee."
					fmt.Println(msg)
					won = true

					break GameLoop
				} else {
					msg := fmt.Sprintf("Direct Hit! Queen took 10 hit points. %d HP left.", bees[beeIndex].Health)
					fmt.Println(msg)
				}
			case WorkerBee:
				bees[beeIndex].Health -= 25
				if bees[beeIndex].Health <= 0 {
					bees = slices.Delete(bees, beeIndex, beeIndex+1)

					msg := fmt.Sprintf("You just took out one of the worker bees. %d bees left.", len(bees))
					fmt.Println(msg)
				} else {
					msg := fmt.Sprintf("Direct Hit! Worker took 25 hit points. %d HP left.", bees[beeIndex].Health)
					fmt.Println(msg)
				}
			case DroneBee:
				bees[beeIndex].Health -= 30
				if bees[beeIndex].Health <= 0 {
					bees = slices.Delete(bees, beeIndex, beeIndex+1)

					msg := fmt.Sprintf("You just took out one of the drone bees. %d bees left.", len(bees))
					fmt.Println(msg)
				} else {
					msg := fmt.Sprintf("Direct Hit! Drone took 30 hit points. %d HP left.", bees[beeIndex].Health)
					fmt.Println(msg)
				}
			}
		} else {
			msg := "Miss! You just missed the hive, better luck next time!"
			fmt.Println(msg)
		}

		// Hive's turn. Choose a random bee that will attack the player.
		beeIndex = rand.IntN(len(bees))

		// Check if bee actually manages to sting the player.
		if rand.IntN(101) > bees[beeIndex].MissChance {
			stings += 1

			// Determine the kind of bee that just stung the player.
			switch bees[beeIndex].Type {
			case QueenBee:
				player.Health -= 10
				if player.Health <= 0 {
					msg := "The Queen Bee just killed you!"
					fmt.Println(msg)
					break GameLoop
				} else {
					msg := fmt.Sprintf("Sting! You just got stun by the queen bee. You have %d HP left.", player.Health)
					fmt.Println(msg)
				}
			case WorkerBee:
				player.Health -= 5
				if player.Health <= 0 {
					msg := "A worker bee just killed you!"
					fmt.Println(msg)
					break GameLoop
				} else {
					msg := fmt.Sprintf("Sting! You just got stun by a worker bee. You have %d HP left.", player.Health)
					fmt.Println(msg)
				}
			case DroneBee:
				player.Health -= 1
				if player.Health <= 0 {
					msg := "A drone bee just killed you!"
					fmt.Println(msg)
					break GameLoop
				} else {
					msg := fmt.Sprintf("Sting! You just got stun by a drone bee. You have %d HP left.", player.Health)
					fmt.Println(msg)
				}
			}
		} else {
			msg := fmt.Sprintf("Buzz! That was close! The %s just missed you!", bees[beeIndex].Type)
			fmt.Println(msg)
		}
	}

	// Print summary information about the game.

	fmt.Println("\nGame summary:")

	if won {
		fmt.Println("You won the game!")
	} else {
		fmt.Println("You lost the game.")
	}

	fmt.Printf("You played for %d rounds\n", rounds)
	fmt.Printf("You were stung %d times\n", stings)
	fmt.Printf("You managed to hit the hive %d times\n", hits)
}
