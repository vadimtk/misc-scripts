package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "log"
import "time"
import "math/rand"
import "runtime"
import "strconv"
import "os"
import "flag"


func BumpMySQL(db *sql.DB, buf chan int, tab int) {

var k int
	// Open doesn't open a connection. Validate DSN data:
	err := db.Ping()
	if err != nil { panic(err.Error()) }

	for {
    select {
        case i := <-buf:
		rows, err := db.Query("select k from sbtest"+strconv.Itoa(tab+1)+" where id = "+strconv.Itoa(i))
			if err != nil { log.Fatal(err) }
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&k)
	     		if err != nil { log.Fatal(err) }
		}
		err = rows.Err()
			if err != nil { log.Fatal(err) }

	}
    }

}

var threads = flag.Int("threads", 8, "help message for flagname")

func main() {

flag.Parse()
log.Println("User threads:",*threads)

runtime.GOMAXPROCS(runtime.NumCPU())



var buffer = make(chan int, 1000) // 1000 just looks big enough for this case

maxtime := 300 * time.Second // time to run benchmark

db, err := sql.Open("mysql", "root:@/sbtest")

if err != nil {
    panic(err.Error())
}
defer db.Close()

db.SetMaxIdleConns(10000)

for j:=0; j<*threads; j++ {
	go BumpMySQL(db, buffer, j)
}

req:=0

t0 := time.Now()


go func() { // Daniel told me to write this handler this way. 
timer := time.NewTimer(maxtime)
for {
	select {
	case <-time.After(time.Second * 10):
		ts := time.Since(t0)
        	log.Println("time: ",ts," requests: ", req," throughput:",float64(req)/ts.Seconds())
	case <-timer.C:
        	log.Println("Finish!")
		ts := time.Since(t0)
        	log.Println("Final result: ",ts," requests: ", req," throughput:",float64(req)/ts.Seconds())
		os.Exit(1) // this is not a quite nice way to exit, but I do not care.
	}	
}
}()

for {

	buffer <- rand.Intn(1000000)+1 
	req = req + 1
}

time.Sleep(10000 * time.Millisecond)

}
