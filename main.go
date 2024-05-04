package main

import (
	"encoding/json"
	"fmt"
	"os"

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
	FPS           = 60.00
)

var (
	mainPlayer *Player

	entities []Drawable

	//game frame
	gameFrame int32 = 0
	camera    rl.Camera2D

	//bg music
	bgMusic rl.Music

	//settings
	MusicOn bool = false
	SoundOn bool = true
)

type Drawable interface {
	Update()
}

type Player struct {
	Speed          float32
	_OriginalSpeed float32
	Direction      direction
	CurrentFrame   int32
	CurrentRow     int32
	SpriteRowCount int32
	SpriteColCount int32
	FrameWidth     int32
	FrameHeight    int32
	FrameSpeed     int32

	Source      rl.Rectangle
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

	source := rl.NewRectangle(float32(currentFrame*frameWidth), float32(currentRow*frameHeight), float32(frameWidth), float32(frameHeight))
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

		Source:      source,
		Destination: destination,

		Texture: texture,
		Running: false,
	}
}

func (p *Player) Draw() {
	//draw player

	//draw player border
	x := int32(p.Destination.X - p.Source.Width/2)
	y := int32(p.Destination.Y - p.Source.Height/2)
	//rl.DrawRectangleLines(x, y, p.Destination.ToInt32().Width, p.Destination.ToInt32().Height, rl.Green)
	//Debug(x, y)
	rl.DrawTexturePro(p.Texture, p.Source, p.Destination, rl.NewVector2(p.Destination.Width, p.Destination.Height), 0, rl.White)
	rl.DrawCircleLines(x, y, p.Destination.Width/3, rl.Green)
}

func DebugPlayer(p *Player) {

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
}

func (p *Player) Update() {

	p.HandleControls()

	p.Draw()

	if p.CurrentFrame >= p.SpriteColCount {
		p.CurrentFrame = 0
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

	} else if gameFrame % 20 == 1 { // change charecter frame every 45 frames of the game
		p.CurrentFrame++
	}

	//Idle animation
	if !p.Running && p.CurrentFrame > 1 {
		p.CurrentFrame = 0
	}

	p.Source.X = float32(p.CurrentFrame * p.FrameWidth)
	p.Source.Y = float32(p.CurrentRow * p.FrameHeight)
}


//source of the texture
// Tileset represents a tileset in the tilemap JSON
type Tileset struct {
	Columns     int    `json:"columns"`
	FirstGID    int    `json:"firstgid"`
	Image       string `json:"image"`
	ImageHeight int    `json:"imageheight"`
	ImageWidth  int    `json:"imagewidth"`
	Margin      int    `json:"margin"`
	Name        string `json:"name"`
	Class 		string `json:"class"`
	Spacing     int    `json:"spacing"`
	TileCount   int    `json:"tilecount"`
	TileHeight  int    `json:"tileheight"`
	TileWidth   int    `json:"tilewidth"`
}

// Chunk represents a chunk of tiles in a tile layer
type Chunk struct {
	Data   []int `json:"data"`
	Height int   `json:"height"`
	Width  int   `json:"width"`
	X      int   `json:"x"`
	Y      int   `json:"y"`
}

// TileLayer represents a tile layer in the tilemap JSON
type TileLayer struct {
	Chunks  []Chunk `json:"chunks"`
	Height  int     `json:"height"`
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Class 	string  `json:"class"`
	Width   int     `json:"width"`
	X       int     `json:"x"`
	Y       int     `json:"y"`
	StartX  int     `json:"startx"`
	StartY  int     `json:"starty"`
	Visible bool    `json:"visible"`
	Opacity float32 `json:"opacity"`
}

func LoadTilesetTextures(tilesets []Tileset) map[string]rl.Texture2D {
	textures := make(map[string]rl.Texture2D)

	for _, tileset := range tilesets {
		texture := rl.LoadTexture(tileset.Image)
		textures[tileset.Class] = texture
	}

	return textures
}

type TileMap struct {
	Tilesets []Tileset   `json:"tilesets"`
	Layers   []TileLayer `json:"layers"`
	Textures map[string]rl.Texture2D
}

func NewTileMap(filePath string) *TileMap {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tilemap := TileMap{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tilemap); err != nil {
		panic(err)
	}

	tilemap.Textures = make(map[string]rl.Texture2D)
	tilemap.Textures = LoadTilesetTextures(tilemap.Tilesets)

	return &tilemap
}

/*
func Draw(layer TileLayer, textures map[int]rl.Texture2D, tileSize rl.Vector2) {
    for _, chunk := range layer.Chunks {
        for y := 0; y < chunk.Height; y++ {
            for x := 0; x < chunk.Width; x++ {
                tileIndex := chunk.Data[y*chunk.Width+x]
                if tileIndex != 0 { // Assuming 0 means no tile
                    texture := textures[tileIndex]
                    //dest := rl.NewVector2(float32(chunk.X+x)*tileSize.X, float32(chunk.Y+y)*tileSize.Y)
					sourceRect := rl.NewRectangle(float32(x), float32(y), tileSize.X, tileSize.Y)
					dest := rl.NewRectangle(float32(chunk.X+x)*tileSize.X, float32(chunk.Y+y)*tileSize.Y, tileSize.X, tileSize.Y)
                    //rl.DrawTextureRec(texture, rl.NewRectangle(float32(x), float32(y), tileSize.X, tileSize.Y), dest, rl.White)
					rl.DrawTexturePro(texture, sourceRect, dest, rl.NewVector2(float32(tileSize.X), float32(tileSize.Y)), 0, rl.White)
					//draw a solid rectangle to show the tile
					rl.DrawRectangleLines(int32(dest.X), int32(dest.Y), int32(tileSize.X), int32(tileSize.Y), rl.Red)
                }
            }
        }
    }
}
*/

