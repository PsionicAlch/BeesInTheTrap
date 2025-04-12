package game

import (
	"fmt"
	"testing"
)

func TestBeeTypeToString(t *testing.T) {
	scenarios := []struct {
		beeType     BeeType
		expectedStr string
	}{
		{QueenBee, "Queen bee"},
		{WorkerBee, "worker bee"},
		{DroneBee, "drone bee"},
		{BeeType(4), ""},
	}

	for _, scenario := range scenarios {
		receivedStr := scenario.beeType.String()

		if scenario.expectedStr != receivedStr {
			t.Errorf("Unexpected string received. Expected: \"%s\". Received: \"%s\"\n", scenario.expectedStr, receivedStr)
		}
	}
}

func TestTakeDamage(t *testing.T) {
	var died bool

	queenBee := &Bee{
		Type:   QueenBee,
		Health: 2 * damageAgainstQueen,
	}

	died = queenBee.takeDamage()

	if died {
		t.Errorf("Queen bee died when she should not have. Her health: %d\n", queenBee.Health)
	}

	died = queenBee.takeDamage()

	if !died {
		t.Errorf("Queen bee didn't die when she should have. Her health: %d\n", queenBee.Health)
	}

	workerBee := &Bee{
		Type:   WorkerBee,
		Health: 2 * damageAgainstWorker,
	}

	died = workerBee.takeDamage()

	if died {
		t.Errorf("Worker bee died when she should not have. Her health: %d\n", workerBee.Health)
	}

	died = workerBee.takeDamage()

	if !died {
		t.Errorf("Worker bee didn't die when she should have. Her health: %d\n", workerBee.Health)
	}

	droneBee := &Bee{
		Type:   DroneBee,
		Health: 2 * damageAgainstDrone,
	}

	died = droneBee.takeDamage()

	if died {
		t.Errorf("Drone bee died when she should not have. Her health: %d\n", droneBee.Health)
	}

	died = droneBee.takeDamage()

	if !died {
		t.Errorf("Drone bee didn't die when she should have. Her health: %d\n", droneBee.Health)
	}
}

func TestGenerateHitMessage(t *testing.T) {
	queenBee := Bee{
		Type:   QueenBee,
		Health: 100,
	}

	workerBee := Bee{
		Type:   WorkerBee,
		Health: 75,
	}

	droneBee := Bee{
		Type:   DroneBee,
		Health: 60,
	}

	unknownBee := Bee{
		Type: BeeType(3),
	}

	scenarios := []struct {
		bee  Bee
		msg  string
		died bool
	}{
		{queenBee, "You killed the Queen bee.", true},
		{queenBee, fmt.Sprintf("Direct Hit! Queen took %d hit points. %d HP left.", damageAgainstQueen, queenBee.Health), false},
		{workerBee, "You killed a worker bee.", true},
		{workerBee, fmt.Sprintf("Direct Hit! Worker took %d hit points. %d HP left.", damageAgainstWorker, workerBee.Health), false},
		{droneBee, "You killed a drone bee.", true},
		{droneBee, fmt.Sprintf("Direct Hit! Drone took %d hit points. %d HP left.", damageAgainstDrone, droneBee.Health), false},
		{unknownBee, "", true},
		{unknownBee, "", false},
	}

	for _, scenario := range scenarios {
		testBee := scenario.bee
		msg := (&testBee).generateHitMessage(scenario.died)

		if msg != scenario.msg {
			t.Errorf("Generated hit message does not match expected hit message.\n \tExpected: \"%s\"\n\tGot: \"%s\"\n", scenario.msg, msg)
		}
	}
}

func TestCreateBees(t *testing.T) {
	scenarios := []struct {
		beeType            BeeType
		count              uint
		withMissChance     bool
		expectedHealth     int
		expectedMissChance uint
	}{
		{QueenBee, 5, true, 100, 35},
		{QueenBee, 5, false, 100, 10},
		{WorkerBee, 5, true, 75, 35},
		{WorkerBee, 5, false, 75, 15},
		{DroneBee, 5, true, 60, 35},
		{DroneBee, 5, false, 60, 20},
	}

	for _, scenario := range scenarios {
		var bees []Bee

		if scenario.withMissChance {
			bees = createBees(scenario.beeType, scenario.count, scenario.expectedMissChance)
		} else {
			bees = createBees(scenario.beeType, scenario.count)
		}

		if len(bees) != int(scenario.count) {
			t.Errorf("createBees gave incorrect number of bees. Requested: %d. Received: %d\n", scenario.count, len(bees))
		}

		for _, bee := range bees {
			if bee.Type != scenario.beeType {
				t.Errorf("createBees gave incorrect bee type. Requested bee type: %s. Received bee type: %s\n", scenario.beeType, bee.Type)
			}

			if bee.Health != scenario.expectedHealth {
				t.Errorf("createBees gave bee with incorrect health. Expected health: %d. Received health: %d\n", scenario.expectedHealth, bee.Health)
			}

			if bee.MissChance != scenario.expectedMissChance {
				t.Errorf("createBees gave bee with incorrect miss chance. Expected miss chance: %d. Received miss chance: %d\n", scenario.expectedMissChance, bee.MissChance)
			}
		}
	}

	bees := createBees(BeeType(5), 5)

	if bees != nil {
		t.Error("createBees returned a slice when nil was expected.")
	}
}

func TestCreateHive(t *testing.T) {
	numQueens := 5
	numWorkers := 10
	numDrones := 15

	hive := createHive(uint(numQueens), uint(numWorkers), uint(numDrones))

	if len(hive) != numQueens+numWorkers+numDrones {
		receivedQueens := 0
		receivedWorkers := 0
		receivedDrones := 0

		for _, bee := range hive {
			if bee.Type == QueenBee {
				receivedQueens++
			} else if bee.Type == WorkerBee {
				receivedWorkers++
			} else if bee.Type == DroneBee {
				receivedDrones++
			}
		}

		if receivedQueens != numQueens {
			t.Errorf("createHive gave incorrect number of queens. Requested: %d. Received: %d\n", numQueens, receivedQueens)
		}

		if receivedWorkers != numWorkers {
			t.Errorf("createHive gave incorrect number of workers. Requested: %d. Received: %d\n", numWorkers, receivedWorkers)
		}

		if receivedDrones != numDrones {
			t.Errorf("createHive gave incorrect number of drones. Requested: %d. Received: %d\n", numDrones, receivedDrones)
		}
	}
}
