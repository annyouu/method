package main

// import "fmt"

func f() {
	defer func() {}() // OK

	for range 10 {
		defer func() {}() // want "NG"
		if false {
			// 難しい(CFGを使わないと無理)
			defer func() {}() // want "NG"
		}
	}
}