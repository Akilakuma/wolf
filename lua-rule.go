package main

import (
	"errors"
	"log"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type Rule struct {
	luaPool *lStatePool
	m       *sync.RWMutex
	rulesFn map[string]map[string]lua.LValue
}

// NewRuleHandler 把從檔案抓出來的lua script
// 丟到套件的Compiler處理
func NewRuleHandler(rules map[string][]byte) *Rule {

	rFn := make(map[string]map[string]lua.LValue)
	for k, v := range rules {
		m, err := Compile(k, v)
		if err != nil {
			log.Printf("%s rule err:%s\n", k, err)
			continue
		}
		log.Println("解析到的rule名稱是：", k)
		rFn[k] = m

	}
	r := &Rule{
		luaPool: &lStatePool{
			saved: make(chan *lua.LState, 1000),
		},
		rulesFn: rFn,
		m:       &sync.RWMutex{},
	}
	return r

}

// Compile 一次處理一個lua檔
func Compile(rule string, script []byte) (m map[string]lua.LValue, err error) {
	// 打開lua 通道
	ls := lua.NewState()
	defer ls.Close()

	// 先驗證對於lua來講，是否可以順利轉成string
	err = ls.DoString(string(script))
	if err != nil {
		log.Println(" syntax error ", err)
		return
	}

	// 取得lua檔案內的全域變數
	// .lua 檔裡面的變數，要跟檔名一樣
	t := ls.GetGlobal(rule)

	switch r := t.(type) {
	case *lua.LTable:

		m = make(map[string]lua.LValue)

		// parser_string : lua method in .lua file
		// parser_content : LValue(lua method) in golang
		m["parser_content"] = r.RawGet(lua.LString("parser_string"))

	case *lua.LNilType:
		err = errors.New("no rule nil")
	default:
		err = errors.New("no rule")
	}
	return
}

// Parser 丟給lua處理的method接口
func (r *Rule) Parser(rule string, content string) (parsedResult string, err error) {

	r.m.RLock()
	defer r.m.RUnlock()

	l := r.luaPool.Get()
	defer r.luaPool.Put(l)

	fnList, ok := r.rulesFn[rule]
	if !ok {
		log.Println("找不到設定檔")
	}
	fn, ok := fnList["parser_content"]
	if !ok {
		log.Println("找不到method")
	}
	if err := l.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}, lua.LString(content)); err != nil {
		log.Println(err.Error())
	}
	l.Get(-1)
	ret := l.Get(-1)
	l.Pop(1)

	switch v := ret.(type) {
	case lua.LString:
		parsedResult = string(v)
	}
	return
}
