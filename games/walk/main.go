package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/krile136/testGame/funcs/picture"

	"github.com/hajimehoshi/ebiten/v2/inpututil"

	// 行列を使うライブラリ
	"gonum.org/v1/gonum/mat"
)

const (
	screenWidth  = 240
	screenHeight = 240
)

const (
	// タイルのサイズは16px(正方形)
	tileSize = 16
	// マップチップの最大の横の大きさは、25タイル分( 16px x 25 = 400px)
	tileXNum = 25
)

const (
	frameOX     = 0  // フレーム開始時点のX座標
	frameOY     = 0  // フレーム開始時点のY座標
	frameWidth  = 32 // 1フレームで表示する横幅
	frameHeight = 32 // 1フレームで表示する縦幅
	frameNum    = 4  // 表示させる画像の数

)

var (
	tilesImage *ebiten.Image
)

// 基準速度
var basic_velocity float64 = 2
var velocity_array = []float64{basic_velocity, 0}        // 速度ベクトル生成用の行列を定義
var velocity_vector = mat.NewDense(2, 1, velocity_array) // 速度ベクトル（２行１列の行列）を生成　X軸の正方向にbasic_velocityの大きさを持つベクトル

// 移動方向回転行列 デフォルトでは下向き
var basic_angle float64 = -math.Pi / 2
var basic_posture_array = []float64{math.Cos(basic_angle), -math.Sin(basic_angle), math.Sin(basic_angle), math.Cos(basic_angle)}
var posture_rotate_matrix = mat.NewDense(2, 2, basic_posture_array) // 移動方向回転行列（2行2列）を生成

// キャラクターの初期位置
var positionX float64 = 100
var positionY float64 = 100

// キャラクターの向き  0:下 1:左 2:右 3:上　　これは読み込む画像による
var character_direction int = 0

// 斜めの時キャラクターの向きを決める変数 0:下 1:左 2:右 3:上 4:未入力
var pushed_arrow_key_num = 4
var released_arrow_key_num = 0
var pushed_arrow_key ebiten.Key

// 向きを保存した配列
var arrow_key_array [4]ebiten.Key = [4]ebiten.Key{ebiten.KeyDown, ebiten.KeyLeft, ebiten.KeyRight, ebiten.KeyUp}

//　キャラクター用の変数宣言
var character *ebiten.Image

// 最初に画像を読み込む
func init() {

	var err error
	tilesImage, _, err = ebitenutil.NewImageFromFile("tiles.png")

	character, _, err = ebitenutil.NewImageFromFile("character.png")

	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	// 毎フレーム毎にgを増加させていく
	count  int
	layers [][]int
}

