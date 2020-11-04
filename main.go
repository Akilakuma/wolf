package main

import "log"

func main() {

	// 先讀檔,取得檔案內的lua script
	scriptValue := InitParser()

	// 取得一個將lua都compiled好的主體，rule handler
	ruleInstance := NewRuleHandler(scriptValue)

	// 告訴rule handler，要使用的rule、參數
	r1, e1 := ruleInstance.Parser("test", "hero tanoki")

	r2, e2 := ruleInstance.Parser("smart", "uncle Tom")

	log.Println("r1:", r1, "e1:", e1)
	log.Println("r2:", r2, "e2:", e2)

}
