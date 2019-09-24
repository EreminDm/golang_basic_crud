package mariadb

import (
	"context"
	"fmt"

	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/pkg/errors"
)

// personalData description.
type personalData struct {
	ID          string
	Name        string
	LastName    string
	Phone       string
	Email       string
	YearOfBirth int
}

// receive returns mongo personal data construction.
func receive(ep entity.PersonalData) personalData {
	return personalData{
		ID:          ep.DocumentID,
		Name:        ep.Name,
		LastName:    ep.LastName,
		Phone:       ep.Phone,
		Email:       ep.Email,
		YearOfBirth: ep.YearOfBirth,
	}
}

// transmit returns entity data construction.
func (p personalData) transmit() entity.PersonalData {
	return entity.PersonalData{
		DocumentID:  p.ID,
		Name:        p.Name,
		LastName:    p.LastName,
		Phone:       p.Phone,
		Email:       p.Email,
		YearOfBirth: p.YearOfBirth,
	}
}

// One returns personal data for a given id,
// id params to make filtration.
func (m MariaDB) One(ctx context.Context, id string) (entity.PersonalData, error) {
	sqlQuery := "SELECT * FROM person WHERE id = ?"
	var p personalData

	row := m.Person.QueryRowContext(ctx, sqlQuery, id)
	err := row.Scan(&p.ID, &p.Name, &p.LastName, &p.Phone, &p.Email, &p.YearOfBirth)
	if err != nil {
		return entity.PersonalData{}, errors.Wrap(err, "could not scan row")
	}

	return p.transmit(), nil
}

// Insert is a function which adding data to database.
func (m MariaDB) Insert(ctx context.Context, document entity.PersonalData) (entity.PersonalData, error) {
	p := receive(document)
	sqlQuery := "INSERT INTO person (id, name, last_name, phone, email, year_od_birth ) VALUES (?,?,?,?,?,?);"
	_, err := m.Person.ExecContext(ctx, sqlQuery, p.ID, p.Name, p.LastName, p.Phone, p.Email, p.YearOfBirth)
	if err != nil {
		return entity.PersonalData{}, errors.Wrap(err, "could not exec query statement")
	}
	return document, nil
}

// All selects all documents from database.
func (m MariaDB) All(ctx context.Context) ([]entity.PersonalData, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM person")
	var p personalData
	var persons []entity.PersonalData
	rows, err := m.Person.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, errors.Wrap(err, "could not make query")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&p.ID, &p.Name, &p.LastName, &p.Phone, &p.Email, &p.YearOfBirth)
		if err != nil {
			return nil, errors.Wrap(err, "could not scan rows")
		}
		persons = append(persons, p.transmit())
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows error")
	}
	return persons, nil
}

// Remove deletes document from Mongo.
func (m MariaDB) Remove(ctx context.Context, id string) (int64, error) {
	sqlQuery := "DELETE FROM person WHERE id = ?"
	rslt, err := m.Person.ExecContext(ctx, sqlQuery, id)
	if err != nil {
		return 0, errors.Wrap(err, "could not remove data")
	}
	count, err := rslt.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "rows are not affected")
	}
	return count, nil
}

// Update rewrites information in db by user id filtration.
func (m MariaDB) Update(ctx context.Context, ep entity.PersonalData) (int64, error) {
	p := receive(ep)
	sqlQuery := "UPDATE person SET name=?, last_name=?, phone=?, email=?, year_od_birth=? where id= ?"

	rslt, err := m.Person.ExecContext(ctx, sqlQuery, p.Name, p.LastName, p.Phone, p.Email, p.YearOfBirth, p.ID)
	if err != nil {
		return 0, errors.Wrap(err, "could not update data")
	}
	count, err := rslt.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "rows are not affected")
	}
	return count, nil
}
