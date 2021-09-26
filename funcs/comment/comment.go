package comment 

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func Print(screen *ebiten.Image, text string) {
	ebitenutil.DebugPrint(screen, text)
}

