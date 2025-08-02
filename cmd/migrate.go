package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/MadManJJ/cms-api/config"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	cfg := config.New()

	port, err := strconv.Atoi(cfg.Database.Port)
	if err != nil {
		log.Fatalf("Invalid DB port: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, port, cfg.Database.Username, cfg.Database.Password, cfg.Database.DatabaseName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Enable uuid-ossp extension (optional but recommended if your SQL uses uuid_generate_v4())
	_, err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		log.Fatalf("Failed to create uuid-ossp extension: %v", err)
	}

	migrationFile := flag.String("file", "cmd/migration/default.up.sql", "Migration SQL file to run")
	flag.Parse()

	// Read migration SQL file (change path to your .up.sql file)
	sqlBytes, err := os.ReadFile(*migrationFile)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	sqlStatements := string(sqlBytes)

	isFunctionFile := strings.Contains(*migrationFile, "create_trigger_set_timestamp_function.up.sql")

	var commands []string
	if isFunctionFile {
		// ถ้าเป็นไฟล์ function, execute ทั้งหมดเป็นคำสั่งเดียว
		// ตรวจสอบให้แน่ใจว่าไฟล์ function ไม่มี semicolon ที่ท้ายสุดของ CREATE FUNCTION statement
		commands = append(commands, strings.TrimSpace(sqlStatements))
	} else {
		rawCommands := strings.Split(sqlStatements, ";")
		for _, cmd := range rawCommands {
			trimmedCmd := strings.TrimSpace(cmd)
			if trimmedCmd != "" {
				commands = append(commands, trimmedCmd)
			}
		}
	}

	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}

		_, err := db.Exec(cmd)
		if err != nil {
			log.Fatalf("Failed to execute migration command: %v\nSQL: %s", err, cmd)
		}
	}

	log.Println("Migration applied successfully.")
}
