/*
builder模式
*/

package codepattern

type MysqlConfBuilder struct {
	MysqlConf
}

func (b *MysqlConfBuilder) Create(host string) *MysqlConfBuilder {
	s := MysqlConf{
		Host:    host,
		Port:    3306,
		User:    "demo",
		Passwd:  "demo123",
		Dbname:  "test",
		Charset: "utf8mb4",
	}

	b.MysqlConf = s
	return b
}

func (b *MysqlConfBuilder) WithPort(port int) *MysqlConfBuilder {
	b.MysqlConf.Port = port
	return b
}

func (b *MysqlConfBuilder) WithUser(user string) *MysqlConfBuilder {
	b.MysqlConf.User = user
	return b
}

func (b *MysqlConfBuilder) WithPasswd(passwd string) *MysqlConfBuilder {
	b.MysqlConf.Passwd = passwd
	return b
}

func (b *MysqlConfBuilder) WithDbname(dbname string) *MysqlConfBuilder {
	b.MysqlConf.Dbname = dbname
	return b
}

func (b *MysqlConfBuilder) WithCharset(charset string) *MysqlConfBuilder {
	b.MysqlConf.Charset = charset
	return b
}

func (b *MysqlConfBuilder) Build() MysqlConf {
	return b.MysqlConf
}
