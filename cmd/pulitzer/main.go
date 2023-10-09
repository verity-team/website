package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/verity-team/dws/internal/pulitzer"
)

func main() {
	s := gocron.NewScheduler(time.UTC)
	s.Every("1m").Do(getETHPrice)
	s.StartBlocking()
}

func calculateAveragePrice(prices ...decimal.Decimal) (decimal.Decimal, error) {
	if len(prices) == 0 {
		return decimal.Zero, errors.New("no prices provided")
	}

	// Calculate the sum of prices
	sum := decimal.Zero
	for _, price := range prices {
		sum = sum.Add(price)
	}

	// Calculate the average by dividing the sum by the number of prices
	average := sum.Div(decimal.NewFromInt(int64(len(prices))))

	return average, nil
}

func getETHPrice() (decimal.Decimal, error) {
	// Create a wait group to synchronize the goroutines
	var wg sync.WaitGroup

	// Create channels to receive the results from the goroutines
	binanceResultCh := make(chan decimal.Decimal, 1)
	krakenResultCh := make(chan decimal.Decimal, 1)
	bitfinexResultCh := make(chan decimal.Decimal, 1)

	// Start the first goroutine to fetch Ethereum price from Binance
	wg.Add(1)
	go func() {
		defer wg.Done()
		ethereumPrice, err := pulitzer.GetWeightedAvgPriceFromBinance()
		if err != nil {
			log.Errorf("Error fetching Ethereum price from Binance:", err)
			return
		}
		binanceResultCh <- ethereumPrice
		close(binanceResultCh)
	}()

	// Start the second goroutine to fetch Ethereum price from Kraken
	wg.Add(1)
	go func() {
		defer wg.Done()
		ethereumPrice, err := pulitzer.GetKrakenETHPrice()
		if err != nil {
			log.Errorf("Error fetching Ethereum price from Kraken:", err)
			return
		}
		krakenResultCh <- ethereumPrice
		close(krakenResultCh)
	}()

	// Start the third goroutine to fetch Ethereum price from bitfinex
	wg.Add(1)
	go func() {
		defer wg.Done()
		ethereumPrice, err := pulitzer.GetBitfinexETHPrice()
		if err != nil {
			log.Errorf("Error fetching Ethereum price from bitfinex:", err)
			return
		}
		bitfinexResultCh <- ethereumPrice
		close(bitfinexResultCh)
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Receive and print the results from the channels
	binancePrice, binanceOk := <-binanceResultCh
	krakenPrice, krakenOk := <-krakenResultCh
	bitfinexPrice, bitfinexOk := <-bitfinexResultCh

	if binanceOk {
		log.Infof("Ethereum Price from Binance: $%s USD\n", binancePrice.StringFixed(2))
	} else {
		err := errors.New("failed to get Ethereum Price from Binance")
		log.Error(err)
		return decimal.Zero, err
	}

	if krakenOk {
		log.Infof("Ethereum Price from Kraken: $%s USD\n", krakenPrice.StringFixed(2))
	} else {
		err := errors.New("failed to get Ethereum Price from Kraken")
		log.Error(err)
		return decimal.Zero, err
	}

	if bitfinexOk {
		log.Infof("Ethereum Price from bitfinex: $%s USD\n", bitfinexPrice.StringFixed(2))
	} else {
		err := errors.New("failed to get Ethereum Price from bitfinex")
		log.Error(err)
		return decimal.Zero, err
	}

	av, err := calculateAveragePrice(binancePrice, krakenPrice, bitfinexPrice)

	if err != nil {
		log.Error(err)
		return decimal.Zero, err
	} else {
		log.Infof("average price: $%s\n", av.StringFixed(2))
	}
	return av, nil
}

func checkPriceDeviation(prices ...decimal.Decimal) error {
	if len(prices) < 3 {
		return errors.New("at least 3 prices are required for deviation check")
	}

	// Check deviation between all pairs of prices
	for i := 0; i < len(prices); i++ {
		for j := i + 1; j < len(prices); j++ {
			price1 := prices[i]
			price2 := prices[j]

			// Calculate the percentage deviation
			deviation := price1.Sub(price2).Div(price1).Abs()

			if deviation.GreaterThan(decimal.NewFromFloat(0.10)) {
				return fmt.Errorf("price deviation between %s and %s is greater than 10%%", price1.StringFixed(2), price2.StringFixed(2))
			}
		}
	}

	return nil
}
