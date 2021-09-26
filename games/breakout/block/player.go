package block

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/krile136/testGame/funcs/input"
)

var thePlayer = &Player{320, 20, 5}

func Current() *Player {
	return thePlayer
}

type Player struct {
	y      int
	width  int
	height int
}

func (p *Player) Show(screen *ebiten.Image) {
	px, _ := input.Current().GetPosition()

	s, err := ebiten.NewShader([]byte(`package main
  var Color vec4
  func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
    return Color
  }`))
	if err != nil {
		return
	}

	op := &ebiten.DrawRectShaderOptions{}

	var cr, cb, cg, ca float32 = 1, 1, 1, 1
	op.Uniforms = map[string]interface{}{
		"Color": []float32{cr, cg, cb, ca},
	}

	op.GeoM.Translate(float64(px-p.width/2), float64(p.y-p.height/2))
	screen.DrawRectShader(p.width, p.height, s, op)
}

func (p *Player) IsReflect(moved_x float64, moved_y float64, r float64) bool {
	var reflect bool = false
	px, _ := input.Current().GetPosition()
	if moved_y+r > float64(p.y) {
		if float64(px-p.width/2) < moved_x && moved_x < float64(px+p.width/2) {
			reflect = true
		}
	}
	return reflect
}
