package persistence

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ortymid/projector/projector"
)

type bsonUser struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

type MongoUserRepo struct {
	db *mongo.Database
}

func NewMongoUserRepo(db *mongo.Database) MongoUserRepo {
	return MongoUserRepo{db: db}
}

func (repo MongoUserRepo) All(ctx context.Context) ([]projector.User, error) {
	cur, err := repo.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(ctx)

	var users []projector.User

	for cur.Next(ctx) {
		var result bsonUser

		err := cur.Decode(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// log.Println(result)
		user, err := projector.NewUser(result.Name)

		users = append(users, user)
	}
	if err = cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

func (repo MongoUserRepo) Create(ctx context.Context, user projector.User) (err error) {
	_, err = repo.db.Collection("users").InsertOne(ctx, bson.M{
		"name": user.Name,
	})

	return err
}
