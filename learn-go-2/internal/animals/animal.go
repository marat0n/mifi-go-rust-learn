package animals

// 1. Value Pointer | 8 byte
// 2. Functions table Pointer | 8 byte
type Animal interface{
	GetName() string
	SetBreed(string)
	SetName(string)
	SetAge(int)
	SetHeight(int)
}
