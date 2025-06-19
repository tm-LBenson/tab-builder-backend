package db

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Song struct {
	ID        string         `json:"id"`
	OwnerUID  string         `json:"ownerUid,omitempty"`
	Title     string         `json:"title"`
	IsPublic  bool           `json:"isPublic"`
	Payload   map[string]any `json:"payload"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type SongIn struct {
	Title    string         `json:"title"`
	IsPublic bool           `json:"isPublic"`
	Payload  map[string]any `json:"payload"`
}

func scanSong(row pgx.Row) (Song, error) {
	var s Song
	err := row.Scan(&s.ID, &s.OwnerUID, &s.Title, &s.IsPublic, &s.Payload, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

type Store struct{ Pool *pgxpool.Pool }

func NewStore(ctx context.Context) (*Store, error) {
	url := os.Getenv("DATABASE_URL")
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 8
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	return &Store{Pool: pool}, err
}

func (s *Store) Close() { s.Pool.Close() }

func (s *Store) ListSongs(ctx context.Context, uid string) ([]map[string]any, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT id,title,is_public,payload,updated_at
		   FROM songs WHERE owner_uid=$1 OR is_public`,
		uid)
	if err != nil {
		return nil, err
	}
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
func (s *Store) CreateSong(ctx context.Context, uid string, in SongIn) (Song, error) {
	row := s.Pool.QueryRow(ctx,
		`INSERT INTO songs (owner_uid,title,is_public,payload)
             VALUES ($1,$2,$3,$4)
         RETURNING id,title,is_public,payload,created_at,updated_at`,
		uid, in.Title, in.IsPublic, in.Payload)
	return scanSong(row)
}

func (s *Store) GetSong(ctx context.Context, id string) (Song, error) {
	row := s.Pool.QueryRow(ctx,
		`SELECT id,owner_uid,title,is_public,payload,created_at,updated_at
           FROM songs WHERE id=$1`,
		id)
	return scanSong(row)
}

func (s *Store) UpdateSong(ctx context.Context, id string, in SongIn) (Song, error) {
	row := s.Pool.QueryRow(ctx,
		`UPDATE songs
            SET title=$2, is_public=$3, payload=$4, updated_at=NOW()
          WHERE id=$1
          RETURNING id,title,is_public,payload,created_at,updated_at`,
		id, in.Title, in.IsPublic, in.Payload)
	return scanSong(row)
}

func (s *Store) DeleteSong(ctx context.Context, id string) error {
	_, err := s.Pool.Exec(ctx, `DELETE FROM songs WHERE id=$1`, id)
	return err
}

func (s *Store) CloneSong(ctx context.Context, id, newUID string) (Song, error) {
	row := s.Pool.QueryRow(ctx,
		`INSERT INTO songs (owner_uid,title,is_public,payload)
             SELECT $2,title,is_public,payload FROM songs WHERE id=$1
         RETURNING id,title,is_public,payload,created_at,updated_at`,
		id, newUID)
	return scanSong(row)
}
