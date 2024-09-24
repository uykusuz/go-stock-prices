package repos

import (
	"context"
	"encoding/csv"
	"fmt"
	"go-stock-prices/model"
	"go-stock-prices/ports"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type financialDataRepository struct {
}

func NewFinancialDataRepository() ports.FinancialDataRepository {
	return &financialDataRepository{}
}

func (r *financialDataRepository) GetCandles(
	ctx context.Context,
	symbol string,
	from *model.Date,
	to *model.Date,
	mDirection model.SortDirection,
) ([]*model.Candle, error) {
	data, err := r.requestCandles(symbol, from, to, mDirection)
	if err != nil {
		return nil, err
	}

	return parseCandles(data)
}

func parseCandles(data string) ([]*model.Candle, error) {
	reader := csv.NewReader(strings.NewReader(data))

	var candles []*model.Candle

	lineNo := 0

	for {
		lineNo++

		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, errors.Wrap(err, "failed to read response line")
		}

		if len(record) < 1 {
			continue
		}

		if lineNo == 1 {
			// skip header
			continue
		}

		if len(record) != 7 {
			skippedError(lineNo, err, "invalid number of columns")
			continue
		}

		date, err := model.ParseDateFromString(record[0])
		if err != nil {
			skippedError(lineNo, err, "date column is not a valid date")
			continue
		}

		open, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			skippedError(lineNo, err, "open column is not a valid float")
			continue
		}

		high, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			skippedError(lineNo, err, "high column is not a valid float")
			continue
		}

		low, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			skippedError(lineNo, err, "low column is not a valid float")
			continue
		}

		closePrice, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			skippedError(lineNo, err, "close column is not a valid float")
			continue
		}

		volume, err := strconv.ParseInt(record[5], 10, 64)
		if err != nil {
			skippedError(lineNo, err, "volume column is not a valid float")
			continue
		}

		unadjustedVolume, err := strconv.ParseInt(record[6], 10, 64)
		if err != nil {
			skippedError(lineNo, err, "unadjustedVolume column is not a valid float")
			continue
		}

		candle := &model.Candle{
			Date:             date,
			Open:             open,
			High:             high,
			Low:              low,
			Close:            closePrice,
			Volume:           volume,
			UnadjustedVolume: unadjustedVolume,
		}

		candles = append(candles, candle)
	}

	return candles, nil
}

func (*financialDataRepository) requestCandles(
	symbol string,
	from *model.Date,
	to *model.Date,
	mDirection model.SortDirection,
) (string, error) {
	baseURLString := "https://fn44uqlbc5ulyst6l4treh5tem0klihr.lambda-url.eu-central-1.on.aws/historical/"
	baseURLString += symbol

	dateFrom := dateFromModel(from)
	dateTo := dateFromModel(to)
	direction, err := sortDirectionFromModel(mDirection)
	if err != nil {
		return "", err
	}

	baseURL, err := url.Parse(baseURLString)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse base URL of backend")
	}

	queryParams := baseURL.Query()
	queryParams.Set("date_from", dateFrom)
	queryParams.Set("date_to", dateTo)
	queryParams.Set("direction", direction)

	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create request for backend")
	}

	req.Header.Add("Content-Type", "application/csv")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status. Expected %v, but is %v", http.StatusOK, resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}

	return string(responseBody), nil
}

func sortDirectionFromModel(direction model.SortDirection) (string, error) {
	switch direction {
	case model.SortDirectionAscending:
		return "asc", nil
	case model.SortDirectionDescending:
		return "desc", nil
	}

	return "asc", fmt.Errorf("invalid sort direction: %v", direction)
}

func dateFromModel(d *model.Date) string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func skippedError(lineNo int, err error, message string) {
	fmt.Fprintf(os.Stderr, "[line %v] failed parsing line: %v. Details: %v", lineNo, message, err)
}
