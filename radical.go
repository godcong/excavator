package excavator

// RadicalCharacter ...
type RadicalCharacter struct {
	Hash   string `json:"hash" xorm:"hash,p"`
	Zi     string `json:"zi"`
	PinYin string `json:"pinyin"`
	BiHua  string `json:"bihua"`
	BuShou string `json:"bushou"`
	Num    string `json:"num"`
	URL    string `json:"url"`
}
