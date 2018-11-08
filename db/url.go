package db

import (
	"database/sql"
	"encoding/base64"
	"fmt"

	// M
	_ "github.com/go-sql-driver/mysql"
)

// URLService is the
type URLService interface {
	//Init() error

	CreateNoPassword(string, string) error
	CreatePassword(string, string, string, string) error
	GetLong(string) (bool, string, string, string, error)
	PresentShort(string) (bool, error)
}

type urlsDb struct {
	createNoPassword *sql.Stmt
	createPassword   *sql.Stmt
	getLong          *sql.Stmt
	checkShort       *sql.Stmt
	db               *sql.DB
}

// Init creates all the prepared statements
func (u *urlsDb) Init() error {
	var prepareStatementError error

	u.checkShort, prepareStatementError = u.db.Prepare("select `id` from `short_urls` where `short` = ?")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	u.createNoPassword, prepareStatementError = u.db.Prepare("insert into `short_urls` (`short`,`long`) values (?,?)")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	u.createPassword, prepareStatementError = u.db.Prepare("insert into `short_urls` (`short`,`long`, `password`, `salt`) values (?,?,?,?)")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	u.getLong, prepareStatementError = u.db.Prepare("select `long`,`password`, `salt` from `short_urls` where `short` = ?")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	return nil
}

func (u *urlsDb) CreateNoPassword(short, long string) error {
	_, err := u.createNoPassword.Exec(short, long)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (u *urlsDb) CreatePassword(short, long, password, salt string) error {
	salt = base64.URLEncoding.EncodeToString([]byte(salt))
	long = base64.URLEncoding.EncodeToString([]byte(long))
	fmt.Println(salt)
	_, err := u.createPassword.Exec(short, long, password, salt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (u *urlsDb) GetLong(short string) (bool, string, string, string, error) {
	var longURL string
	var password string
	var salt string
	err := u.getLong.QueryRow(short).Scan(&longURL, &password, &salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", "", "", nil
		}
		return false, "", "", "", err
	}
	return true, longURL, password, salt, nil
}

func (u *urlsDb) PresentShort(short string) (bool, error) {
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
