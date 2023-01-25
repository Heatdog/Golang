package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gitlab.com/vk-go/lectures-2022-2/pkg/middleware"
	itemdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData"
	itemdatamongo "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData/itemDataMongo"
	userdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData"
	"gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData/userDataMySQL"
	"gitlab.com/vk-go/lectures-2022-2/pkg/server"
	"gitlab.com/vk-go/lectures-2022-2/pkg/service"
	"gitlab.com/vk-go/lectures-2022-2/pkg/session"
	"gitlab.com/vk-go/lectures-2022-2/pkg/session/sessionManagerMySQL"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile)
	key := []byte("super secret")
	ctx := context.Background()
	db := mySQLInit("root", "123", "dbMySQL", "3306", "webDB", logger)
	collection, client := mongoInit(ctx, "dbMongo", "27017", "webDB", "posts")
	defer closeDB(client)
	var usData userdata.UserData = userdatamysql.NewUserDataMySql(db)
	var itmData itemdata.ItemData = itemdatamongo.NewItemDataMongo(collection, ctx)
	var sesManager session.SesManager = sessionmanagermysql.NewSessionManagerMySQL(db)
	serv := service.NewService(usData, itmData, sesManager, key)
	srv := server.NewServer(serv, logger)
	mid := middleware.NewMiddleware(key, sesManager, logger)
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./static/html/")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.Use(mid.Panic)
	r.Use(mid.AccessLog)

	routerSub := r.PathPrefix("/api").Subrouter()
	routerSub.HandleFunc("/register", srv.Register).Methods("POST")
	routerSub.HandleFunc("/login", srv.Login).Methods("POST")
	routerSub.HandleFunc("/posts/", srv.GetPosts).Methods("GET")
	routerSub.HandleFunc("/posts/{category}", srv.GetCategory).Methods("GET")
	routerSub.HandleFunc("/user/{user_login}", srv.GetUser).Methods("GET")
	routerSub.HandleFunc("/post/{post_id}", srv.GetPostID).Methods("GET")

	routerPost := r.PathPrefix("/api").Subrouter()
	routerPost.Use(mid.Auth)
	routerPost.HandleFunc("/posts", srv.CreatePost).Methods("POST")
	routerPost.HandleFunc("/post/{post_id}", srv.CreateComment).Methods("POST")
	routerPost.HandleFunc("/post/{post_id}/{comment_id}", srv.DeleteComment).Methods("DELETE")
	routerPost.HandleFunc("/post/{post_id}/upvote", srv.Upvote).Methods("GET")
	routerPost.HandleFunc("/post/{post_id}/downvote", srv.Downvote).Methods("GET")
	routerPost.HandleFunc("/post/{post_id}/unvote", srv.Unvote).Methods("GET")
	routerPost.HandleFunc("/post/{post_id}", srv.DeletePost).Methods("DELETE")

	r.NotFoundHandler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "static/html/index.html")
	})
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Fatal(err.Error())
	}
}

func mySQLInit(login, password, host, port, dataBase string, logger *log.Logger) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", login, password, host, port, dataBase)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal(err.Error())
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func mongoInit(ctx context.Context, host, port, dataBase, collection string) (*mongo.Collection, *mongo.Client) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))
	if err != nil {
		log.Fatal(err.Error())
	}
	if err = client.Ping(nil, readpref.Primary()); err != nil {
		log.Fatal(err.Error())
	}
	res := client.Database(dataBase).Collection(collection)
	return res, client
}

func closeDB(client *mongo.Client) {
	if err := client.Disconnect(nil); err != nil {
		log.Fatal(err.Error())
	}
}
