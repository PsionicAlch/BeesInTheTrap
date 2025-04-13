# 🐝 Bees In The Trap

A turn-based command-line game written in Go, where you face off against a hive of dangerous bees. Your objective: destroy the hive before they sting you to death.

Each game runs in its own goroutine, simulating server-like behavior in a local environment. The core loop is built to handle both player-controlled and automated simulations, showcasing concurrency via Go channels.

This project was built to the specifications of the technical assessment, a copy of which can be found in [OBJECTIVE.md]("https://github.com/PsionicAlch/BeesInTheTrap/blob/main/OBJECTIVE.md").

## 🧩 Game Overview

You play as a lone wanderer with 100 HP, facing a hive composed of:

- **1 Queen Bee** — 100 HP, deals 10 damage
- **5 Worker Bees** — 75 HP each, deal 5 damage
- **25 Drone Bees** — 60 HP each, deal 1 damage

Every turn:

- You can choose to `hit` the hive.
- A bee is selected randomly (weighted by population) and might be hit — or missed.
- The hive retaliates: a random bee stings you — or misses.
- The game ends when either **all bees are dead** or **you are**.

---

## 📦 Build & Run

### Prerequisites

This project was built using Golang v1.24.1.

### Build

```bash
make build
```

This compiles the project into a binary at ```./tmp/BeesInTheTrap```.

### Run

```bash
make run
```

This builds and runs the game in one step.

---

## 🧪 Running Tests

```bash
make test
```

To run tests with coverage:

```bash
make coverage
```

A coverage report will be generated and opened in your browser.

---

## 🧹 Clean Build Artifacts

```bash
make clean
```

Removes all files from the ./tmp/ directory.