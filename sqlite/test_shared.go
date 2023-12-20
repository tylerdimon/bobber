package sqlite

import "log"
import "context"

func initDB() *DB {
	db := &DB{
		DSN: ":memory:",
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	err := db.Open()
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	return db
}
