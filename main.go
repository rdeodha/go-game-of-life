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

var xC int
var yC int
var gameState [32][32]bool
var isPaused bool = true
var isSet bool = false
var resetBox image.Rectangle = text.BoundString(basicfont.Face7x13, "Reset")
var startBox image.Rectangle = text.BoundString(basicfont.Face7x13, "Start")
var updates int

func (g *Game) Update() error {
	updates++
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		xC, yC = ebiten.CursorPosition()

		var hBlock int = int(xC / 20)
		var vBlock int = int(yC / 20)
		if !(hBlock >= 32 || hBlock < 0 || vBlock >= 32 || vBlock < 0) {
			isSet = true
			gameState[hBlock][vBlock] = !gameState[hBlock][vBlock]
		} else if isSet && xC >= 330-2 && xC <= 330+resetBox.Dx()+2 && yC >= 700-7 && yC <= 700+resetBox.Dy()+7 {
			var resetGame [32][32]bool
			gameState = resetGame
			isPaused = true
			isSet = false
		} else if isPaused && xC >= 250-2 && xC <= 250+startBox.Dx()+2 && yC >= 700-7 && yC <= 700+startBox.Dy()+7 {
			isPaused = false
		} else if !isPaused && xC >= 250-2 && xC <= 250+startBox.Dx()+2 && yC >= 700-7 && yC <= 700+startBox.Dy()+7 {
			isPaused = true
		}

	}

	if !isPaused && updates%7 == 0 {
		var copyGameState = gameState
		for i := 0; i < 32; i++ {
			for j := 0; j < 32; j++ {
				liveCount := 0
				if i-1 >= 0 && j-1 >= 0 && gameState[i-1][j-1] {
					fmt.Println("NW")
					liveCount++
				}

				if j-1 >= 0 && gameState[i][j-1] {
					fmt.Println("N")
					liveCount++
				}

				if i+1 < 32 && j-1 >= 0 && gameState[i+1][j-1] {
					fmt.Println("NE")
					liveCount++
				}

				if i+1 < 32 && gameState[i+1][j] {
					fmt.Println("E")
					liveCount++
				}

				if i+1 < 32 && j+1 < 32 && gameState[i+1][j+1] {
					fmt.Println("SE")
					liveCount++
				}

				if j+1 < 32 && gameState[i][j+1] {
					fmt.Println("S")
					liveCount++
				}

				if i-1 >= 0 && j+1 < 32 && gameState[i-1][j+1] {
					fmt.Println("SW")
					liveCount++
				}

				if i-1 >= 0 && gameState[i-1][j] {
					fmt.Println("W")
					liveCount++
				}

				if copyGameState[i][j] && (liveCount == 2 || liveCount == 3) {
					copyGameState[i][j] = true
				} else if !copyGameState[i][j] && (liveCount == 3) {
					copyGameState[i][j] = true
				} else {
					copyGameState[i][j] = false
				}

			}
		}

		gameState = copyGameState
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < 32; i++ {
		for j := 0; j < 32; j++ {
			if gameState[i][j] {
				ebitenutil.DrawRect(screen, float64((i * 20)), float64((j * 20)), 20, 20, color.White)
			}
			ebitenutil.DrawRect(screen, float64(i*20), float64(j*20), 1, 20, color.RGBA{105, 105, 105, 255})
			ebitenutil.DrawRect(screen, float64(i*20), float64(j*20), 20, 1, color.RGBA{105, 105, 105, 255})
		}
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
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
