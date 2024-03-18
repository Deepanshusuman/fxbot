package functions

import "rsifxbot/capital"

type CandleStick struct {
	Date    string  `json:"date"`
	Closing float64 `json:"close"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Volume  int64   `json:"volume"`
}
type CandleSticks struct {
	Candles []CandleStick `json:"candles"`
}

func NewCandle(date string) (c *CandleStick) {
	return &CandleStick{
		Date:    date,
		Closing: 0,
		High:    0,
		Low:     0,
		Volume:  0,
	}
}

func (c *CandleStick) MergeCandle(newcandle capital.Candle) {
	c.Closing = newcandle.ClosePrice.Bid
	if c.High == 0 {
		c.High = newcandle.HighPrice.Bid
	} else if newcandle.HighPrice.Bid > c.High {
		c.High = newcandle.HighPrice.Bid
	}

	if c.Low == 0 {
		c.Low = newcandle.LowPrice.Bid
	} else if newcandle.LowPrice.Bid < c.Low {
		c.Low = newcandle.LowPrice.Bid
	}

	if c.Volume == 0 {
		c.Volume = newcandle.LastTradedVolume
	} else {
		c.Volume += newcandle.LastTradedVolume
	}

}