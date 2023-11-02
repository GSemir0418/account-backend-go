package controller

import (
	"account/api"
	queries "account/config/sqlc"
	"account/internal/jwt_helper"
	"encoding/json"
	"fmt"
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

func TestTagUpdateWithUser(t *testing.T) {
	cleanup := setUpTestCase(t)
	defer cleanup(t)

	tc := TagController{}
	tc.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()

	u, _ := q.CreateUser(c, "1@qq.com")
	tag, err := q.CreateTag(c, queries.CreateTagParams{
		Kind:   "in_come",
		Name:   "test",
		Sign:   "😈",
		UserID: u.ID,
	})
	if err != nil {
		log.Fatalln(err)
		return
	}

	reqBody := gin.H{
		"name": "newTest",
		"sign": "👟",
	}
	bytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(
		"PATCH",
		fmt.Sprintf("/api/v1/tags/%d", tag.ID),
		strings.NewReader(string(bytes)),
	)

	auth, err := jwt_helper.GenerateJWT(int(u.ID))
	if err != nil {
		log.Fatalln(err)
	}
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
	assert.Equal(t, "newTest", resBody.Resource.Name)
	assert.Equal(t, "👟", resBody.Resource.Sign)
	assert.Equal(t, "in_come", resBody.Resource.Kind)
	assert.Nil(t, resBody.Resource.DeletedAt)
}

func TestTagDelete(t *testing.T) {
	cleanup := setUpTestCase(t)
	defer cleanup(t)

	tc := TagController{}
	tc.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()

	u, _ := q.CreateUser(c, "1@qq.com")
	tag, err := q.CreateTag(c, queries.CreateTagParams{
		Kind:   "in_come",
		Name:   "testDelete",
		Sign:   "😈",
		UserID: u.ID,
	})
	if err != nil {
		log.Fatalln(err)
		return
	}

	req, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/tags/%d", tag.ID),
		nil,
	)

	auth, err := jwt_helper.GenerateJWT(int(u.ID))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + auth},
	}

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestPagedTags(t *testing.T) {
	cleanup := setUpTestCase(t)
	defer cleanup(t)

	it := TagController{}
	it.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/tags?page=3&page_size=5",
		nil,
	)

	// 登录
	u, _ := q.CreateUser(c, "1@qq.com")
	logIn(t, u.ID, req)

	for i := 0; i < int(13); i++ {
		if _, err := q.CreateTag(c, queries.CreateTagParams{
			UserID: u.ID,
			Name:   "testPaged",
			Kind:   "expenses",
		}); err != nil {
			t.Error(err)
		}
	}

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resBody api.GetPagedTagsResponse
	if err := json.Unmarshal([]byte(w.Body.String()), &resBody); err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, 3, len(resBody.Resources))
}
