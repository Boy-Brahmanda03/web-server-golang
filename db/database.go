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

func AddInbox(id_pesan int, id_user int, uname string, text string)(int64, error){
	result, err := db.Exec("INSERT INTO inbox (id_pesan, id_user, username, message, tanggal) VALUES (?, ?, ?, ?, ?)", id_pesan, id_user, uname, text, time.Now())
    if err != nil {
        return 0, fmt.Errorf("addInbox: %v", err)
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("addInbox: %v", err)
    }
    return id, nil
}

func AddOutbox(id_pesan int,id_user int, uname string, text string)(int64, error){
	result, err := db.Exec("INSERT INTO outbox (id_pesan, id_user, username, message, tanggal) VALUES (?, ?, ?, ?, ?)", id_pesan, id_user, uname, text, time.Now())
    if err != nil {
        return 0, fmt.Errorf("addOutbox: %v", err)
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("addOutbox: %v", err)
    }
    return id, nil
}


type Menu struct {
	No int64
	Label string
	Deskripsi string
}

func ShowMenu()([]Menu, error){
	rows, err := db.Query("Select no,label,deskripsi from tb_menu")
	if err != nil {
        return nil, err
    }
    defer rows.Close()

	var menu []Menu

	for rows.Next(){
		var each = Menu{}
		var err = rows.Scan(&each.No, &each.Label, &each.Deskripsi)

		if err != nil {
			fmt.Println(err.Error())
            return nil, err
        }

		menu = append(menu, each)
	}

	if err = rows.Err(); err != nil {
        fmt.Println(err.Error())
        return nil, err
    }

	return menu, nil
}


func GetStateMessage(idMsg int) int{
	var state int
	rows := db.QueryRow("select state from inbox where id_pesan = ?", idMsg)
	rows.Scan(&state)
	return state
}

func UpdateState(idMsg int, state int) error {
	nState :=  state + 1
	_, err := db.Exec("Update inbox set state = ? where id_pesan = ?", nState, idMsg)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

type Mahasiswa struct {
	ID int64
	NIM string
	Nama string
}

func CariMahasiswa(nim string)([]Mahasiswa, error) {
	rows, err := db.Query("Select * from tb_mhs where nim = ?", nim)
	if err != nil {
        return nil, err
    }
    defer rows.Close()

	var mahasiswa []Mahasiswa

	for rows.Next() {
		var mhs Mahasiswa
		err := rows.Scan(&mhs.ID, &mhs.NIM, &mhs.Nama); 
		if err != nil {
            fmt.Println(err.Error())
            return nil, err
        }

        mahasiswa = append(mahasiswa, mhs)
	}
	return mahasiswa, err
} 