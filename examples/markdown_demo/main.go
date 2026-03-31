package main

import (
	"fmt"
	"log"
)

func main() {
	demos := []struct {
		name string
		run  func() error
	}{
		{name: "数学公式示例", run: runMathFormulaDemo},
		{name: "软换行示例", run: runSoftLinebreakDemo},
		{name: "表格与任务列表示例", run: runTableAndTasklistDemo},
	}

	for _, demo := range demos {
		fmt.Printf("\n=== 开始执行: %s ===\n", demo.name)
		if err := demo.run(); err != nil {
			log.Fatalf("执行失败（%s）: %v", demo.name, err)
		}
		fmt.Printf("=== 执行完成: %s ===\n", demo.name)
	}

	fmt.Println("\n所有 Markdown 示例执行完成。")
}
