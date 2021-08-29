package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

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
	emptyImage  = ebiten.NewImage(3, 3)
	vertices    []ebiten.Vertex
	ballCenterX float64 = screenWidth / 2
	ballCenterY float64 = screenHeight / 2
	velAngle    float64 = math.Pi / 4
	velocity    float64 = 5
)

func init() {
	// 空の画像の色を白で初期化
	emptyImage.Fill(color.White)

	// 頂点数を初期化
	vertices = genVertices(ngon)
}

// 頂点を生成する関数（サンプルより拝借）
func genVertices(num int) []ebiten.Vertex {
	const (
		r = radius
	)

	var (
		centerX = float32(ballCenterX)
		centerY = float32(ballCenterY)
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
	// 進行方向の角度による回転行列を生成
	basicPostureArray := []float64{math.Cos(velAngle), -math.Sin(velAngle), math.Sin(velAngle), math.Cos(velAngle)}
	postureRotateMatrix := mat.NewDense(2, 2, basicPostureArray)

	// 速度ベクトルを生成(Y軸方向に正の方向がvelocityとなるようなベクトル)
	basicVelocityArray := []float64{0, -velocity}
	velocityVector := mat.NewDense(2, 1, basicVelocityArray)

	// 移動ベクトルを生成
	moveVector := mat.NewDense(2, 1, nil)
	moveVector.Product(postureRotateMatrix, velocityVector)

	// matPrint(moveVector)

	// ボールの中心位置を移動させる
	// 壁の橋にあたったときにX,Yそれぞれ角度を反転させる
	prevBallCenterX := ballCenterX + moveVector.At(0, 0)
	if prevBallCenterX-radius < 0 || prevBallCenterX+radius > screenWidth {
		ballCenterX -= moveVector.At(0, 0)
		velAngle = math.Pi*2 - velAngle
	} else {
		ballCenterX += moveVector.At(0, 0)
	}

	prevBallCenterY := ballCenterY + moveVector.At(1, 0)
	if prevBallCenterY-radius < 0 || prevBallCenterY+radius > screenHeight {
		ballCenterY -= moveVector.At(1, 0)
		velAngle = math.Pi - velAngle
	} else {
		ballCenterY += moveVector.At(1, 0)
	}

	// ボールの描画位置を計算
	vertices = genVertices(ngon)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// ボールを表示させる
	op := &ebiten.DrawTrianglesOptions{}
	op.Address = ebiten.AddressUnsafe
	indices := []uint16{}
	for i := 0; i < ngon; i++ {
		indices = append(indices, uint16(i), uint16(i+1)%uint16(ngon), uint16(ngon))
	}
	screen.DrawTriangles(vertices, indices, emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)

	// 現在のTPSを表示させる
	msg := fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)
	// msg := fmt.Sprintf("centerX %0.2f", ballCenterX)
	// ebitenutil.DebugPrint(screen, msg)
	// fmt.Printf("%f\n", ballCenterX)
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
