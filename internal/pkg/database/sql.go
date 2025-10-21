package database

import (
	"context"
	"database/sql"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	logging "fiber-boilerplate/internal/pkg/logging"
	"fiber-boilerplate/internal/pkg/setting"
	"fiber-boilerplate/internal/pkg/util"

	"gopkg.in/guregu/null.v4"

	"github.com/jmoiron/sqlx"
	// postgres driver
	_ "github.com/lib/pq"
)

const (
	socketDir = "/cloudsql"
)

// SQL :
type SQL struct {
	db *sqlx.DB
}

func (x *SQL) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	query, args = confirmQuery(query, args...)
	if setting.Runtime.Env == "local" {
		logging.TraceSQL("%s / %v", query, args)
	}
	return x.db.ExecContext(ctx, query, args...)
}

func (x *SQL) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if setting.Runtime.Env == "local" {
		logging.TraceSQL("%s", query)
	}
	return x.db.PrepareContext(ctx, query)
}

func (x *SQL) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	query, args = confirmQuery(query, args...)
	if setting.Runtime.Env == "local" {
		logging.TraceSQL("%s / %v", query, args)
	}
	return x.db.QueryContext(ctx, query, args...)
}

func (x *SQL) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	query, args = confirmQuery(query, args...)
	if setting.Runtime.Env == "local" {
		logging.TraceSQL("%s / %v", query, args)
	}
	return x.db.QueryRowContext(ctx, query, args...)
}

// BeginxContext : Begin a transaction with context propagation
func (x *SQL) BeginxContext(ctx context.Context) (*SQLTX, context.Context, error) {
	tx, err := x.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	if setting.Runtime.Env == "local" {
		logging.TraceSQL("TRANSACTION %p : BEGIN", tx)
	}

	return &SQLTX{
		Tx: tx,
		Commit: func() error {
			if setting.Runtime.Env == "local" {
				logging.TraceSQL("TRANSACTION %p : COMMIT", tx)
			}
			return tx.Commit()
		},

		Rollback: func() error {
			if setting.Runtime.Env == "local" {
				logging.TraceSQL("TRANSACTION %p : ROLLBACK", tx)
			}
			return tx.Rollback()
		},
	}, ctx, nil
}

// SQLTX :
type SQLTX struct {
	Tx       *sql.Tx
	Commit   func() error
	Rollback func() error
}

func (x *SQLTX) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	query, args = confirmQuery(query, args...)
	if setting.Runtime.Env == "local" {
		logging.TraceSQL("TRANSACTION %p : %s / %v", x.Tx, query, args)
	}
	return x.Tx.ExecContext(ctx, query, args...)
}

func (x *SQLTX) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if setting.Runtime.Env == "local" {
		logging.TraceSQL("TRANSACTION %p : %s", x.Tx, query)
	}
	return x.Tx.PrepareContext(ctx, query)
}

func (x *SQLTX) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	query, args = confirmQuery(query, args...)
	if setting.Runtime.Env == "local" {
		logging.TraceSQL("TRANSACTION %p : %s / %v", x.Tx, query, args)
	}
	return x.Tx.QueryContext(ctx, query, args...)
}

func (x *SQLTX) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	query, args = confirmQuery(query, args...)
	if setting.Runtime.Env == "local" {
		logging.TraceSQL("TRANSACTION %p : %s / %v", x.Tx, query, args)
	}
	return x.Tx.QueryRowContext(ctx, query, args...)
}

func confirmQuery(query string, args ...interface{}) (string, []interface{}) {
	// Read query line by line
	querySplit := strings.Split(query, "\n")
	for i, line := range querySplit {
		cut := strings.ToLower(strings.TrimLeft(line, " "))
		switch {
		case strings.HasPrefix(cut, "?"): // ? 1 = $4::text
			// 정규식으로 $ 뒤에 오는 숫자 추출
			reg := regexp.MustCompile(`\$(\d+)`)
			matches := reg.FindStringSubmatch(line)
			if len(matches) == 0 || len(matches[0]) == 0 {
				break
			}

			optionPos, err := strconv.Atoi(strings.TrimPrefix(matches[0], "$"))
			if err != nil {
				logging.Warn(err, "Failed to parse parameter position in query: %s", line)
				break
			}

			if len(args) > optionPos-1 {
				if v, ok := args[optionPos-1].(null.String); ok {
					if v.Valid {
						querySplit[i] = v.String
					} else {
						querySplit[i] = ""
					}
				}
				args = append(args[:optionPos-1], args[optionPos:]...)
			}
		}
	}

	// merge split query
	query = strings.Join(querySplit, "\n")
	return query, args
}

// Connect :
func (x *SQL) Connect(driverID DriverEnum) (err error) {
	config := driverConfigs[driverID]

	var dbHost string
	if len(config.Host) > 0 {
		dbHost = util.String.Concat(
			"host=", config.Host, " ",
			"port=", config.Port, " ",
			"sslmode=", "disable",
		)
	} else {
		dbHost = util.String.Concat(
			"host=", path.Join(socketDir, config.Conn),
		)
	}
	dbURI := util.String.Concat(
		dbHost, " ",
		"user=", config.User, " ",
		"password=", config.Password, " ",
		"database=", config.Name,
	)

	x.db, err = sqlx.Connect(driverID.String(), dbURI)
	if err != nil {
		logging.Warn(err, "")
		return
	}

	x.db.SetMaxIdleConns(config.MaxIdleConns)
	x.db.SetMaxOpenConns(config.MaxOpenConns)
	x.db.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)

	return
}
