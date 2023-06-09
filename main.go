package main

import (
	"fmt"
	"github.com/ares0516/snake/pkg/component"
	"github.com/ares0516/snake/pkg/define"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math/rand"
	"sync"
	"time"
)

type GreedySnake struct {
	screenWidth  int
	screenHeight int

	snakeHead    *component.Square
	snakeBodyLen int                 // 蛇的长度
	snakeBody    []*component.Square // 蛇的身体

	awards []*component.Square // 奖励

	mutex   sync.RWMutex
	running bool
}

func NewGreedySnake() *GreedySnake {
	return &GreedySnake{
		screenWidth:  640,
		screenHeight: 480,
		snakeBodyLen: 10, // 初始长度为3
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

	if g.snakeHead.Eat(&g.awards) {
		g.snakeBodyLen++
	}

	if g.IsRunning() {
		g.snakeHead.Transparent(float64(g.screenWidth), float64(g.screenHeight))
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.snakeHead.SetDirection(define.LEFT)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.snakeHead.SetDirection(define.RIGHT)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.snakeHead.SetDirection(define.UP)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.snakeHead.SetDirection(define.DOWN)
	}

	return nil
}

// Draw 在屏幕上绘制游戏内容
func (g *GreedySnake) Draw(screen *ebiten.Image) {
	// 绘制黑色背景
	screen.Fill(color.RGBA{0, 0, 0, 255})
	// 绘制蛇的头
	screen.DrawImage(g.snakeHead.Image, g.snakeHead.Opts)
	// 绘制蛇的身体
	for _, body := range g.snakeBody {
		screen.DrawImage(body.Image, body.Opts)
	}

	// 绘制奖励
	for _, award := range g.awards {
		screen.DrawImage(award.Image, award.Opts)
	}

	// 绘制计分牌
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.snakeBodyLen))
}

func (g *GreedySnake) AwardGenerator() {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		if g.IsRunning() {
			select {
			case <-ticker.C:
				//log.Printf("award len: %d", len(g.awards))

				if len(g.awards) < 5 {
					g.awards = append(g.awards, component.NewSquare(define.Yellow, 5, 5, float64(rand.Intn(300)+10), float64(rand.Intn(200)+10), 0))
				}
				if len(g.awards) > 5 {
					g.awards = g.awards[1:5]
				}
			default:
				// do nothing
			}
		}
	}
}

func (g *GreedySnake) BodyGenerator(dir define.Position) {
	if g.IsRunning() {
		g.snakeBody = append(g.snakeBody, component.NewSquare(define.Green, 5, 5, dir.X, dir.Y, 5))
		if len(g.snakeBody) > g.snakeBodyLen {
			g.snakeBody = g.snakeBody[1:]
		}
	}
}

func (g *GreedySnake) GetLevel() int {
	switch g.snakeBodyLen {
	case 10:
		return 1
	case 20:
		return 2
	case 30:
		return 3
	case 40:
		return 4
	case 50:
		return 5
	default:
		return 5
	}
}

// XUpdate 更新游戏的逻辑,通过外部协程来更新snake的移动
// TODO: 1.根据游戏等级来控制蛇的移动速度
//
//	2.根据分数绘制蛇的身体
func XUpdate(game *GreedySnake) {
	ticker := time.NewTicker(200 * time.Millisecond)
	pos := define.Position{}
	for {
		if game.IsRunning() {
			select {
			case <-ticker.C:
				pos = game.snakeHead.Move()
				game.BodyGenerator(pos)
				td := time.Duration(600-game.GetLevel()*100) * time.Millisecond
				ticker.Reset(td)
			default:
				// do nothing
			}
		}
	}
}

func main() {
	// 1. 初始化游戏
	game := NewGreedySnake()
	game.snakeHead = component.NewSquare(define.Green, 5, 5, 320, 240, 5)
	// 初始化蛇的身体，优先生成尾部
	//game.snakeBody = append(game.snakeBody, component.NewSquare(define.Green, 5, 5, 305, 240, 5))
	//game.snakeBody = append(game.snakeBody, component.NewSquare(define.Green, 5, 5, 310, 240, 5))
	//game.snakeBody = append(game.snakeBody, component.NewSquare(define.Green, 5, 5, 315, 240, 5))
	go XUpdate(game)
	go game.AwardGenerator()
	// 2. 设置游戏窗口大小
	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	// 3. 设置游戏窗口标题
	ebiten.SetWindowTitle("贪吃蛇")
	// 4. 设置游戏运行时的更新函数和绘制函数
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
