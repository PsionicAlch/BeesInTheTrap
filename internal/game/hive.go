package game

import "fmt"

// BeeType represents the category of a bee in the hive.
type BeeType uint

const (
	QueenBee BeeType = iota
	WorkerBee
	DroneBee
)

// String returns the string representation of a BeeType.
func (beeType BeeType) String() string {
	switch beeType {
	case QueenBee:
		return "Queen bee"
	case WorkerBee:
		return "worker bee"
	case DroneBee:
		return "drone bee"
	default:
		return ""
	}
}

const (
	damageAgainstQueen  = 10
	damageAgainstWorker = 25
	damageAgainstDrone  = 30
)

// Bee represents an individual bee in the hive, including its type,
// remaining health, and chance to miss an attack.
type Bee struct {
	Type       BeeType
	Health     int
	MissChance uint
}

// takeDamage applies damage to the bee based on its type.
// It returns true if the bee's health drops to zero or below (i.e., the bee is dead).
func (bee *Bee) takeDamage() bool {
	damage := 0

	switch bee.Type {
	case QueenBee:
		damage = damageAgainstQueen
	case WorkerBee:
		damage = damageAgainstWorker
	case DroneBee:
		damage = damageAgainstDrone
	}

	bee.Health -= damage

	return bee.Health <= 0
}

// generateHitMessage returns a descriptive string about the result
// of a hit on the bee, including whether it was fatal.
func (bee *Bee) generateHitMessage(died bool) string {
	msg := ""

	switch {
	case bee.Type == QueenBee && died:
		msg = "You killed the Queen bee."
	case bee.Type == QueenBee && !died:
		msg = fmt.Sprintf("Direct Hit! Queen took %d hit points. %d HP left.", damageAgainstQueen, bee.Health)
	case bee.Type == WorkerBee && died:
		msg = "You killed a worker bee."
	case bee.Type == WorkerBee && !died:
		msg = fmt.Sprintf("Direct Hit! Worker took %d hit points. %d HP left.", damageAgainstWorker, bee.Health)
	case bee.Type == DroneBee && died:
		msg = "You killed a drone bee."
	case bee.Type == DroneBee && !died:
		msg = fmt.Sprintf("Direct Hit! Drone took %d hit points. %d HP left.", damageAgainstDrone, bee.Health)
	}

	return msg
}

// createBees returns a slice of Bee instances of the given type and quantity.
// Optionally, a custom miss chance can be provided.
func createBees(beeType BeeType, num uint, missChance ...uint) []Bee {
	var health int
	var miss uint

	switch beeType {
	case QueenBee:
		health = 100
		miss = 10
	case WorkerBee:
		health = 75
		miss = 15
	case DroneBee:
		health = 60
		miss = 20
	default:
		return nil
	}

	if len(missChance) > 0 {
		miss = missChance[0]
	}

	bees := make([]Bee, num)
	for index := range bees {
		bees[index] = Bee{
			Type:       beeType,
			Health:     health,
			MissChance: miss,
		}
	}

	return bees
}

// createHive returns a full hive composed of queens, workers, and drones.
// It constructs the hive using the given number of each type of bee.
func createHive(numQueens, numWorkers, numDrones uint) []Bee {
	queens := createBees(QueenBee, numQueens)
	workers := createBees(WorkerBee, numWorkers)
	drones := createBees(DroneBee, numDrones)

	bees := make([]Bee, 0, numQueens+numWorkers+numDrones)
	bees = append(bees, queens...)
	bees = append(bees, workers...)
	bees = append(bees, drones...)

	return bees
}
