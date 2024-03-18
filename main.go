package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"rsifxbot/capital"
	"rsifxbot/functions"
	"rsifxbot/keys"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/cinar/indicator"
	"github.com/robfig/cron/v3"
)

var previous_reference string

// lucky number = 1.618
func setEnv() {
	os.Setenv("AWS_ACCESS_KEY_ID", keys.AWS_ACCESS_KEY_ID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", keys.AWS_SECRET_ACCESS_KEY)
	os.Setenv("AWS_REGION", keys.AWS_REGION)
}

// UTC TIME when market is open
// MON-TUE-WED-THU 12:00 am -9:58pm THEN 10:05 pm - 12:00 am
// FRI 12:00 am-9:58pm
// SUN 10:05 pm - 12:00 am
func isMarketOpen() bool {
	t := time.Now()
	//fmt.Println(t)
	if t.Weekday() == 0 {
		if t.Hour() == 22 && t.Minute() >= 5 {
			return true
		}
		if t.Hour() == 23 {
			return true
		}
		return false
	}
	if t.Weekday() == 5 {
		if t.Hour() >= 0 && t.Hour() < 22 {
			if t.Hour() == 21 && t.Minute() >= 58 {
				return false
			}
			return true
		}
		return false
	}

	if t.Hour() >= 0 && t.Hour() < 22 {
		if t.Hour() == 21 && t.Minute() >= 58 {
			return false
		}
		return true
	}
	if t.Hour() == 22 && t.Minute() >= 5 {
		return true
	}
	if t.Hour() == 23 {
		return true
	}
	return false
}

func main() {
	setEnv()
	c := cron.New(cron.WithSeconds())
	session, err := capital.GetSession()
	if err != nil {
		panic(err)
	}

	c.AddFunc("0 0 22 * * 1-5", func() {
		notifyPipsPerDay(session)
	})
	c.AddFunc("0 * * * * *", func() {
		capital.PingSession(session)
	})
	c.AddFunc("*/30 * * * * *", func() {
		updateStopLoss(session)
	})
	c.AddFunc("1 * * * * *", func() {
		notify(session)
	})

	c.AddFunc("0 */9 * * * *", func() {
		if isMarketOpen() {
			min9rsi(session)
		}

	})

	// c.AddFunc("0 */5 * * * *", func() {
	// 	if isMarketOpen() {
	// 		min5rsib(session)
	// 	}
	// })

	c.AddFunc("0 */5 * * * *", func() {
		if isMarketOpen() {
			//min5rsia(session)
		}
	})

	c.AddFunc("0 */6 * * * *", func() {
		if isMarketOpen() {
			//min6rsi(session)
		}
	})

	// once a week on friday at 22:00
	c.AddFunc("0 0 22 * * 5", func() {
		notifyPipsPerWeek(session)
	})

	// end of every day at 9:30 pm
	c.AddFunc("0 30 21 * * 1-5", func() {
		closePofitablePositions(session)
	})

	c.Start()
	select {}

}

func closePofitablePositions(session capital.Session) {
	positions, err := capital.GetPositions(session)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, position := range positions.Positions {
		if position.Position.Direction == "BUY" {
			if position.Market.Bid > position.Position.Level {
				closePosition, err := capital.ClosePosition(session, position.Position.DealId)
				if err != nil {
					fmt.Println(err)
				}
				if closePosition.ErrorCode != "" {
					fmt.Println("EOD Error Closing Buy: ", closePosition.ErrorCode)
				} else {
					fmt.Println("EOD Closed Buy for " + position.Market.Epic)
				}
			}
		}
		if position.Position.Direction == "SELL" {
			if position.Market.Offer < position.Position.Level {
				closePosition, err := capital.ClosePosition(session, position.Position.DealId)
				if err != nil {
					fmt.Println(err)
				}
				if closePosition.ErrorCode != "" {
					fmt.Println("EOD Error Closing Sell: ", closePosition.ErrorCode)
				} else {
					fmt.Println("EOD Closed Sell for " + position.Market.Epic)
				}
			}
		}
	}

}

// AUDUSD 9m
func min9rsi(session capital.Session) {
	candles, err := capital.GetPrices(session, "AUDUSD", "MINUTE", 991)
	if err != nil {
		fmt.Println("Use Min9 RSI Error 1: ", err)
		return
	}
	if len(candles.Candles) == 0 {
		backoff := time.Duration(rand.Intn(5-1)+1) * time.Second
		fmt.Println("API Limit Reached, Retrying in ", backoff, " seconds")
		time.Sleep(backoff)
		min9rsi(session)
		return
	}
	snapshotTimeUTC, err := time.Parse("2006-01-02T15:04:05", candles.Candles[len(candles.Candles)-1].SnapshotTimeUTC)
	if err != nil {
		fmt.Println("Use Min9 RSI Error 2: ", err)
		return
	}
	var newCandles []capital.Candle
	var start int
	if snapshotTimeUTC.Minute() == time.Now().UTC().Minute() {
		newCandles = candles.Candles[:len(candles.Candles)-1]
		start = 0
	} else {
		newCandles = candles.Candles
		start = 1
	}
	var close []float64
	var high []float64
	var low []float64
	for i := start; i < len(newCandles); i = i + 9 {
		candleStick := functions.NewCandle(newCandles[i].SnapshotTimeUTC)
		for j := i; j < i+9; j++ {
			candleStick.MergeCandle(newCandles[j])
		}
		close = append(close, candleStick.Closing)
		high = append(high, candleStick.High)
		low = append(low, candleStick.Low)

	}
	_, rsi := indicator.RsiPeriod(14, close)
	k, d := indicator.StochasticOscillator(high, low, close)
	length := len(rsi)
	if rsi[length-1] > 28 && rsi[length-2] < 28 && k[length-2] < 20 && d[length-2] < 20 {
		if err != nil {
			fmt.Println("Use Strategy Error 3: ", err)
			return
		}
		createBuy(session, "AUDUSD", 2000, 0.00500, 0.01020)
	}
	if rsi[length-1] < 72 && rsi[length-2] > 72 && k[length-2] > 80 && d[length-2] > 80 {
		if err != nil {
			fmt.Println("Use Strategy Error 4: ", err)
			return
		}
		createSell(session, "AUDUSD", 2000, 0.00500, 0.01020)
	}
}

// EURUSD 5m
func min5rsia(session capital.Session) {
	candles, err := capital.GetPrices(session, "EURUSD", "MINUTE_5", 501)
	if err != nil {
		fmt.Println("Use Strategy Error 1: ", err)
		return
	}
	if len(candles.Candles) == 0 {
		backoff := time.Duration(rand.Intn(5-1)+1) * time.Second
		fmt.Println("API Limit Reached, Retrying in ", backoff, " seconds")
		time.Sleep(backoff)
		min5rsia(session)
		return
	}

	snapshotTimeUTC, err := time.Parse("2006-01-02T15:04:05", candles.Candles[len(candles.Candles)-1].SnapshotTimeUTC)
	if err != nil {
		fmt.Println("Use Strategy Error 2: ", err)
		return
	}
	var newCandles []capital.Candle
	var start int
	if snapshotTimeUTC.Minute() == time.Now().UTC().Minute() {
		newCandles = candles.Candles[:len(candles.Candles)-1]
		start = 0
	} else {
		newCandles = candles.Candles
		start = 1
	}

	var close []float64
	var high []float64
	var low []float64
	//var volume []int64
	//var candleSticks capital.CandleSticks
	//var tradeAt []string

	// create a new candle on the fist and and add all the candles to it
	for i := start; i < len(newCandles); i++ {
		candleStick := functions.NewCandle(newCandles[i].SnapshotTimeUTC)
		candleStick.MergeCandle(newCandles[i])
		//volume = append(volume, candleSticks.Volume)
		//candleSticks.Candles = append(candleSticks.Candles, *candleStick)
		// tradeAt = append(tradeAt, candleStick.Date)
		close = append(close, candleStick.Closing)
		high = append(high, candleStick.High)
		low = append(low, candleStick.Low)

	}

	_, rsi := indicator.RsiPeriod(14, close)
	k, d := indicator.StochasticOscillator(high, low, close)
	length := len(rsi)

	if rsi[length-1] > 32 && rsi[length-2] < 32 && k[length-2] < 20 && d[length-2] < 20 {
		if err != nil {
			fmt.Println("Use Strategy Error 3: ", err)
			return
		}

		createBuy(session, "EURUSD", 1000, 0.00500, 0.01020)
	}
	if rsi[length-1] < 72 && rsi[length-2] > 72 && k[length-2] > 80 && d[length-2] > 80 {
		if err != nil {
			fmt.Println("Use Strategy Error 4: ", err)
			return
		}

		createSell(session, "EURUSD", 1000, 0.00500, 0.01020)
	}
}

// EURCHF 5m
func min5rsib(session capital.Session) {
	candles, err := capital.GetPrices(session, "EURCHF", "MINUTE_5", 501)
	if err != nil {
		fmt.Println("Use Strategy Error 1: ", err)
		return
	}
	if len(candles.Candles) == 0 {
		backoff := time.Duration(rand.Intn(5-1)+1) * time.Second
		fmt.Println("API Limit Reached, Retrying in ", backoff, " seconds")
		time.Sleep(backoff)
		min5rsib(session)
		return
	}

	snapshotTimeUTC, err := time.Parse("2006-01-02T15:04:05", candles.Candles[len(candles.Candles)-1].SnapshotTimeUTC)
	if err != nil {
		fmt.Println("Use Strategy Error 2: ", err)
		return
	}
	var newCandles []capital.Candle
	var start int
	if snapshotTimeUTC.Minute() == time.Now().UTC().Minute() {
		newCandles = candles.Candles[:len(candles.Candles)-1]
		start = 0
	} else {
		newCandles = candles.Candles
		start = 1
	}

	var close []float64
	var high []float64
	var low []float64
	//var volume []int64
	//var candleSticks capital.CandleSticks
	//var tradeAt []string

	// create a new candle on the fist and and add all the candles to it
	for i := start; i < len(newCandles); i = i + 8 {
		candleStick := functions.NewCandle(newCandles[i].SnapshotTimeUTC)
		candleStick.MergeCandle(newCandles[i])
		//volume = append(volume, candleSticks.Volume)
		//candleSticks.Candles = append(candleSticks.Candles, *candleStick)
		// tradeAt = append(tradeAt, candleStick.Date)
		close = append(close, candleStick.Closing)
		high = append(high, candleStick.High)
		low = append(low, candleStick.Low)

	}

	_, rsi := indicator.RsiPeriod(14, close)
	k, d := indicator.StochasticOscillator(high, low, close)
	length := len(rsi)

	if rsi[length-1] > 28 && rsi[length-2] < 28 && k[length-2] < 20 && d[length-2] < 20 {
		if err != nil {
			fmt.Println("Use Strategy Error 3: ", err)
			return
		}
		createBuy(session, "EURCHF", 1000, 0.00500, 0.01020)
	}
	if rsi[length-1] < 72 && rsi[length-2] > 72 && k[length-2] > 80 && d[length-2] > 80 {
		if err != nil {
			fmt.Println("Use Strategy Error 4: ", err)
			return
		}

		createSell(session, "EURCHF", 1000, 0.00500, 0.01020)
	}
}

// EURUSD 6m
func min6rsi(session capital.Session) {
	candles, err := capital.GetPrices(session, "EURUSD", "MINUTE", 901)
	if err != nil {
		fmt.Println("Use Strategy Error 1: ", err)
		return
	}
	if len(candles.Candles) == 0 {
		backoff := time.Duration(rand.Intn(5-1)+1) * time.Second
		fmt.Println("API Limit Reached, Retrying in ", backoff, " seconds")
		time.Sleep(backoff)
		min6rsi(session)
		return
	}

	snapshotTimeUTC, err := time.Parse("2006-01-02T15:04:05", candles.Candles[len(candles.Candles)-1].SnapshotTimeUTC)
	if err != nil {
		fmt.Println("Use Strategy Error 2: ", err)
		return
	}
	var newCandles []capital.Candle
	var start int
	if snapshotTimeUTC.Minute() == time.Now().UTC().Minute() {
		newCandles = candles.Candles[:len(candles.Candles)-1]
		start = 0
	} else {
		newCandles = candles.Candles
		start = 1
	}

	var close []float64
	var high []float64
	var low []float64
	for i := start; i < len(newCandles); i = i + 6 {
		candleStick := functions.NewCandle(newCandles[i].SnapshotTimeUTC)
		for j := i; j < i+6; j++ {
			candleStick.MergeCandle(newCandles[j])
		}
		close = append(close, candleStick.Closing)
		high = append(high, candleStick.High)
		low = append(low, candleStick.Low)

	}

	_, rsi := indicator.RsiPeriod(14, close)
	k, d := indicator.StochasticOscillator(high, low, close)
	length := len(rsi)

	if rsi[length-1] > 28 && rsi[length-2] < 28 && k[length-2] < 20 && d[length-2] < 20 {
		if err != nil {
			fmt.Println("Use Strategy Error 3: ", err)
			return
		}
		createBuy(session, "EURUSD", 2000, 0.00500, 0.01020)
	}
	if rsi[length-1] < 72 && rsi[length-2] > 72 && k[length-2] > 80 && d[length-2] > 80 {
		if err != nil {
			fmt.Println("Use Strategy Error 4: ", err)
			return
		}

		createSell(session, "EURUSD", 2000, 0.00500, 0.01020)
	}
}

func updateStopLoss(session capital.Session) {
	positions, err := capital.GetPositions(session)
	if err != nil {
		return
	}

	mapofhigh := make(map[string]float64)
	mapoflow := make(map[string]float64)
	for _, position := range positions.Positions {
		mapofhigh[position.Market.Epic] = 0
		mapoflow[position.Market.Epic] = 0
	}

	for key := range mapofhigh {
		price, err := capital.GetPrices(session, key, "MINUTE", 2)
		if err != nil {
			fmt.Println("Update Stop Loss Error 2: ", err)
			return
		}
		if len(price.Candles) == 0 {
			backoff := time.Duration(rand.Intn(5-1)+1) * time.Second
			fmt.Println("API Limit Reached, Retrying in 2 seconds")
			time.Sleep(backoff)
			updateStopLoss(session)
			return
		}
		high := price.Candles[0].HighPrice.Bid
		if high < price.Candles[1].HighPrice.Bid {
			high = price.Candles[1].HighPrice.Bid
		}
		mapofhigh[key] = high

		low := price.Candles[0].LowPrice.Bid
		if low > price.Candles[1].LowPrice.Bid {
			low = price.Candles[1].LowPrice.Bid
		}
		mapoflow[key] = low

	}

	for _, position := range positions.Positions {
		if position.Position.Direction == "BUY" {
			created_at, err := time.Parse("2006-01-02T15:04:05.000", position.Position.CreatedDateUTC)
			if err != nil {
				fmt.Println("Update Stop Loss Error 1: ", err)
				continue
			}

			if time.Now().UTC().Sub(created_at) > 2*time.Minute {
				high := mapofhigh[position.Market.Epic]
				k := high - position.Position.Level
				var ts float64

				if position.Position.StopLevel == 0 {
					ts = capital.TRAILINGSTOP
					if position.Market.Epic == "EURUSD" && position.Position.Size == 2000 {
						ts = 0.0050
					}

				} else {
					ts = functions.TrailingStop(k)
					if position.Market.Epic == "EURUSD" && position.Position.Size == 2000 {
						ts = functions.BigtrailingStop(k)
					}
				}

				sl := high - ts
				truncated := strconv.FormatFloat(sl, 'f', 5, 64)
				sl, err = strconv.ParseFloat(truncated, 64)
				if err != nil {
					fmt.Println("Update Stop Loss Error 3: ", err)
					return
				}
				if sl > position.Position.StopLevel || position.Position.StopLevel == 0 {
					updatePosition, err := capital.UpdatePosition(session, position.Position.DealId, sl, position.Position.ProfitLevel)
					if err != nil {
						fmt.Println("Update Stop Loss Error 4: ", err)
						continue
					}
					if updatePosition.ErrorCode != "" {
						fmt.Println("Updating Stop Loss Error: ", updatePosition.ErrorCode)
					} else {
						fmt.Println("Stop Loss Updated for "+position.Market.Epic, " to", sl)
					}
				}

			}
		}
		if position.Position.Direction == "SELL" {

			created_at, err := time.Parse("2006-01-02T15:04:05.000", position.Position.CreatedDateUTC)
			if err != nil {
				fmt.Println("Update Stop Loss Error 5: ", err)
				continue
			}

			if time.Now().UTC().Sub(created_at) > 2*time.Minute {
				low := mapoflow[position.Market.Epic]
				k := position.Position.Level - low

				var ts float64

				if position.Position.StopLevel == 0 {
					ts = capital.TRAILINGSTOP
					if position.Market.Epic == "EURUSD" && position.Position.Size == 2000 {
						ts = 0.0050
					}

				} else {
					ts = functions.TrailingStop(k)
					if position.Market.Epic == "EURUSD" && position.Position.Size == 2000 {
						ts = functions.BigtrailingStop(k)
					}
				}

				sl := low + ts
				truncated := strconv.FormatFloat(sl, 'f', 5, 64)
				sl, err = strconv.ParseFloat(truncated, 64)
				if err != nil {
					fmt.Println("Update Stop Loss Error 7: ", err)
					return
				}
				if sl < position.Position.StopLevel || position.Position.StopLevel == 0 {
					updatePosition, err := capital.UpdatePosition(session, position.Position.DealId, sl, position.Position.ProfitLevel)
					if err != nil {
						fmt.Println("Update Stop Loss Error 8: ", err)
						continue
					}
					if updatePosition.ErrorCode != "" {
						fmt.Println("Updating Stop Loss Error: ", updatePosition.ErrorCode)
					} else {
						fmt.Println("Stop Loss Updated for "+position.Market.Epic, " to", sl)
					}
				}

			}
		}

	}

}

//	func rejectedOrders(session capital.Session) {
//		startAt := time.Now().UTC().Format("2006-01-02")
//		activity, err := capital.GetActivity(session, startAt, startAt, "REJECTED")
//		if err != nil {
//			fmt.Println("RejectedOrders Error 1: ", err)
//			return
//		}
//		// fmt.Println("Rejected Orders: ", activity)
//		for _, deal := range activity.Activities {
//			fmt.Println("Rejected Order: ", deal.DateUTC, deal.DealId, deal.Epic, deal.Details.Direction, deal.Details.Size)
//		}
//	}
func notifyPipsPerDay(session capital.Session) {
	type Position struct {
		BuyAt  float64 `json:"buyAt"`
		SellAt float64 `json:"sellAt"`
		Size   float64 `json:"size"`
	}
	startAt := time.Now().UTC().Format("2006-01-02")
	activity, err := capital.GetActivity(session, startAt, startAt, "ACCEPTED")
	if err != nil {
		fmt.Println("NotifyPipsPerDay Error 1: ", err)
		return
	}

	mapofDealId := make(map[string]Position)
	for _, deal := range activity.Activities {
		if deal.Epic != "AUDUSD" && deal.Epic != "USDCAD" && deal.Epic != "NZDUSD" && deal.Epic != "GBPUSD" && deal.Epic != "USDCHF" && deal.Epic != "EURCHF" {
			continue
		}
		if deal.Details.Direction == "BUY" {
			mapofDealId[deal.DealId] = Position{BuyAt: deal.Details.Level, SellAt: mapofDealId[deal.DealId].SellAt, Size: deal.Details.Size}
		}
		if deal.Details.Direction == "SELL" {
			mapofDealId[deal.DealId] = Position{SellAt: deal.Details.Level, BuyAt: mapofDealId[deal.DealId].BuyAt, Size: deal.Details.Size}
		}
	}

	var totalRealizedPips float64
	for _, deal := range mapofDealId {
		if deal.BuyAt != 0 && deal.SellAt != 0 {
			totalRealizedPips += (deal.SellAt - deal.BuyAt) * (deal.Size / 1000)
		}
	}

	// 	else {

	// 		activityByDealId, err := capital.GetActivityByDealId(session, dealId)

	// 		if err != nil {
	// 			fmt.Println("NotifyPipsPerDay Error 2: ", err)
	// 			return
	// 		}

	// 		if len(activityByDealId.Activities) == 2 {
	// 			if activityByDealId.Activities[0].Details.Direction == "BUY" {
	// 				totalRealizedPips += activityByDealId.Activities[1].Details.Level - activityByDealId.Activities[0].Details.Level
	// 			} else {
	// 				totalRealizedPips += activityByDealId.Activities[0].Details.Level - activityByDealId.Activities[1].Details.Level
	// 			}
	// 		} else {
	// 			var price float64
	// 			if activityByDealId.Activities[0].Details.StopLevel == 0 {
	// 				prices, err := capital.GetPrices(session, activityByDealId.Activities[0].Epic, "MINUTE", 1)
	// 				if err != nil {
	// 					fmt.Println("NotifyPipsPerDay Error 3: ", err)
	// 					return
	// 				}
	// 				price = prices.Candles[0].ClosePrice.Bid
	// 			} else {
	// 				price = activityByDealId.Activities[0].Details.StopLevel
	// 			}

	// 			if activityByDealId.Activities[0].Details.Direction == "BUY" {
	// 				totalUnrealizedPips += activityByDealId.Activities[0].Details.Level - price

	// 			} else {
	// 				totalUnrealizedPips += price - activityByDealId.Activities[0].Details.Level

	// 			}
	// 		}

	// 	}

	// }
	totalRealizedPips = totalRealizedPips * 10000
	if totalRealizedPips != 0 {
		msg := "Today Realized Pips: " + strconv.FormatFloat(totalRealizedPips, 'f', 1, 64)
		sendNotification(msg)
	} else {
		fmt.Println("No Pips Today")
	}

}

func notifyPipsPerWeek(session capital.Session) {
	type Position struct {
		BuyAt  float64 `json:"buyAt"`
		SellAt float64 `json:"sellAt"`
		Size   float64 `json:"size"`
	}
	startAt := time.Now().UTC().AddDate(0, 0, -int(time.Now().UTC().Weekday())).Format("2006-01-02")
	endAt := time.Now().UTC().Format("2006-01-02")
	activity, err := capital.GetActivity(session, startAt, endAt, "ACCEPTED")
	if err != nil {
		fmt.Println("notifyPipsPerWeek Error 1: ", err)
		return
	}

	mapofDealId := make(map[string]Position)
	for _, deal := range activity.Activities {
		if deal.Epic != "AUDUSD" && deal.Epic != "USDCAD" && deal.Epic != "NZDUSD" && deal.Epic != "GBPUSD" && deal.Epic != "USDCHF" && deal.Epic != "EURCHF" {
			continue
		}
		if deal.Details.Direction == "BUY" {
			mapofDealId[deal.DealId] = Position{BuyAt: deal.Details.Level, SellAt: mapofDealId[deal.DealId].SellAt, Size: deal.Details.Size}
		}
		if deal.Details.Direction == "SELL" {
			mapofDealId[deal.DealId] = Position{SellAt: deal.Details.Level, BuyAt: mapofDealId[deal.DealId].BuyAt, Size: deal.Details.Size}
		}
	}

	var totalRealizedPips float64
	for _, deal := range mapofDealId {
		if deal.BuyAt != 0 && deal.SellAt != 0 {
			totalRealizedPips += (deal.SellAt - deal.BuyAt) * (deal.Size / 1000)
		}
	}
	totalRealizedPips = totalRealizedPips * 10000
	if totalRealizedPips != 0 {
		msg := "Weekly Realized Pips: " + strconv.FormatFloat(totalRealizedPips, 'f', 1, 64)
		sendNotification(msg)
	} else {
		fmt.Println("No Pips Today")
	}

}

func notify(session capital.Session) {
	transactions, err := capital.GetTransactions(session)
	if err != nil {
		return
	}

	if len(transactions.Transactions) > 0 {
		var changed float64 = 0
		for _, transaction := range transactions.Transactions {
			if transaction.Reference == previous_reference {
				break
			}
			if transaction.TransactionType == "TRADE" {
				size, err := strconv.ParseFloat(transaction.Size, 64)
				if err != nil {
					fmt.Println("Notify Error 2: ", err)
					continue
				}
				changed += size
			}
			if transaction.TransactionType == "SWAP" {
				size, err := strconv.ParseFloat(transaction.Size, 64)
				if err != nil {
					fmt.Println("Notify Error 3: ", err)
					continue
				}
				changed -= math.Abs(size)
			}

		}
		if changed != 0 {
			previous_reference = transactions.Transactions[0].Reference
			sendNotification(strconv.FormatFloat(changed, 'f', 2, 64) + " USD")
		}
	}
}

func createBuy(session capital.Session, symbol string, size float64, sl float64, tp float64) {
	p, err := capital.GetPositions(session)
	if err != nil {
		fmt.Println(err)
		return
	}

	lengthofpositions := len(p.Positions)
	ifpresent := false
	for _, position := range p.Positions {
		if position.Market.Epic == symbol {
			ifpresent = true
		}
		if position.Position.Direction == "SELL" && position.Market.Epic == symbol && position.Position.Size == size {
			if position.Market.Offer < position.Position.Level {
				closePosition, err := capital.ClosePosition(session, position.Position.DealId)
				if err != nil {
					fmt.Println(err)
				}
				if closePosition.ErrorCode != "" {
					fmt.Println("Error Closing Sell: ", closePosition.ErrorCode)
				} else {
					lengthofpositions--
					fmt.Println("Closed Sell for " + position.Market.Epic)
				}
			}
		}
	}

	spread := 0.0
	if ifpresent {
		for _, position := range p.Positions {
			if symbol == position.Market.Epic {
				spread = position.Market.Offer - position.Market.Bid
			}
		}
	} else {
		c, err := capital.GetPrices(session, symbol, "MINUTE", 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(c.Candles) == 0 {
			backoff := time.Duration(rand.Intn(5-1)+1) * time.Second
			fmt.Println("API Limit Reached, Retrying in ", backoff, " seconds for symbol", symbol, "and size", size)
			time.Sleep(backoff)
			createBuy(session, symbol, size, sl, tp)
			return
		}
		spread = c.Candles[0].OpenPrice.Ask - c.Candles[0].OpenPrice.Bid
	}
	if spread > 0.00006 {
		fmt.Println("Spread too high, not opening a position for ", symbol)
		return
	}

	if lengthofpositions > 30 {
		sendNotification("BUY (" + symbol + ") : Too many open positions, for size " + strconv.FormatFloat(size, 'f', 6, 64) + ", please close some or manually open a position")
		return
	}

	createPosition, err := capital.CreatePosition(session, symbol, "BUY", size, sl, tp)
	if err != nil {
		fmt.Println(err)
	}
	if createPosition.ErrorCode != "" {
		fmt.Println("Error Creating Buy: ", createPosition.ErrorCode)
	} else {
		fmt.Println("Created Buy for " + symbol)
	}
}

// Opens a sell position and closes all buy positions on the same symbol
func createSell(session capital.Session, symbol string, size float64, sl float64, tp float64) {
	p, err := capital.GetPositions(session)
	if err != nil {
		fmt.Println(err)
		return
	}
	lengthofpositions := len(p.Positions)
	ifpresent := false
	for _, position := range p.Positions {
		// if symbol is already present
		if position.Market.Epic == symbol {
			ifpresent = true
		}
		if position.Position.Direction == "BUY" && position.Market.Epic == symbol && position.Position.Size == size {
			if position.Market.Bid > position.Position.Level {
				closePosition, err := capital.ClosePosition(session, position.Position.DealId)
				if err != nil {
					fmt.Println(err)
				}
				if closePosition.ErrorCode != "" {
					fmt.Println("Error Closing Buy: ", closePosition.ErrorCode)
				} else {
					lengthofpositions--
					fmt.Println("Closed Buy for " + position.Market.Epic)
				}
			}
		}
	}

	spread := 0.0
	if ifpresent {
		for _, position := range p.Positions {
			if symbol == position.Market.Epic {
				spread = position.Market.Offer - position.Market.Bid
			}
		}
	} else {
		c, err := capital.GetPrices(session, symbol, "MINUTE", 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(c.Candles) == 0 {
			backoff := time.Duration(rand.Intn(5-1)+1) * time.Second
			fmt.Println("API Limit Reached, Retrying in ", backoff, " seconds for symbol", symbol, "and size", size)
			time.Sleep(backoff)
			createSell(session, symbol, size, sl, tp)
			return
		}
		spread = c.Candles[0].OpenPrice.Ask - c.Candles[0].OpenPrice.Bid
	}

	if spread > 0.00006 {
		fmt.Println("Spread too high, not opening a position for ", symbol)
		return
	}
	if lengthofpositions > 30 {
		sendNotification("SELL (" + symbol + ") : Too many open positions, for size " + strconv.FormatFloat(size, 'f', 6, 64) + ", please close some or manually open a position")
		return
	}

	createPosition, err := capital.CreatePosition(session, symbol, "SELL", size, sl, tp)
	if err != nil {
		fmt.Println(err)

	}
	if createPosition.ErrorCode != "" {
		fmt.Println("Error Creating Sell: ", createPosition.ErrorCode)
	} else {
		fmt.Println("Created Sell for " + symbol)
	}

}

func sendNotification(message string) {
	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_REGION := os.Getenv("AWS_REGION")
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(AWS_REGION),
		Credentials: credentials.NewStaticCredentials(
			AWS_ACCESS_KEY_ID,
			AWS_SECRET_ACCESS_KEY,
			""),
	}))

	svc := sns.New(sess)
	params := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String("+919999999999"),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"AWS.SNS.SMS.SMSType": {StringValue: aws.String("Transactional"), DataType: aws.String("String")},
		},
	}
	_, err := svc.Publish(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
