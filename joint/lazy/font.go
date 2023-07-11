package lazy

import (
	"image"
	"log"

	"github.com/codeation/impress/driver"
	lru "github.com/hashicorp/golang-lru/v2"
)

type keySplit struct {
	text   string
	edge   int
	indent int
}

type font struct {
	driver.Fonter
	cacheSplit *lru.Cache[keySplit, []string]
	cacheSize  *lru.Cache[string, image.Point]
}

func (a *app) NewFont(height int, attributes map[string]string) driver.Fonter {
	cacheSplit, err := lru.New[keySplit, []string](128)
	if err != nil {
		log.Printf("lru.New: %v", err)
	}
	cacheSize, err := lru.New[string, image.Point](1024)
	if err != nil {
		log.Printf("lru.New: %v", err)
	}
	return &font{
		Fonter:     a.Driver.NewFont(height, attributes),
		cacheSplit: cacheSplit,
		cacheSize:  cacheSize,
	}
}

func (f *font) Split(text string, edge int, indent int) []string {
	key := keySplit{
		text:   text,
		edge:   edge,
		indent: indent,
	}
	if value, ok := f.cacheSplit.Get(key); ok {
		return value
	}
	output := f.Fonter.Split(text, edge, indent)
	f.cacheSplit.Add(key, output)
	return output
}

func (f *font) Size(text string) image.Point {
	if value, ok := f.cacheSize.Get(text); ok {
		return value
	}
	output := f.Fonter.Size(text)
	f.cacheSize.Add(text, output)
	return output
}

func (f *font) ID() int {
	type ider interface{ ID() int }
	if id, ok := f.Fonter.(ider); ok {
		return id.ID()
	}
	return 0
}
