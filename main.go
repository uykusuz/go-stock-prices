package main

import (
	"context"
	"flag"
	"fmt"
	"go-stock-prices/algos"
	"go-stock-prices/model"
	"go-stock-prices/ports"
	"go-stock-prices/repos"
	"go-stock-prices/repos/rtwitter"
	"os"

	"github.com/pkg/errors"
)

func main() {
	ctx := context.Background()

	help := flag.Bool("help", false, "print a help message and exit")
	symbol := flag.String("symbol", "AAPL", "stock symbol to query")
	from := flag.String("from", "2014-02-21", "start date")
	to := flag.String("to", "2015-02-20", "end date")
	useMockFinData := flag.Bool("mockfindata", false, "use mock financial data as input")

	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	err := run(ctx, *symbol, *from, *to, *useMockFinData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n\n", err)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf("Usage: %v [options]\n\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func run(ctx context.Context, symbol string, from string, to string, useMockFinData bool) error {
	var finData ports.FinancialDataRepository
	if useMockFinData {
		finData = repos.NewMockFinancialDataRepository()
	} else {
		finData = repos.NewFinancialDataRepository()
	}

	twitterClient, err := rtwitter.NewTwitterClient()
	if err != nil {
		return err
	}

	notifications := rtwitter.NewTwitterNotificationRepository(twitterClient)

	if len(symbol) < 1 {
		return fmt.Errorf("empty symbol")
	}

	fromDate, err := model.ParseDateFromString(from)
	if err != nil {
		return errors.Wrap(err, "failed to parse from time")
	}

	toDate, err := model.ParseDateFromString(to)
	if err != nil {
		return errors.Wrap(err, "failed to parse to time")
	}

	candles, err := finData.GetCandles(ctx, symbol, fromDate, toDate, model.SortDirectionAscending)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve candles")
	}

	performance, err := algos.CalculatePerformance(candles, symbol, fromDate, toDate)
	if err != nil {
		return err
	}

	if err := notifications.NotifyPerformance(ctx, performance); err != nil {
		return err
	}

	fmt.Println("Successfully notified user of performance.")
	return nil
}
