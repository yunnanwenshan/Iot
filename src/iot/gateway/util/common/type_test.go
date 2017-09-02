package common

import "fmt"

func typeTest(t interface{})  {
	switch t := t.(type) {
	default:
		fmt.Printf("unexpected type %T", t)       // %T prints whatever type t has
	case bool:
		fmt.Printf("boolean %t\n", t)             // t has type bool
	case int:
		fmt.Printf("integer %d\n", t)             // t has type int
	case *bool:
		fmt.Printf("pointer to boolean %t\n", *t) // t has type *bool
	case *int:
		fmt.Printf("pointer to integer %d\n", *t) // t has type *int
	}
}

