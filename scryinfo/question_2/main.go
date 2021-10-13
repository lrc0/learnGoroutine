package main

import (
	"fmt"
	"os"
	"strconv"
)

//数字倒排，例如：1234倒排成4321
func dealString(s int) (int, error) {
	numS := fmt.Sprintf("%d", s)
	numL := len(numS)

	var result string
	for i := numL - 1; i >= 0; i-- {
		result += string(numS[i])
	}

	rs, err := strconv.Atoi(result)
	if err != nil {
		return -1, err
	}
	return rs, nil
}

func main() {
	numStr := os.Args[1]
	l := len(numStr)
	if l >= 19 {
		fmt.Println("数字位数不能超过19位")
		return
	}
	num, err := strconv.Atoi(numStr)
	if err != nil {
		fmt.Println("请输入数字, err: ", err)
		return
	}
	if num < 0 {
		fmt.Println("只支持正数")
		return
	}

	rs, err := dealString(num)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println(rs)
}
