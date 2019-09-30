package excavator

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

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

func RadicalReader(reader io.ReadCloser) (*Radical, error) {
	bytes, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return UnmarshalRadical(bytes)
}

// Radical ...
type Radical []RadicalUnion

// UnmarshalRadical ...
func UnmarshalRadical(data []byte) (*Radical, error) {
	var r Radical
	err := json.Unmarshal(data, &r)
	return &r, err
}


// Marshal ...
func (r *Radical) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// RadicalUnion ...
type RadicalUnion struct {
	String                *string
	RadicalCharacterArray []RadicalCharacter
}

// UnmarshalJSON ...
func (x *RadicalUnion) UnmarshalJSON(data []byte) error {
	x.RadicalCharacterArray = nil
	_, err := unmarshalUnion(data, nil, nil, nil, &x.String, true, &x.RadicalCharacterArray, false, nil, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	return nil
}

// MarshalJSON ...
func (x *RadicalUnion) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, x.String, x.RadicalCharacterArray != nil, x.RadicalCharacterArray, false, nil, false, nil, false, nil, false)
}
