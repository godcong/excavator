package excavator

import (
	"bytes"
	"encoding/json"
	"errors"
)

// const error list ...
const (
	UnionMustNotBeNull            = "union must not be null"
	UnparsableNumber              = "unparsable number"
	UnionDoesNotContainNumber     = "union does not contain number"
	DecoderShouldNotReturnFloat64 = "decoder should not return float64"
	UnionDoesNotContainBool       = "union does not contain bool"
	UnionDoesNotContainString     = "union does not contain string"
	UnionDoesNotContainNull       = "union does not contain null"
	UnionDoesNotContainObject     = "union does not contain object"
	UnionDoesNotContainArray      = "union does not contain array"
	CannotHandleDelimiter         = "cannot handle delimiter"
	CannotUnmarshalUnion          = "cannot unmarshal union"
)

// marshalUnion ...
func marshalUnion(pi *int64, pf *float64, pb *bool, ps *string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) ([]byte, error) {
	if pi != nil {
		return json.Marshal(*pi)
	}
	if pf != nil {
		return json.Marshal(*pf)
	}
	if pb != nil {
		return json.Marshal(*pb)
	}
	if ps != nil {
		return json.Marshal(*ps)
	}
	if haveArray {
		return json.Marshal(pa)
	}
	if haveObject {
		return json.Marshal(pc)
	}
	if haveMap {
		return json.Marshal(pm)
	}
	if haveEnum {
		return json.Marshal(pe)
	}
	if nullable {
		return json.Marshal(nil)
	}
	return nil, errors.New(UnionMustNotBeNull)
}

// unmarshalUnion ...
func unmarshalUnion(data []byte, pi **int64, pf **float64, pb **bool, ps **string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) (bool, error) {
	if pi != nil {
		*pi = nil
	}
	if pf != nil {
		*pf = nil
	}
	if pb != nil {
		*pb = nil
	}
	if ps != nil {
		*ps = nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
		return false, err
	}

	switch v := tok.(type) {
	case json.Number:
		if pi != nil {
			i, err := v.Int64()
			if err == nil {
				*pi = &i
				return false, nil
			}
		}
		if pf != nil {
			f, err := v.Float64()
			if err == nil {
				*pf = &f
				return false, nil
			}
			return false, errors.New(UnparsableNumber)
		}
		return false, errors.New(UnionDoesNotContainNumber)
	case float64:
		return false, errors.New(DecoderShouldNotReturnFloat64)
	case bool:
		if pb != nil {
			*pb = &v
			return false, nil
		}
		return false, errors.New(UnionDoesNotContainBool)
	case string:
		if haveEnum {
			return false, json.Unmarshal(data, pe)
		}
		if ps != nil {
			*ps = &v
			return false, nil
		}
		return false, errors.New(UnionDoesNotContainString)
	case nil:
		if nullable {
			return false, nil
		}
		return false, errors.New(UnionDoesNotContainNull)
	case json.Delim:
		if v == '{' {
			if haveObject {
				return true, json.Unmarshal(data, pc)
			}
			if haveMap {
				return false, json.Unmarshal(data, pm)
			}
			return false, errors.New(UnionDoesNotContainObject)
		}
		if v == '[' {
			if haveArray {
				return false, json.Unmarshal(data, pa)
			}
			return false, errors.New(UnionDoesNotContainArray)
		}
		return false, errors.New(CannotHandleDelimiter)
	}
	return false, errors.New(CannotUnmarshalUnion)

}
