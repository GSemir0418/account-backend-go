package controller

import (
	"account/api"
	queries "account/config/sqlc"
	"account/internal/jwt_helper"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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
		"/api/v1/items",
		nil,
	)

	// 登录
	u, _ := q.CreateUser(c, "1@qq.com")
	logIn(t, u.ID, req)

	for i := 0; i < int(20); i++ {
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
	assert.Equal(t, 10, len(resBody.Resources))
}
