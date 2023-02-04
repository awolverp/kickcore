package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/awolverp/kickcore/cache"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteCacheDriver struct {
	conn *sql.DB
}

func (db *SQLiteCacheDriver) execTx(ctx context.Context, isolationLevel sql.IsolationLevel, callback func(*sql.Tx) error) error {
	tx, err := db.conn.BeginTx(ctx, &sql.TxOptions{Isolation: isolationLevel})
	if err != nil {
		return err
	}

	err = callback(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *SQLiteCacheDriver) Init() error {
	err := db.execTx(context.Background(), sql.LevelSerializable, func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`CREATE TABLE IF NOT EXISTS cache(key TEXT PRIMARY KEY, value BLOB, date BIGINT);`,
		)
		return err
	})
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(`VACUUM`)
	return err
}

func (db *SQLiteCacheDriver) PingContext(ctx context.Context) error { return db.conn.PingContext(ctx) }

func (db *SQLiteCacheDriver) Insert(key string, value []byte, date int64) (bool, error) {
	var result bool = false

	err := db.execTx(context.Background(), sql.LevelReadCommitted, func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO cache(key,value,date) VALUES(?,?,?);`, key, value, date)
		if err != nil {
			return nil
		}

		result = true
		return nil
	})

	return result, err
}

func (db *SQLiteCacheDriver) Select(key string, value *[]byte) error {
	err := db.conn.QueryRow(`SELECT value FROM cache WHERE key=? LIMIT 1;`, key).Scan(value)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	return err
}

func (db *SQLiteCacheDriver) SelectExpiredValues(expireAfter int64) ([]string, error) {
	rows, err := db.conn.Query(
		`SELECT key FROM cache WHERE (?-date > ?);`, time.Now().Unix(), expireAfter-1,
	)
	if err != nil {
		return nil, err
	}

	var keys []string

	for rows.Next() {
		var result string
		rows.Scan(&result)
		keys = append(keys, result)
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	return keys, nil
}

func (db *SQLiteCacheDriver) Delete(key string) (bool, error) {
	var result bool = false

	err := db.execTx(context.Background(), sql.LevelReadCommitted, func(tx *sql.Tx) error {
		sqlresult, err := tx.Exec(`DELETE FROM cache WHERE key=?;`, key)
		if err != nil {
			return err
		}

		n, _ := sqlresult.RowsAffected()
		result = (n != 0)
		return nil
	})

	return result, err
}

func stringListToQuery(keys []string) string {
	length := len(keys)

	if length == 0 {
		return "()"
	} else if length == 1 {
		return `('` + keys[0] + `')`
	}

	var s string = "("

	for i, v := range keys {
		s += `'` + v + `'`
		if i < length-1 {
			s += ","
		}
	}

	return s + ")"
}

func (db *SQLiteCacheDriver) DeleteMany(keys []string) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	var result int64

	err := db.execTx(context.Background(), sql.LevelReadCommitted, func(tx *sql.Tx) error {
		sqlresult, err := tx.Exec(
			`DELETE FROM cache WHERE key IN ` + stringListToQuery(keys),
		)
		if err != nil {
			return err
		}

		result, _ = sqlresult.RowsAffected()
		return nil
	})

	return result, err
}

func (db *SQLiteCacheDriver) Len() (int64, error) {
	var result int64

	err := db.conn.QueryRow(`SELECT COUNT(key) FROM cache;`).Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (db *SQLiteCacheDriver) Close() error { return db.conn.Close() }

func Connect(dsn string, timeout time.Duration) (cache.CacheDriver, error) {
	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if timeout == 0 {
		timeout = time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		conn.Close()
		return nil, err
	}

	db := new(SQLiteCacheDriver)
	db.conn = conn
	return db, nil
}
