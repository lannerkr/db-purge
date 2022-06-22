package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func purgedb(coll *mongo.Collection) {

	loc := time.FixedZone("KST", +9*60*60)
	currentTime := time.Now().In(loc)
	filter := bson.D{}

	sub := bson.M{"$subtract": []interface{}{currentTime, "$last_login"}}
	day := bson.M{"$divide": []interface{}{sub, 1000 * 60 * 60 * 24}}
	trunc := bson.M{"$trunc": day}
	updays := bson.M{"days": trunc}
	update := bson.D{{"$set", updays}}

	_, err := coll.UpdateMany(context.TODO(), filter, mongo.Pipeline{update})
	if err != nil {
		panic(err)
	}

	log.Println("days since last login are updated at : " + currentTime.Format("2006-01-02 15:04:05"))

}

func oldusers(coll *mongo.Collection) {
	days := configuration.CheckDays
	f1 := bson.D{{"days", bson.D{{"$gte", days}}}, {"enabled", "True"}}

	filter := bson.D{{"$match", f1}}

	result, err := coll.Aggregate(context.TODO(), mongo.Pipeline{filter})
	if err != nil {
		panic(err)
	}

	var olduser []bson.D
	if err := result.All(context.TODO(), &olduser); err != nil {
		panic(err)
	}

	disabledUsers := updateUser(olduser)
	//= []DisabledUsers{{"store-101", "True"}}
	//disabledUsers = updateUser(olduser)

	dbStatusUpdate(coll, disabledUsers)

}

func dbStatusUpdate(coll *mongo.Collection, disabledUsers []DisabledUsers) {

	for _, users := range disabledUsers {
		user := users.users
		status := users.status

		filter := bson.M{"user_name": user}
		upStatus := bson.M{"enabled": status}
		update := bson.M{"$set": upStatus}

		_, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			panic(err)
		}

		log.Println("user [" + user + "] enabled status is updated to [" + status + "] !!!")
	}
}
