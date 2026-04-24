package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"order-crm/config"
	"order-crm/internal/model"
	"order-crm/internal/router"
	"order-crm/pkg/database"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testRouter http.Handler
var testToken = os.Getenv("TEST_TOKEN")

func TestMain(m *testing.M) {
	config.InitEnv()
	db := database.InitDB()
	testRouter = router.InitGinRouter(db)

	// Запуск тестов
	code := m.Run()

	os.Exit(code)
}

func TestGetUsers(t *testing.T) {
	// Создаем запрос
	req, err := http.NewRequest("GET", "/api/users", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+testToken)

	// Создаем рекордер и вызываем хендлер
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUsersUnauthorized(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/users", nil)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUserById(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/users/1", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

var createdUserId int

func TestCreateUser(t *testing.T) {
	body := `{
		"login":"testuser_create",
		"fio":"Test User Create",
		"password":"12345678",
		"id_role":3
	}`

	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code) // или http.StatusOK

	// распарсим ответ
	var response struct {
		ID int `json:"id"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	createdUserId = response.ID
}

func TestDeleteUser(t *testing.T) {
	if createdUserId == 0 {
		t.Fatal("createdUserId is empty")
	}
	// Удаляем созданного ранее пользователя
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/users/%d", createdUserId), nil)
	req.Header.Set("Authorization", "Bearer "+testToken)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetClients(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/clients", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

var createdOrder *model.Order

func TestCreateOrder(t *testing.T) {
	body := `{
		"label":"Test Order",
		"id_client":1,
		"items":[{"label":"Горох","amount":20},{"label":"Бананы","amount":40}]
	}`

	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err := json.Unmarshal(w.Body.Bytes(), &createdOrder)
	assert.NoError(t, err)

	assert.Equal(t, float64(60), createdOrder.Amount, "Значение поля amount не совпадает")
}

func TestAddOrderItem(t *testing.T) {
	if createdOrder == nil {
		t.Fatal("createdOrder is empty")
	}
	body := `{
		"label":"Test Order Item",
		"amount":999
	}`

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("/api/orders/%d/items", createdOrder.ID),
		bytes.NewBufferString(body),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteOrder(t *testing.T) {
	if createdOrder == nil {
		t.Fatal("createdOrder is empty")
	}
	// Удаляем созданный ранее заказ
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/orders/%d", createdOrder.ID), nil)
	req.Header.Set("Authorization", "Bearer "+testToken)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
