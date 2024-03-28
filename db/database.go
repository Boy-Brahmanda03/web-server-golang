package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func DatabaseConnection(){
	connection := mysql.Config {
		User:   "root",
        Passwd: "",
        Net:    "tcp",
        Addr:   "localhost:3306",
        DBName: "latihan_golang",
		AllowNativePasswords: true,
	}
	 // Get a database handle.
	 var err error
	 db, err = sql.Open("mysql", connection.FormatDSN())
	 if err != nil {
		 log.Fatal(err)
	 }
 
	 pingErr := db.Ping()
	 if pingErr != nil {
		 log.Fatal(pingErr)
	 }
	 fmt.Println("Connected!")
}

func AddInbox(uname string, text string)(int64, error){
	result, err := db.Exec("INSERT INTO inbox (user, message, tanggal) VALUES (?, ?, ?)", uname, text, time.Now())
    if err != nil {
        return 0, fmt.Errorf("addInbox: %v", err)
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("addInbox: %v", err)
    }
	fmt.Println("id: ", id, uname, text)
    return id, nil
}

func AddOutbox(uname string, text string)(int64, error){
	result, err := db.Exec("INSERT INTO outbox (user, message, tanggal) VALUES (?, ?, ?)", uname, text, time.Now())
    if err != nil {
        return 0, fmt.Errorf("addOutbox: %v", err)
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("addOutbox: %v", err)
    }
	fmt.Println("id: ", id, uname, text)
    return id, nil
}


func ShowMenu()(int64, error){
	result, err := db.Exec("Select no,label,deskripsi from tb_menu")
	if err != nil {
        return 0, fmt.Errorf("showMenu: %v", err)
    }
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("showMenu: %v", err)
	}

	fmt.Println(id)
	return id, nil
}