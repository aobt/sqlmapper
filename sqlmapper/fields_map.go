package sqlmapper

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
)

// Field db field
// describe struct mapping in DB like:
// type DemoRow struct {
// 	FieldKey string  `sql:"field_key"`
// 	FieldOne string  `sql:"field_one"`
// 	FieldTwo bool    `sql:"field_two"`
// 	FieldThr int64   `sql:"field_thr"`
// 	FieldFou float64 `sql:"field_fou"`
// }
//
type Field struct {
	Name       string
	Tag        string
	Type       string
	Addr       interface{}
	IntSave    sql.NullInt64
	StringSave sql.NullString
	FloatSave  sql.NullFloat64
	BoolSave   sql.NullBool
}

// FieldsMap hold Field
type FieldsMap interface {

	// GetFields Fields
	GetFields() []Field

	// GetFieldNamesInDB get Names in db from Fields
	GetFieldNamesInDB() []string

	// GetFieldValues get Values in Object(struct)
	GetFieldValues() []interface{}

	// GetFieldValue get Value in Object(struct)
	GetFieldValue(idx int) interface{}

	// GetFieldSaveAddrs get Pointers of Values in Object(struct)
	GetFieldSaveAddrs() []interface{}

	// GetFieldSaveAddr get Pointer of Value in Object(struct)
	GetFieldSaveAddr(idx int) interface{}

	// MapBackToObject mapping back to the original object
	MapBackToObject() interface{}

	////////////////////////////////////////////////////////////////
	// generate SQL string
	// SQLFieldsStr generate sqlstr in db from Fields
	SQLFieldsStr() string

	// SQLFieldsStrForSet generate sqlstr in db from Fields for set
	SQLFieldsStrForSet() string

	////////////////////////////////////////////////////////////////
	// generate statement
	// PrepareStmt prepare statement
	// Must Close after Stmt used
	PrepareStmt(ctx context.Context, tx *sql.Tx, db *sql.DB,
		sqlstr string) (*sql.Stmt, error)

	// SQLSelectStmt generate statement for SELECT
	SQLSelectStmt(ctx context.Context, tx *sql.Tx, db *sql.DB,
		extStr string) (*sql.Stmt, error)

	// SQLInsertStmt generate statement for INSERT
	SQLInsertStmt(ctx context.Context, tx *sql.Tx, db *sql.DB) (*sql.Stmt, error)

	// SQLUpdateStmt generate statement for UPDATE
	SQLUpdateStmt(ctx context.Context, tx *sql.Tx, db *sql.DB,
		extStr string) (*sql.Stmt, error)

	// SQLDeleteStmt generate statement for DELETE
	SQLDeleteStmt(ctx context.Context, tx *sql.Tx, db *sql.DB,
		extStr string) (*sql.Stmt, error)

	////////////////////////////////////////////////////////////////
	// exec sql
	// SQLLockByPriKey by primary key (field[0])
	SQLLockByPriKey(ctx context.Context, tx *sql.Tx,
		db *sql.DB) (interface{}, error)

	// SQLSelectByPriKey by primary key (field[0])
	SQLSelectByPriKey(ctx context.Context, tx *sql.Tx,
		db *sql.DB) (interface{}, error)

	// SQLSelectRowsByFieldNameInDB by field name in DB
	SQLSelectRowsByFieldNameInDB(ctx context.Context, tx *sql.Tx,
		db *sql.DB, nameInDB string) ([]interface{}, error)

	// SQLSelectAllRows
	SQLSelectAllRows(ctx context.Context, tx *sql.Tx,
		db *sql.DB) ([]interface{}, error)

	// SQLInsert
	SQLInsert(ctx context.Context, tx *sql.Tx, db *sql.DB) error

	// SQLUpdateByPriKey by primary key (field[0])
	SQLUpdateByPriKey(ctx context.Context, tx *sql.Tx, db *sql.DB) error

	// SQLDeleteByPriKey by primary key (field[0])
	SQLDeleteByPriKey(ctx context.Context, tx *sql.Tx, db *sql.DB) error
}

////////////////////////////////////////////////////////////////

