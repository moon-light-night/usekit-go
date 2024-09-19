package main

import "fmt"

func main() {
	// defer позволяет выполнить функцию в конце выполнения программы
	// самый первый defer выполнится самым последним
	defer finish()
	defer fmt.Println("Program has been almost finished")
	println(divide(1, 2))
	arr := [...]int{1, 2, 3, 4, 5}
	for i := 0; i <= len(arr); i++ {
		fmt.Println(i * i)
	}
	hello()
	sumOfPair := sum(2, 2)
	println(sumOfPair)
	add(1, 2, 3, 4, 5)
	age, name := nameAndAge(10, 15, "John", "Doe")
	println(age, name)
	println("-----------------")
	var fSum func(int, int) int = sum
	// или так fSum := sum
	println(fSum(10, 10))
	println("----------")
	println(action(2, 5, sum))
	// --------------

	// анонимная функция
	f := func(x, y int) int {
		return x + y
	}
	println(f(20, 20))
	println("----------")

	// ----------
	println(factorial(3))
	println("----------")

	// ----------
	println(fibbonachi(10))
	println("----------")

	// ----------
	// срезы
	// var users []string
	users2 := []string{"John", "Sam", "Peter"}
	var users3 = make([]string, 2)
	users3[0] = "John"
	users3[1] = "Sam"
	// добавление элемента в срез, возвращает новый срез
	users3 = append(users3, "Peter")
	println(users2[1], users3[2])
	println("----------")

	// ----------
	// оператор среза
	initUsersIdx := 2
	initUsers := [5]string{"John", "Sam", "Peter", "Mike", "Donald"}
	// операции удаления, добавления элемента можно производить только со слайсами, не с массивами
	initUsersSlice := []string{"John", "Sam", "Peter", "Mike", "Donald"}
	fmt.Println(initUsers[1:2])
	fmt.Println(initUsers[:2])
	fmt.Println(initUsers[2:])
	initUsersSlice = append(initUsersSlice[:initUsersIdx], initUsersSlice[initUsersIdx+1:]...)
	println("---")
	fmt.Println(initUsersSlice)
	println("----------")

	// ---------- объекты/карты/отображения/мапы
	peopleAges := map[string]int {
		"John": 30,
		"Peter": 25,
		"Donald": 31,
	}

	fmt.Println(peopleAges)
	fmt.Println(peopleAges["John"])
	for key, value := range peopleAges {
		fmt.Println(key, value)
	}
	// Функция make представляет альтернативный вариант создания отображения. Она создает пустую хеш-таблицу:
	people := make(map[string]int)
	people["Kate"] = 28
	people["Sam"] = 20
	delete(people, "Sam")
	fmt.Println(people)
}

func hello() {
	fmt.Println("hello, go")
}

func sum(x, y int) (z int) {
	z = x + y
	return z
}

func action(x, y int, operation func(int, int) int) int {
	return operation(x, y)
}

func add(numbers ...int) {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	fmt.Println(sum)
}

func nameAndAge(x, y int, fName, sName string) (age int, name string) {
	age = x + y
	name = fName + " " + sName
	return age, name
}

// функция square возвращает анонимную функцию, которая имеет доступ к среде выполнения анонимной функции и получает доступ к переменной x
// и манипулирует переменной x
func square() func() int {
	x := 2
	return func() int {
		x++
		return x * x
	}
}

func factorial(n uint) uint {
	if n == 0 {
		return 1
	}

	return n * factorial(n-1)
}

// рекурсивная функция для чисел фиббоначи
func fibbonachi(n uint) uint {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	return fibbonachi(n-1) + fibbonachi(n-2)
}

func finish() {
	println("Program is finished")
}

func divide(a, b float64) float64 {
	if b == 0 {
		panic("divide by zero!")
	}

	return a / b
}