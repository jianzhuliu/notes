/*
状态模式 + 简单工厂模式
内部修改状态
*/
package godesign

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

//状态类型
type StateType int

const (
	StateTypeYidong StateType = iota
	StateTypeLiantong
	StateTypeDianxin
)

func (s StateType) String() string {
	switch s {
	case StateTypeYidong:
		return "移动"
	case StateTypeLiantong:
		return "联通"
	case StateTypeDianxin:
		return "电信"
	default:
		return ""
	}
}

//状态上下文
type StateContext struct {
	Mobile  int    //手机号
	Message string //短信内容
}

var stateTypeMap = []StateType{StateTypeYidong, StateTypeLiantong, StateTypeDianxin}

func NewStateContext(mobile int, message string) *StateContext {
	return &StateContext{Mobile: mobile, Message: message}
}

type IState interface {
	Send(*StateContext) error
}

//状态管理者
type StateManager struct {
	currentState     IState        //当前状态对象
	currentStateType StateType     //当前状态类型
	setStateDuration time.Duration //切换状态的时间间隔
	mu               sync.Mutex
}

func NewStateManager(ctx context.Context, duration time.Duration) *StateManager {
	m := &StateManager{
		setStateDuration: duration,
	}

	//初始化状态对象
	m.setState()

	go func() {
		for {
			select {
			case <-time.NewTicker(duration).C:
				m.setState()
			case <-ctx.Done():
				return
			}
		}
	}()

	return m
}

func (m *StateManager) setState() {
	m.mu.Lock()
	defer m.mu.Unlock()

	//随机生成一个状态
	rand.Seed(time.Now().Unix())
	randIndex := rand.Intn(len(stateTypeMap))
	stateType := stateTypeMap[randIndex]

	m.currentStateType = stateType

	//简单工厂模式
	switch stateType {
	case StateTypeYidong:
		m.currentState = &YidongService{}
	case StateTypeLiantong:
		m.currentState = &LiantongService{}
	case StateTypeDianxin:
		m.currentState = &DianxinService{}
	}

	log.Println("切换状态", stateType)
}

func (m *StateManager) GetCurrentState() IState {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.currentState
}

///////////////////////////状态接口具体实现
//移动
type YidongService struct {
}

func (s *YidongService) Send(c *StateContext) error {
	log.Println("移动发送短信", c.Mobile, c.Message)
	return nil
}

//联通
type LiantongService struct {
}

func (s *LiantongService) Send(c *StateContext) error {
	log.Println("联通发送短信", c.Mobile, c.Message)
	return nil
}

//电信
type DianxinService struct {
}

func (s *DianxinService) Send(c *StateContext) error {
	log.Println("电信发送短信", c.Mobile, c.Message)
	return nil
}
