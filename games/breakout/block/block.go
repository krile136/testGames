package block

import "github.com/hajimehoshi/ebiten/v2"

// 変数名を大文字にしないとimportしたファイルからインポートできない
type Block struct {
	x      float64
	y      float64
	width  int
	height int
	hit    int
	show   bool
}

// Blockを初期化する
func NewBlock(x float64, y float64, width int, height int, hit int) (b *Block) {
	b = new(Block)
	b.x = x
	b.y = y
	b.width = width
	b.height = height
	b.hit = hit
	b.show = true

	return b
}

func (b *Block) AngleCoodinates() (float64, float64, float64, float64) {
	lux := b.x
	luy := b.y
	ldy := b.y + float64(b.height)
	rux := b.x + float64(b.width)

	return lux, luy, ldy, rux

}

// Blockを表示する
func (b *Block) Show(screen *ebiten.Image) {
	if b.show {
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

		var cr, cb, cg, ca float32
		switch b.hit {
		case 3:
			cr = 1
			cg = 0
			cb = 0
			ca = 1
		case 2:
			cr = 0
			cg = 1
			cb = 0
			ca = 1
		case 1:
			cr = 0
			cg = 0
			cb = 1
			ca = 1
		}
		op.Uniforms = map[string]interface{}{
			"Color": []float32{cr, cg, cb, ca},
		}
		op.GeoM.Translate(b.x, b.y)
		screen.DrawRectShader(b.width, b.height, s, op)
	}
}

func (b *Block) GetTouchedBorder(x float64, y float64, moved_x float64, moved_y float64, r float64) (bool, bool) {
	var left_or_right_touched bool = false
	var upper_or_down_touched bool = false

	if b.show {
		lux, luy, ldy, rux := b.AngleCoodinates()
		// ボールがブロックの範囲内に侵入したかどうか
		if moved_x+r > lux && moved_x-r < rux && moved_y+r > luy && moved_y-r < ldy {
			if lux < x && x < rux {
				// 移動前のボールの位置がブロックの横幅範囲内のとき、上下枠のどちらかに触れた
				upper_or_down_touched = true
			} else if lux < y && y < ldy {
				// 移動前のボールの位置がブロックの縦幅範囲のとき、左右枠のどちらかに触れた
				left_or_right_touched = true
			} else {
				// それ以外のとき、斜めから触れているので上下左右を反転
				upper_or_down_touched = true
				left_or_right_touched = true
			}
			// ボールの残ヒット数を減らし、もし0以下になったら非表示にする
			b.hit -= 1
			if b.hit <= 0 {
				b.show = false
			}
		}
	}
	return left_or_right_touched, upper_or_down_touched
}

// 半径rの円がブロックのどの線の内側にいるのかを計算する
func (b *Block) CalcBorderTouch(x float64, y float64, r float64) []string {
	lux := b.x
	luy := b.y
	ldy := b.y + float64(b.height)
	rux := b.x + float64(b.width)

	var result []string

	// 左の線と触れたか
	if calcCrossProduct(lux, ldy, lux, luy, x+r, y) > 0 {
		result = append(result, "left")
	}
	// 上の線と触れたか
	if calcCrossProduct(lux, luy, rux, luy, x, y+r) > 0 {
		result = append(result, "upper")
	}
	// 右の線と触れたか
	if calcCrossProduct(rux, luy, rux, ldy, x-r, y) > 0 {
		result = append(result, "right")
	}
	// 下の線と触れたか
	if calcCrossProduct(rux, ldy, lux, ldy, x, y-r) > 0 {
		result = append(result, "down")
	}

	return result
}

// 外積を使って線分ACが線分ABの左右どちらにいるかを計算する
// A->Bの向きに対して、A->Cが左の場合 -1、右の場合 1 を返す
func calcCrossProduct(Ax float64, Ay float64, Bx float64, By float64, Cx float64, Cy float64) int {
	ABx := Bx - Ax
	ABy := By - Ay
	ACx := Cx - Ax
	ACy := Cy - Ay

	result := ABx*ACy - ABy*ACx
	if result > 0 {
		return 1
	} else {
		return -1
	}
}
