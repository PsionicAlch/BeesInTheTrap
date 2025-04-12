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

	playerWithoutMissChance := createPlayer(uint(missChance))

	if playerWithoutMissChance.Health != 100 {
		t.Errorf("createPlayer generated player with incorrect health. Expected: %d. Received: %d.\n", 100, playerWithoutMissChance.Health)
	}

	if playerWithoutMissChance.MissChance != uint(missChance) {
		t.Errorf("createPlayer generated player with incorrect miss chance. Expected: %d. Received: %d.\n", missChance, playerWithoutMissChance.MissChance)
	}
}
