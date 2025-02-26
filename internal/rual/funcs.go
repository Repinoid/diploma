package rual

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/theplant/luhn"
)

type Tovar struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
type Buyback struct {
	Match      string `json:"match"`
	Reward     int    `json:"reward"`
	RewardType string `json:"RewardType"`
}
type orda struct {
	Order string  `json:"order"`
	Goods []Tovar `json:"goods"`
}
type OrderStatus struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

var Accrualhost = "localhost:8089"
var Time429 time.Time

var marks = []Buyback{
	{Match: "Acer", Reward: 20, RewardType: "pt"},
	{Match: "Bork", Reward: 10, RewardType: "%"},
	{Match: "Asus", Reward: 20, RewardType: "pt"},
	{Match: "Samsung", Reward: 25, RewardType: "%"},
	{Match: "Apple", Reward: 35, RewardType: "%"},
}

func LoadGood(num int, goodIdx int, price float64) error {
	ord := orda{Order: strconv.Itoa(Luhner(num)), Goods: []Tovar{
		{Description: "Smth " + marks[goodIdx].Match + " " + strconv.Itoa(num), Price: price}}}
	buyM, _ := json.Marshal(ord)
	err := poster("/api/orders", buyM)
	return err
}

func InitAccrualForTests() error {
	for _, r := range marks { // load to accrual good's type and buybacks
		buyM, err := json.Marshal(r)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		err = poster("/api/goods", buyM)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	for idx := range 999 { // затарим ордерами
		err := LoadGood(idx+1, int(rand.Int63n(5)), 1000)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}

func poster(postCMD string, wts []byte) error {
	httpc := resty.New() //
	httpc.SetBaseURL("http://" + Accrualhost)
	req := httpc.R().
		SetHeader("Content-Type", "application/json").
		SetBody(wts)
	_, err := req.
		SetDoNotParseResponse(false).
		Post(postCMD) //
	return err
}

func Luhner(numb int) int {
	// if luhn.Valid(numb) {	// если возвращать неизменённым, возникнут коллизии, типа у 2 Лун 26, и у 26 тоже 26
	// 	return numb
	// }
	return 10*numb + luhn.CalculateLuhn(numb)
}

// OrderStatus - {номер заказа; статус расчёта начисления; рассчитанные баллы к начислению}
func GetFromAccrual(number string) (OrderStatus, int, error) {

	wait429 := time.Until(Time429) // время до разморозки
	time.Sleep(wait429)

	httpc := resty.New() //
	httpc.SetBaseURL(Accrualhost)
	//	httpc.SetBaseURL("http://" + Accrualhost)
	getReq := httpc.R()

	orderStat := &OrderStatus{}
	resp, err := getReq.
		SetResult(&orderStat).
		SetDoNotParseResponse(false).
		SetHeader("Content-Type", "application/json").
		Get("/api/orders/" + number)

	if err != nil {
		return *orderStat, http.StatusInternalServerError, err // 500
	}

	contentType := resp.Header().Get("Content-Type")

	if resp.StatusCode() == http.StatusTooManyRequests && contentType == "text/plain" { // http.StatusTooManyRequests 429
		delayTime := resp.Header().Get("Retry-After")
		dTime, err := strconv.Atoi(delayTime)
		if err == nil {
			var mutter sync.Mutex // установка wait429 - everybody sleeps until this
			mutter.Lock()
			Time429 = time.Now().Add(time.Duration(dTime) * time.Second)
			mutter.Unlock()
			//		time.Sleep(time.Duration(dTime) * time.Second)
		}
	}
	return *orderStat, resp.StatusCode(), nil
}