func (g *Game) Update() error {
	// 毎フレーム毎にgを増加させていく
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	const xNum = screenWidth / tileSize

	for _, l := range g.layers {
		for i, t := range l {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(1.3, 1.3)
			op.GeoM.Translate(float64((i%xNum)*tileSize)*1.3, float64((i/xNum)*tileSize)*1.3)

			sx := (t % tileXNum) * tileSize
			sy := (t / tileXNum) * tileSize
			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}

	// 斜め移動時の向きの優先度を保存
	for i, v := range arrow_key_array {
		if pressedKey(v) {
			savePressedArrowKey(v, i)
		}
	}

	// どの方向を向くか決定する
	if pushed_arrow_key_num != 4 {
		character_direction = pushed_arrow_key_num
	} else {
		character_direction = released_arrow_key_num
	}

	// 優先度が高いキーが離されたかを監視
	checkReleasedArrowKey(pushed_arrow_key, pushed_arrow_key_num)

	// 基準速度を現在の速度に代入
	current_velocity := basic_velocity

	// 入力したキーで移動方向の角度と速度を決定
	if pressedKey(ebiten.KeyArrowUp) {
		if pressedKey(ebiten.KeyArrowRight) {
			// 右上の時 -45度
			basic_angle = -math.Pi / 4
		} else if pressedKey(ebiten.KeyArrowLeft) {
			// 左上の時 -135度
			basic_angle = -math.Pi / 4 * 3
		} else {
			// 上の時 -90度
			basic_angle = -math.Pi / 2
		}
		if pressedKey(ebiten.KeyArrowDown) {
			// 上下キーを同時押しした時、その場にとどまらせる
			current_velocity = 0
		}
	} else if pressedKey(ebiten.KeyArrowDown) {
		if pressedKey(ebiten.KeyArrowRight) {
			// 右下の時 45度
			basic_angle = math.Pi / 4
		} else if pressedKey(ebiten.KeyArrowLeft) {
			// 左下の時 135度
			basic_angle = math.Pi / 4 * 3
		} else {
			// 下の時 90度
			basic_angle = math.Pi / 2
		}
	} else if pressedKey(ebiten.KeyArrowRight) {
		// 右の時 0度
		basic_angle = 0
		if pressedKey(ebiten.KeyArrowLeft) {
			// 左右キーを同時押しした時、その場に留まらせる
			current_velocity = 0
		}
	} else if pressedKey(ebiten.KeyArrowLeft) {
		// 左の時 180度
		basic_angle = math.Pi
	} else {
		// 十字キーの入力がない時に、現在の速度を0にする
		current_velocity = 0
	}

	// 現在の入力に合わせた移動方向回転行列を更新
	posture_rotate_matrix.Set(0, 0, math.Cos(basic_angle))
	posture_rotate_matrix.Set(0, 1, -math.Sin(basic_angle))
	posture_rotate_matrix.Set(1, 0, math.Sin(basic_angle))
	posture_rotate_matrix.Set(1, 1, math.Cos(basic_angle))

	// 速度ベクトルの成分を更新
	velocity_vector.Set(0, 0, current_velocity)

	// 移動ベクトルを移動方向回転行列と速度ベクトルから計算
	move_vector := mat.NewDense(2, 1, nil)
	move_vector.Product(posture_rotate_matrix, velocity_vector)

	// 13フレームに一度画像を更新する
	i := (g.count / 13) % frameNum

	// i が3 すなわち右手→ニュートラル→左手→ここ　の時にニュートラルを表示させる
	if i == 3 {
		i = 1
	}
	// 現在の速度が0の時、ニュートラルの姿勢にする
	if current_velocity == 0 {
		i = 1
	}

	sx, sy := frameOX+i*frameWidth, frameOY+32*character_direction

	// 表示する画像を切り出す
	var this_frame_img *ebiten.Image = character.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image)

	// キャラクターの表示位置をmove_vectorの分だけ移動させる
	positionX += move_vector.At(0, 0)
	positionY += move_vector.At(1, 0)

	// 画像を表示させる
	picture.Show(screen, this_frame_img, 1, positionX, positionY, 0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth * 1.3, screenHeight * 1.3
}

func main() {
	g := &Game{
		layers: [][]int{
			{
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 244, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 219, 243, 243, 243, 219, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 26, 27, 28, 29, 30, 31, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 51, 52, 53, 54, 55, 56, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 76, 77, 78, 79, 80, 81, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 101, 102, 103, 104, 105, 106, 0, 0, 0, 0,

				0, 0, 0, 0, 0, 126, 127, 128, 129, 130, 131, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 303, 303, 245, 242, 303, 303, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,

				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			},
		},
	}

	// tileとキャラクターの大きさを合わせるためにtilesを1.3倍している
	// windowSizeを1.5倍してしまうとでかいので1.5倍にしている
	ebiten.SetWindowSize(screenWidth*1.3*1.5, screenHeight*1.3*1.5)
	ebiten.SetWindowTitle("Tiles (Ebiten Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// 行列を標準出力する
func matPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n", fa)
}

// 特定のキーが押されているかをチェックする
func pressedKey(str ebiten.Key) bool {
	inputArray := inpututil.PressedKeys()
	for _, v := range inputArray {
		if v == str {
			return true
		}
	}
	return false
}

// 優先的に押されたキーを保持する
func savePressedArrowKey(key ebiten.Key, key_num int) {
	if pushed_arrow_key_num == 4 {
		pushed_arrow_key = key
		pushed_arrow_key_num = key_num
	}
}

// 優先的に押されていたキーが離されたかどうかをチェックする
func checkReleasedArrowKey(key ebiten.Key, key_num int) {
	if inpututil.IsKeyJustReleased(key) {
		pushed_arrow_key_num = 4
		released_arrow_key_num = key_num
	}
}
