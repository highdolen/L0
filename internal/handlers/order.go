package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/highdolen/L0/internal/cache"
	"github.com/highdolen/L0/internal/database"
)

type OrderHandler struct {
	cache *cache.OrderCache
	repo  *database.OrderRepository
}

func NewOrderHandler(cache *cache.OrderCache, repo *database.OrderRepository) *OrderHandler {
	return &OrderHandler{
		cache: cache,
		repo:  repo,
	}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["order_uid"]

	// Сначала проверяем кэш
	if cachedOrder, found := h.cache.Get(uid); found {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cachedOrder); err != nil {
			http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
		}
		return
	}

	// Если в кэше нет, идем в базу
	order, err := h.repo.GetOrderByUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Ошибка при получении заказа: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if order == nil {
		http.Error(w, "Заказ не найден", http.StatusNotFound)
		return
	}

	// Добавляем в кэш для будущих запросов
	h.cache.Set(uid, *order)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}
