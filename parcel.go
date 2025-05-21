package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (number, client, status, address, created_at) values (:number, :client, :status, :address, :created_at)",
		sql.Named("number", p.Number),
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = $1", number)

	if row == nil {
		return Parcel{}, sql.ErrNoRows
	}

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = $1", client)
	if err != nil {
		return nil, err
	}

	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		if rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return res, err
		}
		res = append(res, p)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = $1 WHERE number = $2", status, number)
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	parcel, err := s.Get(number)

	if err != nil {
		return nil
	}

	if parcel.Status == ParcelStatusRegistered {
		_, err = s.db.Exec("UPDATE parcel SET address = $1 WHERE number = $2", address, number)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	parcel, err := s.Get(number)

	if err != nil {
		return nil
	}

	if parcel.Status == ParcelStatusRegistered {
		_, err = s.db.Exec("Delete from parcel where number = $1", number)
		if err != nil {
			return err
		}
	}

	return nil
}
