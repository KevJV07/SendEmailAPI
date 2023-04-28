package transactions

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

// Transaction represents a single transaction with a date and value.
type Transaction struct {
	Date  string
	Value float64
}

// MonthlyStat represents the monthly statistics for a set of transactions.
type MonthlyStat struct {
	Balance      float64
	Transactions int
	CreditTotal  float64
	DebitTotal   float64
	CreditCount  int
	DebitCount   int
}

// PerformTransactionCalculations reads data from the provided string, parses the transactions, calculates the monthly stats, and total balance.
func PerformTransactionCalculations(data string) (map[string]*MonthlyStat, float64, error) {
	// Convert the data string to an io.Reader
	dataReader := strings.NewReader(data)
	transactionsList, err := readAndParseCSVFromReader(dataReader)
	if err != nil {
		return nil, 0, err
	}

	monthlyStats := getMonthlyStats(transactionsList)

	// Calculate the total balance
	totalBalance := 0.0
	for _, transaction := range transactionsList {
		totalBalance += transaction.Value
	}

	return monthlyStats, totalBalance, nil
}

// readAndParseCSVFromReader reads and parses CSV data from the provided io.Reader.
func readAndParseCSVFromReader(reader io.Reader) ([]Transaction, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var transactions []Transaction
	for _, record := range records[1:] {
		value, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, Transaction{
			Date:  record[1],
			Value: value,
		})
	}

	return transactions, nil
}

// getMonthlyStats calculates the monthly statistics for a list of transactions.
func getMonthlyStats(transactions []Transaction) map[string]*MonthlyStat {
	monthlyStats := make(map[string]*MonthlyStat)

	for _, transaction := range transactions {
		month := getMonth(transaction.Date)
		stats, ok := monthlyStats[month]
		if !ok {
			stats = &MonthlyStat{}
			monthlyStats[month] = stats
		}

		stats.Balance += transaction.Value
		stats.Transactions++

		if transaction.Value > 0 {
			stats.CreditTotal += transaction.Value
			stats.CreditCount++
		} else {
			stats.DebitTotal += transaction.Value
			stats.DebitCount++
		}
	}

	return monthlyStats
}

// getMonth extracts the month from a given date string in the format "MM/DD/YYYY".
func getMonth(date string) string {
	return strings.Split(date, "/")[0]
}
