package game

import (
	"fmt"
	"math/rand/v2"
	"slices"
)

// GameState represents the current status of a running game.
// It holds information about the player, the hive, and turn statistics.
type GameState struct {
	Player Player
	Hive   []Bee
	Round  uint
	Hits   uint
	Stings uint
}

// GameServer manages the lifecycle of the game.
// It coordinates turns, updates state, and communicates with the client via a protocol.
type GameServer struct {
	finished      bool
	state         GameState
	communication Protocol
}

// StartupServer initializes the game server and returns a communication channel for the client.
// This function also spawns a goroutine to run the main game loop.
func StartupServer() *CommunicationProtocol {
	communication := createCommunicationProtocol()

	server := &GameServer{
		finished: false,
		state: GameState{
			Round:  0,
			Hits:   0,
			Stings: 0,
			Player: createPlayer(),
			Hive:   createHive(1, 5, 25),
		},
		communication: communication,
	}

	go server.run()

	return communication
}

// run starts the main game loop, alternating turns between the player and the hive.
func (server *GameServer) run() {
	for {
		server.state.Round += 1

		server.playersTurn()

		if server.finished {
			return
		}

		server.hivesTurn()

		if server.finished {
			return
		}
	}
}

// playersTurn handles the player's action phase.
// It waits for input, applies damage to a random bee, and checks for win conditions.
func (server *GameServer) playersTurn() {
	// Wait for the player's input.
	server.communication.WaitForPlayer()

	// Get a reference to a random bee from the hive.
	beeIndex := rand.IntN(len(server.state.Hive))
	selectedBee := &server.state.Hive[beeIndex]

	// Check to see if player misses their shot.
	if rand.UintN(101) < server.state.Player.MissChance {
		server.communication.HitResponse("Miss! You just missed the hive, better luck next time!", server.state)
		return
	}

	// Player managed to successfully hit a bee. Increment counter.
	server.state.Hits += 1

	// Determine the damage the bee took and generate a witty message.
	beeDied := selectedBee.takeDamage()
	hitMsg := selectedBee.generateHitMessage(beeDied)

	// If the queen died the player won.
	if beeDied && selectedBee.Type == QueenBee {
		server.finished = true

		server.communication.GameFinishedResponse(hitMsg, server.state)
		return
	}

	// If the bee that was just hit died (and wasn't the queen) we remove it from the hive.
	if beeDied {
		server.state.Hive = slices.Delete(server.state.Hive, beeIndex, beeIndex+1)
	}

	server.communication.HitResponse(hitMsg, server.state)
}

// hivesTurn handles the hive's action phase.
// A random bee attempts to sting the player. Death or miss is resolved accordingly.
func (server *GameServer) hivesTurn() {
	// Select a random bee from the hive.
	beeIndex := rand.IntN(len(server.state.Hive))
	selectedBee := &server.state.Hive[beeIndex]
	player := &server.state.Player

	// Check to see if the bee misses their shot.
	if rand.UintN(101) <= selectedBee.MissChance {
		msg := fmt.Sprintf("Buzz! That was close! The %s just missed you!", selectedBee.Type)
		server.communication.StingResponse(msg, server.state)
		return
	}

	// Bee managed to successfully sting the player. Increment counter.
	server.state.Stings += 1

	// Determine the damage dealt to the player and generate a witty message.
	playerDied := player.takeDamage(selectedBee.Type)
	stingMsg := player.generateHitMessage(selectedBee.Type, playerDied)

	// If the player died the game is over.
	if playerDied {
		server.finished = true

		server.communication.GameFinishedResponse(stingMsg, server.state)
		return
	}

	server.communication.StingResponse(stingMsg, server.state)
}
