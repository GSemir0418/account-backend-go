package controller

import (
	"account/api"
	queries "account/config/sqlc"
	"account/internal/jwt_helper"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nav-inc/datetime"
	"github.com/stretchr/testify/assert"
)

func TestItemController(t *testing.T) {
	cleanup := setUpTestCase(t)
	defer cleanup(t)
	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	reqBody := gin.H{
		"user_id":     1,
		"amount":      100,
		"kind":        "in_come",
		"happened_at": time.Now(),
		"tag_ids":     []int32{1, 2, 3},
	}
	bytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(
		"POST",
		"/api/v1/items",
		strings.NewReader(string(bytes)),
	)
	r.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func TestItemControllerWithUser(t *testing.T) {
	cleanup := setUpTestCase(t)
	defer cleanup(t)
	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()

	u, err := q.CreateUser(c, "1@qq.com")
	if err != nil {
		log.Fatalln(err)
	}
	auth, err := jwt_helper.GenerateJWT(int(u.ID))
	if err != nil {
		log.Fatalln(err)
	}

	reqBody := gin.H{
		"amount":      100,
		"kind":        "in_come",
		"happened_at": time.Now(),
		"tag_ids":     []int32{1, 2, 3},
	}
	bytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(
		"POST",
		"/api/v1/items",
		strings.NewReader(string(bytes)),
	)
	// 也可以直接写
	// req, _ := http.NewRequest(
	// 	"POST",
	// 	"/api/v1/items",
	// 	strings.NewReader(`{
	// 		"amount": 100,
	// 		"kind": "expenses",
	// 		"happened_at": "2020-01-01T00:00:00Z",
	// 		"tag_ids": [1, 2, 3]
	// 	}`),
	// )

	req.Header = http.Header{
		"Authorization": []string{"Bearer " + auth},
	}

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resBody struct {
		Resource queries.Item
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, u.ID, resBody.Resource.UserID)
}

func TestPagedItems(t *testing.T) {
	cleanup := setUpTestCase(t)
	defer cleanup(t)
	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/items?page=3&page_size=5",
		nil,
	)

	// 登录
	u, _ := q.CreateUser(c, "1@qq.com")
	logIn(t, u.ID, req)

	for i := 0; i < int(13); i++ {
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: time.Now(),
		}); err != nil {
			t.Error(err)
		}
	}

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resBody api.GetPagedItemsResponse
	if err := json.Unmarshal([]byte(w.Body.String()), &resBody); err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, 3, len(resBody.Resources))
}

func TestBalance(t *testing.T) {
	cleanup := setUpTestCase(t)
	defer cleanup(t)

	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/items/balance",
		nil,
	)

	u, _ := q.CreateUser(c, "1@qq.com")
	logIn(t, u.ID, req)

	for i := 0; i < int(10); i++ {
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: time.Now(),
		}); err != nil {
			t.Error(err)
		}
	}

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	var j api.GetBalanceResponse
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, 10000*10, j.Expenses)
	assert.Equal(t, 0, j.Income)
	assert.Equal(t, -10000*10, j.Balance)
}

func TestBalanceWithTime(t *testing.T) {
	cleanup := setUpTestCase(t)
	defer cleanup(t)

	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/items/balance?happened_after="+url.QueryEscape("2023-09-29T00:00:00+0800")+
			"&happened_before="+url.QueryEscape("2023-10-01T00:00:00+0800"),
		nil,
	)

	u, _ := q.CreateUser(c, "1@qq.com")
	logIn(t, u.ID, req)

	// 构造不同时间的数据
	for i := 0; i < int(3); i++ {
		d, _ := datetime.Parse("2023-09-29T00:00:00+08:00", time.Local)
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: d,
		}); err != nil {
			t.Error(err)
		}
	}
	for i := 0; i < int(3); i++ {
		d, _ := datetime.Parse("2023-09-30T23:59:59+08:00", time.Local)
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "in_come",
			TagIds:     []int32{1},
			HappenedAt: d,
		}); err != nil {
			t.Error(err)
		}
	}
	for i := 0; i < int(3); i++ {
		d, _ := datetime.Parse("2023-10-30T23:59:59+08:00", time.Local)
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "in_come",
			TagIds:     []int32{1},
			HappenedAt: d,
		}); err != nil {
			t.Error(err)
		}
	}

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	var j api.GetBalanceResponse
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, 10000*3, j.Expenses)
	assert.Equal(t, 10000*3, j.Income)
	assert.Equal(t, 0, j.Balance)
}

func TestSummary(t *testing.T) {
	cleanup := setUpTestCase(t)
	defer cleanup(t)

	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))

	qs := url.Values{
		// "happened_after":  []string{"2023-09-01T00:00:00+08:00"},
		// "happened_before": []string{"2023-10-31T00:00:00+08:00"},
		"kind":     []string{"xxx"},
		"group_by": []string{"happened_at"},
	}.Encode()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/items/summary?"+qs,
		nil,
	)

	u, _ := q.CreateUser(c, "1@qq.com")
	logIn(t, u.ID, req)

	for i := 0; i < 3; i++ {
		d, _ := datetime.Parse("2023-10-01T00:00:00+08:00", time.Local)
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: d,
		}); err != nil {
			t.Error(err)
		}
	}

	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	// body := w.Body.String()
	// var j api.ErrorResponse
	// if err := json.Unmarshal([]byte(body), &j); err != nil {
	// 	t.Error("json.Unmarshal fail", err)
	// }
	// assert.Equal(t, 3, len(j.Errors))
	// assert.Equal(t, 10000*3, j.Income)
	// assert.Equal(t, 0, j.Balance)
}
