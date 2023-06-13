package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type transactionsResponse struct {
	Transactions []*transaction `json:"data"`
}

type transaction struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		AccountID            int32   `json:"accountId"`
		CreatedAt            string  `json:"createdAt"`
		Currency             string  `json:"transactionCurrency"`
		Amount               float64 `json:"transactionAmount"`
		Direction            string  `json:"transactionDirection"`
		Type                 string  `json:"transactionType"`
		ClientReference      string  `json:"clientReference"`
		TransactionReference string  `json:"transactionReference"`
	} `json:"attributes"`
}

func (c *Client) GetAllTransactions(ctx context.Context, accountID string) ([]*transaction, error) {
	var transactions []*transaction

	for page := 0; ; page++ {
		pagedTransactions, err := c.getTransactions(ctx, accountID, page, pageSize)
		if err != nil {
			return nil, err
		}

		if len(pagedTransactions) == 0 {
			break
		}

		transactions = append(transactions, pagedTransactions...)

		if len(pagedTransactions) < pageSize {
			break
		}
	}

	return transactions, nil
}

func (c *Client) getTransactions(ctx context.Context, accountID string, page, pageSize int) ([]*transaction, error) {
	endpoint := fmt.Sprintf("%s/accounts/%s/transactions/find", c.endpoint, accountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("pagesize", strconv.Itoa(pageSize))
	q.Add("pagenumber", fmt.Sprint(page))
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	var transactions transactionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal login response body: %w", err)
	}

	return transactions.Transactions, nil
}