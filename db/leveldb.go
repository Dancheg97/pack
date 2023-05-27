// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package db

import (
	"encoding/json"
	"path"

	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	*leveldb.DB
}

func GetLevelDB(p string) (Database, error) {
	db, err := leveldb.OpenFile(path.Join(p, "users.db"), nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{DB: db}, nil
}

func (d *LevelDB) List() ([]User, error) {
	var users []User
	iter := d.NewIterator(nil, nil)
	for iter.Next() {
		u := User{}
		b := iter.Value()
		err := json.Unmarshal(b, &u)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (d *LevelDB) Update(u User) error {
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return d.Put([]byte(u.Name), b, nil)
}
