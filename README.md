#### `sqlmapper` is a light mapper between `golang struct` and `table rows` in db

### example
We need to read/write a table in db, like:
```sql
CREATE TABLE `test_table` (
  `field_key` varchar(64) NOT NULL DEFAULT '',
  `field_one` varchar(64) DEFAULT NULL,
  `field_two` tinyint(1) DEFAULT NULL,
  `field_thr` int(12) DEFAULT NULL,
  `field_fou` float DEFAULT NULL,
  PRIMARY KEY (`field_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```
In golang, we create a struct corresponding to the table, like:
```go
// struct in go such as:
type DemoRow struct {
	FieldKey string  `sql:"field_key"`
	FieldOne string  `sql:"field_one"`
	FieldTwo bool    `sql:"field_two"`
	FieldThr int64   `sql:"field_thr"`
	FieldFou float64 `sql:"field_fou"`
}
```
Then, we can execute `SELECT`/`INSERT`/`UPDATE`/`DELETE` 

without long `Hard-Code` sql string which is easy to make mistakes.

sample (follow [fields_map_test.go](https://github.com/arthas29/sqlmapper/blob/master/fields_map_test.go) for more):
```go

// select single row
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
    
```
