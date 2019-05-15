package pgx

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	log "gopkg.in/inconshreveable/log15.v2"
	"math/rand"
	"os"
	"sync"
	"time"
)

var pool *pgx.ConnPool

// afterConnect creates the prepared statements that this application uses
func afterConnect(conn *pgx.Conn) (err error) {
	_, err = conn.Prepare("read", `
    SELECT key FROM config_parameters where user_id=$1
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("read2", `
    SELECT value FROM config_parameters where user_id=$1 AND key=$2
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("write", `
INSERT INTO config_parameters(user_id, key, value) VALUES ($1, $2, $3)
  `)
	return
}

func RunBatch(){
	//logger := log15adapter.NewLogger(log.New("module", "pgx"))

	var err error
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "127.0.0.1",
			User:     "nasim",
			Password: "nasim",
			Database: "nasim",
		},
		MaxConnections: 5,
		AfterConnect:   afterConnect,
	}




	pool, err = pgx.NewConnPool(connPoolConfig)
	if err != nil {
		log.Crit("Unable to create connection pool", "error", err)
		os.Exit(1)
	}

	runBatchTest()
}


func runBatchTest(){
	start := time.Now()
	var waitGroup sync.WaitGroup
	for i:=0 ; i< 20; i++ {
		waitGroup.Add(1)
		go batchTest(&waitGroup, 20000)
	}

	waitGroup.Wait()
	fmt.Println("Elapsed : ", time.Since(start))
}

func runTest(){
	var waitGroup sync.WaitGroup

	rand.Seed(time.Now().UnixNano())
	for i:=0 ; i< 150000; i++ {
		waitGroup.Add(1)
		go batchQuery(1, &waitGroup)
	}

	start := time.Now()
	waitGroup.Wait()
	fmt.Println("Elapsed : ", time.Since(start))
}

func batchTest(w *sync.WaitGroup, n int) {
	defer w.Done()
	batch := pool.BeginBatch()


	for i:=0 ; i< n; i++ {
		batch.Queue("read2",
			[]interface{}{0, "s"},
			[]pgtype.OID{pgtype.Int4OID, pgtype.VarcharOID},
			nil,
		)
	}


	err := batch.Send(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	for j:=0 ; j < n; j++ {
		rows, err := batch.QueryResults()
		if err != nil {
			fmt.Println(err)
		}

		for rows.Next() {
			var key string
			if e := rows.Scan(&key); e!=nil {
				fmt.Println(e)
			}
			fmt.Println(key)

		}

	}

	_ = batch.Close()

}

func batchQuery(n int, waitGroup *sync.WaitGroup){
	defer waitGroup.Done()
	for i:=0; i< n; i++ {
		query()
	}
}

func query(){
	var k string

	err := pool.QueryRow("read", rand.Int31()).Scan(&k)
	switch err {
	case nil:
		fmt.Println("===", k)
	case pgx.ErrNoRows:

	default:
		fmt.Println(err)
	}

}
