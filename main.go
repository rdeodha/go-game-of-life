package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"golang.org/x/image/font/basicfont"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct{}
type Cell struct {
	x         int
	y         int
	alive     bool
	count     int
	makeAlive bool
	makeKill  bool
}
type Range struct {
	min int
	max int
}

var xC int
var yC int
var isPaused bool = true
var isSet bool = false
var resetBox image.Rectangle = text.BoundString(basicfont.Face7x13, "Reset")
var startBox image.Rectangle = text.BoundString(basicfont.Face7x13, "Start")
var updates int
var cellSize int = 20
var horizontalRange Range = Range{0, 32}
var verticleRange Range = Range{0, 32}
var activeCells = make(map[string]*Cell)

func (g *Game) Update() error {
	updates++
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		xC, yC = ebiten.CursorPosition()

		if yC >= 0 && yC < 640 {
			isSet = true
			var absoluteX int = horizontalRange.min + int(xC/cellSize)
			var absoluteY int = verticleRange.min + int(yC/cellSize)

			var key string = fmt.Sprintf("%d-%d", absoluteX, absoluteY)
			// var keyN string = fmt.Sprintf("%d-%d", absoluteX, absoluteY-1)
			// var keyNE string = fmt.Sprintf("%d-%d", absoluteX+1, absoluteY-1)
			// var keyE string = fmt.Sprintf("%d-%d", absoluteX+1, absoluteY)
			// var keySE string = fmt.Sprintf("%d-%d", absoluteX+1, absoluteY+1)
			// var keyS string = fmt.Sprintf("%d-%d", absoluteX, absoluteY+1)
			// var keySW string = fmt.Sprintf("%d-%d", absoluteX-1, absoluteY+1)
			// var keyW string = fmt.Sprintf("%d-%d", absoluteX-1, absoluteY)
			// var keyNW string = fmt.Sprintf("%d-%d", absoluteX-1, absoluteY-1)

			var cell = activeCells[key]

			if cell != nil && cell.alive {
				// Valid && Alive  -> Make cell !alive. Subtract count from cell and surrounding. If count == 0 of cell and surrounding, delete.
				cell.alive = false
				cell.count--

				if cell.count == 0 {
					delete(activeCells, key)
				}

				// N
				deleteCellHelper(absoluteX, absoluteY-1)
				// NE
				deleteCellHelper(absoluteX+1, absoluteY-1)
				// E
				deleteCellHelper(absoluteX+1, absoluteY)
				// SE
				deleteCellHelper(absoluteX+1, absoluteY+1)
				// S
				deleteCellHelper(absoluteX, absoluteY+1)
				// SW
				deleteCellHelper(absoluteX-1, absoluteY+1)
				// W
				deleteCellHelper(absoluteX-1, absoluteY)
				// NW
				deleteCellHelper(absoluteX-1, absoluteY-1)

			} else if cell == nil {
				// !Valid -> Make new alive cell, set count to 1. Create not alive cells around it, if they dont already exist, increment count
				var cell Cell = Cell{absoluteX, absoluteY, true, 1, false, false}
				activeCells[key] = &cell

				// N
				notValidCellHelper(absoluteX, absoluteY-1)
				// NE
				notValidCellHelper(absoluteX+1, absoluteY-1)
				// E
				notValidCellHelper(absoluteX+1, absoluteY)
				// SE
				notValidCellHelper(absoluteX+1, absoluteY+1)
				// S
				notValidCellHelper(absoluteX, absoluteY+1)
				// SW
				notValidCellHelper(absoluteX-1, absoluteY+1)
				// W
				notValidCellHelper(absoluteX-1, absoluteY)
				// NW
				notValidCellHelper(absoluteX-1, absoluteY-1)

			} else {
				// Valid && !Alive -> Make cell alive, add count, Create not alive cells around it, if they dont already exist, increment count
				cell.alive = true
				cell.count++

				// N
				notValidCellHelper(absoluteX, absoluteY-1)
				// NE
				notValidCellHelper(absoluteX+1, absoluteY-1)
				// E
				notValidCellHelper(absoluteX+1, absoluteY)
				// SE
				notValidCellHelper(absoluteX+1, absoluteY+1)
				// S
				notValidCellHelper(absoluteX, absoluteY+1)
				// SW
				notValidCellHelper(absoluteX-1, absoluteY+1)
				// W
				notValidCellHelper(absoluteX-1, absoluteY)
				// NW
				notValidCellHelper(absoluteX-1, absoluteY-1)
			}

		} else if isSet && xC >= 330-2 && xC <= 330+resetBox.Dx()+2 && yC >= 700-7 && yC <= 700+resetBox.Dy()+7 {
			activeCells = make(map[string]*Cell)
			isPaused = true
			isSet = false
		} else if isPaused && xC >= 250-2 && xC <= 250+startBox.Dx()+2 && yC >= 700-7 && yC <= 700+startBox.Dy()+7 {
			isPaused = false
		} else if !isPaused && xC >= 250-2 && xC <= 250+startBox.Dx()+2 && yC >= 700-7 && yC <= 700+startBox.Dy()+7 {
			isPaused = true
		}
	}

	if !isPaused && updates%7 == 0 {
		for _, value := range activeCells {
			liveCount := 0

			// if value.alive {
			// 	fmt.Printf("X: %d, Y: %d\n", value.x, value.y)
			// }

			var absoluteX = value.x
			var absoluteY = value.y

			// North
			if liveCountHelper(absoluteX, absoluteY-1) {
				liveCount++
				//fmt.Println("N")
			}
			// North-East
			if liveCountHelper(absoluteX+1, absoluteY-1) {
				liveCount++
				//fmt.Println("NE")
			}
			// East
			if liveCountHelper(absoluteX+1, absoluteY) {
				liveCount++
				//fmt.Println("E")
			}
			// South-East
			if liveCountHelper(absoluteX+1, absoluteY+1) {
				liveCount++
				//fmt.Println("SE")
			}
			// South
			if liveCountHelper(absoluteX, absoluteY+1) {
				liveCount++
				//fmt.Println("S")
			}
			// South-West
			if liveCountHelper(absoluteX-1, absoluteY+1) {
				liveCount++
				//fmt.Println("SW")
			}
			// West
			if liveCountHelper(absoluteX-1, absoluteY) {
				liveCount++
				//fmt.Println("W")
			}
			// North-West
			if liveCountHelper(absoluteX-1, absoluteY-1) {
				liveCount++
				// fmt.Println("NW")
			}

			if value.alive && (liveCount == 2 || liveCount == 3) {
				// Do nothing, cell survives
			} else if !value.alive && liveCount == 3 {
				value.makeAlive = true
			} else if value.alive {
				value.makeKill = true
			}

			// if value.alive {
			// 	fmt.Println(fmt.Sprintf("%t, X: %d, Y: %d, Live: %d, MakeAlive: %t, MakeKill: %t\n", value.alive, value.x, value.y, liveCount, value.makeAlive, value.makeKill))
			// }
		}

		for _, value := range activeCells {
			if value.makeAlive {
				value.alive = true
				value.count++
				var absoluteX = value.x
				var absoluteY = value.y

				// N
				notValidCellHelper(absoluteX, absoluteY-1)
				// NE
				notValidCellHelper(absoluteX+1, absoluteY-1)
				// E
				notValidCellHelper(absoluteX+1, absoluteY)
				// SE
				notValidCellHelper(absoluteX+1, absoluteY+1)
				// S
				notValidCellHelper(absoluteX, absoluteY+1)
				// SW
				notValidCellHelper(absoluteX-1, absoluteY+1)
				// W
				notValidCellHelper(absoluteX-1, absoluteY)
				// NW
				notValidCellHelper(absoluteX-1, absoluteY-1)

				value.makeAlive = false
			} else if value.makeKill {
				var absoluteX = value.x
				var absoluteY = value.y
				value.alive = false
				value.makeKill = false
				value.count--
				if value.count == 0 {
					delete(activeCells, fmt.Sprintf("%d-%d", value.x, value.y))
				}

				// N
				deleteCellHelper(absoluteX, absoluteY-1)
				// NE
				deleteCellHelper(absoluteX+1, absoluteY-1)
				// E
				deleteCellHelper(absoluteX+1, absoluteY)
				// SE
				deleteCellHelper(absoluteX+1, absoluteY+1)
				// S
				deleteCellHelper(absoluteX, absoluteY+1)
				// SW
				deleteCellHelper(absoluteX-1, absoluteY+1)
				// W
				deleteCellHelper(absoluteX-1, absoluteY)
				// NW
				deleteCellHelper(absoluteX-1, absoluteY-1)
			}

		}
		// fmt.Println("\n\n\n\n")
	}

	return nil
}

