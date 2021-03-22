package goim

import "sync"

//用户管理
type UserMgr struct {
	users map[string]*User
	//涉及添加修改需要加锁
	userLock sync.RWMutex
}

//创建用户管理对象
func NewUserMgr() *UserMgr {
	return &UserMgr{
		users: make(map[string]*User),
	}
}

//添加一个用户
func (m *UserMgr) AddUser(user *User) {
	m.userLock.Lock()
	defer m.userLock.Unlock()

	m.users[user.Name] = user
}

//移除一个用户
func (m *UserMgr) RemoveUser(user *User) {
	m.userLock.Lock()
	defer m.userLock.Unlock()

	delete(m.users, user.Name)
}

//获取所有用户
func (m *UserMgr) GetAllUsers() []*User {
	m.userLock.RLock()
	defer m.userLock.RUnlock()

	//构建结果
	users := make([]*User, 0, len(m.users))

	for _, user := range m.users {
		users = append(users, user)
	}

	return users
}

//根据用户名，获取用户信息
func (m *UserMgr) GetUserByName(name string) (*User, bool) {
	m.userLock.RLock()
	defer m.userLock.RUnlock()

	user, ok := m.users[name]
	return user, ok
}

//广播消息
func (m *UserMgr) Broadcast(msg string) {
	users := m.GetAllUsers()
	for _, user := range users {
		user.Notify(msg)
	}
}
