package picture 

import (
	"math"
	"github.com/hajimehoshi/ebiten/v2"
)

/*
引数解説
      screen : *ebiten.Image  描画するスクリーン
         img : *ebiten.Image  描画したい画像
 coefficient : float64  画像の縮小/拡大の係数
x_coodinates : float64  画像のX座標
y_coodinates : float64  画像のY座標
       angle : float64  画像の回転角度（degree 度 )
*/
func Show(screen *ebiten.Image, img *ebiten.Image, coefficient float64, x_coodinates float64, y_coodinates float64, angle float64) {
	// 画像のサイズを取得
	w, h := img.Size()

	// 係数で画像を拡大/縮小したときの大きさを計算しておく
	var sw, sh float64 = float64(w) * coefficient, float64(h) * coefficient

	// オプションを宣言
	op := &ebiten.DrawImageOptions{}

	// 画像を拡大/縮小する
	op.GeoM.Scale(coefficient, coefficient)

	// 縮小したサイズに合わせて、画面の左上に縦横半分めり込む形にする
	op.GeoM.Translate(-sw/2, -sh/2)

	// 画像を画面の左上を中心に回転させる（縦横半分めり込んでいるので、中心で回転することになる)
	op.GeoM.Rotate(angle / 180 * math.Pi)

	// 好きな位置へ移動させる
	op.GeoM.Translate(x_coodinates, y_coodinates)

	// 画像を描画する
	screen.DrawImage(img, op)
}

