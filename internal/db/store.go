package db

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct{ Pool *pgxpool.Pool }

func NewStore(ctx context.Context) (*Store, error) {
	url := os.Getenv("DATABASE_URL")
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil { return nil, err }
	cfg.MaxConns = 8
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	return &Store{Pool: pool}, err
}

func (s *Store) Close() { s.Pool.Close() }

// example query
func (s *Store) ListSongs(ctx context.Context, uid string) ([]map[string]any, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT id,title,is_public,payload,updated_at
		   FROM songs WHERE owner_uid=$1 OR is_public`,
		uid)
	if err != nil { return nil, err }
	defer rows.Close()

	var out []map[string]any
	for rows.Next() {
		var id, title string
		var pub bool
		var payload map[string]any
		var updated time.Time
		if err := rows.Scan(&id, &title, &pub, &payload, &updated); err != nil {
			return nil, err
		}
		out = append(out, map[string]any{
			"id": id, "title": title, "public": pub,
			"payload": payload, "updated": updated,
		})
	}
	return out, rows.Err()
}
