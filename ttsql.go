package main

import (
	//"bufio"
	"database/sql"
	"flag"
	"fmt"
	//"io"
	//"os"
	"time"

	_ "gomssqldb"
)

func main() {
	var (
		userid   = flag.String("U", "sa", "login_id")
		password = flag.String("P", "gecx1057@123", "password")
		server   = flag.String("S", "192.168.11.112", "server_name[\\instance_name]")
		database = flag.String("d", "LabConsole", "db_name")
	)
	flag.Parse()

	dsn := "server=" + *server + ";user id=" + *userid + ";password=" + *password + ";database=" + *database
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		return
	}
	
	cmd := "select * from labcode" //zhiling
	
	err = exec(db, cmd)
		if err != nil {
			fmt.Println(err)
		}
}

func exec(db *sql.DB, cmd string) error {
	rows, err := db.Query(cmd)
	if err != nil {
		return err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	if cols == nil {
		return nil
	}
	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
		if i != 0 {
			fmt.Print("\t")
		}
		fmt.Print(cols[i])
	}
	fmt.Println()
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for i := 0; i < len(vals); i++ {
			if i != 0 {
				fmt.Print("\t")
			}
			printValue(vals[i].(*interface{}))
		}
		fmt.Println()

	}
	if rows.Err() != nil {
		return rows.Err()
	}
	return nil
}

func printValue(pval *interface{}) {
	switch v := (*pval).(type) {
	case nil:
		fmt.Print("NULL")
	case bool:
		if v {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
	case []byte:
		fmt.Print(string(v))
	case time.Time:
		fmt.Print(v.Format("2006-01-02 15:04:05.999"))
	default:
		fmt.Print(v)
	}
}
	
	
	