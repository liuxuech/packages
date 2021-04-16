package huobi

import (
	"github.com/huobirdcenter/huobi_golang/config"
	"github.com/huobirdcenter/huobi_golang/pkg/client"
	"github.com/huobirdcenter/huobi_golang/pkg/model/market"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"strings"
)

const (
	AccessKey = "a492fa4d-1bb9a53f-5b17b5d6-fr2wer5t6y"
	SecretKey = "67acbffd-378db9f5-16e491c3-9614d"
)

type HuoBi struct {
	Host      string
	MarketCli *client.MarketClient

	History map[string]*Ticker
}

type Ticker struct {
	Symbol     string          // 名称
	FirPrice   decimal.Decimal // 第一次价格
	PrePrice   decimal.Decimal // 上一次价格
	CurPrice   decimal.Decimal // 当前价格
	CurrentPer float64         // 本次增长百分比
	TotalPer   float64         // 总的增长百分比
	Weigh      int             // 权重
}

func NewHuoBi() *HuoBi {
	var (
		hb   HuoBi
		mCli client.MarketClient
	)

	hb.Host = config.Host
	hb.MarketCli = mCli.Init(hb.Host)

	hb.History = make(map[string]*Ticker)

	return &hb
}

// GetAllTickers 记录最高值
func (hb *HuoBi) GetAllTickers() (tickers []market.SymbolCandlestick, err error) {
	// 保留usdt

	results, err := hb.MarketCli.GetAllSymbolsLast24hCandlesticksAskBid()
	if err != nil {
		return nil, errors.Wrap(err, "获取所有交易失败")
	}

	for _, item := range results {
		if strings.HasSuffix(item.Symbol, "usdt") {
			tickers = append(tickers, item)
		}
	}

	return
}

func (hb *HuoBi) UpdateHistory(tickers []market.SymbolCandlestick) {
	for _, item := range tickers {
		if tk, ok := hb.History[item.Symbol]; !ok {
			hb.History[item.Symbol] = &Ticker{
				Symbol:     item.Symbol,
				FirPrice:   item.Close,
				PrePrice:   item.Close,
				CurPrice:   item.Close,
				CurrentPer: 0,
				TotalPer:   0,
				Weigh:      1,
			}
		} else {
			tk.PrePrice = tk.CurPrice
			tk.CurPrice = item.Close

			v1, _ := tk.CurPrice.Sub(tk.PrePrice).Mul(tk.PrePrice).Value()
			v2, _ := tk.CurPrice.Sub(tk.FirPrice).Mul(tk.FirPrice).Value()

			tk.CurrentPer = v1.(float64)
			tk.TotalPer = v2.(float64)
		}
	}
}

// CalculateWeigh 计算权重
func (hb *HuoBi) CalculateWeigh() {

}

// 买入、卖出
