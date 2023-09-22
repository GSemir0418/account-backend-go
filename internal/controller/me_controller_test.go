package controller

import (
	"account/internal/jwt_helper"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeController(t *testing.T) {
	teardown := setUpTestCase(t)
	defer teardown(t)
	mc := MeController{}
	mc.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/me",
		strings.NewReader(""),
	)
	r.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

}

func TestMeWithJWT(t *testing.T) {
	teardown := setUpTestCase(t)
	defer teardown(t)

	u, err := q.CreateUser(c, "1@qq.com")
	if err != nil {
		log.Fatalln(err)
	}
	auth, err := jwt_helper.GenerateJWT(int(u.ID))
	if err != nil {
		log.Fatalln(err)
	}

	mc := MeController{}
	mc.RegisterRoutes(r.Group("/api"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/me",
		strings.NewReader(""),
	)
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + auth},
	}
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resBody struct {
		Resource struct {
			ID    int32  `json:"id"`
			Email string `json:"email"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, u.ID, resBody.Resource.ID)
	assert.Equal(t, u.Email, resBody.Resource.Email)
}
