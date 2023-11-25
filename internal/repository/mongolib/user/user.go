package user

import (
	"context"
	"fmt"
	"slices"

	"github.com/mallvielfrass/fmc"
	"github.com/ponyCorp/rebornPony/internal/repository/mongolib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	db *mongolib.IMongo
	ct *mongo.Collection
}

func Init(driver *mongolib.IMongo) (*User, error) {
	collectionName := "user"
	names, err := driver.Driver.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	if !slices.Contains(names, collectionName) {
		command := bson.M{"create": collectionName}
		var result bson.M
		if err := driver.Driver.RunCommand(context.TODO(), command).Decode(&result); err != nil {
			fmt.Printf("%s: %+v\n", fmc.WhoCallerIs(), err)
			return nil, err
		}
	}

	return &User{
		db: driver,
		ct: driver.Driver.Collection(collectionName),
	}, nil
}
func (u *User) GetName() string {
	return u.db.DriverType()
}
func (u *User) GetUserByID(id string) (*User, error) {
	return nil, nil
}

type UserScheme struct {
	Name string
	Id   string
}

// Create
func (u *User) CreateUser(name string) (UserScheme, error) {
	res, err := u.ct.InsertOne(context.TODO(), bson.M{"name": name})
	if err != nil {
		return UserScheme{}, err
	}
	usr, isExist, err := u.getUserById(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return UserScheme{}, err
	}
	if !isExist {
		return UserScheme{}, fmt.Errorf("User not created")
	}
	return UserScheme{
		Name: usr.Name,
		Id:   usr.Id,
	}, nil
}

func (u *User) getUserById(id primitive.ObjectID) (UserScheme, bool, error) {

	var user UserScheme
	err := u.ct.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return UserScheme{}, false, err
	}
	return user, true, nil
}
