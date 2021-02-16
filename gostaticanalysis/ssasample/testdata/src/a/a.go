package a

func f(n, m int) { // want "unused"
	n = 20
	println(n, m)
}
