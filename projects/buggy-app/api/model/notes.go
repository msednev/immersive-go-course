package model

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Note struct {
	Id       string    `json:"id"`
	Owner    string    `json:"owner"`
	Content  string    `json:"content"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Tags     []string  `json:"tags"`
}

type Notes []Note

type dbConn interface {
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func GetNotesForOwner(ctx context.Context, conn dbConn, owner string) (Notes, error) {
	if owner == "" {
		return nil, errors.New("model: owner not supplied")
	}

	queryRows, err := conn.Query(ctx, "SELECT id, owner, content, created, modified FROM public.note")
	if err != nil {
		return nil, fmt.Errorf("model: could not query notes: %w", err)
	}
	defer queryRows.Close()

	notes := []Note{}
	for queryRows.Next() {
		note := Note{}
		err = queryRows.Scan(&note.Id, &note.Owner, &note.Content, &note.Created, &note.Modified)
		if err != nil {
			return nil, fmt.Errorf("model: query scan failed: %w", err)
		}
		if note.Owner == owner {
			note.Tags = extractTags(note.Content)
			notes = append(notes, note)
		}
	}

	if queryRows.Err() != nil {
		return nil, fmt.Errorf("model: query read failed: %w", queryRows.Err())
	}

	return notes, nil
}

func GetNoteByIdForOwner(ctx context.Context, conn dbConn, id string, owner string) (Note, error) {
	var note Note
	if id == "" {
		return note, errors.New("model: id not supplied")
	}

	if owner == "" {
		return note, errors.New("model: owner not supplied")
	}

	row := conn.QueryRow(ctx, "SELECT id, owner, content, created, modified FROM public.note WHERE id = $1 AND owner = $2", id, owner)

	err := row.Scan(&note.Id, &note.Owner, &note.Content, &note.Created, &note.Modified)
	if err != nil {
		return note, fmt.Errorf("model: query scan failed: %w", err)
	}
	note.Tags = extractTags(note.Content)
	return note, nil
}

// Extract tags from the note. We're looking for #something. There could be
// multiple tags, so we FindAll.
func extractTags(input string) []string {
	re := regexp.MustCompile(`#([^# ]+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	tags := make([]string, 0, len(matches))
	for _, f := range matches {
		tags = append(tags, strings.TrimSpace(f[1]))
	}
	return tags
}
