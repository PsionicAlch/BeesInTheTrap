package game

import "testing"

func TestCreatePlayer(t *testing.T) {
	missChance := 35

	playerWithMissChance := createPlayer(uint(missChance))

	if playerWithMissChance.Health != 100 {
		t.Errorf("createPlayer generated player with incorrect health. Expected: %d. Received: %d.\n", 100, playerWithMissChance.Health)
	}

	if playerWithMissChance.MissChance != uint(missChance) {
		t.Errorf("createPlayer generated player with incorrect miss chance. Expected: %d. Received: %d.\n", missChance, playerWithMissChance.MissChance)
	}

	missChance = 10

	playerWithoutMissChance := createPlayer()

	if playerWithoutMissChance.Health != 100 {
		t.Errorf("createPlayer generated player with incorrect health. Expected: %d. Received: %d.\n", 100, playerWithoutMissChance.Health)
	}

	if playerWithoutMissChance.MissChance != uint(missChance) {
		t.Errorf("createPlayer generated player with incorrect miss chance. Expected: %d. Received: %d.\n", missChance, playerWithoutMissChance.MissChance)
	}
}

func TestPlayerTakeDamage(t *testing.T) {
	scenarios := []struct {
		playerHealth   int
		expectedHealth int
		beeType        BeeType
		shouldHaveDied bool
	}{
		{damageFromQueen + 5, 5, QueenBee, false},
		{damageFromQueen - 5, -5, QueenBee, true},
		{damageFromWorker + 5, 5, WorkerBee, false},
		{damageFromWorker - 5, -5, WorkerBee, true},
		{damageFromDrone + 5, 5, DroneBee, false},
		{damageFromDrone - 5, -5, DroneBee, true},
	}

	for _, scenario := range scenarios {
		player := &Player{
			Health: scenario.playerHealth,
		}

		died := player.takeDamage(scenario.beeType)

		if player.Health != scenario.expectedHealth {
			t.Errorf("Player's health is different from what was expected after an attack. Expected health: %d. Current health: %d. Bee type: %s\n", scenario.expectedHealth, player.Health, scenario.beeType)
		}

		if died && !scenario.shouldHaveDied {
			t.Errorf("Player died when they should have lived. Bee type: %s\n", scenario.beeType)
		}

		if !died && scenario.shouldHaveDied {
			t.Errorf("Player didn't die when they should have died. Bee type: %s\n", scenario.beeType)
		}
	}
}

func TestPlayerGenerateHitMessage(t *testing.T) {
	scenarios := []struct {
		playerHealth int
		beeType      BeeType
		died         bool
		expectedMsg  string
	}{
		{0, QueenBee, true, "The Queen bee just killed you!"},
		{5, QueenBee, false, "Sting! You just got stun by the Queen bee. You have 5 HP left."},
		{0, WorkerBee, true, "A worker bee just killed you!"},
		{5, WorkerBee, false, "Sting! You just got stun by a worker bee. You have 5 HP left."},
		{0, DroneBee, true, "A drone bee just killed you!"},
		{5, DroneBee, false, "Sting! You just got stun by a drone bee. You have 5 HP left."},
		{5, BeeType(5), false, ""},
	}

	for _, scenario := range scenarios {
		player := &Player{
			Health: scenario.playerHealth,
		}

		msg := player.generateHitMessage(scenario.beeType, scenario.died)

		if msg != scenario.expectedMsg {
			t.Errorf("Unexpected message generated from generateHitMessage function. Expected: \"%s\". Generated: \"%s\"\n", scenario.expectedMsg, msg)
		}
	}
}
