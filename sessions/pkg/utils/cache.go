package utils

import (
	"sync"
	"time"
)

// CacheEntry representa una entrada en la caché con valor y expiración
type CacheEntry struct {
	Value      interface{}
	Expiration int64
}

// Cache es la estructura que mantiene el mapa de la caché
type Cache struct {
	data       map[string]CacheEntry
	mu         sync.Mutex
	cleanupInt time.Duration // Intervalo de limpieza automática
	maxSize    int           // Tamaño máximo de la caché
}

// NewCache crea una nueva instancia de Cache con intervalos de limpieza y tamaño máximo
func NewCache(cleanupInterval time.Duration, maxSize int) *Cache {
	cache := &Cache{
		data:       make(map[string]CacheEntry),
		cleanupInt: cleanupInterval,
		maxSize:    maxSize,
	}

	// Lanzar una goroutine para limpiar las entradas expiradas periódicamente
	go cache.cleanupExpiredEntries()

	return cache
}

// Set agrega un valor a la caché con un tiempo de expiración
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Si la caché excede el tamaño, eliminar la entrada más antigua (o aplicar una política como LRU)
	if len(c.data) >= c.maxSize {
		c.evictEntry()
	}

	expiration := time.Now().Add(duration).Unix()
	c.data[key] = CacheEntry{
		Value:      value,
		Expiration: expiration,
	}
}

// Get obtiene un valor de la caché, si no está expirado
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.data[key]
	if !exists || time.Now().Unix() > entry.Expiration {
		delete(c.data, key)
		return nil, false
	}

	return entry.Value, true
}

// Delete elimina un valor de la caché
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Clear limpia toda la caché
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]CacheEntry)
}

// cleanupExpiredEntries elimina las entradas expiradas en intervalos regulares
func (c *Cache) cleanupExpiredEntries() {
	for {
		time.Sleep(c.cleanupInt)
		c.mu.Lock()
		for key, entry := range c.data {
			if time.Now().Unix() > entry.Expiration {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

// evictEntry elimina la entrada más antigua o aplica una política de reemplazo
func (c *Cache) evictEntry() {
	// Implementa la lógica para eliminar una entrada, por ejemplo LRU, FIFO, etc.
	// Aquí podrías eliminar simplemente la primera entrada encontrada (FIFO)
	for key := range c.data {
		delete(c.data, key)
		break
	}
}
