package qiita

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// QiitaItem represents a minimal Qiita article structure
// 必要に応じてフィールドを追加
type QiitaItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Tags  []struct {
		Name string `json:"name"`
	} `json:"tags"`
}

// QiitaClient handles API requests
type QiitaClient struct {
	Token  string
	UserID string
}

func NewQiitaClient(token, userID string) *QiitaClient {
	return &QiitaClient{Token: token, UserID: userID}
}

// FetchStocks fetches stock articles from Qiita API
func (c *QiitaClient) FetchStocks() ([]QiitaItem, error) {
	// パラメータをつけたい。per_page=100&page=1
	url := fmt.Sprintf("https://qiita.com/api/v2/users/%s/stocks", c.UserID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	q := req.URL.Query()
	q.Set("per_page", "100")
	q.Set("page", "1")
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Qiita API error: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var items []QiitaItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}
	return items, nil
}
