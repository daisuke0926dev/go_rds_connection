package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DSN string `json:"dsn"`
}

func main() {
	// 設定ファイルを読み込む
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// データベースに接続
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// データベース接続を確認
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// SQLクエリ: テーブルからデータを選択
	rows, err := db.Query("SELECT User FROM user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 結果を読み出す
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("user: %s\n", user)
	}

	// エラーチェック
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
