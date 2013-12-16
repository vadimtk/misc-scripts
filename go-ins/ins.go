package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "strconv"
import "os"

func main() {
db, err := sql.Open("mysql", "root:@/t1")

if err != nil {
    panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
}
defer db.Close()

// Open doesn't open a connection. Validate DSN data:
err = db.Ping()
if err != nil {
    panic(err.Error()) // proper error handling instead of panic in your app
}

_,err = db.Exec("TRUNCATE TABLE tab1")
if err != nil {
    panic(err.Error()) // proper error handling instead of panic in your app
}

for i := 0; i < 1000000; i++ {

str := strconv.Itoa(i)
_,err = db.Exec("INSERT INTO tab1 (id,val) VALUES (" + str + ",'randomvalue')")
if err != nil {
    panic(err.Error()) // proper error handling instead of panic in your app
}

_,err = db.Exec("COMMIT")
if err != nil {
    panic(err.Error()) // proper error handling instead of panic in your app
}

os.Stdout.Write([]byte("loop: "+str+"\n"))
os.Stdout.Sync()

}


}
