package main

import "fmt"

func panicf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

func _panic(err error) {
	if err != nil {
		panic(err)
	}
}

func _panicf(err error, format string, args ...interface{}) {
	if err != nil {
		panic(fmt.Errorf(format+"\n\terr: %v", append(args, err)...))
	}
}
