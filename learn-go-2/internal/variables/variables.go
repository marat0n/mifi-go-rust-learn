package variables

func createVariables() int64 {
	// 1. Declared & Initiated
	var declaredAndInitiatedVar int64 = 1

	// 2.1. Declaration
	var myVariable int

	// 2.2. Initiatian | Initialization
	myVariable = 2

	// 3. Short declaration and initiation
	myNewVariable := 3

	return declaredAndInitiatedVar +
		int64(myVariable) +
		int64(myNewVariable)
}
