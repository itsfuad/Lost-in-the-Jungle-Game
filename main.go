package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type direction int

const (
	idle direction = iota
	up
	down
	left
	right
)

const (
	WINDOW_WIDTH  = 1500
	WINDOW_HEIGHT = 800
	FPS           = 120.00
)

var (
	mainPlayer 	*Player

	//grass tiles
	grassTile 	*Map

	entities 	[]Entity

	//game frame
	gameFrame 	int32 = 0
	camera 		rl.Camera2D

	//bg music
	bgMusic 	rl.Music

	//settings
	MusicOn 	bool = true
	SoundOn 	bool = true
)

type Entity interface {
	Draw()
	Update()
}

type Player struct {
	Speed          	float32
	_OriginalSpeed 	float32
	Direction      	direction
	CurrentFrame   	int32
	CurrentRow     	int32
	SpriteRowCount 	int32
	SpriteColCount 	int32
	FrameWidth     	int32
	FrameHeight    	int32
	FrameSpeed     	int32

	Source rl.Rectangle
	Destination rl.Rectangle

	//sprite
	Texture rl.Texture2D
	Running bool
}

func NewPlayer(texturePath string, position rl.Vector2, speed float32, spriteRowCount int32, spriteColCount int32, frameSpeed int32, currentRow int32) *Player {


	var currentFrame int32 = 0
	texture := rl.LoadTexture(texturePath)
	frameWidth := texture.Width / spriteColCount
	frameHeight := texture.Height / spriteRowCount

	fmt.Printf("Frame Width: %v\n", frameWidth)
	fmt.Printf("Frame Height: %v\n", frameHeight)

	source :=  rl.NewRectangle(float32(currentFrame * frameWidth), float32(currentRow * frameHeight), float32(frameWidth), float32(frameHeight))
	destination := rl.NewRectangle(position.X, position.Y, float32(frameWidth), float32(frameHeight))

	//declare player
	return &Player{
		Speed:          speed,
		_OriginalSpeed: speed,
		Direction:      idle,
		CurrentFrame:   currentFrame,
		CurrentRow:     currentRow,
		SpriteRowCount: spriteRowCount,
		SpriteColCount: spriteColCount,
		FrameWidth:     frameWidth,
		FrameHeight:    frameHeight,
		FrameSpeed:     frameSpeed,

		Source: source,
		Destination: destination,

		Texture:        rl.LoadTexture(texturePath),
		Running:        false,
	}
}

func (p *Player) Draw() {
	//draw player
	rl.DrawTexturePro(p.Texture, p.Source, p.Destination, rl.NewVector2(p.Destination.Width, p.Destination.Height), 0, rl.White)
	p.Update()
}

func (p *Player) HandleControls() {

	p.Running = false

	updown := false
	leftRight := false
	
	if mainPlayer == p {
		if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
			p.Destination.Y -= p.Speed
			p.Direction = up
			p.Running = true
			updown = true
			//fmt.Printf("Moving Up %v\n", p.Destination.Y)
		} else if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
			p.Destination.Y += p.Speed
			p.Direction = down
			p.Running = true
			updown = true
			//fmt.Printf("Moving Down %v\n", p.Destination.Y)
		} 
		if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
			p.Destination.X -= p.Speed
			p.Direction = left
			p.Running = true
			leftRight = true
			//fmt.Printf("Moving Left %v\n", p.Destination.X)
		} else if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
			p.Destination.X += p.Speed
			p.Direction = right
			p.Running = true
			leftRight = true
			//fmt.Printf("Moving Right %v\n", p.Destination.X)
		}
	}

	if p.Running && updown && leftRight {
		//diagonal movement should be slower.
		p.Speed = p._OriginalSpeed / 1.5
	} else {
		p.Speed = p._OriginalSpeed
	}

	if p.Running {
		//update player
		switch p.Direction {
		case up:
			p.CurrentRow = 1
		case down:
			p.CurrentRow = 0
		case left:
			p.CurrentRow = 2
		case right:
			p.CurrentRow = 3
		}

		if gameFrame % 8 == 1 { // change charecter frame every 8 frames of the game
			p.CurrentFrame++
		}
		
	} else {
		p.CurrentFrame = 0
	}

	p.Source.X = float32(p.CurrentFrame * p.FrameWidth)
	p.Source.Y = float32(p.CurrentRow * p.FrameHeight)
}

