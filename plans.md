# Future Plans for the Riichi Mahjong Project

## Short-Term Goals

- ~~Implement tests for all existing modules to ensure code reliability.~~ (Completed)
- ~~Add additional tests for edge cases to patch oversights.~~ (Completed)
- ~~Finalize the yaku detection module to cover all standard yaku.~~ (Completed)
- Complete the scoring module to calculate hand scores based on detected yaku.
- [Important] Rework hand/set/tile/win context window implementations to accommodate additional information.
  - eg. red fives, dora indicators, riichi status, wait types, etc.
- Modify parse to tolerate precompleted Sets from calls during gameplay. (should help with scoring)
- Develop a complete progression that utilizes all currently implemented features.
  - eg. input hand, detect validity, calculate yaku, score hand.
  - Result should be:
    - Input: Hand as a list of tiles.
    - Output: Validity (T/F), Yaku List, Han, Fu, Score.

## Long-Term Goals

- Explore and implement the Quadtree Algorithm for deficiency calculation.
- Explore and implement the Block Deficiency Algorithm for improved hand analysis.
- [Target] Explore and implement the Hierarchical Branch and Bound Algorithm for optimal tile selection.
- Develop a solution that utilizes the above algorithms to create an CPU opponent.
- Consider possible graphical implementations for user interaction. (eg. web app, desktop app)
- Expand documentation to cover all modules and functions comprehensively.
