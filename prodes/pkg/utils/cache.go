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
	data map[string]CacheEntry
	mu   sync.Mutex
}

// NewCache crea una nueva instancia de Cache
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheEntry),
	}
}

// Set agrega un valor a la caché con un tiempo de expiración en segundos
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

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
	if !exists {
		return nil, false
	}

	// Verificar si la entrada ha expirado
	if time.Now().Unix() > entry.Expiration {
		delete(c.data, key) // Limpiar la entrada caducada
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
