package db

import (
	"database/sql"
	"fmt"

	// M
	_ "github.com/go-sql-driver/mysql"
)

// UserAuthService is the
type URLService interface {
	Init() error

	CreateNoPassword(string, string) error
	CreatePassword(string, string, string) error
	GetLong(string) (bool, string, error)
	PresentShort(string) (bool, error)
}

type urls struct {
	createNoPassword *sql.Stmt
	createPassword   *sql.Stmt
	getLong          *sql.Stmt
	checkShort       *sql.Stmt
	db               *sql.DB
}

// Init creates all the prepared statements
func (u *urls) Init() error {
	var prepareStatementError error

	u.checkShort, prepareStatementError = u.db.Prepare("select `id` from `short_urls` where `short` = ?")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	u.createNoPassword, prepareStatementError = u.db.Prepare("insert into `short_urls` (`short`,`long`) values (?,?)")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	u.createPassword, prepareStatementError = u.db.Prepare("insert into `short_urls` (`short`,`long`, `password`) values (?,?,?)")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	u.getLong, prepareStatementError = u.db.Prepare("select `long`,`password` from `short_urls` where `short` = ?")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	return nil
}

func (u *urls) CreateNoPassword(short, long string) error {
	_, err := u.createNoPassword.Exec(short, long)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (u *urls) CreatePassword(short, long, password string) error {
	_, err := u.createPassword.Exec(short, long, password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (u *urls) GetLong(short string) (bool, string, error) {
	var longURL string
	err := u.getLong.QueryRow(short).Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil
		}
		return false, "", err
	}
	return true, longURL, nil
}

func (u *urls) PresentShort(short string) (bool, error) {
	var id int
	err := u.checkShort.QueryRow(short).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
