package models

// 1. Value Pointer | 8 byte
// 2. Functions table Pointer | 8 byte
type Animal interface{
	GetName() string
	GetBreed() string
	GetAge() int
	GetHeight() int

	GetImageUrl() string

	SetBreed(string)
	SetName(string)
	SetAge(int)
	SetHeight(int)
}
