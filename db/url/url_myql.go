package url

import (
	"database/sql"
	"fmt"

	// Only way to get it working
	_ "github.com/go-sql-driver/mysql"
)

// UrlsDb is used to i
type UrlsDb struct {
	create     *sql.Stmt
	getLong    *sql.Stmt
	checkShort *sql.Stmt
	DB         *sql.DB
}

// Init creates all the prepared statements
func (u *UrlsDb) Init() error {
	var prepareStatementError error

	u.checkShort, prepareStatementError = u.DB.Prepare("select `id` from `short_urls` where `short` = ?")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	u.create, prepareStatementError = u.DB.Prepare("insert into `short_urls` (`short`,`long`, `password`, `nonce`, `salt`) values (?,?,?,?,?)")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	u.getLong, prepareStatementError = u.DB.Prepare("select `nonce`, `salt`, `long`, `password` from `short_urls` where `short` = ?")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	return nil
}

// GetLong returns all the info related to the given short URL
func (u *UrlsDb) GetLong(short string) (map[string]string, error) {
	var mp map[string]string
	err := u.getLong.QueryRow(short).Scan(&mp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mp, nil
}

// Present checks the presence of given short URL
func (u *UrlsDb) Present(short string) (bool, error) {
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

// Create is used to create an entry in the datastore
func (u *UrlsDb) Create(short, nonce, salt, encrypted, passwordHash string) error {
	_, err := u.create.Exec(short, encrypted, passwordHash, nonce, salt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
