package excavator

import "encoding/json"

// Radical ...
type Radical []RadicalUnion

// UnmarshalRadical ...
func UnmarshalRadical(data []byte) (Radical, error) {
	var r Radical
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal ...
func (r *Radical) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// RadicalClass ...
type RadicalClass struct {
	Zi     string `json:"zi"`
	Pinyin string `json:"pinyin"`
	Bushou string `json:"bushou"`
	Num    string `json:"num"`
	URL    string `json:"url"`
}

// RadicalUnion ...
type RadicalUnion struct {
	String            *string
	RadicalClassArray []RadicalClass
}

// UnmarshalJSON ...
func (x *RadicalUnion) UnmarshalJSON(data []byte) error {
	x.RadicalClassArray = nil
	object, err := unmarshalUnion(data, nil, nil, nil, &x.String, true, &x.RadicalClassArray, false, nil, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
	}
	return nil
}

// MarshalJSON ...
func (x *RadicalUnion) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, x.String, x.RadicalClassArray != nil, x.RadicalClassArray, false, nil, false, nil, false, nil, false)
}
