package test

import "time"

type NamedInt int
type NamedString string

type Test1 struct {
	V1  int64
	V2  string
	V3  NamedInt
	V4  string
	V5  string
	V6  *int
	V7  *NamedInt
	V8  **NamedInt
	V9  time.Time
	v10 int32
	V11 []int
	V12 []*string
}

type Test2 struct {
	V1  int32
	V2  NamedString
	V3  int
	V4  NamedString
	V5  *string
	V6  int
	V7  *int
	V8  *NamedString
	v9  int32
	v10 time.Time
	V11 []int8
	V12 []*NamedString
}

type Test3 struct {
	//V1 time.Time
	V2 int32
}

type Test4 struct {
	//V1 int32
	V2 time.Time
}
