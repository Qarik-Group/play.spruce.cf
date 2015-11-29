package main

import (
	"fmt"
	"log"
	"math/rand"
)

var Memory map[string]*Merge

func init() {
	Memory = make(map[string]*Merge, 20)
}

func Store(m *Merge) (string, error) {
	key := fmt.Sprintf("%x", rand.Int())
	Memory[key] = m
	log.Printf("stored %s in memory for sharing", key)
	return key, nil
}

func Retrieve(key string) (*Merge, error) {
	if m, ok := Memory[key]; ok {
		return m, nil
	}
	return nil, fmt.Errorf("not found")
}
