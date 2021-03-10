package chrome

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Condo) insert(ctx context.Context, collection *mongo.Collection) error {
	d := bson.D{
		{
			"name", c.name,
		},
		{
			"address", c.address,
		},
		{
			"district", c.district,
		},
		{
			"tenure", c.tenure,
		},
		{
			"developer", c.developer,
		},
		{
			"url", c.url,
		},
		{
			"facility", c.facility,
		},
		{
			"facString", c.facString,
		},
	}
	_, err := collection.InsertOne(ctx, d)
	if err != nil {
		return err
	}
	return nil
}
