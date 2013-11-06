package main

import (
	"./parser"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type FeedCache struct {
	HandleIterator  int
	RecycledHandles []int
	Feeds           map[string][]int
}

func (me *FeedCache) NextHandle() int {
	if len(me.RecycledHandles) > 0 {
		retval := me.RecycledHandles[0]
		me.RecycledHandles = me.RecycledHandles[1:]
		return retval
	} else {
		me.HandleIterator++
		return me.HandleIterator
	}
}

func (me *FeedCache) RecycleHandle(h int) {
	me.RecycledHandles = append(me.RecycledHandles, h)
}

func (me *FeedCache) RecycleHandles(h []int) {
	for _, v := range h {
		me.RecycledHandles = append(me.RecycledHandles, v)
	}
}

func main() {
	in, _ := ioutil.ReadFile("feeds.txt")
	feeds, err := parser.Parse(string(in))
	if err != nil {
		fmt.Println(err)
	}

	var cache FeedCache
	cache.Feeds = make(map[string][]int)

	cacheIn, _ := ioutil.ReadFile("FeedCache.json")
	println(string(cacheIn))
	json.Unmarshal(cacheIn, cache.Feeds)

	for _, v := range feeds {
		handles := cache.Feeds[v.Name]
		if v.FamilySize > len(handles) {
			println("Add feed")
			diff := v.FamilySize - len(handles)
			for i := 0; i < diff; i++ {
				handles = append(handles, cache.NextHandle())
			}
		} else if v.FamilySize < len(handles) {
			println("Removed feed")
			diff := len(handles) - v.FamilySize
			cache.RecycleHandles(handles[:diff])
			handles = handles[diff:]
		}
		cache.Feeds[v.Name] = handles
	}

	data, _ := json.MarshalIndent(cache, "", "   ")
	ioutil.WriteFile("FeedCache.json", data, 0644)
}
