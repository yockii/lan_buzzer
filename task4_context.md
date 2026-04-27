You are implementing Task 4: Implement Game Server Logic

## Task Description

**Files:**
- Create: `backend/game/server.go`

This task implements the core game server logic with methods for managing game state, players, and buzz recording. Following TDD methodology.

## Full Task Specification

See: D:/projects/github.com/yockii/lan_qr/docs/superpowers/plans/2026-04-27-lan-buzzer-implementation.md
Lines 352-570

## Key Requirements

1. **Write tests first** (TDD approach)
2. **Implement all methods:**
   - NewGameServer() - Constructor
   - AddPlayer(player) - Thread-safe player addition
   - RemovePlayer(playerID) - Thread-safe player removal
   - GetPlayers() - Thread-safe player list retrieval
   - StartGame() - Transition to ready state, set start time
   - RecordBuzz(playerID) - Record buzz, transition to locked, return bool
   - ResetGame() - Reset to waiting state
   - GetState() - Thread-safe state getter
   - GetWinner() - Thread-safe winner retrieval

3. **All operations must be thread-safe** using Mutex
4. **Follow TDD:** Write tests → Run (fail) → Implement → Run (pass)
5. **Commit** with message "feat: implement game server state management"

## Context

This task builds on Tasks 2 and 3 (types and colors). The GameServer struct is already defined in types.go with Mutex field. This task adds methods to manipulate the game state.

Dependencies: Tasks 2 and 3 must be complete.

## Work From

D:\projects\github.com\yockii\lan_qr

## Report Format

When done, report:
- **Status:** DONE | DONE_WITH_CONCERNS | BLOCKED | NEEDS_CONTEXT
- What you implemented
- Test results
- Files changed
- Self-review findings
- Any issues or concerns
