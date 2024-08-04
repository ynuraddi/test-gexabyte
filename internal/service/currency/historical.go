package currency

import (
	"context"
	"gexabyte/internal/model"
	"math"
	"strconv"
	"time"
)

/*
Чтобы получать корректные значения мне нужно для начала понять по какому принципу берутся интервалы на бинансе.
Они округляются в зависимости от выбранного интервала, нужно потыкать чтобы понять как правильно округлять значения.

Лимиты я буду брать стандартные, максимальный будет как и на бинансе 1000 элементов.
Пагинацию  я собираюсь вычислять путем startTime + interval * (page-1) что-то типа такого

Нужно еще добавить значение максимального кол-ва страниц в ДТО.

Потыкал их интервалы, понял как работают:
На примере легче обьяснить startTime=00:01 endTime=01:00 interval=30m, то выйдет только 1 запись, потому что она начнется с 00:30
Аналогично с часами, месяцами и т.д., только нужно обязательно вычислять все в UTC-0

# Недели начинаются с понедельников, а месяцы с первых чисел

Так, секунды, минуты и часы считать будет легко, потому что можно просто используя ceil
*/

// TODO: написать тесты
func (s *Currency) GetPriceHistorical(ctx context.Context, req model.GetCurrencyPriceHistoricalDTOReq) (*model.GetCurrencyPriceHistoricalDTORes, error) {
	st, mp := s.solvePagination(req.StartTime, req.EndTime, req.Limit, req.Page, req.Interval)
	req.StartTime = st

	res, err := s.binanceClient.KlineService(ctx, req.Symbol, req.Interval, req.StartTime, req.EndTime, req.Limit)
	if err != nil {
		return nil, err
	}

	prices := make([]model.CurrencyPriceInterval, 0, len(res))
	for i := 0; i < len(res); i++ {
		r := res[i]

		openPrice, err := strconv.ParseFloat(r.Open, 64)
		if err != nil {
			return nil, err
		}
		closePrice, err := strconv.ParseFloat(r.Close, 64)
		if err != nil {
			return nil, err
		}
		highPrice, err := strconv.ParseFloat(r.High, 64)
		if err != nil {
			return nil, err
		}
		lowPrice, err := strconv.ParseFloat(r.Low, 64)
		if err != nil {
			return nil, err
		}

		prices = append(prices, model.CurrencyPriceInterval{
			OpenPrice:  openPrice,
			ClosePrice: closePrice,

			HighPrice: highPrice,
			LowPrice:  lowPrice,

			OpenTime:  int64(r.OpenTime),
			CloseTime: int64(r.CloseTime),
		})
	}

	return &model.GetCurrencyPriceHistoricalDTORes{
		Page:    req.Page,
		MaxPage: mp,

		Prices: prices,
	}, nil
}

func (s *Currency) solvePagination(startTime, endTime int64, limit, page int, interval string) (sTime int64, maxPage int) {
	countSkipIntervals := (page - 1) * limit
	i := time.Second.Milliseconds() // interval

	startTime = ceil(startTime, i) * i
	endTime = ceil(endTime, i) * i
	switch interval {
	case
		"1s",
		"1m", "3m", "5m", "15m", "30m",
		"1h", "2h", "4h", "6h", "12h", "1d", "3d":
		i = int64(model.KlineInterval.GetDuration(interval).Milliseconds())
		startTime, endTime = ceil(startTime, i)*i, ceil(endTime, i)*i

		maxPage = int(ceil(endTime-startTime, i*int64(limit)))
		sTime = startTime + int64(countSkipIntervals)*i
		return sTime, maxPage

	case "1w":
		i = int64(model.KlineInterval.GetDuration(interval))
		startTime, endTime = ceilWeek(startTime), ceilWeek(endTime)

		maxPage = int(ceil(endTime-startTime, i))
		sTime = startTime + int64(countSkipIntervals)*i
		return sTime, maxPage

	case "1M":
		startTime, endTime = ceilMonth(startTime), ceilMonth(endTime)
		st, et := time.UnixMilli(startTime), time.UnixMilli(endTime)

		for st.Before(et) {
			maxPage += 1
			st = st.AddDate(0, 1, 0)
		}

		st = time.UnixMilli(startTime)
		for i := 0; i < countSkipIntervals; i++ {
			st = st.AddDate(0, 1, 0)
		}
		return st.UnixMilli(), maxPage
	}

	return
}

func ceil(x, div int64) int64 {
	return int64(math.Ceil(float64(x) / float64(div)))
}

func ceilDay(in int64) int64 {
	t := time.UnixMilli(in)
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)

	if t.After(d) {
		return d.AddDate(0, 0, 1).UnixMilli() // return next day
	}
	return t.UnixMilli() // return 00:00 current day
}

func ceilWeek(in int64) int64 {
	d := time.UnixMilli(ceilDay(in))
	for d.Weekday() != time.Monday {
		d = d.AddDate(0, 0, 1)
	}

	return d.UnixMilli()
}

func ceilMonth(in int64) int64 {
	d := time.UnixMilli(ceilDay(in))
	m := time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.Local)

	if d.After(m) {
		return m.AddDate(0, 1, 0).UnixMilli()
	}
	return m.UnixMilli()
}
