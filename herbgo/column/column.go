package column

import (
	"errors"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herbgo/name"
)

var Drivers = map[string]func() ColumnsLoader{}

type Column struct {
	Field      string
	ColumnType string
	AutoValue  bool
	PrimayKey  bool
	NotNull    bool
}

func (c *Column) Name() string {
	n := name.MustNewOrExitWithoutParents(c.Field)
	return n.Pascal
}

type Columns []Column

type ColumnsLoader interface {
	Columns() ([]Column, error)
	Load(conn db.Database, table string) error
}
type ModelColumns struct {
	Columns     []Column
	Name        *name.Name
	Database    string
	PrimaryKeys []Column
	HasTime     bool
}

func (m *ModelColumns) FirstPrimayKey() Column {
	return m.Columns[0]
}

func (m *ModelColumns) CanCreate() bool {
	if len(m.PrimaryKeys) == 1 {
		if m.PrimaryKeys[0].ColumnType == "string" {
			return true
		}
		if m.PrimaryKeys[0].AutoValue {
			return true
		}
	}
	return false
}

func (m *ModelColumns) HasPrimayKey() bool {
	return len(m.Columns) > 0
}

func (m *ModelColumns) PrimaryKeyType() string {
	output := "//" + m.Name.Pascal + "PrimaryKey : table " + m.Name.Raw + " primary key type\n"
	switch len(m.PrimaryKeys) {
	case 0:
		output = output + "type " + m.Name.Pascal + "PrimaryKey map[string]interface{}\n"
	case 1:
		output = output + "type " + m.Name.Pascal + "PrimaryKey "
		if !m.PrimaryKeys[0].NotNull {
			output = output + "*"
		}
		output = output + m.Columns[0].ColumnType + "\n"
	default:
		output = output + "type " + m.Name.Pascal + "PrimaryKey struct{\n"
		for _, v := range m.PrimaryKeys {
			output = output + "    " + v.Name() + " "
			if !v.NotNull {
				output = output + "*"
			}
			output = output + v.ColumnType + "\n"
		}
		output = output + "}\n"
	}
	return output
}

func (m *ModelColumns) BuildByPKQuery() string {
	output := "//BuildByPKQuery : build by pk query for table " + m.Name.Raw + "\n"
	output = output + "func (mapper *" + m.Name.Pascal + "Mapper) BuildByPKQuery(pk interface{}) *querybuilder.PlainQuery {\n"
	output = output + "var query= mapper.QueryBuilder()\n"
	output = output + "    var q = query.New(\"\")\n"
	switch len(m.PrimaryKeys) {
	case 0:
		output = output + "    for k,v :=range pk {\n"
		output = output + " q.And(query.Equal( " + m.Name.Pascal + ".FieldAlias(k),v))\n"
		output = output + "    }"
	case 1:
		output = output + "    q.And(query.Equal(" + m.Name.Pascal + "FieldAlias" + m.Columns[0].Name() + ",pk))\n"
	default:
		for _, v := range m.PrimaryKeys {
			output = output + "    q.And(query.Equal(" + m.Name.Pascal + "FieldAlias" + v.Name() + ",pk." + v.Name() + "))\n"
		}
	}
	output = output + "    return q\n}\n"
	return output
}
func (m *ModelColumns) ModelPrimaryKey() string {
	output := "//ModelPrimaryKey :  get primary key from model.\n"
	output = output + "func (mapper *" + m.Name.Pascal + "Mapper) ModelPrimaryKey( model *" + m.Name.Pascal + "Model ) *" + m.Name.Pascal + "PrimaryKey {\n"
	switch len(m.PrimaryKeys) {
	case 0:
		output = output + "    return nil\n"
	case 1:
		output = output + "    var pk " + m.Name.Pascal + "PrimaryKey\n"
		output = output + "    pk=" + m.Name.Pascal + "PrimaryKey(model." + m.Columns[0].Name() + ")\n"
		output = output + "    return &pk\n"
	default:
		output = output + "    pk:=" + m.Name.Pascal + "PrimaryKey{}\n"
		for _, v := range m.PrimaryKeys {
			output = output + "    pk." + v.Name() + " = model." + v.Name() + "\n"
		}
		output = output + "    return &pk\n"
	}
	output = output + "}\n"
	return output
}
func (m *ModelColumns) ColumnsToModelStruct() string {
	output := "//" + m.Name.Pascal + "Model :" + m.Name.Raw + " model.\n"
	output = output + "type " + m.Name.Pascal + "Model struct{\n"
	for _, v := range m.Columns {
		output = output + "    " + v.Name() + " "
		if !v.NotNull {
			output = output + "*"
		}
		output = output + v.ColumnType + "\n"
	}
	output = output + "}\n"
	return output
}

