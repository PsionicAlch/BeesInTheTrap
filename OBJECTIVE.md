# BeesInTheTrap

Write a simple game in Go using the following rules. You should consider using Go design principles. The code should also include tests.

- Demonstrate concurrency (goroutines)
- Demonstrate Go Data structures such as maps and slices

## Rules

- The game must run from the command line.
- The game must be compiled to an executable binary for Windows and Linux.
- The game should be turn based.
- There is a player, and a hive of bees, each player starts with 100 HP.
- The player can enter “hit” during their turn to play. Once a player has had their turn, the bees can take their turn. The bees are the computer against which players play.
- The game should show a message to the player after each turn, both computers and player turns, with the outcome of that turn. e.g. "Direct Hit! You took 12 hit points from a Drone bee" or "Sting! You just got stun by a drone bee" or "Miss! You just missed the hive, better luck next time!" or "Buzz! That was close! The Queen Bee just missed you!" ect…
- The game should be single player.
- The game should exit on game over, with a message informing the user of how many hits were needed to destroy the hive, or if they are dead, how many stings it took for the hive to kill them. As well as a summary of the game.
- The object of the game is to destroy a hive of bees before the bees sting you to death.
- The hive has three types of bee.
- The game is over when all bees are dead, or you are dead first.
- You are dead when your HP reaches 0.
- The player should sometimes be able to miss the hive entirely.
- The bees should sometimes be able to miss the player entirely.
- The bee a player hits should be picked at random, but should be statistically correct given the ratio of each bee type and miss chances.
- The bee that strings a player should be picked at random, but should be statistically correct given the ratio of each bee type and miss chances.
- You should include a way for the player to auto spin and let the game run through to an ending on their behalf.

### Queen Bee

- The Queen Bee has a lifespan of 100 Hit Points.
- When the Queen Bee is hit, 10 Hit Points are deducted from her lifespan.
- If/When the Queen Bee has run out of Hit Points, All remaining alive Bees automatically run out of hit points.
- When the Queen bee stings a player, the player loses 10 Hit Points.
- There is only 1 Queen Bee.

### Worker Bee

- Worker Bees have a lifespan of 75 Hit Points.
- When a Worker Bee is hit, 25 Hit Points are deducted from their lifespan.
- When a Worker Bee stings a player, the player loses 5 Hit Points.
- There are 5 Worker Bees.

### Drone Bee

- Drone Bees have a lifespan of 60 Hit Points.
- When a Drone Bee is hit, 30 Hit Points are deducted from his lifespan.
- When a Drone Bee stings a player, the player loses 1 Hit Point.
- There are 25 Drone Bees.