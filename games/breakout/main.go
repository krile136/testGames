package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/krile136/testGame/funcs/input"
	"github.com/krile136/testGame/games/breakout/block"

	// 行列を使うライブラリ
	"gonum.org/v1/gonum/mat"
)

const (
	screenWidth  = 240
	screenHeight = 360
	ngon         = 25
	radius       = 3
)

// 空の画像を定義
var (
	emptyImage     = ebiten.NewImage(3, 3)
	vertices       []ebiten.Vertex
	vertices_mouse []ebiten.Vertex
	ballCenterX    float64 = screenWidth / 2
	ballCenterY    float64 = screenHeight / 2
	velAngle       float64 = math.Pi / 4
	velocity       float64 = 5
)

var blks Blocks

type Blocks []*block.Block

func init() {
	// 空の画像の色を白で初期化
	emptyImage.Fill(color.White)

	// 頂点数を初期化
	vertices = genVertices(ngon, radius, float32(ballCenterX), float32(ballCenterY))

	// ブロックを生成
	// blks = append(blks, block.NewBlock(screenWidth/2, screenHeight/2, 16, 16, 1))
	// blks = append(blks, block.NewBlock(screenWidth/2+17, screenHeight/2, 16, 16, 2))
	// blks = append(blks, block.NewBlock(screenWidth/2+32+2, screenHeight/2, 16, 16, 3))
	// blks = append(blks, block.NewBlock(0, 0, 16, 16, 3))
	blks = append(blks, block.NewBlock(screenWidth/2+20, screenHeight/2-20, 50, 50, 3))
}

// 頂点を生成する関数（サンプルより拝借）
func genVertices(num int, rad float64, x float32, y float32) []ebiten.Vertex {
	var (
		r       = rad
		centerX = x
		centerY = y
	)

	vs := []ebiten.Vertex{}
	for i := 0; i < num; i++ {
		rate := float64(i) / float64(num)
		cr := 0.0
		cg := 0.0
		cb := 0.0
		if rate < 1.0/3.0 {
			cb = 2 - 2*(rate*3)
			cr = 2 * (rate * 3)
		}
		if 1.0/3.0 <= rate && rate < 2.0/3.0 {
			cr = 2 - 2*(rate-1.0/3.0)*3
			cg = 2 * (rate - 1.0/3.0) * 3
		}
		if 2.0/3.0 <= rate {
			cg = 2 - 2*(rate-2.0/3.0)*3
			cb = 2 * (rate - 2.0/3.0) * 3
		}
		vs = append(vs, ebiten.Vertex{
			DstX:   float32(r*math.Cos(2*math.Pi*rate)) + centerX,
			DstY:   float32(r*math.Sin(2*math.Pi*rate)) + centerY,
			SrcX:   0,
			SrcY:   0,
			ColorR: float32(cr),
			ColorG: float32(cg),
			ColorB: float32(cb),
			ColorA: 1,
		})
	}

	vs = append(vs, ebiten.Vertex{
		DstX:   centerX,
		DstY:   centerY,
		SrcX:   0,
		SrcY:   0,
		ColorR: 1,
		ColorG: 1,
		ColorB: 1,
		ColorA: 1,
	})

	return vs
}

type Game struct{}

func (g *Game) Update() error {
	px, py := input.Current().GetPosition()
	vertices_mouse = genVertices(ngon, radius, float32(px), float32(py))
	// str := blks[0].CalcBorderTouch(float64(px), float64(py), radius)
	// log.Print(str)
	// 進行方向の角度による回転行列を生成
	basicPostureArray := []float64{math.Cos(velAngle), -math.Sin(velAngle), math.Sin(velAngle), math.Cos(velAngle)}
	postureRotateMatrix := mat.NewDense(2, 2, basicPostureArray)

	// 速度ベクトルを生成(Y軸方向に正の方向がvelocityとなるようなベクトル)
	basicVelocityArray := []float64{0, -velocity}
	velocityVector := mat.NewDense(2, 1, basicVelocityArray)

	// 移動ベクトルを生成
	moveVector := mat.NewDense(2, 1, nil)
	moveVector.Product(postureRotateMatrix, velocityVector)

	var xReverse, yReverse int = 1, 1
	// 左右の端に当たったときにボールを反転させる処理
	movedBallCenterX := ballCenterX + moveVector.At(0, 0)
	if movedBallCenterX-radius < 0 || movedBallCenterX+radius > screenWidth {
		xReverse = -1
	}
	// 上下の端に当たったときにボールを反転させる処理
	movedBallCenterY := ballCenterY + moveVector.At(1, 0)
	if movedBallCenterY-radius < 0 || movedBallCenterY+radius > screenHeight {
		yReverse = -1
	}

	// ブロックにぶつかったときの処理
	for _, v := range blks {
		isTouch_left_or_right, isTouch_upper_or_down := v.GetTouchedBorder(ballCenterX, ballCenterY, movedBallCenterX, movedBallCenterY, radius)
		if isTouch_left_or_right {
			xReverse = -1
		}
		if isTouch_upper_or_down {
			yReverse = -1
		}

	}

	// プレイヤーで反射するかの処理
	reflect := block.Current().IsReflect(movedBallCenterX, movedBallCenterY, radius)
	if reflect {
		yReverse = -1
	}

	// 最終的な次のボール位置と角度の計算
	ballCenterX += moveVector.At(0, 0) * float64(xReverse)
	ballCenterY += moveVector.At(1, 0) * float64(yReverse)
	if xReverse < 0 {
		velAngle = math.Pi*2 - velAngle
	}
	if yReverse < 0 {
		velAngle = math.Pi - velAngle
	}

	// ボールの描画位置を計算
	vertices = genVertices(ngon, radius, float32(ballCenterX), float32(ballCenterY))

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// プレイヤーのブロックを表示させる
	block.Current().Show(screen)

	// ボールを表示させる
	op := &ebiten.DrawTrianglesOptions{}
	op.Address = ebiten.AddressUnsafe
	indices := []uint16{}
	for i := 0; i < ngon; i++ {
		indices = append(indices, uint16(i), uint16(i+1)%uint16(ngon), uint16(ngon))
	}
	screen.DrawTriangles(vertices, indices, emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)

	indices_mouse := []uint16{}
	for i := 0; i < ngon; i++ {
		indices_mouse = append(indices_mouse, uint16(i), uint16(i+1)%uint16(ngon), uint16(ngon))
	}
	// screen.DrawTriangles(vertices_mouse, indices_mouse, emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)

	// 現在のTPSを表示させる
	msg := fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)

	// ブロックを表示させる
	for _, v := range blks {
		v.Show(screen)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("ball")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

// 行列を標準出力する
func matPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n", fa)
	fmt.Printf("----------------------\n")
}
