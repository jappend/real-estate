package tests

import (
	"database/sql"
	"fmt"
	"jappend/real_estate/internal/database"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"testing"
  _ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type TestRunner struct {
  DefDB *sql.DB
  DB *database.Queries
  TestDB string
}

type testCases struct {
	name string
	test func(*testing.T)
}

var testRunner TestRunner

func TestMain(m *testing.M) {
  tr, err := setup()
  if err != nil {
    log.Fatal("Error setting up tests: ", err)
  }
  testRunner = *tr

  // Running tests
  exitCode := m.Run()

  defer func() {
    if err := teardown(tr); err != nil {
      log.Fatal("Error tearing down tests: ", err)
    }

    // Finishing tests
    os.Exit(exitCode)
  }()
}

func setup() (*TestRunner, error) {
  // Initializing env variables
  if _, err := os.Stat("../../.env"); err == nil {
    if err := godotenv.Load("../../.env"); err != nil {
      log.Fatal("Error initializing environment variables: ", err)
    }
  } else {
    		log.Println(".env file not found, skipping loading of environment variables from file") // We do this for the CI testing
  }

  // Setting host to localhost, because the db would be local in the test goroutine
  if err := os.Setenv("PGHOST", "localhost"); err != nil {
    log.Fatal("Error setting PGHOST environment variable: ", err)
  }

  // Generating a random string to be the test database name.
  const letterBytes string = "abcdefghijklmnopqrstuvwxyz"
	b := []byte{'t', 'e', 's', 't', '_'}
	for len(b) < 15 {
    b = append(b, []byte(letterBytes)[rand.IntN(len(letterBytes))]) 
	}
	testDB := string(b)

	// Connect to the default postgres db
	defaultDB, err := connectToDefaultDB()
	if err != nil {
		return nil, fmt.Errorf("error connecting to default database: %v", err)
	}

	// Create test database
	if _, err := defaultDB.Exec("CREATE DATABASE " + testDB + ";"); err != nil {
		return nil, fmt.Errorf("error creating test database: %v", err)
	}

  // Run migrations on test database
	err = runMigrations(testDB)
	if err != nil {
		return nil, fmt.Errorf("error running migrations: %v", err)
	}

	// Create Test Runner
  testDBConn := dbInitializer(testDB)
  queries := database.New(testDBConn)
  defer testDBConn.Close()

	tr := TestRunner{
		DefDB:        defaultDB,
		DB:           queries,
		TestDB:       testDB,
	}

	return &tr, nil
}

func teardown(tr *TestRunner) error {
	// Close connection to test database
	log.Println("Running teardown after tests...")

	// Drop test database
	if _, err := tr.DefDB.Exec("DROP DATABASE " + tr.TestDB + ";"); err != nil {
		return fmt.Errorf("error dropping test database: %v", err)
	}

	// Close the default connection
	if err := tr.DefDB.Close(); err != nil {
		return fmt.Errorf("error closing default DB connection: %v", err)
	}

	return nil
}

// Utilities
func dbInitializer(dbName string) *sql.DB {
	// Opening a connection to the database
	var (
		user     string = os.Getenv("PGUSER")
		password string = os.Getenv("PGPASSWORD")
		host     string = os.Getenv("PGHOST")
		port     string = os.Getenv("PGPORT")
	)
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	log.Printf("Opening connection with database %s on port %s...\n", dbName, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
	}
	log.Println("Database connection opened successefuly!")

	return db
}

func runMigrations(dbName string) error {
	// Construct connection string dynamically
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		dbName,
	)

	// Run goose migrations
	cmd := exec.Command("goose", "-dir", "../database/migrations", "postgres", connStr, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running migrations: %v", err)
	}

	return nil
}

func connectToDefaultDB() (*sql.DB, error) {
  // Connecting to the default "postgres" database so we can create the test db
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
	)

	log.Println("Connecting to default 'postgres' database...")
	defaultDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to default DB: %v", err)
	}

	if err := defaultDB.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging default DB: %v", err)
	}

	log.Println("Connected to default 'postgres' database successfully!")
	return defaultDB, nil
}

func runTestCasesInParallel(t *testing.T, tests []testCases) {
	// Create a channel and loop all test cases to run then in parallel
	done := make(chan struct{})
	for _, tc := range tests {
		go func() {
			t.Run(tc.name, tc.test)
			done <- struct{}{}
		}()
	}

	// Wait for all tests to complete
	for range tests {
		<-done
	}
}

