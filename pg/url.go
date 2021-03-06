package pg

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zirius/url-shortener/models"
)

func GetURL(db *sqlx.DB, longUrl, slug string) (*models.URL, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sb := psql.Select("*").
		From("urls")
	if longUrl != "" {
		sb = sb.Where(squirrel.Eq{"url": longUrl})
	}
	if slug != "" {
		sb = sb.Where(squirrel.Eq{"slug": slug})
	}

	sqlStr, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.Queryx(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var url models.URL
		if err := rows.StructScan(&url); err != nil {
			return nil, err
		}
		return &url, nil
	}
	return nil, sql.ErrNoRows
}

func GetURLs(db *sqlx.DB, clauses map[string]interface{}) ([]models.URL, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sb := psql.Select("*").
		From("urls").OrderBy("created desc")

	if slug, ok := clauses["slug"].(string); ok {
		sb = sb.Where(squirrel.Eq{"slug": slug})
	}

	if url, ok := clauses["url"].(string); ok {
		sb = sb.Where(squirrel.Eq{"url": url})
	}

	if status, ok := clauses["status"].(string); ok {
		sb = sb.Where(squirrel.Eq{"status": status})
	}

	if status, ok := clauses["status"].(string); ok {
		sb = sb.Where(squirrel.Eq{"status": status})
	}

	if limit, ok := clauses["_limit"].(uint64); ok {
		sb = sb.Limit(limit)
	}

	if offset, ok := clauses["_offset"].(uint64); ok {
		sb = sb.Offset(offset)
	}

	sqlStr, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	var urls []models.URL

	rows, err := db.Queryx(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url models.URL
		if err := rows.StructScan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func CreateURL(db *sqlx.DB, url *models.URL) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sb := psql.Insert("urls").Columns("url, slug, ip, counter, created, updated, password, expired, mindful").
		Values(url.Url, url.Slug, url.IP, url.Counter, url.Created, url.Updated, url.Password, url.Expired, url.Mindful)
	sqlStr, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	if _, err = db.Exec(sqlStr, args...); err != nil {
		return err
	}
	return nil
}

func UpdateURL(db *sqlx.DB, url *models.URL) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	clauses := make(map[string]interface{})
	clauses["counter"] = url.Counter
	clauses["status"] = url.Status
	clauses["updated"] = time.Now()
	sb := psql.Update("urls").SetMap(clauses).Where(squirrel.Eq{"slug": url.Slug})
	sqlStr, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	if _, err = db.Exec(sqlStr, args...); err != nil {
		return err
	}
	return nil
}
