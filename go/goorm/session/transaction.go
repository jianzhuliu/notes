/*
事务的ACID属性
1、原子性(Atomicity)：事务中的全部操作在数据库中是不可分割的，要么全部完成，要么全部不执行。
2、一致性(Consistency): 几个并行执行的事务，其执行结果必须与按某一顺序 串行执行的结果相一致。
3、隔离性(Isolation)：事务的执行不受其他事务的干扰，事务执行的中间结果对其他事务必须是透明的。
4、持久性(Durability)：对于任意已提交事务，系统必须保证该事务对数据库的改变不被丢失，即使数据库出现故障。

*/

package session

import (
	"goorm/log"
)

//开启事务
func (s *Session) Begin() (err error) {
	log.Info("begin to transaction")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
	}

	return
}

//回滚事务
func (s *Session) Rollback() (err error) {
	log.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}

	return
}

//提交事务
func (s *Session) Commit() (err error) {
	log.Info("transaction commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}

	return
}
