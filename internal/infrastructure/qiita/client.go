package qiita

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
	token  string
	userId string
}

func NewQiitaClient(token, userID string) *QiitaClient {
	return &QiitaClient{token: token, userId: userID}
}

// FetchStocksAll fetches stock articles from Qiita API
func (c *QiitaClient) FetchStocksAll(ctx context.Context) ([]QiitaItem, error) {
	allItems := make([]QiitaItem, 0)
	page := 1
	for {
		items, hasMore, err := c.fetchStocks(ctx, page)
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, items...)
		if !hasMore {
			break
		}
		page++
	}
	return allItems, nil
}

// fetchStockPage fetches a single page of stock articles from Qiita API
func (c *QiitaClient) fetchStocks(ctx context.Context, page int) ([]QiitaItem, bool, error) {
	url := fmt.Sprintf("https://qiita.com/api/v2/users/%s/stocks", c.userId)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, false, err
	}

	// パラメータを設定
	q := req.URL.Query()
	q.Set("per_page", "100")
	q.Set("page", fmt.Sprintf("%d", page))
	req.URL.RawQuery = q.Encode()

	// 認証ヘッダーを設定
	req.Header.Set("Authorization", "Bearer "+c.token)

	// APIリクエストを実行
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	// ステータスコードのチェック(200-299の範囲外はエラーとする)
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		body, _ := io.ReadAll(resp.Body)
		return nil, false, fmt.Errorf("Qiita API error: %s, body: %s", resp.Status, string(body))
	}
	// レスポンスボディの読み込み
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}
	// JSON解析
	var items []QiitaItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, false, err
	}
	return items, len(items) == 100, nil
}
