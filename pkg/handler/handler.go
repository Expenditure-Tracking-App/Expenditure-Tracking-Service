package handler

import (
	"github.com/patrickmn/go-cache"
	"time"
)

// --- Global Cache ---
var c = cache.New(5*time.Minute, 10*time.Minute)
