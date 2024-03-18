package functions

import "rsifxbot/capital"

// after 5 pips update the trailing stop loss to 26.5 pips
// after 10 pips update the trailing stop loss to 25.5 pips
// after 15 pips update the trailing stop loss to 24.5 pips
// after 20 pips update the trailing stop loss to 23.5 pips
// after 25 pips update the trailing stop loss to 22.5 pips
// after 30 pips update the trailing stop loss to 21.5 pips
// after 35 pips update the trailing stop loss to 20.5 pips
// after 40 pips update the trailing stop loss to 19.5 pips
// after 45 pips update the trailing stop loss to 18.5 pips
// after 50 pips update the trailing stop loss to 17.5 pips
// after 55 pips update the trailing stop loss to 16.5 pips
// after 60 pips update the trailing stop loss to 15.5 pips
// after 65 pips update the trailing stop loss to 14.5 pips
// after 70 pips update the trailing stop loss to 13.5 pips
// after 75 pips update the trailing stop loss to 12.5 pips
// after 80 pips update the trailing stop loss to 11.5 pips
// after 85 pips update the trailing stop loss to 10.5 pips
// after 90 pips update the trailing stop loss to 9.5 pips
// after 95 pips update the trailing stop loss to 8.5 pips
// after 100 pips update the trailing stop loss to 7.5 pips
// after 105 pips update the trailing stop loss to 6.5 pips
// after 110 pips update the trailing stop loss to 5.5 pips
// after 115 pips update the trailing stop loss to 4.5 pips
// after 120 pips update the trailing stop loss to 3.5 pips
// after 125 pips update the trailing stop loss to 2.5 pips
// after 130 pips update the trailing stop loss to 1.5 pips
// after 135 pips update the trailing stop loss to 0.5 pips
func TrailingStop(k float64) float64 {
	switch {
	case k > 0.0005:
		return 0.00265
	case k > 0.001:
		return 0.00255
	case k > 0.0015:
		return 0.00245
	case k > 0.002:
		return 0.00235
	case k > 0.0025:
		return 0.00225
	case k > 0.003:
		return 0.00215
	case k > 0.0035:
		return 0.00205
	case k > 0.004:
		return 0.00195
	case k > 0.0045:
		return 0.00185
	case k > 0.005:
		return 0.00175
	case k > 0.0055:
		return 0.00165
	case k > 0.006:
		return 0.00155
	case k > 0.0065:
		return 0.00145
	case k > 0.007:
		return 0.00135
	case k > 0.0075:
		return 0.00125
	case k > 0.008:
		return 0.00115
	case k > 0.0085:
		return 0.00105
	case k > 0.009:
		return 0.00095
	case k > 0.0095:
		return 0.00085
	case k > 0.01:
		return 0.00075
	case k > 0.0105:
		return 0.00065
	case k > 0.011:
		return 0.00055
	case k > 0.0115:
		return 0.00045
	case k > 0.012:
		return 0.00035
	case k > 0.0125:
		return 0.00025
	case k > 0.013:
		return 0.00015
	case k > 0.0135:
		return 0.00005
	default:
		// return 0.0005
		return capital.TRAILINGSTOP
	}
}

// trailing stop loss for EURUSD WITH 45 PIPs
func BigtrailingStop(k float64) float64 {
	//0.0005 => 5 pips
	// for every 5 pips decrease the trailing stop distance by 10 pips
	switch {
	case k > 0.0005:
		return 0.00440
	case k > 0.001:
		return 0.00430
	case k > 0.0015:
		return 0.00420
	case k > 0.002:
		return 0.00410
	case k > 0.0025:
		return 0.00400
	case k > 0.003:

		return 0.00390
	case k > 0.0035:
		return 0.00380
	case k > 0.004:
		return 0.00370
	case k > 0.0045:
		return 0.00360
	case k > 0.005:
		return 0.00350
	case k > 0.0055:
		return 0.00340
	case k > 0.006:
		return 0.00330
	case k > 0.0065:
		return 0.00320
	case k > 0.007:
		return 0.00310
	case k > 0.0075:
		return 0.00300
	case k > 0.008:
		return 0.00290
	case k > 0.0085:
		return 0.00280
	case k > 0.009:
		return 0.00270
	case k > 0.0095:
		return 0.00260
	case k > 0.01:
		return 0.00250
	case k > 0.0105:
		return 0.00240
	case k > 0.011:
		return 0.00230
	case k > 0.0115:
		return 0.00220
	case k > 0.012:
		return 0.00210
	case k > 0.0125:
		return 0.00200
	case k > 0.013:
		return 0.00190
	case k > 0.0135:
		return 0.00180
	case k > 0.014:
		return 0.00170
	case k > 0.0145:
		return 0.00160
	case k > 0.015:
		return 0.00150
	case k > 0.0155:
		return 0.00140
	case k > 0.016:
		return 0.00130
	case k > 0.0165:
		return 0.00120
	case k > 0.017:
		return 0.00110
	case k > 0.0175:
		return 0.00100
	case k > 0.018:
		return 0.00090
	case k > 0.0185:
		return 0.00080
	case k > 0.019:
		return 0.00070
	case k > 0.0195:
		return 0.00060
	case k > 0.02:
		return 0.00050
	case k > 0.0205:
		return 0.00040
	case k > 0.021:
		return 0.00030
	case k > 0.0215:
		return 0.00020
	default:
		return 0.005
	}

}

// MACD BOLLINGER BANDS AND MA
