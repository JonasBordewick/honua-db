package honuadb

import (
	"database/sql"
	"errors"

	"github.com/JonasBordewick/honua-db/models"
)

/*
CREATE TABLE IF NOT EXISTS honua_services (
    id INTEGER NOT NULL,
    identity TEXT NOT NULL,
    PRIMARY KEY(id, identity),
    CONSTRAINT fk_identity FOREIGN KEY(identity) REFERENCES identities(identifier) ON DELETE CASCADE,
    domain TEXT NOT NULL,
    name TEXT NOT NULL
);
*/

func (hdb *HonuaDB) AddHonuaService(service *models.HonuaService) error {
	const query = `INSERT INTO honua_services(id, identity, domain, name) VALUES($1, $2, $3, $4);`
	id, err := hdb.get_honua_service_id(service.Identity)
	if err != nil {
		return err
	}

	_, err = hdb.psqlDB.Exec(query, id, service.Identity, service.Domain, service.Name)
	return err
}

func (hdb *HonuaDB) GetHonuaServices(identity string) ([]*models.HonuaService, error) {
	const query = "SELECT * from honua_services WHERE identity = $1;"
	rows, err := hdb.psqlDB.Query(query, identity)
	if err != nil {
		return nil, err
	}

	var services []*models.HonuaService = []*models.HonuaService{}

	for rows.Next() {
		service, err := hdb.make_honua_service(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		services = append(services, service)
	}

	rows.Close()

	return services, nil
}

func (hdb *HonuaDB) DeleteHonuaService(identity string, serviceID int32) error {
	const query = "DELETE FROM honua_services WHERE identity=$1 AND id = $2;"

	_, err := hdb.psqlDB.Exec(query, identity, serviceID)
	return err
}

func (hdb *HonuaDB) get_honua_service_id(identity string) (int, error) {
	query := "SELECT CASE WHEN EXISTS ( SELECT * FROM honua_services WHERE identity = $1) THEN true ELSE false END"

	rows, err := hdb.psqlDB.Query(query, identity)
	if err != nil {
		return -1, err
	}

	var exist_identity bool = false

	for rows.Next() {
		err = rows.Scan(&exist_identity)
		if err != nil {
			rows.Close()
			return -1, err
		}
	}

	rows.Close()

	if !exist_identity {
		return 0, nil
	}

	query = "SELECT MAX(id) FROM honua_services WHERE identity = $1;"

	rows, err = hdb.psqlDB.Query(query, identity)
	if err != nil {
		return -1, err
	}

	var id int = -1

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			rows.Close()
			return -1, err
		}
	}
	rows.Close()

	if id == -1 {
		return -1, errors.New("something went wrong during getting id of entity")
	}

	id = id + 1

	return id, nil
}

func (hdb *HonuaDB) make_honua_service(rows *sql.Rows) (*models.HonuaService, error) {
	var id int32
	var identity string
	var domain string
	var name string

	err := rows.Scan(&id, &identity, &domain, &name)
	if err != nil {
		return nil, err
	}

	return &models.HonuaService{ID: id, Identity: identity, Domain: domain, Name: name}, nil
}