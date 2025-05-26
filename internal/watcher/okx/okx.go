package okx

import (
	"context"
	"errors"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"strconv"
	"time"
)

// GetOkxTrxUsdtRate 获取 Okx TRX 汇率 https://www.okx.com/zh-hans/trade-spot/trx-usdt
func (w *Redis) GetOkxTrxUsdtRate() (float64, error) {
	var url = "https://www.okx.com/priapi/v5/market/candles?instId=TRX-USDT&before=1727143156000&bar=4H&limit=1&t=" + cast.ToString(time.Now().UnixNano())
	var client = http.Client{Timeout: time.Second * 5}
	var req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("referer", "https://www.okx.com/zh-hans/trade-spot/trx-usdt")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {

		return 0, errors.New("okx resp error:" + err.Error())
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {

		return 0, errors.New("okx resp status code:" + strconv.Itoa(resp.StatusCode))
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {

		return 0, errors.New("okx resp read body error:" + err.Error())
	}

	result := gjson.ParseBytes(all)
	if result.Get("data").Exists() {
		var data = result.Get("data").Array()
		if len(data) > 0 {
			w.Redis.Set(context.Background(), "usdt-cny", data[0].Get("1").Float(), 0)

			return data[0].Get("1").Float(), nil
		}
	}

	return 0, errors.New("okx resp json data not found")
}

// GetOkxUsdtCnySellPrice   Okx  C2C快捷交易 USDT出售 实时汇率
func (w *Redis) GetOkxUsdtCnySellPrice() (float64, error) {
	var t = strconv.Itoa(int(time.Now().Unix()))
	var okxApi = "https://www.okx.com/v4/c2c/express/price?crypto=USDT&fiat=CNY&side=sell&t=" + t
	client := http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", okxApi, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {

		return 0, errors.New("okx resp error:" + err.Error())
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {

		return 0, errors.New("okx resp status code:" + strconv.Itoa(resp.StatusCode))
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {

		return 0, errors.New("okx resp read error:" + err.Error())
	}

	result := gjson.ParseBytes(all)
	if result.Get("error_code").Int() != 0 {

		return 0, errors.New("json parse error:" + result.Get("error_message").String())
	}

	if result.Get("data.price").Exists() {
		var _ret = result.Get("data.price").Float()
		if _ret <= 0 {
			return 0, errors.New("okx resp json data.price <= 0")
		}

		w.Redis.Set(context.Background(), "usdt-cny", cast.ToFloat64(_ret), 0)

		return cast.ToFloat64(_ret), nil
	}

	return 0, errors.New("okx resp json data.price not found")
}
