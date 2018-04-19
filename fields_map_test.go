package sqlmapper

import (
	"context"
	"database/sql"
	"testing"
)

var (
	table = "test_table"
)

// DemoRow for `test_table`
type DemoRow struct {
	FieldKey string  `sql:"field_key"`
	FieldOne string  `sql:"field_one"`
	FieldTwo bool    `sql:"field_two"`
	FieldThr int64   `sql:"field_thr"`
	FieldFou float64 `sql:"field_fou"`
}

func TestSqlmapper(t *testing.T) {
	t.Log("====>TestSqlmapper")

	var db *sql.DB
	// db = GetDB() // get db in your own way
	if db == nil {
		t.Log("db is nil.")
		return
	}
	ctx := context.Background()

	// test Query
	row0, _ := QueryByKey(ctx, nil, db, "key001")
	rowArr1, _ := QueryByFieldOne(ctx, nil, db, "one")
	rowArrAll, _ := QueryAll(ctx, nil, db)
	t.Log(rowArr1)
	t.Log(rowArrAll)

	// test Update
	row0.FieldOne = "one123"
	row0.FieldTwo = true
	row0.FieldThr = 1234
	row0.FieldFou = 123.45
	_ = Update(ctx, nil, db, row0)

	// test Insert
	newRow0 := DemoRow{
		FieldKey: "key002",
		FieldOne: "one456",
		FieldTwo: false,
		FieldThr: 5678,
		FieldFou: 0.01,
	}
	newRow1 := DemoRow{
		FieldKey: "key003",
		FieldOne: "one789",
		FieldTwo: true,
		FieldThr: 5678,
		FieldFou: 0.02,
	}
	_ = Insert(ctx, nil, db, newRow0, newRow1)

	// test Remove
	_ = Remove(ctx, nil, db, "key001")

	t.Log("====>End")
}

// Query by primary key (field[0])
func QueryByKey(ctx context.Context, tx *sql.Tx, db *sql.DB, fieldKey string) (
	*DemoRow, error) {

	var row DemoRow
	row.FieldKey = fieldKey
	fm, err := NewFieldsMap(table, &row)
	if err != nil {
		return nil, err
	}

	objptr, err := fm.SQLSelectByPriKey(ctx, tx, db)
	if err != nil {
		return nil, err
	}

	return objptr.(*DemoRow), nil

}

// Query by `field_one`
func QueryByFieldOne(ctx context.Context, tx *sql.Tx, db *sql.DB, fieldOne string) (
	[]DemoRow, error) {

	var row DemoRow
	row.FieldOne = fieldOne
	fm, err := NewFieldsMap(table, &row)
	if err != nil {
		return nil, err
	}

	objptrs, err := fm.SQLSelectRowsByFieldNameInDB(ctx, tx, db, "field_one")
	if err != nil {
		return nil, err
	}

	var objs []DemoRow
	for i, olen := 0, len(objptrs); i < olen; i++ {
		objs = append(objs, *objptrs[i].(*DemoRow))
	}

	return objs, nil
}

// Query all
func QueryAll(ctx context.Context, tx *sql.Tx, db *sql.DB) ([]DemoRow, error) {

	var row DemoRow
	fm, err := NewFieldsMap(table, &row)
	if err != nil {
		return nil, err
	}

	objptrs, err := fm.SQLSelectAllRows(ctx, tx, db)
	if err != nil {
		return nil, err
	}

	var objs []DemoRow
	for i, olen := 0, len(objptrs); i < olen; i++ {
		objs = append(objs, *objptrs[i].(*DemoRow))
	}

	return objs, nil
}

// Insert
func Insert(ctx context.Context, tx *sql.Tx, db *sql.DB, rows ...DemoRow) error {

	for i, tlen := 0, len(rows); i < tlen; i++ {

		fm, err := NewFieldsMap(table, &rows[i])
		if err != nil {
			return err
		}

		err = fm.SQLInsert(ctx, tx, db)
		if err != nil {
			return err
		}
	}

	return nil
}

// Update by primary key (field[0])
func Update(ctx context.Context, tx *sql.Tx, db *sql.DB, row *DemoRow) error {

	fm, err := NewFieldsMap(table, row)
	if err != nil {
		return err
	}

	err = fm.SQLUpdateByPriKey(ctx, tx, db)
	if err != nil {
		return err
	}

	return nil
}

// Remove by primary key (field[0])
func Remove(ctx context.Context, tx *sql.Tx, db *sql.DB, fieldKey string) error {

	var row DemoRow
	row.FieldKey = fieldKey
	fm, err := NewFieldsMap(table, &row)
	if err != nil {
		return err
	}

	err = fm.SQLDeleteByPriKey(ctx, tx, db)
	if err != nil {
		return err
	}

	return nil
}
