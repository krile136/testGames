package input

import (
	"fmt"
	_ "log"

	"github.com/hajimehoshi/ebiten/v2"
)

// "github.com/hajimehoshi/ebiten/v2/inpututil"

// Inputの構造体をtheInputに格納
var theInput = &Input{}

// CurrentでtheInputを返すことで、importしたファイルからinput.Current().hogehoge() のように
// Input構造体が持つ関数を呼ぶことができる
func Current() *Input {
	return theInput
}

// Inputの構造体を宣言
type Input struct {
}

// タッチ可能ならmobile、そうでないならPCを返す
func (i *Input) Platform() string {
	if isTouchEnabled() {
		return "mobile"
	}
	return "PC"
}

// ウィンドウ上のカーソル位置を返す
func (i *Input) GetPosition() (x, y int) {
	px, py := ebiten.CursorPosition()
	fmt.Printf("x=%v y=%v\n", px, py)
	return px, py
}
