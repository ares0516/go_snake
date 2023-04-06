package component

import (
	"github.com/ares0516/snake/pkg/define"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

type Square struct {
	bgc   color.RGBA
	h     float64 //	高度
	w     float64 //	宽度
	x     float64 //	锚点坐标x
	y     float64 //	锚点坐标y
	step  float64
	stepX float64
	stepY float64
	Image *ebiten.Image
	Opts  *ebiten.DrawImageOptions
}

func NewSquare(bgc color.RGBA, h, w int, x, y, step float64) *Square {
	image := ebiten.NewImage(w, h)
	image.Fill(bgc)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)
	return &Square{
		bgc:   bgc,
		h:     float64(h),
		w:     float64(w),
		x:     x,
		y:     y,
		step:  step,
		stepX: step,
		stepY: 0,
		Image: image,
		Opts:  opts,
	}
}

func NewSquareWithImage(image *ebiten.Image, x, y, step float64) *Square {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)
	w, h := image.Bounds().Dx(), image.Bounds().Dy()
	return &Square{
		h:     float64(h),
		w:     float64(w),
		x:     x,
		y:     y,
		step:  step,
		stepX: step,
		stepY: step,
		Image: image,
		Opts:  opts,
	}
}

// Move 移动到新坐标并返回之前的坐标
func (s *Square) Move() define.Position {
	log.Printf("current squre x[%f]y[%f]  stepX[%f]stepY[%f]", s.x, s.y, s.stepX, s.stepY)
	s.x += s.stepX
	s.y += s.stepY
	s.Opts.GeoM.Translate(s.stepX, s.stepY)
	return define.Position{X: s.x, Y: s.y}
}

func (s *Square) SetDirection(dir define.Direction) {
	log.Printf("----------------------dir[%v]", dir)
	// 1. 不需要转向的场景
	if (s.stepX == 0 && dir.X == 0) || (s.stepY == 0 && dir.Y == 0) {
		return
	}
	// 2. 需要转向的场景
	log.Printf("000000-----stepX[%f]stepY[%f]", s.stepX, s.stepY)
	s.stepX = dir.X * s.step
	s.stepY = dir.Y * s.step
	log.Printf("111111-----stepX[%f]stepY[%f]", s.stepX, s.stepY)
	return
}

func (s *Square) Transparent(w, h float64) {
	x, y := s.x+s.stepX, s.y+s.stepY // 小球下一步位置坐标
	if x <= 0 {
		s.x = w - s.w
	} else if x+s.w >= w {
		s.x = 0
	} else if y <= 0 {
		s.y = h - s.h
	} else if y+s.h >= h { // 小球落地
		s.y = 0
	}
	s.Opts.GeoM.Translate(s.x, s.y)
}
