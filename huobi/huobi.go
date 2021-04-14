package main

import (
	"fmt"
	"github.com/huobirdcenter/huobi_golang/config"
	"github.com/huobirdcenter/huobi_golang/logging/applogger"
	"github.com/huobirdcenter/huobi_golang/pkg/client"
	"github.com/huobirdcenter/huobi_golang/pkg/client/marketwebsocketclient"
	"github.com/huobirdcenter/huobi_golang/pkg/model/market"
	"strconv"
	"strings"
	"time"
)

type Calculate struct {
	Inc bool
	Val string
}

func main() {

	var store = make(map[string]*Calculate)

	getAllTickers(store)
}

func getAllTickers(m map[string]*Calculate) {
	var c client.MarketClient

	marketClient := c.Init(config.Host)

	ticker := time.NewTicker(time.Second * 1)

	all, _ := marketClient.GetAllSymbolsLast24hCandlesticksAskBid()
	for _, item := range all {
		v, _ := item.Close.Value()
		vv := v.(string)
		m[item.Symbol] = &Calculate{
			Inc: false,
			Val: vv,
		}
	}

	counter := make(map[string]int)

	go func() {
		ticker := time.NewTicker(time.Second * 5)

		for _ = range ticker.C {

			kvSlice := make([]string, 0)

			for k, v := range counter {
				if v > 5 || v < -5 {
					kvSlice = append(kvSlice, fmt.Sprintf("%s(%d)", k, v))
				}
			}
			fmt.Println(strings.Join(kvSlice, " / "))
		}
	}()

	for _ = range ticker.C {
		all, _ := marketClient.GetAllSymbolsLast24hCandlesticksAskBid()

		for _, item := range all {

			if !strings.HasSuffix(item.Symbol, "usdt") {
				continue
			}

			v, _ := item.Close.Value()
			vv := v.(string)

			res := m[item.Symbol]

			resNumber, _ := strconv.ParseFloat(res.Val, 12)
			vvNumber, _ := strconv.ParseFloat(vv, 12)

			if resNumber-vvNumber > 0 {
				if ress, ok := counter[item.Symbol]; ok {
					counter[item.Symbol] = ress + 1
				} else {
					counter[item.Symbol] = 1
				}
			} else if resNumber-vvNumber < 0 {
				if ress, ok := counter[item.Symbol]; ok {
					counter[item.Symbol] = ress - 1
				}
			}
			res.Inc = resNumber-vvNumber > 0
			res.Val = vv
		}
	}
}

func GetCurrencys() {
	c := new(client.CommonClient).Init(config.Host)
	resp, _ := c.GetCurrencys()

	for _, item := range resp {
		fmt.Println(item)
	}
}

func GetSymbols() {
	c := new(client.CommonClient).Init(config.Host)
	resp, _ := c.GetSymbols()

	for _, item := range resp {
		fmt.Println(item.Symbol)
	}
}
func http() {
	//config.Host = "api-aws.huobi.pro"
	//// Get the timestamp from Huobi server and print on console
	//c := new(client.CommonClient).Init(config.Host)
	//resp, err := c.GetCurrencys()
	//
	//if err != nil {
	//	applogger.Error("Get timestamp error: %s", err)
	//} else {
	//	applogger.Info("Get timestamp: %+v", resp)
	//}

	// 获取账户列表
	//c := new(client.AccountClient).Init(huobi.AccessKey, huobi.SecretKey, config.Host)
	//resp, err := c.GetAccountInfo()
	//if err != nil {
	//	applogger.Error("Get account error: %s", err)
	//} else {
	//	applogger.Info("Get account, count=%d", len(resp))
	//	for _, result := range resp {
	//		applogger.Info("account: %+v", result)
	//	}
	//}
	//
	//// 获取账户余额
	//for _, ac := range resp {
	//	balance, _ := c.GetAccountBalance(strconv.Itoa(int(ac.Id)))
	//	fmt.Printf("%+v\n", balance)
	//	fmt.Println("-----------------------")
	//}

	//// 获取币余额
	//balance, err := c.GetAccountBalance("20184357")
	//fmt.Println(len(balance.List), err)
	//
	//mC := new(client.MarketClient).Init(config.Host)
	//
	//candlestick, err := mC.GetCandlestick("btcusdt", market.GetCandlestickOptionalRequest{
	//	Period: "1day",
	//	Size:   200,
	//})
	//fmt.Println(candlestick, err)
}

func ws() {
	client := new(marketwebsocketclient.Last24hCandlestickWebSocketClient).Init(config.Host)

	client.SetHandler(
		// Connected handler
		func() {
			client.Request("bttusdt", "1608")

			client.Subscribe("bttusdt", "1608")
		},
		// Response handler
		func(resp interface{}) {
			candlestickResponse, ok := resp.(market.SubscribeLast24hCandlestickResponse)
			if ok {
				if &candlestickResponse != nil {
					if candlestickResponse.Tick != nil {
						t := candlestickResponse.Tick
						applogger.Info("Candlestick update, id: %d, count: %v, volume: %v [%v-%v-%v-%v]",
							t.Id, t.Count, t.Vol, t.Open, t.Close, t.Low, t.High)
					}

					if candlestickResponse.Data != nil {
						t := candlestickResponse.Data
						applogger.Info("Candlestick data, id: %d, count: %v, volume: %v [%v-%v-%v-%v]",
							t.Id, t.Count, t.Vol, t.Open, t.Close, t.Low, t.High)
					}
				}
			} else {
				applogger.Warn("Unknown response: %v", resp)
			}
		},
	)
	// Connect to the server and wait for the handler to handle the response
	client.Connect(true)
	select {}
}
