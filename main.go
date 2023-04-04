package main

import (
	"github.com/ares0516/snake/pkg/component"
	"github.com/ares0516/snake/pkg/define"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
	"sync"
	"time"
)

type GreedySnake struct {
	screenWidth  int
	screenHeight int

	snake *component.Square

	mutex   sync.RWMutex
	running bool
}

func NewGreedySnake() *GreedySnake {
	return &GreedySnake{
		screenWidth:  640,
		screenHeight: 480,
	}
}

func (g *GreedySnake) IsRunning() bool {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.running
}

func (g *GreedySnake) SetRunning(running bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.running = running
}

// Layout 设置游戏窗口的大小
func (g *GreedySnake) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

// Update 更新游戏的逻辑 默认每秒60次
func (g *GreedySnake) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) { // 按下空格键开始游戏
		if !g.IsRunning() {
			g.SetRunning(true)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.snake.SetDirection(define.LEFT)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.snake.SetDirection(define.RIGHT)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.snake.SetDirection(define.UP)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.snake.SetDirection(define.DOWN)
	}

	return nil
}

// Draw 在屏幕上绘制游戏内容
func (g *GreedySnake) Draw(screen *ebiten.Image) {
	// 绘制黑色背景
	screen.Fill(color.RGBA{0, 0, 0, 255})
	// 绘制图像
	screen.DrawImage(g.snake.Image, g.snake.Opts)
}

func (g *GreedySnake) SnakeMove() {
	for g.IsRunning() { // 如果游戏开始，蛇开始移动
		log.Printf("snake move")
		g.snake.Move()
		time.After(1 * time.Second) // 每隔1秒移动一次
	}
}

// XUpdate 更新游戏的逻辑,通过外部协程来更新snake的移动
// TODO: 根据游戏等级来控制蛇的移动速度
func XUpdate(game *GreedySnake) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			game.snake.Move()
		default:
			// do nothing
		}
	}
}

func main() {
	// 1. 初始化游戏
	game := NewGreedySnake()
	game.snake = component.NewSquare(define.Green, 5, 5, 320, 240, 5)
	go XUpdate(game)
	// 2. 设置游戏窗口大小
	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	// 3. 设置游戏窗口标题
	ebiten.SetWindowTitle("贪吃蛇")
	// 4. 设置游戏运行时的更新函数和绘制函数
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