// NewFieldsMap new Fields
func NewFieldsMap(table string, objptr interface{}) (FieldsMap, error) {

	elem := reflect.ValueOf(objptr).Elem()
	reftype := elem.Type()

	var fields []Field
	for i, flen := 0, reftype.NumField(); i < flen; i++ {

		var field Field
		field.Type = reftype.Field(i).Type.String()
		if field.Type != "int64" && field.Type != "string" &&
			field.Type != "float64" && field.Type != "bool" {
			return nil, errors.New("Unsupported Type: " + field.Type)
		}

		field.Name = reftype.Field(i).Name
		field.Tag = reftype.Field(i).Tag.Get("sql")
		field.Addr = elem.Field(i).Addr().Interface()
		fields = append(fields, field)
	}

	return &_FieldsMap{
		objptr:  objptr,
		reftype: reftype,
		fields:  fields,
		table:   table,
	}, nil
}

////////////////////////////////////////////////////////////////

var _ FieldsMap = &_FieldsMap{}

type _FieldsMap struct {
	objptr  interface{}
	reftype reflect.Type
	fields  []Field
	table   string
}

// GetFields get Fields for an Object(struct)
func (fds *_FieldsMap) GetFields() []Field {

	return fds.fields
}

// GetFieldNamesInDB get Names in db from Fields
// example:
// type DemoRow struct {
// 	FieldKey string  `sql:"field_key"`
// 	FieldOne string  `sql:"field_one"`
// 	FieldTwo bool    `sql:"field_two"`
// 	FieldThr int64   `sql:"field_thr"`
// 	FieldFou float64 `sql:"field_fou"`
// }
//
// return ["field_key", "field_one", "field_two","field_thr","field_fou"]
//
func (fds *_FieldsMap) GetFieldNamesInDB() []string {

	var tags []string
	for i, flen := 0, len(fds.fields); i < flen; i++ {
		tags = append(tags, fds.fields[i].Tag)
	}

	return tags
}

// GetFieldValues get Values in Object(struct)
func (fds *_FieldsMap) GetFieldValues() []interface{} {

	var values []interface{}
	for i, flen := 0, len(fds.fields); i < flen; i++ {
		values = append(values, fds.GetFieldValue(i))
	}

	return values
}

// GetFieldValue get Values in Object(struct)
func (fds *_FieldsMap) GetFieldValue(idx int) interface{} {

	switch fds.fields[idx].Type {
	case "int64":
		return *fds.fields[idx].Addr.(*int64)
	case "string":
		return *fds.fields[idx].Addr.(*string)
	case "float64":
		return *fds.fields[idx].Addr.(*float64)
	case "bool":
		return *fds.fields[idx].Addr.(*bool)
	default:
	}

	return nil
}

// GetFieldSaveAddrs get Pointers of Values in Object(struct)
func (fds *_FieldsMap) GetFieldSaveAddrs() []interface{} {

	var addrs []interface{}
	for i, flen := 0, len(fds.fields); i < flen; i++ {
		addrs = append(addrs, fds.GetFieldSaveAddr(i))
	}

	return addrs
}

// GetFieldSaveAddr get Pointers of Values in Object(struct)
func (fds *_FieldsMap) GetFieldSaveAddr(idx int) interface{} {

	switch fds.fields[idx].Type {
	case "int64":
		return &fds.fields[idx].IntSave
	case "string":
		return &fds.fields[idx].StringSave
	case "float64":
		return &fds.fields[idx].FloatSave
	case "bool":
		return &fds.fields[idx].BoolSave
	default:
	}

	return nil
}

// MapBackToObject mapping back to the original object
func (fds *_FieldsMap) MapBackToObject() interface{} {

	for i, flen := 0, len(fds.fields); i < flen; i++ {
		switch fds.fields[i].Type {
		case "int64":
			if fds.fields[i].IntSave.Valid {
				*fds.fields[i].Addr.(*int64) = fds.fields[i].IntSave.Int64
			}
			break
		case "string":
			if fds.fields[i].StringSave.Valid {
				*fds.fields[i].Addr.(*string) = fds.fields[i].StringSave.String
			}
			break
		case "float64":
			if fds.fields[i].FloatSave.Valid {
				*fds.fields[i].Addr.(*float64) = fds.fields[i].FloatSave.Float64
			}
			break
		case "bool":
			if fds.fields[i].BoolSave.Valid {
				*fds.fields[i].Addr.(*bool) = fds.fields[i].BoolSave.Bool
			}
			break
		default:
		}
	}

	return fds.objptr
}

////////////////////////////////////////////////////////////////
// generate SQL string

