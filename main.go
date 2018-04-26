package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	
	"image/color"
	"image/png"
	"math/rand"
	"fmt"
	"bytes"
	"time"
	"log"
	"strconv"
)

const (
	screenWidth = 256
	screenHeight = 224
)

type course struct {
	curveTo int
	positionX int
	isFirst bool
	data []float64
}

type player struct {
	positionX float64
}

var raceCourse course
var racePlayer player
var playerImage *ebiten.Image
var courseImage *ebiten.Image
var mainFont font.Face
var score int
var isGameOver bool
var bigFont font.Face
var isDrawScore bool
var count int
var waitSpaceKey bool

func init() {
	rand.Seed(time.Now().UnixNano())

	var err error
	r := bytes.NewReader(carBytes)
	carPNG, _ := png.Decode(r)
	playerImage, err = ebiten.NewImageFromImage(carPNG, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	courseImage, err = ebiten.NewImage(64, 1, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	courseImage.Fill(color.RGBA{168, 171, 169, 255})

	raceCourse = course{
		curveTo: rand.Intn(2),
		positionX: 88,
		isFirst: true,
	}

	mainFont = decodeFont(fontBytes, 8)
	bigFont = decodeFont(fontBytes, 24)
	score = 0

	isGameOver = false
	isDrawScore = false
	count = 0

}

func createCourse() {
	for {
		for i := 0; !isGameOver; i ++ {
			if i == 180 && count == 0 {
				racePlayer.positionX = raceCourse.data[60] + 32
			}
			curve := []int{1, 0, -1}
			if i % 46 == 0 {//52 72
				raceCourse.curveTo = 0
			}
			if i % 56 == 0 && score != 0 {
				raceCourse.curveTo = curve[rand.Intn(3)]
			}

			if i % 1 == 0 {
				if raceCourse.positionX < 20 && raceCourse.curveTo == curve[2] {
					raceCourse.curveTo = 0
				}else if raceCourse.positionX > screenWidth - 64 - 20 && raceCourse.curveTo == curve[0] {
					raceCourse.curveTo = 0
				}else {
					raceCourse.positionX += raceCourse.curveTo	
				}
			}

			if count == 0 {
				if i > 180 {
					time.Sleep(time.Millisecond * 16)
					raceCourse.data = append(raceCourse.data[1:], float64(raceCourse.positionX))	
				}else {
					raceCourse.data = append(raceCourse.data, float64(raceCourse.positionX))
				}
			}else {
				time.Sleep(time.Millisecond * 16)
				raceCourse.data = append(raceCourse.data[1:], float64(raceCourse.positionX))	
			}
			
			if i % 110 == 0 && i > 180 {
				score ++
				isDrawScore = true
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
	
}

func (c *course) draw(screen *ebiten.Image) {
	// for i := 0; i < len(c.data); i ++ {
	for i := 0; i < 181; i ++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(c.data[i], float64(210-i))
		screen.DrawImage(courseImage, op)
	}
}

func collision() {
	for i := 30; i < 66 && !isGameOver; i ++ {
		if racePlayer.positionX <= raceCourse.data[i] {
			isGameOver = true
		}
		if racePlayer.positionX + 24 >= raceCourse.data[i] + 64 {
			isGameOver = true
		}
	}
}

func update(screen *ebiten.Image) error {
	screen.Fill(colornames.Darkgreen)
	
	raceCourse.draw(screen)

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(racePlayer.positionX, 144)
	screen.DrawImage(playerImage, op)

	text.Draw(screen, "SCORE " + strconv.Itoa(score), mainFont, 172, 18, colornames.White)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	
	collision()

	if score <= 2 && count == 0 {
		text.Draw(screen, "セブン レース", bigFont, 48, 50, colornames.White)
	}

	if isGameOver {
		text.Draw(screen, "GAME OVER", bigFont, 20, 77, colornames.White)
		text.Draw(screen, "PRESS SPACE KEY", mainFont, 67, 90, colornames.White)
	}

	return nil
}

func keyEvent() {
	for {
		if isGameOver {
			time.Sleep(time.Millisecond * 10)
		}
		// if inpututil.IsKeyJustPressed(ebiten.KeySpace) && isGameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) && (count == 0 || isGameOver) {
			for i := 0; i < len(raceCourse.data)-40; i ++ {
				raceCourse.data[i] = raceCourse.data[len(raceCourse.data)-40]
			}
			racePlayer.positionX = raceCourse.data[len(raceCourse.data)-40] + 32
			isGameOver = false
			count ++

			score = 0
		}else if !isGameOver{
			time.Sleep(time.Millisecond * 10)
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				racePlayer.positionX -= 0.5
			}
			if ebiten.IsKeyPressed(ebiten.KeyRight) {
				racePlayer.positionX += 0.5
			}
		}
		
	}
}

func main() {
	go createCourse()	
	go keyEvent()
	
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "セブン レース"); err != nil {
		log.Fatal(err)
	}
}
