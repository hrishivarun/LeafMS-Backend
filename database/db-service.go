package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const URI = "mongodb+srv://jitesh:hrishi3688@cluster0.l6a8xav.mongodb.net/?retryWrites=true&w=majority"

type Database struct {
	Context            context.Context
	Database           *mongo.Database
	EmployeeCollection *mongo.Collection
	LeaveCollection    *mongo.Collection
	PublicHolidays     *mongo.Collection
}

func ConnectDB() *Database {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(URI).SetServerAPIOptions(serverAPI)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*100000000)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err = client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()

	database := client.Database("test")
	employeeCollection := database.Collection("employees")
	publicHolidaysCollection := database.Collection("publicHolidays")
	leaveCollection := database.Collection("leaves")

	// Send a ping to confirm a successful connection
	if err := database.RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return &Database{
		Context:            ctx,
		Database:           database,
		EmployeeCollection: employeeCollection,
		LeaveCollection:    leaveCollection,
		PublicHolidays:     publicHolidaysCollection,
	}
}

func (db Database) InsertOne(collectionName string, document any) (*mongo.InsertOneResult, error) {
	collection := db.Database.Collection(collectionName)
	result, err := collection.InsertOne(db.Context, document)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("Could not insert document. Error:-\n\t", err)
		return nil, err
	}
	return result, nil
}

func (db Database) InsertMany(collectionName string, document []any) (*mongo.InsertManyResult, error) {
	collection := db.Database.Collection(collectionName)
	result, err := collection.InsertMany(db.Context, document)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("Could not insert documents. Error:-\n\t", err)
		return nil, err
	}
	return result, nil
}

func (db Database) FindOne(collectionName string, filter bson.D, findOptions *options.FindOneOptions) (bson.Raw, error) {
	var data bson.Raw

	collection := db.Database.Collection(collectionName)
	err := collection.FindOne(db.Context, filter, findOptions).Decode(&data)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("The FindOne query did not return a result. Error:-\n\t", err)
		return nil, err
	}

	return data, nil
}

func (db Database) Find(collectionName string, filter bson.D) ([]bson.Raw, error) {
	var data []bson.Raw

	collection := db.Database.Collection(collectionName)
	resultCursor, err := collection.Find(db.Context, filter)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("The Find query did not return a cursor. Error:-\n\t", err)
		return nil, err
	}

	if err = resultCursor.All(db.Context, &data); err != nil {
		log.Panic("Could not complete the Find query in the database. Error:-\n\t", err)
		return nil, err
	}
	return data, nil
}

func (db Database) UpdateOne(collectionName string, filter bson.D, update any, updateOption *options.UpdateOptions) (*mongo.UpdateResult, error) {
	collection := db.Database.Collection(collectionName)
	res, err := collection.UpdateOne(db.Context, filter, update, updateOption)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("Could not update the entry for filter:- ", filter)
		return nil, err
	}
	return res, nil
}

func (db Database) UpdateMany(collectionName string, filter bson.D, updates any) (*mongo.UpdateResult, error) {
	collection := db.Database.Collection(collectionName)
	res, err := collection.UpdateMany(db.Context, filter, updates)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("Could not update the entries for filter:- ", filter)
		return nil, err
	}
	return res, nil
}

func (db Database) DeleteOne(collectionName string, filter bson.D) (*mongo.DeleteResult, error) {
	collection := db.Database.Collection(collectionName)
	res, err := collection.DeleteOne(db.Context, filter)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("Could not delete entry for filter:- ", filter)
		log.SetPrefix("WARNING: ")
		log.Println("Err:- ", err)
		return nil, err
	}
	return res, nil
}

func (db Database) DeleteMany(collectionName string, filter bson.D) (*mongo.DeleteResult, error) {
	collection := db.Database.Collection(collectionName)
	res, err := collection.DeleteMany(db.Context, filter)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("Could not delete entries for filter:- ", filter)
		log.SetPrefix("WARNING: ")
		log.Println("Err:- ", err)
		return nil, err
	}
	return res, nil
}

func (db Database) Aggregate(collectioName string, pipeline mongo.Pipeline) ([]bson.Raw, error) {
	var data []bson.Raw
	collection := db.Database.Collection((collectioName))
	resultCursor, err := collection.Aggregate(db.Context, pipeline)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("Could not perform aggregate operation with pipeline:- ", pipeline)
		return nil, err
	}

	if err = resultCursor.All(db.Context, &data); err != nil {
		log.Panic("Could not complete the Find query in the database. Error:-\n\t", err)
		return nil, err
	}
	return data, nil
}

var DbConn = ConnectDB()
