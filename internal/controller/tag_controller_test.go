package controller

import (
	queries "account/config/sqlc"
	"account/internal/jwt_helper"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTagCreate(t *testing.T) {
	cleanup := setUpTestCase(t)
	defer cleanup(t)
	tc := TagController{}
	tc.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	reqBody := gin.H{
		"name": "test",
		"kind": "in_come",
		"sign": "😈",
	}
	bytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(
		"POST",
		"/api/v1/tags",
		strings.NewReader(string(bytes)),
	)
	r.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func TestTagCreateWithUser(t *testing.T) {
	cleanup := setUpTestCase(t)
	defer cleanup(t)

	tc := TagController{}
	tc.RegisterRoutes(r.Group("/api"))

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
		"name": "test",
		"kind": "in_come",
		"sign": "😈",
	}
	bytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(
		"POST",
		"/api/v1/tags",
		strings.NewReader(string(bytes)),
	)
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + auth},
	}

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resBody struct {
		Resource queries.Tag
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, u.ID, resBody.Resource.UserID)
	assert.Equal(t, "test", resBody.Resource.Name)
	assert.Nil(t, resBody.Resource.DeletedAt)
}
