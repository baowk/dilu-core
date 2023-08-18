package core

type (
	Mode   string
	DBType string
)

const (
	ModeDev  Mode = "dev"  //开发模式
	ModeTest Mode = "test" //测试模式
	ModeProd Mode = "prod" //生产模式

	Mysql  DBType = "mysql" //mysql数据库标识
	Mssql  DBType = "mssql"
	Pgsql  DBType = "pgsql"
	Sqlite DBType = "sqlite" //sqlite
)

func (e Mode) String() string {
	return string(e)
}

func (e DBType) String() string {
	return string(e)
}
