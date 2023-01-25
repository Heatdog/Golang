package itemdatamap

import (
	itemdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData"
	"sync"
)

var _ itemdata.ItemData = (*itemDataMap)(nil)

type itemDataMap struct {
	data map[string]itemdata.Post
	mux  *sync.RWMutex
}

func NewItemDataMap() *itemDataMap {
	return &itemDataMap{
		data: make(map[string]itemdata.Post, 10),
		mux:  &sync.RWMutex{},
	}
}
