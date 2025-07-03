package main

import "fmt"

func main() {
	// Значение по умолчанию (0)
	var num0 int

	// Значение при инициализации
	var num1 int = 1

	// Пропуск типа
	var num2 = 20
	fmt.Println(num0, num1, num2)

	// Короткое объявление переменной
	num := 30
	// Только для новых переменных
	// no new variables on left side of :=
	// num := 31

	num += 1
	fmt.Println("+=", num)

	// ++num нет
	num++
	fmt.Println("++", num)

	// camelCase - принятый стиль
	userIndex := 10
	// under_score - не принято
	user_index := 10
	fmt.Println(userIndex, user_index)

	// Объявление нескольких переменных
	var weight, height int = 10, 20

	// Присваивание в существующие переменные
	weight, height = 11, 21

	// Короткое присваивание
	// Хотя-бы одна переменная должна быть новой!
	weight, age := 12, 22

	fmt.Println(weight, height, age)
}
