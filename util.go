package main

import (
	lua "github.com/yuin/gopher-lua"
)

// 給lua 呼叫的 function
var exports = map[string]lua.LGFunction{
	//"get_random_by_probability": LuaGetRandomByProbability,
	//"split":                     LuaSplit,
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
