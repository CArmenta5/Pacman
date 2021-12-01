# Pacman 🎮

## Authors 🔥:

- Carlos Armenta
- Mathew Gabriel Lopez
- Francisco Ramos

## Usage: 
- Use the command *make build* to compile the project
- Use the command *make run* to run the game

## Architure

Multithreaded Pacman created with graphical interface and implemented in C, using the libraries SDL (SDL2 image and SDL2 ttf) and OpenMP for GUI and user input and use of threads, respectively.

Flow Diagram: 

![alt text](https://github.com/CArmenta5/Pacman/blob/main/PacmanDiagram.png?raw=true)

## Data Structures
We used the following data structures: 

### Sprite (Player and ghost)
Atributes | DataType | Description | 
--- | --- | --- |
row | int | Coordinate in the X axis |
col | int | Coordinate in the Y axis | 

### Game
Atributes | DataType | Description | 
--- | --- | --- |
layers | int[] | Represent the tile in the Image |


### Config
Atributes | DataType | Description | 
--- | --- | --- |
Player | string | Emoji to represent the player |
Ghost | string | Emoji to represent the ghost | 
Wall | string | Emoji to represent the wall |
Dot | string | Emoji to represent a dot | 
Death | string | Emoji to represent the game over |
Space | string | Representation for white space |

### Functions
* `*loadconfig():*` Reads information from JSON File. (Emojis and settings)

* `*cleanup():*` Clear the console and turn of cbreak mode.
* `*clearscreen():*` Clears the console string in order to set a new frame.
* `*moveCursor(row, col int):*` Moves the cursor to the given coordinates in order to print a change in the board.
* `*moveCursorEmoji(row, col int):*` Moves the cursor to the given coordinates in order to print a change in the board.
* `*readInput():*` Reads from standard input and returns the direction for the movement. This is done in a goroutine.
* `*loadLeve(file string):*` Reads the designed level from the .txt file.
* `*printScreen():*` Prints a new frame of the game in the console.
* `*movePlayer():*` Performs the movement of the player. If the current cell has a normal dot, the score increments 1 unit. If there is a power-up, the score increments 10 units.
* `*directionGhost() string:*` Calculates a random number and sets the corresponding direction for the next move of the ghost.
* `*moveGhost():*` Moves the ghost to the corresponding direction. This is done in a goroutine.
* `*init():*` Decode an image from the image file's byte slice.
* `*Update(g *Game) error:*` Refresh the screen of the library ebitenutil
* `*Draw(screeb *ebiten.Image):*` DDraw each tile with each DrawImage call.

### Assumptions:

- You are running in linux distribution


