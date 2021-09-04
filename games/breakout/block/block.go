package block

// 変数名を大文字にしないとimportしたファイルからインポートできない
type Block struct {
	Width int
	// height int
	// cr     float32
	// cg     float32
	// gb     float32
	// ca     float32
	// hit    int
}

func NewBlock(width int) (b *Block) {
	b = new(Block)
	b.Width = width
	return b
}

func (b *Block) Output() string {
	return "test"
}
