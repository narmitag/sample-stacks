package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Custom metrics
var (
	dbQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "test_app_mysql_query_duration_seconds",
			Help:    "Histogram of the query durations for MySQL.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query"},
	)

	dbQueryErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_app_mysql_query_errors_total",
			Help: "Total number of errors while querying MySQL.",
		},
		[]string{"query"},
	)
)

func init() {
	// Register metrics
	prometheus.MustRegister(dbQueryDuration)
	prometheus.MustRegister(dbQueryErrors)
}

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatalf("MYSQL_DSN environment variable is not set")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Connected to MySQL database")

	// Start a goroutine to periodically run a sample query
	go func() {
		for {
			for i := 1; i <= 10; i++ {
				query := fmt.Sprintf("select count(*)  from sbtest%d group by k;", i)
				runQuery(db, query)
				time.Sleep(2 * time.Second)
			}
		}
	}()

	// Set up Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Prometheus metrics exposed at /metrics")

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func runQuery(db *sql.DB, query string) {
	start := time.Now()
	var count int
	err := db.QueryRow(query).Scan(&count)

	duration := time.Since(start).Seconds()
	dbQueryDuration.WithLabelValues(query).Observe(duration)

	if err != nil {
		log.Printf("Error executing query '%s': %v", query, err)
		dbQueryErrors.WithLabelValues(query).Inc()
		return
	}

	log.Printf("Query '%s' returned count: %d", query, count)
}