func liveCountHelper(absX int, absY int) (live bool) {
	var key = fmt.Sprintf("%d-%d", absX, absY)
	if activeCells[key] != nil && activeCells[key].alive {
		return true
	} else {
		return false
	}
}

func deleteCellHelper(absX int, absY int) {
	var key string = fmt.Sprintf("%d-%d", absX, absY)
	var c = activeCells[key]

	if c != nil {
		c.count--

		if c.count == 0 {
			delete(activeCells, key)
		}
	}
}

func notValidCellHelper(absX int, absY int) {
	var key string = fmt.Sprintf("%d-%d", absX, absY)
	var cell = activeCells[key]
	if cell == nil {
		var c Cell = Cell{absX, absY, false, 1, false, false}
		activeCells[key] = &c
	} else {
		var c = activeCells[key]
		c.count++
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, value := range activeCells {
		var cell = value
		if value.x >= horizontalRange.min && value.x < horizontalRange.max && value.y >= verticleRange.min && value.y < verticleRange.max {
			var spacesRight int = value.x - horizontalRange.min
			var spacesDown int = value.y - horizontalRange.min
			if cell.alive {
				ebitenutil.DrawRect(screen, float64((spacesRight * cellSize)), float64((spacesDown * cellSize)), float64(cellSize), float64(cellSize), color.White)
			} else {
				// ebitenutil.DrawRect(screen, float64((spacesRight * cellSize)), float64((spacesDown * cellSize)), float64(cellSize), float64(cellSize), color.RGBA{105, 105, 105, 255})
			}

			// if cell.alive {
			// 	text.Draw(screen, fmt.Sprintf("%d", cell.count), basicfont.Face7x13, int((spacesRight*cellSize))+10, int((spacesDown*cellSize))+10, color.RGBA{0, 255, 0, 255})
			// } else {
			// 	text.Draw(screen, fmt.Sprintf("%d", cell.count), basicfont.Face7x13, int((spacesRight*cellSize))+10, int((spacesDown*cellSize))+10, color.Black)
			// }
		}

	}

	for i := 0; i < horizontalRange.max; i++ {
		ebitenutil.DrawRect(screen, float64(i*cellSize), float64(0), 1, 640, color.RGBA{105, 105, 105, 255})
	}

	for i := 0; i < verticleRange.max; i++ {
		ebitenutil.DrawRect(screen, float64(0), float64(i*cellSize), 640, 1, color.RGBA{105, 105, 105, 255})
	}

	ebitenutil.DrawRect(screen, float64(639), float64(0), 1, 640, color.RGBA{105, 105, 105, 255})
	ebitenutil.DrawRect(screen, float64(0), float64(640), 640, 1, color.RGBA{105, 105, 105, 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d ---- Y: %d", xC, yC))

	if isPaused {
		text.Draw(screen, "Start", basicfont.Face7x13, 250, 700, color.White)
	} else {
		text.Draw(screen, "Pause", basicfont.Face7x13, 250, 700, color.White)
	}

	if isSet {
		text.Draw(screen, "Reset", basicfont.Face7x13, 330, 700, color.White)
	} else {
		text.Draw(screen, "Reset", basicfont.Face7x13, 330, 700, color.RGBA{105, 105, 105, 255})
	}
	text.BoundString(basicfont.Face7x13, "Reset")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 750
}

func main() {
	ebiten.SetWindowSize(640, 750)
	ebiten.SetWindowTitle("Game of Life")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
