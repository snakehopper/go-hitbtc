package spiral

import (
	"encoding/json"
	"strconv"
)

type KLine struct {
	OpenTs        int64
	Open          float64
	High          float64
	Low           float64
	Close         float64
	Vol           float64
	CloseTs       int64
	RESERVED      string
	NumberOfTrade int64
}

func (r *KLine) UnmarshalJSON(bs []byte) error {
	arr := []interface{}{}
	if err := json.Unmarshal(bs, &arr); err != nil {
		return err
	}

	var err error

	r.OpenTs = int64(arr[0].(float64))
	if r.Open, err = strconv.ParseFloat(arr[1].(string), 64); err != nil {
		return err
	}
	if r.High, err = strconv.ParseFloat(arr[2].(string), 64); err != nil {
		return err
	}
	if r.Low, err = strconv.ParseFloat(arr[3].(string), 64); err != nil {
		return err
	}
	if r.Close, err = strconv.ParseFloat(arr[4].(string), 64); err != nil {
		return err
	}
	if r.Vol, err = strconv.ParseFloat(arr[5].(string), 64); err != nil {
		return err
	}
	r.CloseTs = int64(arr[6].(float64))
	r.RESERVED = arr[7].(string)
	r.NumberOfTrade = int64(arr[8].(float64))

	return nil
}

type KLineReturn struct {
	Data []KLine `json:"data"`
	errorResponse
}
