package main

import (
	"fmt"
	"reflect"
)

func convert[S, T any](tmp *S) *T {
	var target T
	unpack(reflect.ValueOf(&target), reflect.ValueOf(tmp))
	return &target
}

func unpack(to reflect.Value, from reflect.Value) error {
	var emptyValue reflect.Value
	if to == emptyValue || from == emptyValue {
		return nil
	}

	for to.Kind() == reflect.Pointer {
		to = to.Elem()
	}
	for from.Kind() == reflect.Pointer {
		from = from.Elem()
	}

	if to.Kind() != from.Kind() {
		return fmt.Errorf("kinds do not match")
	}
	if to.Kind() == reflect.Struct {
		for i := 0; i < to.NumField(); i++ {
			err := unpack(
				to.Field(i),
				from.FieldByName(to.Type().Field(i).Name),
			)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if to.Kind() == reflect.Slice {
		for i := 0; i < from.Len(); i++ {
			tmp := reflect.New(to.Type().Elem())
			err := unpack(tmp, from.Index(i))
			if err != nil {
				return err
			}
			to.Set(reflect.Append(to, tmp.Elem()))
		}
		return nil
	}

	to.Set(from)

	return nil
}
