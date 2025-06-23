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

func scanSong(row pgx.Row) (Song, error) {
	var s Song
	err := row.Scan(
		&s.ID, &s.OwnerUID, &s.Title, &s.IsPublic,
		&s.Payload, &s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

func (s *Store) EnsureUser(ctx context.Context, uid, display string) error {
	_, err := s.Pool.Exec(ctx, `
	    INSERT INTO users (uid, display_name)
	    VALUES ($1,$2)
	    ON CONFLICT (uid) DO NOTHING
	`, uid, display)
	return err
}

func (s *Store) ListUserSongs(ctx context.Context, uid string) ([]Song, error) {
	rows, err := s.Pool.Query(ctx, `
	    SELECT id,owner_uid,title,is_public,payload,created_at,updated_at
	      FROM songs
	     WHERE owner_uid=$1
	     ORDER BY updated_at DESC
	`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Song
	for rows.Next() {
		song, err := scanSong(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, song)
	}
	return out, rows.Err()
}

func (s *Store) ListPublicSongs(ctx context.Context, limit int) ([]Song, error) {
	rows, err := s.Pool.Query(ctx, `
	    SELECT id,owner_uid,title,is_public,payload,created_at,updated_at
	      FROM songs
	     WHERE is_public
	     ORDER BY updated_at DESC
	     LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Song
	for rows.Next() {
		song, err := scanSong(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, song)
	}
	return out, rows.Err()
}

func (s *Store) CreateSong(ctx context.Context, uid string, in SongIn) (Song, error) {
	row := s.Pool.QueryRow(ctx, `
	    INSERT INTO songs (owner_uid,title,is_public,payload)
	    VALUES ($1,$2,$3,$4)
	    RETURNING id,owner_uid,title,is_public,payload,created_at,updated_at
	`, uid, in.Title, in.IsPublic, in.Payload)
	return scanSong(row)
}

func (s *Store) GetSong(ctx context.Context, id string) (Song, error) {
	row := s.Pool.QueryRow(ctx, `
	    SELECT id,owner_uid,title,is_public,payload,created_at,updated_at
	      FROM songs
	     WHERE id=$1
	`, id)
	return scanSong(row)
}

func (s *Store) UpdateSong(ctx context.Context, id string, in SongIn) (Song, error) {
	row := s.Pool.QueryRow(ctx, `
	    UPDATE songs
	       SET title=$2, is_public=$3, payload=$4, updated_at=NOW()
	     WHERE id=$1
	    RETURNING id,owner_uid,title,is_public,payload,created_at,updated_at
	`, id, in.Title, in.IsPublic, in.Payload)
	return scanSong(row)
}

func (s *Store) DeleteSong(ctx context.Context, id string) error {
	_, err := s.Pool.Exec(ctx, `DELETE FROM songs WHERE id=$1`, id)
	return err
}

func (s *Store) CloneSong(ctx context.Context, id, newUID string) (Song, error) {
	row := s.Pool.QueryRow(ctx, `
	    INSERT INTO songs (owner_uid,title,is_public,payload)
	    SELECT $2,title,is_public,payload FROM songs WHERE id=$1
	    RETURNING id,owner_uid,title,is_public,payload,created_at,updated_at
	`, id, newUID)
	return scanSong(row)
}