func (t *TileMap) Draw() {
	// For every layer

	for _, layer := range t.Layers {
		if !layer.Visible {
			continue
		}
		t.drawLayers(&layer)
	}
}

func (t *TileMap) drawLayers(layer *TileLayer) {
	// If the layer is not visible, skip it.
	texture := t.Textures[layer.Class]

	// For every chunk
	for _, chunk := range layer.Chunks {
		t.drawChunk(texture, &chunk)
	}
}

func (t *TileMap) drawChunk(texture rl.Texture2D, chunk *Chunk) {
	//chunk has data in the form of a 1D array
	for y := 0; y < chunk.Height; y++ {
		for x := 0; x < chunk.Width; x++ {
			tileIndex := chunk.Data[y*chunk.Width+x]
			if tileIndex != 0 { // Assuming 0 means no tile
				//texture := textures[tileIndex]
				//dest := rl.NewVector2(float32(chunk.X+x)*tileSize.X, float32(chunk.Y+y)*tileSize.Y)
				sourceRect := rl.NewRectangle(float32(x), float32(y), float32(t.Tilesets[0].TileWidth), float32(t.Tilesets[0].TileHeight))
				dest := rl.NewRectangle(float32(chunk.X+x)*float32(t.Tilesets[0].TileWidth), float32(chunk.Y+y)*float32(t.Tilesets[0].TileHeight), float32(t.Tilesets[0].TileWidth), float32(t.Tilesets[0].TileHeight))
				//rl.DrawTextureRec(texture, rl.NewRectangle(float32(x), float32(y), tileSize.X, tileSize.Y), dest, rl.White)
				rl.DrawTexturePro(texture, sourceRect, dest, rl.NewVector2(float32(t.Tilesets[0].TileWidth), float32(t.Tilesets[0].TileHeight)), 0, rl.White)
				//draw a solid rectangle to show the tile
				rl.DrawRectangleLines(int32(dest.X), int32(dest.Y), int32(t.Tilesets[0].TileWidth), int32(t.Tilesets[0].TileHeight), rl.Red)
				//rl.DrawRectangle(int32(dest.X), int32(dest.Y), int32(tileSize.X), int32(tileSize.Y), rl.White)
			}
		}
	}
}

func (t *TileMap) Update() {
	//
	t.Draw()
}

func init() {

	defer fmt.Printf("Game Initialized\n")

	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Lost in the Jungle")
	rl.SetExitKey(0)
	rl.SetTargetFPS(FPS)

	mainPlayer = NewPlayer("assets/Characters/player.png", rl.NewVector2(WINDOW_WIDTH/2, WINDOW_HEIGHT/2), 2, 4, 4, 2, 0)

	fmt.Printf("Player: %v\n", mainPlayer)

	camera = rl.NewCamera2D(rl.NewVector2(WINDOW_WIDTH/2, WINDOW_HEIGHT/2), rl.NewVector2(mainPlayer.Destination.X-(mainPlayer.Destination.Width/2), mainPlayer.Destination.Y-(mainPlayer.Destination.Height/2)), 0, 1)

	//load tilesets
	tileMap := NewTileMap("maps/tiles2.json")

	entities = append(entities, tileMap)

	// run entities

	rl.InitAudioDevice()
	bgMusic = rl.LoadMusicStream("assets/audios/bg-music.mp3")
	rl.PlayMusicStream(bgMusic)
}

func ProcessInput() {
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

func UpdateMusic() {
	rl.UpdateMusicStream(bgMusic)

	if MusicOn {
		rl.ResumeMusicStream(bgMusic)
	} else {
		rl.PauseMusicStream(bgMusic)
	}
}

func UpdateEntities() {
	for _, e := range entities {
		e.Update()
	}
}

func UpdateFrames() {
	gameFrame++
	if gameFrame > FPS*2 {
		gameFrame = 0
	}
}

func UpdateCamera() {
	camera.Target = rl.NewVector2(mainPlayer.Destination.X-(mainPlayer.Destination.Width/2), mainPlayer.Destination.Y-(mainPlayer.Destination.Height/2))
	//camera.Zoom = 2
}

func Debug(x, y int32) {
	//show FPS in the top right corner in colors red < 15, green < 30, blue < 45 and green = 60
	fps := float32(rl.GetFPS())
	fpsText := fmt.Sprintf("FPS: %v", fps)
	var color rl.Color
	if fps < FPS/2 {
		color = rl.Red
	} else if fps < fps/1.5 {
		color = rl.Orange
	} else if fps < fps/1.2 {
		color = rl.Yellow
	} else if fps < FPS/1.1 {
		color = rl.Blue
	} else {
		color = rl.Green
	}

	//length of the text
	textLen := rl.MeasureText(fpsText, 30)

	rl.DrawText(fpsText, WINDOW_WIDTH-textLen-x, y, 30, color)
}

func GameLoop() {

	for !rl.WindowShouldClose() {

		//update entities

		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(175, 250, 202, 255))
		rl.BeginMode2D(camera)

		UpdateFrames()

		ProcessInput()
		UpdateMusic()

		UpdateEntities()
		mainPlayer.Update()
		UpdateCamera()

		rl.EndMode2D()
		rl.EndDrawing()
	}

	defer rl.CloseWindow()
	Close(entities)
}

func Close(entities []Drawable) {

	for _, e := range entities {
		//check type
		switch v := e.(type) {
		case *Player:
			//unload player texture
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
