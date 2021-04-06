package a

func f() {
	var n int
	func() {
		n = 10 // want "NG"
	}()
	println(n)
}
