package utils

type evalFunc func(any, any) any

// 惰性计算斐波那契数列
func Fibonacci(num int) int {
	evenFunc := func(v1, v2 any) any {
		r1 := v1.(int)
		r2 := v2.(int)
		if r1 == 0 || r2 == 0 {
			return 1
		}
		return r1 + r2
	}

	ef := BuildLazyEvaluator(evenFunc)
	even := func() int { return ef().(int) }

	for range num - 1 {
		even()
	}
	return even()
}

func BuildLazyEvaluator(evalFunc evalFunc) func() any {
	retValChan := make(chan any)
	loopFunc := func() {
		var v1, v2 any = 0, 0
		for {
			vv := evalFunc(v1, v2)
			v1, v2 = v2, vv
			retValChan <- vv
		}
	}
	go loopFunc()
	return func() any { return <-retValChan }
}
