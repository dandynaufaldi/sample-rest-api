package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	circuit "github.com/rubyist/circuitbreaker"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbPass = "pass"
	dbHost = "localhost"
	dbPort = "3306"
	dbName = "blog"
)

var (
	db     *sql.DB
	cb     *circuit.Breaker
	events <-chan circuit.BreakerEvent
)

var (
	incomingReq = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "request_incoming",
		Help: "Total count of incoming request",
	})
	// responseSuccess = prometheus.NewCounter(prometheus.CounterOpts{
	// 	Name: "response_success_count",
	// 	Help: "Total count of success response",
	// })
	// responseFail = prometheus.NewCounter(prometheus.CounterOpts{
	// 	Name: "response_fail_count",
	// 	Help: "Total count of fail response",
	// })
	responseCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_counter",
			Help: "Total amount of response by status code",
		},
		[]string{"code"},
	)
)

// Post represent the post model in database
type Post struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
	// CreatedAt  time.Time `json:"created_at"`
	// ModifiedAt time.Time `json:"modified_at"`
}

func main() {
	cb = circuit.NewConsecutiveBreaker(10)
	events = cb.Subscribe()
	go func() {
		for {
			e := <-events
			log.Println(e)
		}
	}()

	prometheus.MustRegister(incomingReq)
	prometheus.MustRegister(responseCounterVec)
	// prometheus.MustRegister(responseSuccess)
	// prometheus.MustRegister(responseFail)

	dbSource := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8", dbPass, dbHost, dbPort, dbName)
	var err error
	db, err = sql.Open("mysql", dbSource)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "PONG"})
	})
	http.HandleFunc("/posts", handlePost)
	http.HandleFunc("/posts/", handlePostID)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", Logger(http.DefaultServeMux)); err != nil {
		panic(err)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	incomingReq.Inc()
	var err error
	cb.Call(func() error {
		switch r.Method {
		case "POST":
			err = CreatePost(w, r)
		default:
			err = GetPosts(w, r)
		}
		return err
	}, 0)
}

func handlePostID(w http.ResponseWriter, r *http.Request) {
	incomingReq.Inc()
	var err error
	cb.Call(func() error {
		switch r.Method {
		case "PUT":
			err = UpdatePost(w, r)
		default:
			err = GetPostByID(w, r)
		}
		return err
	}, 0)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// responseSuccess.Inc()
	responseCounterVec.WithLabelValues(strconv.Itoa(code)).Inc()
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	// responseFail.Inc()
	responseCounterVec.WithLabelValues(strconv.Itoa(code)).Inc()
	respondWithJSON(w, code, map[string]string{"message": message})
}

func validateMethod(allowed string, incoming string) (string, error) {
	if allowed != incoming {
		message := "Method " + incoming + " not allowed"
		return message, errors.New(message)
	}
	return "", nil
}

// GetPosts fetch post by ID
func GetPosts(w http.ResponseWriter, r *http.Request) error {
	log.Println("Get Posts")
	msg, err := validateMethod("GET", r.Method)
	if err != nil {
		log.Println(msg)
		respondWithError(w, http.StatusBadRequest, msg)
		return err
	}

	rows, err := db.Query("Select id, title, content, author from post")
	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusInternalServerError, "Server error")
		return err
	}
	defer rows.Close()
	result := make([]*Post, 0)
	for rows.Next() {
		temp := new(Post)
		err := rows.Scan(&temp.ID, &temp.Title, &temp.Content, &temp.Author)
		if err != nil {
			log.Println(err)
			respondWithError(w, http.StatusInternalServerError, "Server error")
			return err
		}
		result = append(result, temp)
	}
	respondWithJSON(w, http.StatusOK, result)
	return err
}

// CreatePost create new post
func CreatePost(w http.ResponseWriter, r *http.Request) error {
	log.Println("Create Post")
	msg, err := validateMethod("POST", r.Method)
	if err != nil {
		log.Println(msg)
		respondWithError(w, http.StatusBadRequest, msg)
		return err
	}

	var post Post
	json.NewDecoder(r.Body).Decode(&post)

	stmt, err := db.Prepare("Insert post SET title=?, content=?, author=?")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Server error")
		return err
	}

	res, err := stmt.Exec(post.Title, post.Content, post.Author)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Server error")
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Server error")
		return err
	}
	post.ID = lastID
	respondWithJSON(w, http.StatusCreated, post)
	return err
}

// GetPostByID fetch post by ID
func GetPostByID(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	msg, err := validateMethod("GET", r.Method)
	if err != nil {
		log.Println(msg)
		respondWithError(w, http.StatusBadRequest, msg)
		return err
	}
	param := r.URL.Path
	param = strings.TrimPrefix(param, "/posts/")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "ID parameter must be integer")
		return err
	}
	result := new(Post)
	err = db.QueryRow("Select id, title, content, author from post where id = ?", id).Scan(&result.ID, &result.Title, &result.Content, &result.Author)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Post with ID "+param+" not exist")
		return err
	}
	respondWithJSON(w, http.StatusOK, result)
	return err
}

// UpdatePost update post
func UpdatePost(w http.ResponseWriter, r *http.Request) error {
	msg, err := validateMethod("PUT", r.Method)
	if err != nil {
		log.Println(msg)
		respondWithError(w, http.StatusBadRequest, msg)
		return err
	}
	param := r.URL.Path
	param = strings.TrimPrefix(param, "/posts/")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "ID parameter must be integer")
		return err
	}
	body, err := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Missing request body")
		return err
	}
	var post Post
	json.Unmarshal(body, &post)
	stmt, err := db.Prepare("Update post set title=?, content=?, author=? where id=?")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Server error")
		return err
	}

	_, err = stmt.Exec(post.Title, post.Content, post.Author, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Post with ID "+param+" not exist")
		return err
	}
	post.ID = int64(id)
	respondWithJSON(w, http.StatusOK, post)
	return err
}

// // DeletePost delete post
// func DeletePost(w http.ResponseWriter, r *http.Request) {
// 	msg, err := validateMethod("DELETE", r.Method)
// 	if err != nil {
// 		log.Println(msg)
// 		respondWithError(w, http.StatusBadRequest, msg)
// 	}

// }

// Logger return log message
func Logger(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now(), r.Method, r.URL)
		handler.ServeHTTP(w, r) // dispatch the request
	})
}
