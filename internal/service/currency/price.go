package currency

import (
	"context"
	"fmt"
	"gexabyte/internal/model"
	"reflect"
	"strconv"
	"time"
)

// TODO: добавить проверку и обработку в хендлере если symbol не существует, написать функцию которая это проверит
func (s *Currency) GetCurrentPrices(ctx context.Context, symbols ...string) ([]model.GetCurrencyPriceDTO, error) {
	dbSymbols, err := s.List(ctx)
	if err != nil {
		return nil, err
	}

	allSymbols := make([]string, 0, len(dbSymbols))
	symbolID := make(map[string]int, len(dbSymbols))
	for _, curr := range dbSymbols {
		allSymbols = append(allSymbols, curr.Symbol)
		symbolID[curr.Symbol] = curr.ID
	}
	for _, symbol := range symbols {
		if _, ok := symbolID[symbol]; !ok {
			allSymbols = append(allSymbols, symbol)
		}
	}

	startReqTime := time.Now().UnixMilli()
	symbolPrice, err := s.fetchCurrentPrices(ctx, allSymbols...)
	if err != nil {
		return nil, err
	}

	{ // update all prices and save to db and update ticker
		saveDB := make([]model.CurrencyPrice, 0, len(symbolPrice))
		for symbol, id := range symbolID { // save only which tracked
			saveDB = append(saveDB, model.CurrencyPrice{
				CurrencyID: id,
				Price:      symbolPrice[symbol],
				Time:       startReqTime,
			})
		}
		if err := s.CreatePrice(ctx, saveDB...); err != nil {
			return nil, err
		}
		s.priceCheckTicker.Reset(s.priceCheckInterval)
	}

	result := make([]model.GetCurrencyPriceDTO, 0, len(symbolPrice))
	for _, symbol := range symbols {
		result = append(result, model.GetCurrencyPriceDTO{
			Symbol: symbol,
			Price:  symbolPrice[symbol],
			Time:   startReqTime,
		})
	}

	return result, err
}

func (s *Currency) fetchCurrentPrices(ctx context.Context, symbols ...string) (map[string]float64, error) {
	prices := make(map[string]float64, len(symbols))

	type task struct {
		symbol string
		price  float64
		err    error
	}

	c, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	taskFuncs := []doTaskFunc{}
	for _, symbol := range symbols {
		taskFuncs = append(taskFuncs, func() interface{} {
			price, err := s.fetchCurrentPrice(ctx, symbol)
			if err != nil {
				return task{err: err}
			}

			return task{
				symbol: symbol,
				price:  price,
			}
		})
	}

	priceStream := s.taskResultStream(ctx, taskFuncs...)
	for i := 0; i < len(symbols); i++ {
		select {
		case <-c.Done():
			return nil, context.DeadlineExceeded
		case out, ok := <-priceStream:
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

			prices[res.symbol] = res.price
		}
	}

	return prices, nil
}

func (s *Currency) fetchCurrentPrice(ctx context.Context, symbol string) (price float64, err error) {
	res, err := s.binanceClient.NewTickerPriceService().Symbol(symbol).Do(ctx)
	if err != nil {
		return 0, err
	}

	fmt.Println("TICKER: ", res, err)

	return strconv.ParseFloat(res.Price, 64)
}

func (s *Currency) CreatePrice(ctx context.Context, rates ...model.CurrencyPrice) error {
	return s.currencyPriceRepo.Create(ctx, rates...)
}
func (s *Currency) ListPrices(ctx context.Context) ([]model.CurrencyPrice, error) {
	return s.currencyPriceRepo.List(ctx)
}
