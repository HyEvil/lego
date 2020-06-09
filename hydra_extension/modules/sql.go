package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"yym/hydra_extension/hydra"
)

func init() {
	hydra.RegisterType("SqlTX", &sqlTX{})
	hydra.RegisterType("SqlDB", NewDB)
}

type sqlExcer struct {
	ext sqlx.Ext
}

type sqlDB struct {
	sqlExcer
	db *sqlx.DB
}

type sqlTX struct {
	*sqlExcer
	tx *sqlx.Tx
}

func NewDB(dbType string, uri string) (*sqlDB, error) {
	s, err := sql.Open(dbType, uri)
	if err != nil {
		return nil, err
	}
	db := sqlx.NewDb(s, dbType)
	return &sqlDB{db: db, sqlExcer: sqlExcer{ext: db}}, nil
}

func (self *sqlDB) BeginTx() (*sqlTX, error) {
	tx, err := self.db.Beginx()
	if err != nil {
		return nil, err
	}
	return &sqlTX{sqlExcer: &sqlExcer{ext: tx}, tx: tx}, nil
}

func (self *sqlTX) Commit() ( error) {
	return self.tx.Commit()
}

func (self *sqlTX) Rollback() ( error) {
	return self.tx.Rollback()
}

func (self *sqlExcer) Query(sql string, args ...interface{}) ([]map[string]interface{}, error) {
	sql = self.ext.Rebind(sql)
	rows, err := self.ext.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := []map[string]interface{}{}
	for rows.Next() {
		m := map[string]interface{}{}
		err = rows.MapScan(m)
		if err != nil {
			return nil, err
		}
		ret = append(ret, m)
	}

	return ret, err
}

func (self *sqlExcer) NamedQuery(sql string, args map[string]interface{}) ([]map[string]interface{}, error) {
	query, argList, err := self.ext.BindNamed(sql, args)
	if err != nil {
		return nil, err
	}
	return self.Query(query, argList...)
}

func (self *sqlExcer) QueryOne(sql string, args ...interface{}) (map[string]interface{}, error) {
	sql = self.ext.Rebind(sql)
	row := self.ext.QueryRowx(sql, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	ret := map[string]interface{}{}
	err := row.MapScan(ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (self *sqlExcer) NamedQueryOne(sql string, args map[string]interface{}) (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	query, argList, err := self.ext.BindNamed(sql, args)
	if err != nil {
		return ret, err
	}
	return self.QueryOne(query, argList...)
}

func (self *sqlExcer) Insert(sql string, args ...interface{}) (int64, error) {
	sql = self.ext.Rebind(sql)
	ret, err := self.ext.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	id, err := ret.LastInsertId()
	return id, err
}

func (self *sqlExcer) NamedInsert(sql string, args map[string]interface{}) (int64, error) {
	query, argList, err := self.ext.BindNamed(sql, args)
	if err != nil {
		return 0, err
	}
	return self.Insert(query, argList...)
}

func (self *sqlExcer) Exec(sql string, args ...interface{}) error {
	sql = self.ext.Rebind(sql)
	_, err := self.ext.Exec(sql, args...)
	return err
}

func (self *sqlExcer) NamedExec(sql string, args map[string]interface{}) error {
	query, argList, err := self.ext.BindNamed(sql, args)
	if err != nil {
		return err
	}
	return self.Exec(query, argList...)
}
