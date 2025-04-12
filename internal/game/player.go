package game

import "fmt"

const (
	damageFromQueen  = 10
	damageFromWorker = 5
	damageFromDrone  = 1
)

// Player represents the player character in the game.
// It tracks the player's health and the probability that their attack will miss.
type Player struct {
	Health     int
	MissChance uint
}

// createPlayer initializes a new Player instance with default or custom miss chance.
// If no miss chance is provided, a default value of 10% is used.
func createPlayer(missChance ...uint) Player {
	var miss uint

	if len(missChance) > 0 {
		miss = missChance[0]
	} else {
		miss = 10
	}

	return Player{
		Health:     100,
		MissChance: miss,
	}
}

// takeDamage applies damage to the player based on the type of bee attacking.
// Returns true if the player's health drops to zero or below (i.e., the player dies).
func (player *Player) takeDamage(beeType BeeType) bool {
	damage := 0

	switch beeType {
	case QueenBee:
		damage = damageFromQueen
	case WorkerBee:
		damage = damageFromWorker
	case DroneBee:
		damage = damageFromDrone
	}

	player.Health -= damage

	return player.Health <= 0
}

// generateHitMessage returns a descriptive string based on the type of bee that attacked
// the player and whether the attack was fatal.
func (player *Player) generateHitMessage(beeType BeeType, died bool) string {
	msg := ""

	switch {
	case beeType == QueenBee && died:
		msg = "The Queen bee just killed you!"
	case beeType == QueenBee && !died:
		msg = fmt.Sprintf("Sting! You just got stun by the Queen bee. You have %d HP left.", player.Health)
	case beeType == WorkerBee && died:
		msg = "A worker bee just killed you!"
	case beeType == WorkerBee && !died:
		msg = fmt.Sprintf("Sting! You just got stun by a worker bee. You have %d HP left.", player.Health)
	case beeType == DroneBee && died:
		msg = "A drone bee just killed you!"
	case beeType == DroneBee && !died:
		msg = fmt.Sprintf("Sting! You just got stun by a drone bee. You have %d HP left.", player.Health)
	}

	return msg
}
