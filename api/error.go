package api

import "fmt"

var (
	MessageErr map[int]string = map[int]string{
		0: "Успешно выполнено",

		10: "Код отправлена на email",
		11: "Не верный код авторизации",
		12: "Не корректный email",
		13: "Ошибка получения access-token",
		14: "Access-token не действителен",
		15: "Access-token истек",
		16: "Пользователь не найден",
		17: "Ошибка авторизации по коду",

		20: "Ошибка получения refresh-token",
		21: "Сессия не найдена ",

		401: "Пользователь не авторизован",
		404: "Страница не найдена",

		500: "Ошибка на сервер попробуйте позже",

		700: "",
	}
)

func GetMessageErr(idError int) string {

	if _, ok := MessageErr[idError]; !ok {
		return fmt.Sprintf("Error: num %i is not found in dict", idError)
	}

	return MessageErr[idError]
}
