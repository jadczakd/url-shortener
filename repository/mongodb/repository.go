package mongo

import (
	"context"
	"time"

	"github.com/jadczakd/url-shortener/shortener"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client     *mongo.Client
	database   string
	timeout    time.Duration
	collection string
}

func newMongoClient(url string, timeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewRedirectRepository(url, db string, timeout int) (shortener.RedirectRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(timeout) * time.Second,
		database: db,
		// collections could be abstracted into a separate wrapper should we need more entities
		collection: "redirects",
	}
	client, err := newMongoClient(url, timeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepository.")
	}
	repo.client = client
	return repo, nil
}

func (r *mongoRepository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	redirect := &shortener.Redirect{}
	collection := r.client.Database(r.database).Collection(r.collection)
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Mongo.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Mongo.Redirect.Find")
	}
	return redirect, nil
}

func (r *mongoRepository) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection(r.collection)
	_, err := collection.InsertOne(ctx, bson.M{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	},
	)
	if err != nil {
		return errors.Wrap(err, "repository.Mongo.Redirect.Store")
	}
	return nil
}