func (p *Player) Update() {

	p.HandleControls()

	if p.CurrentFrame >= p.SpriteColCount {
		p.CurrentFrame = 0
	}

	//fmt.Printf("Current Frame: %v\n", p.CurrentFrame)
}



type Map struct {
	Position rl.Vector2
	Texture  rl.Texture2D
}

func NewTile(texturePath string, position rl.Vector2) *Map {
	return &Map{
		Position: position,
		Texture:  rl.LoadTexture(texturePath),
	}
}

func (t *Map) Draw() {
	//draw tile
	rl.DrawTexture(t.Texture, int32(t.Position.X), int32(t.Position.Y), rl.White)
	t.Update()
}

func (t *Map) Update() {
	//update tile

}

func init() {

	defer fmt.Printf("Game Initialized\n")

	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Lost in the Jungle")
	rl.SetExitKey(0)
	rl.SetTargetFPS(FPS)

	mainPlayer = NewPlayer("assets/Characters/player.png", rl.NewVector2(WINDOW_WIDTH/2, WINDOW_HEIGHT/2), 2, 4, 4, 2, 0)

	fmt.Printf("Player: %v\n", mainPlayer)

	camera = rl.NewCamera2D(rl.NewVector2(WINDOW_WIDTH/2, WINDOW_HEIGHT/2), rl.NewVector2(mainPlayer.Destination.X - (mainPlayer.Destination.Width / 2), mainPlayer.Destination.Y - (mainPlayer.Destination.Height / 2)), 0, 1)
	
	//grass tile
	grassTile = NewTile("assets/Tilesets/Grass.png", rl.NewVector2(0, 0))

	entities = append(entities, grassTile)

	rl.InitAudioDevice()
	bgMusic = rl.LoadMusicStream("assets/audios/bg-music.mp3")
	rl.PlayMusicStream(bgMusic)
}

func Debug() {
	//show FPS in the top right corner in colors red < 15, green < 30, blue < 45 and green = 60
	fps := float32(rl.GetFPS())
	fpsText := fmt.Sprintf("FPS: %v", fps)
	var color rl.Color
	if fps < FPS / 2 {
		color = rl.Red
	} else if fps < fps / 1.5 {
		color = rl.Orange
	} else if fps < fps / 1.2 {
		color = rl.Yellow
	} else if fps < FPS / 1.1 {
		color = rl.Blue
	} else {
		color = rl.Green
	}

	//length of the text
	textLen := rl.MeasureText(fpsText, 30)

	rl.DrawText(fpsText, WINDOW_WIDTH - textLen - 10, 10, 30, color)
}


func GameSettings() {
	//music on/off
	if rl.IsKeyPressed(rl.KeyM) {
		MusicOn = !MusicOn
		fmt.Printf("Music: %v\n", MusicOn)
	}

	//sound on/off
	if rl.IsKeyPressed(rl.KeyN) {
		SoundOn = !SoundOn
		fmt.Printf("Sound: %v\n", SoundOn)
	}
}


func GameLoop() {

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		
		rl.ClearBackground(rl.NewColor(175, 250, 202, 255))
		rl.BeginMode2D(camera)

		gameFrame++

		if gameFrame > FPS * 2 {
			gameFrame = 0
		}
	
		//Debug()
		GameSettings()

		rl.UpdateMusicStream(bgMusic)
		
		if MusicOn {
			rl.ResumeMusicStream(bgMusic)
		} else {
			rl.PauseMusicStream(bgMusic)
		}
	
		//draw entities
		for _, e := range entities {
			e.Draw()
		}

		//draw player
		mainPlayer.Draw()

		camera.Target = rl.NewVector2(mainPlayer.Destination.X - (mainPlayer.Destination.Width / 2), mainPlayer.Destination.Y - (mainPlayer.Destination.Height / 2))
		camera.Zoom = 4.0

		rl.EndMode2D()
		rl.EndDrawing()
	}
	
	defer rl.CloseWindow()
	
	//close entities
	Close(entities)
}

func Close(entities []Entity) {

	for _, e := range entities {
		//check type
		switch v := e.(type) {
		case *Player:
			//unload player texture
			rl.UnloadTexture(v.Texture)
		case *Map:
			//unload tile texture
			rl.UnloadTexture(v.Texture)
		}
	}

	rl.CloseAudioDevice()
	rl.UnloadMusicStream(bgMusic)

	fmt.Printf("Game Closed\n")
}

func main() {

	GameLoop()
}
