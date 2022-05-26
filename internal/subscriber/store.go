package subscriber

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
)

type db struct {
	sql *sqlx.DB
}

func NewStore(d *sqlx.DB) db {
	return db{
		sql: d,
	}
}

func (d db) InsertEvents(ctx context.Context, events []json.RawMessage) error {
	toInsert := make([]map[string]interface{}, len(events))
	for i, v := range events {
		toInsert[i] = map[string]interface{}{"event": v}
	}

	_, err := d.sql.NamedExecContext(ctx, `INSERT INTO event (event) VALUES (:event)`, toInsert)
	return err
}
