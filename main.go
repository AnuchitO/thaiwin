package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"thaichana/logger"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	viper.SetDefault("port", "8000")
	viper.SetDefault("db.conn", "thaichana.db")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	l, _ := zap.NewDevelopment()
	defer l.Sync()
	hostname, _ := os.Hostname()
	l = l.With(zap.String("hostname", hostname))
	zap.ReplaceGlobals(l)

	db, err := sql.Open("sqlite3", viper.GetString("db.conn"))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	r := mux.NewRouter()
	r.Use(logger.LoggerMiddleware(l))
	r.Use(SealMiddleware())

	r.HandleFunc("/recently", Recently).Methods(http.MethodPost)
	r.HandleFunc("/checkin", CheckIn(InFunc(NewInsertCheckIn(db)))).Methods(http.MethodPost)
	r.HandleFunc("/checkout", CheckOut).Methods(http.MethodPost)

	port := viper.GetString("port")
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	zap.L().Info("starting...", zap.String("port", port))
	log.Fatal(srv.ListenAndServe())
}

type Check struct {
	ID      int64 `json:"id"`
	PlaceID int64 `json:"place_id"`
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

// Recently returns currently visited
func Recently(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

type InFunc func(id, placeID int64) error

func (fn InFunc) In(id, placeID int64) error {
	return fn(id, placeID)
}

type Iner interface {
	In(id, placeID int64) error
}

func NewInsertCheckIn(db *sql.DB) func(id, placeID int64) error {
	return func(id, placeID int64) error {
		_, err := db.Exec("INSERT INTO visits VALUES(?, ?);", id, placeID)
		return err
	}
}

// CheckIn check-in to place, returns density (ok, too much)
func CheckIn(check Iner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.L(r.Context()).Info("check-in")
		chk := Check{}
		if err := json.NewDecoder(r.Body).Decode(&chk); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		defer r.Body.Close()

		err := check.In(chk.ID, chk.PlaceID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "density": "ok" }`))
	}
}

// CheckOut check-out from place
func CheckOut(w http.ResponseWriter, r *http.Request) {

}

func SealMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			data, err := base64.StdEncoding.DecodeString(string(b))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			buff := bytes.NewBuffer(data)

			r.Body = ioutil.NopCloser(buff)

			next.ServeHTTP(&EncodeWriter{w}, r)
		})
	}
}

type EncodeWriter struct {
	http.ResponseWriter
}

func (w *EncodeWriter) Write(b []byte) (int, error) {
	body := base64.StdEncoding.EncodeToString(b)
	return w.ResponseWriter.Write([]byte(body))
}

/*
type Interface interface {
	Insert(collection string, doc interface{}) error
}

type Connection struct {
		*mgo.Session
		db string
}

func (conn *Connection) Insert(collection string, doc interface{}) error {
	session := conn.copySession()
	defer session.Close()

	c := session.DB(conn.db).C(collection)
	return c.Insert(doc)
}

type Logger interface {
	StampRequestDateTime()
	StampResponseDateTime()
}

type LogDatabase struct {
	Interface
	Logger Logger
}

func (ld *LogDatabase) Insert(collection string, doc interface{}) error {
	ld.Logger.StampRequestDateTime()
	err := ld.Interface.Insert(collection, doc)
	ld.Logger.StampResponseDateTime()
	ld.Log("Insert", collection, 1, err)
	return err
}

// Factory
func Connect() (Interface, error) {
	if logDatabase != nil {
		return logDatabase, nil
	}
	conn, err := connect()
	if err != nil {
		return nil, err
	}
	if !conf.GetBool("mongodb.log_enable") {
		return conn, nil
	}
	dbLog := newLog(conf.Get("mongodb.log_file"), nil)
	logDatabase = &LogDatabase{Interface: conn, Logger: dbLog}
	return logDatabase, nil
}

func Close() {
	if connection != nil {
		connection.Close()
	}
}

func connect() (*Connection, error) {
	if connection != nil {
		return connection, nil
	}
	session, err := mgo.Dial(conf.Get("mongodb.host"))

	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Strong, true)
	connection = &Connection{Session: session, db: conf.Get("mongodb.dbname")}
	return connection, nil
}

*/
