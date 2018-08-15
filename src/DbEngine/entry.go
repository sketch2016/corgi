package DbEngine

import (
	_ "3rdparty/sql/mysql"
	"Utils"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

var dbEngineTag = "DB Engine"

//DataBase data base
type DataBase struct {
	db *sql.DB
}

//OpenMySQL open database
func OpenMySQL(user string, psw string, host string, port int, dbname string) (result *DataBase) {
	//sql.Open()
	v := new(DataBase)
	config := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", user, psw, host, port, dbname, "utf8")
	db, error := sql.Open("mysql", config)

	if error == nil {
		v.db = db
		return v
	}

	return nil
}

//InsertOnce one data to database
func (p *DataBase) InsertOnce(table string, data interface{}) {
	//we should checke all the table's column
	r1, _ := ormReflect(data)
	sql := "Insert into " + table + " "
	var valueSQL = "("
	var columnSQL = "("

	for key, val := range r1 {
		if val.inout != OUT {
			columnSQL += key + ","
			valueSQL += convertSQLVal(val.value) + ","
		}
	}

	sql += columnSQL[0:len(columnSQL)-1] + ") VALUES "
	sql += valueSQL[0:len(valueSQL)-1] + ")"

	p.db.Exec(sql)
	fmt.Println("sql is ", sql)
}

//InsertList insert list
//insert into tableName (col1,col2,col3,col4) values (value1，value2，value3，value4)，
//(value1，value2，value3，value4)......
func (p *DataBase) InsertList(table string, data []interface{}) {
	//fmt.Println("data type is ", data)

	//we should checke all the table's column
	if len(data) == 0 {
		return
	}

	index := 0
	value := data[index]
	fmt.Println("value is ", value)
	r1, _ := ormReflect(value)

	sql := "Insert into " + table + " "
	var valueSQL = "("
	var columnSQL = "("
	var columnlist = make([]string, 0)

	for key, val := range r1 {
		if val.inout != OUT {
			columnSQL += key + ","
			columnlist = append(columnlist, key)
			valueSQL += convertSQLVal(val.value) + ","
		}
	}

	sql += columnSQL[0:len(columnSQL)-1] + ") VALUES "
	sql += valueSQL[0:len(valueSQL)-1] + "),"
	index++

	for ; index < len(data); index++ {
		rr, _ := ormReflect(data[index])
		valueSQL = "("
		for _, columnkey := range columnlist {
			val, ok := rr[columnkey]
			if ok {
				if val.inout != OUT {
					valueSQL += convertSQLVal(val.value) + ","
				}
			}
		}

		sql += valueSQL[0:len(valueSQL)-1] + "),"
	}
	//fmt.Println("sql is ", sql)
	sql = sql[0 : len(sql)-1]
	fmt.Println("sql is ", sql)
	p.db.Exec(sql)
}

//Query select,data cannot be pointer!!!!
func (p *DataBase) Query(url string, data interface{}) []interface{} {
	//sql.Open()
	Utils.LOGD(dbEngineTag, "query start")
	defer func() {
		Utils.LOGD(dbEngineTag, "query end")
	}()

	r1, r2 := ormReflect(data)
	//fmt.Println("r1,r2 is ", r1, r2)

	var count = 0
	result := make([]interface{}, 0)

	rows, err := p.db.Query(url)
	Utils.LOGD(dbEngineTag, "query result")
	defer rows.Close()
	//fmt.Println("query err is ", err)
	if err == nil {
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		for rows.Next() {
			count++
			if err := rows.Scan(scanArgs...); err != nil {
				log.Fatalln(err)
			}

			//record := make(map[string]string)

			for i, col := range values {
				if col != nil {
					//fmt.Println("columns is ", columns[i])
					//fmt.Println("val is ", string(col.([]byte)))
					//record[columns[i]] = string(col.([]byte))
					val, ok := r1[columns[i]]
					if ok {
						//fmt.Println("val kind is ", val.Kind())
						//fmt.Println("val is ", val)
						setVal(val.value, string(col.([]byte)))
					}
				}
			}
			//start copy one

			finaldata := copyVal(data, &r2)
			result = append(result, finaldata)
		}
	}

	return result[0:count]
}

func setVal(val *reflect.Value, data string) {
	//fmt.Println("set Val is ", data)
	switch val.Kind() {
	case reflect.String:
		val.SetString(data)

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		convIntData, intErr := strconv.ParseInt(data, 10, 64)
		if intErr == nil {
			val.SetInt(convIntData)
		}

	case reflect.Bool:
		convBoolData, boolErr := strconv.ParseBool(data)
		if boolErr == nil {
			val.SetBool(convBoolData)
		}

	case reflect.Float32:
		convFloat32Data, float32Err := strconv.ParseFloat(data, 32)
		if float32Err == nil {
			val.SetFloat(convFloat32Data)
		}

	case reflect.Float64:
		convFloat64Data, float64Err := strconv.ParseFloat(data, 64)
		if float64Err == nil {
			val.SetFloat(convFloat64Data)
		}

	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		convUintData, uintErr := strconv.ParseUint(data, 10, 64)
		if uintErr == nil {
			val.SetUint(convUintData)
		}
	}
}

func copyVal(v interface{}, method *reflect.Value) interface{} {
	field := reflect.TypeOf(v)
	//fmt.Println("copy val field is ", field)
	//fmt.Println("v is ", v)
	if field.Kind() == reflect.Ptr {
		//method := reflect.ValueOf(v).MethodByName("ConvertValue")
		valuelist := method.Call(nil)
		return valuelist[0].Interface()
	}

	return v
}

func convertSQLVal(val *reflect.Value) string {
	//fmt.Println("set Val is ", data)
	switch val.Kind() {
	case reflect.String:
		return "'" + val.String() + "'"

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)

	case reflect.Bool:
		return strconv.FormatBool(val.Bool())

	case reflect.Float32:
		return strconv.FormatFloat(val.Float(), 'f', -1, 32)

	case reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64)

	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10)
	}

	return ""
}