// SQLFieldsStr generate sqlstr in db from Fields
// example:" `field0`, `field1`, `field2`, `field3` "
func (fds *_FieldsMap) SQLFieldsStr() string {

	var tagsStr string
	for i, flen := 0, len(fds.fields); i < flen; i++ {
		if len(tagsStr) > 0 {
			tagsStr += ", "
		}
		tagsStr += "`"
		tagsStr += fds.fields[i].Tag
		tagsStr += "`"
	}
	if len(tagsStr) > 0 {
		tagsStr += " "
		tagsStr = " " + tagsStr
	}

	return tagsStr
}

// SQLFieldsStrForSet generate sqlstr in db from Fields for set
// example:" `field0` = ?, `field1` = ?, `field2` = ?, `field3` = ? "
func (fds *_FieldsMap) SQLFieldsStrForSet() string {

	var tagsStr string
	for i, flen := 0, len(fds.fields); i < flen; i++ {
		if len(tagsStr) > 0 {
			tagsStr += ", "
		}
		tagsStr += "`"
		tagsStr += fds.fields[i].Tag
		tagsStr += "`"
		tagsStr += " = ?"
	}
	if len(tagsStr) > 0 {
		tagsStr += " "
		tagsStr = " " + tagsStr
	}

	return tagsStr
}

////////////////////////////////////////////////////////////////
// generate statement

// PrepareStmt prepare statement
func (fds *_FieldsMap) PrepareStmt(ctx context.Context, tx *sql.Tx, db *sql.DB,
	sqlstr string) (*sql.Stmt, error) {

	if tx != nil {
		return tx.PrepareContext(ctx, sqlstr)
	}

	if db != nil {
		return db.PrepareContext(ctx, sqlstr)
	}

	return nil, errors.New("tx & db both nil")
}

// SQLSelectStmt generate statement for SELECT
func (fds *_FieldsMap) SQLSelectStmt(ctx context.Context, tx *sql.Tx, db *sql.DB,
	extStr string) (*sql.Stmt, error) {

	sqlstr := "SELECT " + fds.SQLFieldsStr() +
		" FROM `" + fds.table + "` " + extStr

	return fds.PrepareStmt(ctx, tx, db, sqlstr)
}

// SQLInsertStmt generate statement for INSERT
func (fds *_FieldsMap) SQLInsertStmt(ctx context.Context, tx *sql.Tx, db *sql.DB) (*sql.Stmt, error) {

	var vs string
	for i, flen := 0, len(fds.fields); i < flen; i++ {
		if len(vs) > 0 {
			vs += ", "
		}
		vs += "?"
	}

	sqlstr := "INSERT INTO `" + fds.table + "` (" + fds.SQLFieldsStr() + ") " +
		"VALUES (" + vs + ")"
	return fds.PrepareStmt(ctx, tx, db, sqlstr)
}

// SQLUpdateStmt generate statement for UPDATE
func (fds *_FieldsMap) SQLUpdateStmt(ctx context.Context, tx *sql.Tx, db *sql.DB,
	extStr string) (*sql.Stmt, error) {

	sqlstr := "UPDATE `" + fds.table + "` SET " + fds.SQLFieldsStrForSet() + extStr
	return fds.PrepareStmt(ctx, tx, db, sqlstr)
}

// SQLDeleteStmt generate statement for DELETE
func (fds *_FieldsMap) SQLDeleteStmt(ctx context.Context, tx *sql.Tx, db *sql.DB,
	extStr string) (*sql.Stmt, error) {

	sqlstr := "DELETE FROM `" + fds.table + "` " + extStr
	return fds.PrepareStmt(ctx, tx, db, sqlstr)
}

////////////////////////////////////////////////////////////////
// exec sql

