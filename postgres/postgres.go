package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"math/rand"
	"strconv"
	"sync"
	"time"
)


const (
	host     = "localhost"
	port     = 5432
	user     = "nasim"
	password = "nasim"
	dbname   = "nasim"
)

var stmt *sql.Stmt
var insertStmt = "INSERT INTO config_parameters(user_id, key, value) VALUES ($1, $2, $3)"
func Run(){

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(1000)
	stmt ,err = db.Prepare("SELECT key FROM config_parameters where user_id=$1")
	if err != nil {
		panic(err)
	}

	runTest(db)


	defer db.Close()

}

func runTest(db *sql.DB){
	var waitGroup sync.WaitGroup

	rand.Seed(time.Now().UnixNano())
	for i:=0 ; i< 20; i++ {
		waitGroup.Add(1)
		go batchQuery(db, 10000, &waitGroup)
	}

	start := time.Now()
	waitGroup.Wait()
	fmt.Println("Elapsed : ", time.Since(start))
}

func batchQuery(db *sql.DB, n int, waitGroup *sync.WaitGroup){
	defer waitGroup.Done()
	for i:=0; i< n; i++ {
		query(db, nil)
	}
}

func batchInsert(db *sql.DB, n int, waitGroup *sync.WaitGroup){
	defer waitGroup.Done()
	for i:=0; i< n; i++ {
		insert(db, rand.Int31(), strconv.Itoa(i), strconv.Itoa(i))
	}
}

func insert(db *sql.DB, userId int32, key string, value string){

	_ , err := db.Exec(insertStmt, userId, key, value)

	if err!= nil {
		fmt.Println(err)
	}


}

func query(db *sql.DB, waitGroup *sync.WaitGroup){
	if waitGroup != nil {
		defer waitGroup.Done()
	}
	rows, err := stmt.Query(rand.Int31())
	if err != nil {
		fmt.Println(err)
	}

	var key string
	for rows.Next() {
		if e := rows.Scan(&key); e!=nil {
			fmt.Println(e)
		}

	}
	fmt.Println(key)

}

