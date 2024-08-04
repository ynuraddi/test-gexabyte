package currency

import (
	"context"
	"fmt"
	"gexabyte/internal/model"
	"reflect"
	"strconv"
	"time"
)

// делаю параллельные запросы потому что клиент бинанса который я использую в кейсе когда len(symbols)>1 все равно возвращает только 1
func (s *Currency) GetStat24H(ctx context.Context, symbols ...string) ([]model.GetCurrencyStat24HDTO, error) {
	return s.fetchStats24H(ctx, symbols...)
}

func (s *Currency) fetchStats24H(ctx context.Context, symbols ...string) ([]model.GetCurrencyStat24HDTO, error) {
	result := make([]model.GetCurrencyStat24HDTO, 0, len(symbols))
	c, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	type task struct {
		item model.GetCurrencyStat24HDTO
		err  error
	}

	taskFuncs := make([]doTaskFunc, 0, len(symbols))
	for _, symbol := range symbols {
		taskFuncs = append(taskFuncs, func() interface{} {
			res, err := s.fetchStat24H(ctx, symbol)
			if err != nil {
				return task{err: err}
			}

			return task{item: res}
		})
	}

	statStream := s.taskResultStream(ctx, taskFuncs...)
	for i := 0; i < len(symbols); i++ {
		select {
		case <-c.Done():
			return nil, context.DeadlineExceeded
		case out, ok := <-statStream:
			if !ok {
				return nil, fmt.Errorf("read from close channel")
			}

			res, ok := out.(task)
			if !ok {
				return nil, fmt.Errorf("incorrect type data: %s", reflect.TypeOf(out))
			}
			if res.err != nil {
				return nil, res.err
			}

			result = append(result, res.item)
		}
	}

	return result, nil
}

func (s *Currency) fetchStat24H(ctx context.Context, symbol string) (model.GetCurrencyStat24HDTO, error) {
	res, err := s.binanceClient.NewTicker24hrService().Symbol(symbol).Do(ctx)
	if err != nil {
		return model.GetCurrencyStat24HDTO{}, err
	}

	openPrice, err := strconv.ParseFloat(res.OpenPrice, 64)
	if err != nil {
		return model.GetCurrencyStat24HDTO{}, err
	}
	lastPrice, err := strconv.ParseFloat(res.LastPrice, 64)
	if err != nil {
		return model.GetCurrencyStat24HDTO{}, err
	}
	highPrice, err := strconv.ParseFloat(res.HighPrice, 64)
	if err != nil {
		return model.GetCurrencyStat24HDTO{}, err
	}
	lowPrice, err := strconv.ParseFloat(res.LowPrice, 64)
	if err != nil {
		return model.GetCurrencyStat24HDTO{}, err
	}

	return model.GetCurrencyStat24HDTO{
		Symbol:    res.Symbol,
		OpenPrice: openPrice,
		LastPrice: lastPrice,
		HighPrice: highPrice,
		LowPrice:  lowPrice,
		OpenTime:  int64(res.OpenTime),
		CloseTime: int64(res.CloseTime),
	}, nil
}
