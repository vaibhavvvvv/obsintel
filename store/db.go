 // Postgres connection pool

 package store

 import (
	 "context"
	 "fmt"
	 "log"
	 "os"
 
	 "github.com/jackc/pgx/v5/pgxpool"
 )
 
 var Pool *pgxpool.Pool
 
 func Init(ctx context.Context) {
	 connStr := fmt.Sprintf(
		 "postgres://%s:%s@%s:%s/%s",
		 os.Getenv("DB_USER"),
		 os.Getenv("DB_PASSWORD"),
		 os.Getenv("DB_HOST"),
		 os.Getenv("DB_PORT"),
		 os.Getenv("DB_NAME"),
	 )
 
	 pool, err := pgxpool.New(ctx, connStr)
	 if err != nil {
		 log.Fatalf("Failed to connect to database: %v", err)
	 }
 
	 // verify connection is actually alive
	 if err := pool.Ping(ctx); err != nil {
		 log.Fatalf("Database ping failed: %v", err)
	 }
 
	 Pool = pool
	 log.Println("Database connected successfully")
 }
 
 func Close() {
	 if Pool != nil {
		 Pool.Close()
	 }
 }