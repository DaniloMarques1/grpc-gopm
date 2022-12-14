package repository

import (
	"database/sql"

	"github.com/danilomarques1/grpc-gopm/server/model"
)

type PasswordRepositoryImpl struct {
	db *sql.DB
}

func NewPasswordRepositoryImpl(db *sql.DB) *PasswordRepositoryImpl {
	return &PasswordRepositoryImpl{db: db}
}

func (p *PasswordRepositoryImpl) Save(password *model.Password) error {
	stmt, err := p.db.Prepare(`
		INSERT INTO passwords(id, key, password) VALUES($1, $2, $3);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(password.Id, password.Key, password.Pwd); err != nil {
		return err
	}
	return nil
}

func (p *PasswordRepositoryImpl) FindAllKeys() ([]string, error) {
	stmt, err := p.db.Prepare(`
		SELECT key FROM passwords;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	for rows.Next() {
		key := new(string)
		if err := rows.Scan(key); err != nil {
			return nil, err
		}
		keys = append(keys, *key)
	}

	return keys, nil
}

func (p *PasswordRepositoryImpl) FindPassword(key string) (*model.Password, error) {
	stmt, err := p.db.Prepare(`
		SELECT id, key, password FROM passwords WHERE key = $1;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	password := &model.Password{}
	if err := stmt.QueryRow(key).Scan(&password.Id, &password.Key, &password.Pwd); err != nil {
		return nil, err
	}
	return password, nil
}

func (p *PasswordRepositoryImpl) Delete(key string) error {
	stmt, err := p.db.Prepare(`
		DELETE FROM passwords where key = $1;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(key); err != nil {
		return err
	}
	return nil
}

func (p *PasswordRepositoryImpl) Update(password *model.Password) error {
	stmt, err := p.db.Prepare(`
		UPDATE passwords SET key = $1, password = $2 WHERE id = $3;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(password.Key, password.Pwd, password.Id); err != nil {
		return err
	}
	return nil
}
