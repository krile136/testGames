package block

import "github.com/hajimehoshi/ebiten/v2"

// 変数名を大文字にしないとimportしたファイルからインポートできない
type Block struct {
	X      float64
	Y      float64
	Width  int
	Height int
	Cr     float32
	Cg     float32
	Cb     float32
	Ca     float32
	Hit    int
}

// Blockを初期化する
func NewBlock(x float64, y float64, width int, height int, cr float32, cg float32, cb float32, ca float32, hit int) (b *Block) {
	b = new(Block)
	b.X = x
	b.Y = y
	b.Width = width
	b.Height = height
	b.Cr = cr
	b.Cg = cg
	b.Cb = cb
	b.Ca = ca
	b.Hit = hit

	return b
}

// Blockを表示する
func (b *Block) Show(screen *ebiten.Image) {
	s, err := ebiten.NewShader([]byte(`package main
	var Color vec4
	func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
		return Color
	}`))
	if err != nil {
		return
	}

	op := &ebiten.DrawRectShaderOptions{}

	op.CompositeMode = ebiten.CompositeModeCopy
	op.Uniforms = map[string]interface{}{
		"Color": []float32{b.Cr, b.Cg, b.Cb, b.Cr},
	}
	op.GeoM.Translate(b.X, b.Y)
	screen.DrawRectShader(b.Width, b.Height, s, op)
}
