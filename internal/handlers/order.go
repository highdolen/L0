package handlers

import (
	"encoding/json"
	"fmt"
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

	// Проверяем параметр для принудительного обновления кеша
	forceRefresh := r.URL.Query().Get("refresh") == "true"

	if forceRefresh {
		// Принудительно обновляем запись в кеше из БД
		if err := h.cache.Refresh(r.Context(), uid, h.repo); err != nil {
			log.Printf("Ошибка обновления кеша для заказа %s: %v", uid, err)
		}
	}

	// Сначала проверяем кэш
	if cachedOrder, found := h.cache.Get(uid); found {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "HIT")
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
	log.Printf("Заказ %s загружен из БД и кеширован", uid)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
	}
}

// GetCacheStats — получить статистику кеша
func (h *OrderHandler) GetCacheStats(w http.ResponseWriter, r *http.Request) {
	stats := h.cache.GetStats()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
	}
}

// InvalidateCache — инвалидировать весь кеш или конкретный заказ
func (h *OrderHandler) InvalidateCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["order_uid"]
	
	if uid != "" {
		// Инвалидируем конкретный заказ
		h.cache.Invalidate(uid)
		log.Printf("Инвалидирован кеш для заказа %s", uid)
		writeJSON(w, map[string]string{"message": fmt.Sprintf("Кеш для заказа %s инвалидирован", uid)}, http.StatusOK)
	} else {
		// Инвалидируем весь кеш
		h.cache.InvalidateAll()
		log.Println("Весь кеш инвалидирован")
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
