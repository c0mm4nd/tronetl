package main

// chk panics when error exists
func chk(err error) {
	if err != nil {
		panic(err)
	}
}
