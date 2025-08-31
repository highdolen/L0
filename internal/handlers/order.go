package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/highdolen/L0/internal/service"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["order_uid"]

	// Проверяем параметр для принудительного обновления кеша
	forceRefresh := r.URL.Query().Get("refresh") == "true"

	var result *service.OrderResult
	var err error

	if forceRefresh {
		// Принудительно обновляем заказ из БД
		result, err = h.orderService.GetOrderByUIDWithRefresh(r.Context(), uid)
	} else {
		// Обычное получение (сначала из кеша, потом из БД)
		result, err = h.orderService.GetOrderByUID(r.Context(), uid)
	}

	if err != nil {
		// Проверяем, является ли ошибка "заказ не найден"
		if errors.Is(err, errors.New("заказ не найден")) {
			http.Error(w, "Заказ не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка при получении заказа: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", "application/json")
	if result.FromCache {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
	}

	// Возвращаем только заказ, без метаданных
	if err := json.NewEncoder(w).Encode(result.Order); err != nil {
		http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
	}
}

// GetCacheStats — получить статистику кеша
func (h *OrderHandler) GetCacheStats(w http.ResponseWriter, r *http.Request) {
	stats := h.orderService.GetCacheStats()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
	}
}

// InvalidateCache — инвалидировать весь кеш или конкретный заказ
func (h *OrderHandler) InvalidateCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["order_uid"]

	var err error
	if uid != "" {
		// Инвалидируем конкретный заказ
		err = h.orderService.InvalidateCache(uid)
		if err != nil {
			http.Error(w, "Ошибка при инвалидации кеша: "+err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{"message": fmt.Sprintf("Кеш для заказа %s инвалидирован", uid)}, http.StatusOK)
	} else {
		// Инвалидируем весь кеш
		err = h.orderService.InvalidateAllCache()
		if err != nil {
			http.Error(w, "Ошибка при полной инвалидации кеша: "+err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{"message": "Весь кеш инвалидирован"}, http.StatusOK)
	}
}

func writeJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}