func (m *ModelColumns) ColumnsToFieldsMethod(query string) string {
	qn, _ := name.MustNew(false, query)
	output := "//" + qn.Pascal + "Fields : map model " + qn.Raw + " fields to database column.\n"
	output = output + "func (mapper *" + m.Name.Pascal + "Mapper) " + qn.Pascal + "Fields(model *" + m.Name.Pascal + "Model)* querybuilder.Fields {\n"
	output = output + "	   if model == nil {\n"
	output = output + "        model = New" + m.Name.Pascal + "Model()\n"
	output = output + "	   }\n"
	output = output + "    return model.BuildFields(true,\n"
	for _, v := range m.Columns {

		output = output + "	//Field \"" + m.Name.Raw + "." + v.Field + "\"\n	" + m.Name.Pascal + "Field" + v.Name() + ","
		output = output + "\n"
	}
	output = output + "    )\n"
	output = output + "}\n"
	return output
}

func (m *ModelColumns) ColumnsToFieldsInsertMethod(query string) string {
	qn, _ := name.MustNew(false, query)
	output := "//" + qn.Pascal + "FieldsInsert : map model " + qn.Raw + " fields to database column used in insert query.\n"
	output = output + "func (mapper *" + m.Name.Pascal + "Mapper) " + qn.Pascal + "FieldsInsert(model *" + m.Name.Pascal + "Model)* querybuilder.Fields {\n"
	output = output + "    return model.BuildFields(false,"
	skiped := ""
	fields := ""
	for _, v := range m.Columns {
		if v.AutoValue {
			skiped = skiped + "\n		//Skip field \"" + v.Field + "\" which should be set by database"
			skiped = skiped + "\n		 //Field \"" + m.Name.Raw + "." + v.Field + "\"\n		//" + m.Name.Pascal + "Field" + v.Name() + ","
		} else {
			fields = fields + "\n		 //Field \"" + m.Name.Raw + "." + v.Field + "\"\n		" + m.Name.Pascal + "Field" + v.Name() + ","
		}
	}

	output = output + skiped + fields + "\n    )\n}\n"
	return output
}
func (m *ModelColumns) ColumnsToFieldsUpdateMethod(query string) string {
	qn, _ := name.MustNew(false, query)
	output := "//" + qn.Pascal + "FieldsUpdate : map model " + qn.Raw + " fields to database column used in update query.\n"
	output = output + "func (mapper *" + m.Name.Pascal + "Mapper) " + qn.Pascal + "FieldsUpdate(model *" + m.Name.Pascal + "Model)* querybuilder.Fields {\n"
	output = output + "    return model.BuildFields(false,"
	skiped := ""
	primaryKey := ""
	fields := ""
	for _, v := range m.Columns {
		if v.AutoValue {
			skiped = skiped + "\n		//Skip field \"" + v.Field + "\" which should be set by database"
			skiped = skiped + "\n		//Field \"" + m.Name.Raw + "." + v.Field + "\"\n		//" + m.Name.Pascal + "Field" + v.Name() + ","
		} else if v.PrimayKey {
			primaryKey = primaryKey + "\n		//Skip primary key field \"" + v.Field + "\""
			primaryKey = primaryKey + "\n		//Field \"" + m.Name.Raw + "." + v.Field + "\"\n		//" + m.Name.Pascal + "Field" + v.Name() + ","

		} else {
			fields = fields + "\n		//Field \"" + m.Name.Raw + "." + v.Field + "\"\n		" + m.Name.Pascal + "Field" + v.Name() + ","
		}
	}

	output = output + skiped + fields + primaryKey + "\n    )\n}\n"
	return output
}
func getLoaderFormDB(conn db.Database) (ColumnsLoader, error) {
	drivername := conn.Driver()
	driver, ok := Drivers[drivername]
	if ok == false {
		return nil, errors.New("unsupported sql driver " + drivername)
	}
	return driver(), nil
}
func New(conn db.Database, database string, table string) (*ModelColumns, error) {
	loader, err := getLoaderFormDB(conn)
	if err != nil {
		return nil, err
	}
	err = loader.Load(conn, table)
	if err != nil {
		return nil, err
	}
	columns, err := loader.Columns()
	if err != nil {
		return nil, err
	}
	pks := []Column{}
	var hasTime bool
	for _, v := range columns {
		if v.PrimayKey {
			pks = append(pks, v)
		}
		if v.ColumnType == "time.Time" {
			hasTime = true
		}
	}
	c := &ModelColumns{
		Columns:     columns,
		Name:        name.MustNewOrExitWithoutParents(table),
		Database:    database,
		PrimaryKeys: pks,
		HasTime:     hasTime,
	}
	return c, nil
}
