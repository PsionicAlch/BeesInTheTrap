package game

import (
	"testing"
	"time"
)

func TestCreateCommunicationProtocol(t *testing.T) {
	communication := createCommunicationProtocol()

	if communication == nil {
		t.Fatal("Expected non-nil CommunicationProtocol.")
	}

	if communication.hitSignal == nil {
		t.Error("Expected hitSignal to be initialized.")
	}

	if communication.eventChannel == nil {
		t.Error("Expected eventChannel to be initialized.")
	}
}

func TestHit(t *testing.T) {
	communication := createCommunicationProtocol()

	expectedEvent := Event{
		Type:    PlayerAttack,
		Message: "Test Message",
		State: GameState{
			Round: 25,
		},
	}

	go func() {
		select {
		case <-communication.hitSignal:
			communication.eventChannel <- expectedEvent
		case <-time.After(time.Second * 3):
			t.Errorf("Timed out waiting for hitSignal")
		}
	}()

	actualEvent := communication.Hit()

	if actualEvent.Type != expectedEvent.Type && actualEvent.Message != expectedEvent.Message && actualEvent.State.Round != expectedEvent.State.Round {
		t.Errorf("Expected event does not match received event.")
	}
}

func TestWaitForCPU(t *testing.T) {
	communication := createCommunicationProtocol()

	expectedEvent := Event{
		Type:    PlayerAttack,
		Message: "Test Message",
		State: GameState{
			Round: 25,
		},
	}

	go func() {
		communication.eventChannel <- expectedEvent
	}()

	actualEvent := communication.WaitForCPU()

	if actualEvent.Type != expectedEvent.Type && actualEvent.Message != expectedEvent.Message && actualEvent.State.Round != expectedEvent.State.Round {
		t.Errorf("Expected event does not match received event.")
	}
}

func TestWaitForPlayer(t *testing.T) {
	communication := createCommunicationProtocol()

	go func() {
		communication.hitSignal <- struct{}{}
	}()

	communication.waitForPlayer()
}

func TestHitResponse(t *testing.T) {
	communication := createCommunicationProtocol()

	msg := "Test Message"
	state := GameState{
		Round: 25,
	}

	go func() {
		communication.hitResponse(msg, state)
	}()

	select {
	case event := <-communication.eventChannel:
		if event.Type != PlayerAttack && event.Message != msg && event.State.Round != state.Round {
			t.Errorf("Expected event does not match received event.")
		}
	case <-time.After(time.Second * 3):
		t.Errorf("Timed out waiting for hitSignal")
	}
}

func TestStingResponse(t *testing.T) {
	communication := createCommunicationProtocol()

	msg := "Test Message"
	state := GameState{
		Round: 25,
	}

	go func() {
		communication.stingResponse(msg, state)
	}()

	select {
	case event := <-communication.eventChannel:
		if event.Type != HiveAttack && event.Message != msg && event.State.Round != state.Round {
			t.Errorf("Expected event does not match received event.")
		}
	case <-time.After(time.Second * 3):
		t.Errorf("Timed out waiting for hitSignal")
	}
}

func TestGameFinishedResponse(t *testing.T) {
	communication := createCommunicationProtocol()

	msg := "Test Message"
	state := GameState{
		Round: 25,
	}

	go func() {
		communication.gameFinishedResponse(msg, state)
	}()

	select {
	case event := <-communication.eventChannel:
		if event.Type != GameFinished && event.Message != msg && event.State.Round != state.Round {
			t.Errorf("Expected event does not match received event.")
		}
	case <-time.After(time.Second * 3):
		t.Errorf("Timed out waiting for hitSignal")
	}
}
