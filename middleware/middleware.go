package middleware

import (
	"context"
	//  I know protocol buffer is more efficient but i didn't have time to set message in an appropriate format
	//  So i used json instead
	"encoding/json"
	"fmt"
	"go-api/models"

	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"bytes"
	"io/ioutil"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// collection object/instance
var collection *mongo.Collection

// create connection with mongo db (Atlas cluster) for demonstration purpose and to get the job done faster :p
func init() {
	loadTheEnv()
	createDBInstance()
}

func loadTheEnv() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func createDBInstance() {
	// DB connection string
	connectionString := os.Getenv("DB_URI")

	// Database Name
	dbName := os.Getenv("DB_NAME")

	// Collection name
	collName := os.Getenv("DB_COLLECTION_NAME")

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database(dbName).Collection(collName)

	fmt.Println("Collection instance created!")
}

// GetAllUser get all the user route
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload := getAllUser()
	json.NewEncoder(w).Encode(payload)
}

// CreateUser create user route
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var user models.UserList
	_ = json.NewDecoder(r.Body).Decode(&user)
	insertOneUser(user)
	json.NewEncoder(w).Encode(user)
}

// Set User Password  route
func SetUserPass(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var user models.UserList
	_ = json.NewDecoder(r.Body).Decode(&user)
	setUserPass(user)
	json.NewEncoder(w).Encode(user)
}

// DeleteUser delete one user route
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	deleteOneUser(params["id"])
	json.NewEncoder(w).Encode(params["id"])
	// json.NewEncoder(w).Encode("User not found")
}

// DeleteAllUser delete all users route
func DeleteAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	count := deleteAllUser()
	json.NewEncoder(w).Encode(count)
	// json.NewEncoder(w).Encode("User not found")

}

// get all user from the DB and return it
func getAllUser() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

// Insert one user in the DB
func insertOneUser(user models.UserList) {
	insertResult, err := collection.InsertOne(context.Background(), user)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a Single Record ", insertResult.InsertedID)
}

// user complete method, update user's profile_id to 0000
func setUserPass(user models.UserList) {
	fmt.Println(user)
	email := user.Email
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"password": user.Password}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

// delete one user from the DB, delete by ID
func deleteOneUser(user string) {
	fmt.Println(user)
	id, _ := primitive.ObjectIDFromHex(user)
	filter := bson.M{"_id": id}
	d, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted Document", d.DeletedCount)
}

// delete all the users from the DB
func deleteAllUser() int64 {
	d, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted Document", d.DeletedCount)
	return d.DeletedCount
}

//  Get Tayara User  route
// Post Body : {"number":"95421449"}
func GetTayaraUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var phone models.Phone
	_ = json.NewDecoder(r.Body).Decode(&phone)
	payload := GetTayaraUserByPhone(phone.Number)
	json.NewEncoder(w).Encode(payload)
}

// Get unprotected user data using  from tayara.tn api and send them back with a formatted version
// phone must start with +216
func GetTayaraUserByPhone(phone string) map[string]string {
	httpposturl := "https://www.tayara.tn/core/marketplace.MarketPlace/GetUserByPhoneNumber"
	var result map[string]string
	result = make(map[string]string)
	// fmt.Println("URL:>", httpposturl)

	var jsonStr = []byte(string(0) + string(0) + string(0) + string(0) + string(14) + string(10) + string(12) + phone)

	req, err := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/grpc-web+proto")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	if bodyString == "" {
		return result
	}

	result["profile_id"] = bodyString[35:72]
	result["info"] = bodyString[91:]

	return result
}

//  Get Tayara User  route
func GetAllTayaraUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetAllTayaraUsers")

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	httpposturl := "https://www.tayara.tn/core/marketplace.MarketPlace/ListUsers"
	fmt.Println("URL:>", httpposturl)
	// This Weird Code is used because i don't know protoc messages structure used by tayara/api , and i did'nt have the time to implement protocol buffer in the proper way
	var BytesStr = []byte(string(0) + string(0) + string(0) + string(0) + "*:" + string(2) + string(8) + "dZ$0bea4865-b40c-4b35-9edb-12b4dcae1ce9")

	req, err := http.NewRequest("POST", httpposturl, bytes.NewBuffer(BytesStr))
	req.Header.Set("Content-Type", "application/grpc-web+proto")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	payload := bodyString
	json.NewEncoder(w).Encode(payload)
}

//  Get Tayara User  route
// Post Body : {"profile_id":"$bla-blaaaa-blaa","password":"xyz123"}
func ChangeUserPass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var user models.UserList
	_ = json.NewDecoder(r.Body).Decode(&user)
	fmt.Println("profile_id")
	fmt.Println(user.Profile_id)
	fmt.Println("user.Password")
	fmt.Println(user.Password)
	newpass := changeUserPass(user.Profile_id, user.Password)
	json.NewEncoder(w).Encode("done new pass : " + newpass)
}

// SHOULD BE USED WITH CAUTION !!!
// Changing USer Password on  tayara.tn  using profile_id obtained in previous method
// No Harm was done to anywone except my own profile (for testing purpose only)
func changeUserPass(profile_id string, new_pass string) string {
	httpposturl := "https://www.tayara.tn/iam/identity.Identity/ChangeUserPassword"

	fmt.Println("URL:>", httpposturl)
	// This Weird Code is used because i don't know protoc messages structure used by tayara/api , and i did'nt have the time to implement protocol buffer in the proper way
	var jsonStr = []byte(string(0) + string(0) + string(0) + string(0) + "0" + string(10) + profile_id + string(18) + string(8) + new_pass)

	req, err := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/grpc-web+proto")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// fmt.Println(resp.Body)

	return new_pass
}
