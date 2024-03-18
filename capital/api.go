package capital

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)



var URL = "https://api-capital.backend-capital.com/api/v1"

// https://demo-api-capital.backend-capital.com/api/v1/positions
// var URL = "https://demo-api-capital.backend-capital.com/api/v1"
var TRAILINGSTOP = 0.00275

func PingSession(session Session) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL+"/ping", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("X-SECURITY-TOKEN", session.SecurityToken)
	req.Header.Add("CST", session.Cst)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}

func CreatePosition(session Session, symbol string, direction string, size float64, sl float64, tp float64) (PositionResponse, error) {
	var result PositionResponse
	client := &http.Client{}
	req, err := http.NewRequest("POST", URL+"/positions", bytes.NewBuffer([]byte(`{
		"epic": "`+symbol+`",
		"direction": "`+direction+`",
		"size": "`+fmt.Sprintf("%.5f", size)+`",
		"stopDistance": `+strconv.FormatFloat(sl, 'f', 5, 64)+`,
		"profitDistance": `+strconv.FormatFloat(tp, 'f', 5, 64)+`
	}`)))

	if err != nil {
		return result, err
	}
	req.Header.Add("X-SECURITY-TOKEN", session.SecurityToken)
	req.Header.Add("CST", session.Cst)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetSession() (Session, error) {
	var session Session
	client := &http.Client{}
	req, err := http.NewRequest("POST", URL+"/session", bytes.NewBuffer([]byte(`{
		"identifier": "`+identifier+`",
		"password": "`+password+`"
	}`)))

	if err != nil {
		return session, err

	}
	req.Header.Add("X-CAP-API-KEY", key)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return session, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return session, err
	}

	err = json.Unmarshal(body, &session)
	if err != nil {
		return session, err
	}
	session.Cst = res.Header.Get("CST")
	session.SecurityToken = res.Header.Get("X-SECURITY-TOKEN")
	return session, nil
}

func GetActivity(session Session, startAt string, endAt string, status string) (Activities, error) {
	var activitys Activities
	client := &http.Client{}
	startAt = startAt + "T00:00:00"
	endAt = endAt + "T23:59:59"
	req, err := http.NewRequest("GET", URL+"/history/activity?from="+startAt+"&to="+endAt+"&detailed=true&filter=type!=EDIT_STOP_AND_LIMIT;status=="+status, nil)
	if err != nil {
		return activitys, err
	}
	req.Header.Add("X-SECURITY-TOKEN", session.SecurityToken)
	req.Header.Add("CST", session.Cst)
	resp, err := client.Do(req)
	if err != nil {
		return activitys, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return activitys, err
	}
	err = json.Unmarshal(body, &activitys)
	if err != nil {
		return activitys, err
	}
	return activitys, err
}

func GetActivityByDealId(session Session, dealId string) (Activities, error) {
	var activitys Activities
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL+"/history/activity?dealId="+dealId+"&detailed=true&filter=type!=EDIT_STOP_AND_LIMIT;status==ACCEPTED", nil)
	if err != nil {
		return activitys, err
	}
	req.Header.Add("X-SECURITY-TOKEN", session.SecurityToken)
	req.Header.Add("CST", session.Cst)
	resp, err := client.Do(req)
	if err != nil {
		return activitys, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return activitys, err
	}
	err = json.Unmarshal(body, &activitys)
	if err != nil {
		return activitys, err
	}
	return activitys, err
}

func GetTransactions(session Session) (Transactions, error) {
	var transactions Transactions
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL+"/history/transactions", nil)
	if err != nil {
		return transactions, err
	}
	req.Header.Add("X-SECURITY-TOKEN", session.SecurityToken)
	req.Header.Add("CST", session.Cst)
	resp, err := client.Do(req)
	if err != nil {
		return transactions, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return transactions, err
	}
	err = json.Unmarshal(body, &transactions)

	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func ClosePosition(session Session, dealId string) (PositionResponse, error) {
	deleteposition := PositionResponse{}
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", URL+"/positions/"+dealId, nil)
	if err != nil {
		return deleteposition, err
	}
	req.Header.Add("X-SECURITY-TOKEN", session.SecurityToken)
	req.Header.Add("CST", session.Cst)
	res, err := client.Do(req)
	if err != nil {
		return deleteposition, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return deleteposition, err
	}
	err = json.Unmarshal(body, &deleteposition)
	if err != nil {
		return deleteposition, err
	}
	return deleteposition, nil
}

func UpdatePosition(session Session, dealId string, sl float64, tp float64) (PositionResponse, error) {
	updatePosition := PositionResponse{}
	client := &http.Client{}
	var jsonStr []byte
	if tp == 0 {
		jsonStr = []byte(`{"stopLevel":` + strconv.FormatFloat(sl, 'f', 5, 64) + `}`)
	} else {
		jsonStr = []byte(`{"stopLevel":` + strconv.FormatFloat(sl, 'f', 5, 64) + `, "profitLevel":` + strconv.FormatFloat(tp, 'f', 5, 64) + `}`)
	}

	req, err := http.NewRequest("PUT", URL+"/positions/"+dealId, bytes.NewBuffer(jsonStr))
	if err != nil {
		return updatePosition, err
	}
	req.Header.Add("X-SECURITY-TOKEN", session.SecurityToken)

	req.Header.Add("CST", session.Cst)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return updatePosition, err
	}
	defer res.Body.Close()
	var result PositionResponse

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return updatePosition, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return updatePosition, err
	}
	return result, nil

}

func GetPrices(session Session, symbol string, resolution string, max int) (Candles, error) {
	var candles Candles
	client := &http.Client{}
	minstr := strconv.Itoa(max)
	req, err := http.NewRequest("GET", URL+"/prices/"+symbol+"?resolution="+resolution+"&max="+minstr, nil)
	if err != nil {

		return candles, err
	}
	req.Header.Add("X-SECURITY-TOKEN", session.SecurityToken)
	req.Header.Add("CST", session.Cst)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Get Prices Error 2: ", err)
		return candles, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Get Prices Error 3: ", err)
		return candles, err
	}
	json.Unmarshal(body, &candles)
	return candles, nil

}

func GetPositions(session Session) (Positions, error) {
	posistions := Positions{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL+"/positions", nil)

	if err != nil {
		return posistions, err
	}
	req.Header.Add("X-SECURITY-TOKEN", session.SecurityToken)
	req.Header.Add("CST", session.Cst)
	resp, err := client.Do(req)
	if err != nil {
		return posistions, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return posistions, err
	}
	err = json.Unmarshal(body, &posistions)
	if err != nil {
		return posistions, err
	}
	return posistions, nil

}
