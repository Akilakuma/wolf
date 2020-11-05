package main

import (
	"github.com/tidwall/gjson"
	lua "github.com/yuin/gopher-lua"
)

// 給lua 呼叫的 function
var exports = map[string]lua.LGFunction{
	"gjson": LuaGJsonFindKey,
}

// LuaGJsonFindKey 利用gjson，找出json key
func LuaGJsonFindKey(l *lua.LState) int {

	// 取出lua第一個參數
	source := string(l.ToString(1))
	// 取出lua第二個參數
	key := string(l.ToString(2))

	jsonValue := gjson.Get(source, key)
	v := jsonValue.String()

	// 轉成lua的型態後，塞回去
	l.Push(lua.LString(v))

	return 1
}

type lStatePool struct {
	saved chan *lua.LState
}

// Get 取得一個lua state 使用
func (pl *lStatePool) Get() (l *lua.LState) {

	select {
	case l = <-pl.saved:
	default:
		l = pl.New()
	}
	return
}

// New 產生一個新的lua state
func (pl *lStatePool) New() (l *lua.LState) {

	l = lua.NewState()
	l.PreloadModule("inmodule", Loader)
	return
}

func Loader(l *lua.LState) int {
	mod := l.SetFuncs(l.NewTable(), exports)
	l.Push(mod)
	return 1
}

// Put 把用完的lua state 塞回去
func (pl *lStatePool) Put(L *lua.LState) {
	select {
	case pl.saved <- L:
	default:
	}
}