// SQLLockByPriKey by primary key (field[0])
func (fds *_FieldsMap) SQLLockByPriKey(ctx context.Context, tx *sql.Tx,
	db *sql.DB) (interface{}, error) {

	extStr := " where `" + fds.fields[0].Tag + "` = ? for update "
	stmt, err := fds.SQLSelectStmt(ctx, tx, db, extStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close() // must close stmt after stmt used

	r := stmt.QueryRowContext(ctx, fds.GetFieldValue(0))
	if r == nil {
		return nil, errors.New("row is nil")
	}

	err = r.Scan(fds.GetFieldSaveAddrs()...)
	if err != nil {
		return nil, err
	}

	return fds.MapBackToObject(), nil
}

// SQLSelectByPriKey by primary key (field[0])
func (fds *_FieldsMap) SQLSelectByPriKey(ctx context.Context, tx *sql.Tx,
	db *sql.DB) (interface{}, error) {

	extStr := " where `" + fds.fields[0].Tag + "` = ? "
	stmt, err := fds.SQLSelectStmt(ctx, tx, db, extStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close() // must close stmt after stmt used

	r := stmt.QueryRowContext(ctx, fds.GetFieldValue(0))
	if r == nil {
		return nil, errors.New("row is nil")
	}

	err = r.Scan(fds.GetFieldSaveAddrs()...)
	if err != nil {
		return nil, err
	}

	return fds.MapBackToObject(), nil
}

// SQLSelectRowsByFieldNameInDB by field name in DB
func (fds *_FieldsMap) SQLSelectRowsByFieldNameInDB(ctx context.Context, tx *sql.Tx,
	db *sql.DB, nameInDB string) ([]interface{}, error) {

	idx := -1
	for i, flen := 0, len(fds.fields); i < flen; i++ {
		if fds.fields[i].Tag == nameInDB {
			idx = i
			break
		}
	}

	if idx < 0 {
		return nil, errors.New("no field match `sql` tag:" + nameInDB)
	}

	extStr := " where `" + fds.fields[idx].Tag + "` = ? "
	stmt, err := fds.SQLSelectStmt(ctx, tx, db, extStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close() // must close stmt after stmt used

	rs, err := stmt.QueryContext(ctx, fds.GetFieldValue(idx))
	if err != nil {
		return nil, err
	}

	var objs []interface{}
	for rs.Next() {
		obj := reflect.New(fds.reftype).Interface()
		fieldsMap, err := NewFieldsMap(fds.table, obj)
		if err != nil {
			return nil, err
		}

		err = rs.Scan(fieldsMap.GetFieldSaveAddrs()...)
		if err != nil {
			return nil, err
		}
		fieldsMap.MapBackToObject()
		objs = append(objs, obj)
	}

	return objs, nil
}

// SQLSelectAllRows
func (fds *_FieldsMap) SQLSelectAllRows(ctx context.Context, tx *sql.Tx,
	db *sql.DB) ([]interface{}, error) {

	stmt, err := fds.SQLSelectStmt(ctx, tx, db, "")
	if err != nil {
		return nil, err
	}
	defer stmt.Close() // must close stmt after stmt used

	rs, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var objs []interface{}
	for rs.Next() {
		obj := reflect.New(fds.reftype).Interface()
		fieldsMap, err := NewFieldsMap(fds.table, obj)
		if err != nil {
			return nil, err
		}

		err = rs.Scan(fieldsMap.GetFieldSaveAddrs()...)
		if err != nil {
			return nil, err
		}
		fieldsMap.MapBackToObject()
		objs = append(objs, obj)
	}

	return objs, nil
}

// SQLInsert
func (fds *_FieldsMap) SQLInsert(ctx context.Context, tx *sql.Tx,
	db *sql.DB) error {

	stmt, err := fds.SQLInsertStmt(ctx, tx, db)
	if err != nil {
		return err
	}
	defer stmt.Close() // must close stmt after stmt used

	_, err = stmt.ExecContext(ctx, fds.GetFieldValues()...)
	if err != nil {
		return err
	}

	return nil
}

// SQLUpdateByPriKey by primary key (field[0])
func (fds *_FieldsMap) SQLUpdateByPriKey(ctx context.Context, tx *sql.Tx,
	db *sql.DB) error {

	extStr := " where `" + fds.fields[0].Tag + "` = ? "
	stmt, err := fds.SQLUpdateStmt(ctx, tx, db, extStr)
	if err != nil {
		return err
	}
	defer stmt.Close() // must close stmt after stmt used

	values := fds.GetFieldValues()
	values = append(values, fds.GetFieldValue(0))
	_, err = stmt.ExecContext(ctx, values...)
	if err != nil {
		return err
	}

	return nil
}

// SQLDeleteByPriKey by primary key (field[0])
func (fds *_FieldsMap) SQLDeleteByPriKey(ctx context.Context, tx *sql.Tx,
	db *sql.DB) error {

	extStr := " where `" + fds.fields[0].Tag + "` = ? "
	stmt, err := fds.SQLDeleteStmt(ctx, tx, db, extStr)
	if err != nil {
		return err
	}
	defer stmt.Close() // must close stmt after stmt used

	_, err = stmt.ExecContext(ctx, fds.GetFieldValue(0))
	if err != nil {
		return err
	}

	return nil
}
