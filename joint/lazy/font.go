package lazy

import (
	"image"

	"github.com/codeation/lru"

	"github.com/codeation/impress/driver"
)

type keySplit struct {
	text   string
	edge   int
	indent int
}

type font struct {
	driver.Fonter
	cacheSplit *lru.SyncLRU[keySplit, []string]
	cacheSize  *lru.SyncLRU[string, image.Point]
}

func (a *app) NewFont(height int, attributes map[string]string) driver.Fonter {
	f := &font{
		Fonter: a.Driver.NewFont(height, attributes),
	}
	f.cacheSplit = lru.NewSyncLRU(128, f.splitByKey)
	f.cacheSize = lru.NewSyncLRU(1024, f.sizeByString)
	return f
}

func (f *font) splitByKey(key keySplit) ([]string, error) {
	return f.Fonter.Split(key.text, key.edge, key.indent), nil
}

func (f *font) sizeByString(text string) (image.Point, error) {
	return f.Fonter.Size(text), nil
}

func (f *font) Split(text string, edge int, indent int) []string {
	output, _ := f.cacheSplit.Get(keySplit{
		text:   text,
		edge:   edge,
		indent: indent,
	})
	return output
}

func (f *font) Size(text string) image.Point {
	output, _ := f.cacheSize.Get(text)
	return output
}

func (f *font) Unwrap() driver.Fonter { return f.Fonter }
