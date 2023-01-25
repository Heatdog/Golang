package itemdatamongo

import (
	"context"
	itemdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ itemdata.ItemData = (*itemDataMongo)(nil)

type itemDataMongo struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewItemDataMongo(collection *mongo.Collection, ctx context.Context) *itemDataMongo {
	return &itemDataMongo{collection: collection, ctx: ctx}
}
