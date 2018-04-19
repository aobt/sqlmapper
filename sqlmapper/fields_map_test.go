package sqlmapper

import (
	"testing"
)

func TestSqlmapper(t *testing.T) {
	t.Log("====>TestSqlmapper")

	// var db *sql.DB
	// db = GetDB() // get db
	// ctx := context.Background()

	// var demo DemoRow
	// demo.FieldKey = "key001"
	// fieldsMap, err := NewFieldsMap("test_table", &demo)
	// if err != nil {
	// 	//...
	// }

	// // Select By PriKey(the first field in struct)
	// objptr, err := fieldsMap.SQLSelectByPriKey(ctx, nil, db)
	// if err != nil {
	// 	//...
	// }

	// // objptr.(*DemoRow) get the pointer to the return obj
	// log.Println(objptr.(*DemoRow))

	// objptrsAll, err := fieldsMap.SQLSelectAllRows(ctx, nil, db)
	// if err != nil {
	// 	t.Error(err)
	// 	return err
	// }

}
