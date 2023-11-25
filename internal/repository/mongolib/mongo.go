package mongolib

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongo struct {
	Client *mongo.Client
	Driver *mongo.Database
}

const (
	MongoType = "mongo"
)

func NewDriver() *IMongo {
	return &IMongo{}
}
func (c *IMongo) Connect(url string, dbName string) error {
	fmt.Println("mongo connect")
	var ctx = context.TODO()
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	c.Client = client
	c.Driver = client.Database(dbName)
	return nil
}
func (c *IMongo) DriverType() string {
	return MongoType
}
func (c *IMongo) GetDriverImplementation() interface{} {
	return c
}

func (c *IMongo) Disconnect() error {
	return c.Client.Disconnect(context.TODO())
}
