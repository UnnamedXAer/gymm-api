package main

import (
	"github.com/unnamedxaer/gymm-api/endpoints/http"
)

var app = http.App{}

// func TestMain(m *testing.M) {
// 	godotenv.Overload(".test.env")
// 	log.Println(" [TestMain] start " + os.Getenv("ENV"))

// 	logger := zerolog.New(os.Stdout)
// 	logger.Info().Msg(time.Now().Local().String() + "-> App starts, env = " + os.Getenv("ENV"))
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		panic("environment variable 'PORT' is not set")
// 	}
// 	dbName := os.Getenv("DB_NAME")
// 	if dbName == "" {
// 		panic("environment variable 'DB_NAME' is not set")
// 	}
// 	mongoURI := os.Getenv("MONGO_URI")
// 	if mongoURI == "" {
// 		panic("environment variable 'MONGO_URI' is not set")
// 	}

// 	db, err := repositories.GetDatabase(&logger, mongoURI, dbName)
// 	if err != nil {
// 		panic(err)
// 	}
// 	repositories.CreateCollections(&logger, db)
// 	usersColl := repositories.GetCollection(&logger, db, "users")
// 	usersRepo := users.NewRepository(&logger, usersColl)
// 	validate := validation.New()
// 	app := http.NewServer(&logger, usersRepo, validate)
// 	app.AddHandlers()

// 	clearTables("mongo")

// 	code := m.Run()

// 	log.Println(" [TestMain] end")
// 	os.Exit(code)
// }

// func clearTables(repoType string) {
// 	log.Println(" [clearTables] start - " + repoType)
// 	switch repoType {
// 	case "mongo":
// 		clearTablesMongo()
// 	default:
// 		log.Fatalln("unsupported repo type")
// 	}
// }

// func clearTablesMongo() {
// 	log.Println(" [clearTablesMongo] start")
// 	db, err := getDb("mongo")
// 	mongoDB := db.(*mongo.Database)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	usersColl := mongoDB.Collection("users")
// 	// results, err := usersColl.DeleteMany(context.Background(), bson.M{})
// 	results, err := usersColl.Find(context.Background(), bson.M{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println(" [clearTablesMongo] results " + fmt.Sprintf("%v", results))
// }

// func getDb(provider string) (interface{}, error) {

// 	switch provider {
// 	case "mongo":
// 		return getMongoDB()
// 	}
// 	return nil, fmt.Errorf("unsupported provider '%s'", provider)
// }

// func getMongoDB() (*mongo.Database, error) {
// 	log.Println("[getMongoDB] start")
// 	uri := os.Getenv("MONGO_URI")
// 	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
// 	if err != nil {
// 		return nil, err

// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	err = client.Connect(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dbName, ok := os.LookupEnv("DB_NAME")
// 	if ok == false {
// 		log.Panic("db name not specified")
// 	}
// 	log.Println("DB name: ", dbName)

// 	return client.Database(dbName), nil
// }
