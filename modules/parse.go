package modules

import (
	"fmt"
)

func Parse(f interface{}) {
	m := f.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:", vv)
			for i, u := range vv {
				if _, ok := u.(map[string]interface{}); ok {
					Parse(u)
				} else {
					fmt.Println(i, u)
				}
			}
		case interface{}:
			fmt.Println(k, "is interface", vv)
			Parse(vv)
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}
