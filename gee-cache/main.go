package main

import "gee-cache/lru"

func main() {
	cache := lru.NewCache(2560, nil)

}
