package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
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

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// SSL/TLS設定の準備
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile("global-bundle.pem")
	if err != nil {
		log.Fatalf("Failed to read CA cert: %v", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatalf("Failed to append PEM.")
	}

	tlsConfig := &tls.Config{
		RootCAs: rootCertPool,
	}

	// TLS設定を登録
	tlsConfigName := "custom-tls"
	if err := mysql.RegisterTLSConfig(tlsConfigName, tlsConfig); err != nil {
		log.Fatal(err)
	}

	// データベースに接続
	db, err := sql.Open("mysql", config.DSN+"&tls="+tlsConfigName)
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
