package main

import (
	"io/ioutil"
	"log"
	"reflect"
	"regexp"
)

// Config
/*
 * Структура нашего будущего PHP массива
 */
type Config struct {
	XServerCorrelationID          string `php:"X-Server-CorrelationId"`
	EntitiesVersion               string `php:"entities-version"`
	ClientCallId                  string `php:"client-call-id"`
	LicenseHash                   string `php:"license-hash"`
	RestrictionsStateHash         string `php:"restrictions-state-hash"`
	ObtainedLicenseConnectionsIds string `php:"obtained-license-connections-ids"`
	XServerPasswordHash           string `php:"X-Server-PasswordHash"`
	XServerLoginName              string `php:"X-Server-LoginName"`
	XServerBackVersion            string `php:"X-Server-BackVersion"`
	XServerAuthType               string `php:"X-Server-AuthType"`
	XServerServerEdition          string `php:"X-Server-ServerEdition"`
	Host                          string `php:"Host"`   // Будет работать и без тега
	Expect                        string `php:"Expect"` // Будет работать и без тега
}

func main() {
	// Объявляем переменные
	var ConfigData Config               // Переменная под конфигурацию
	var tagOrFieldNameIfTagEmpty string // Переменная под имя тега или структуры

	// Читаем данные из входящего файла
	inputData, err := ioutil.ReadFile("input.txt")

	// Если считать файл не удалось, выводим ошибку и завершаем работу
	if err != nil {
		log.Fatal("Не удалось прочитать файл с данными о curl запросе")
	}

	// Получаем информацию о ConfigData
	ConfigDataReflectValues := reflect.ValueOf(ConfigData)

	// Итерируем поля из содержимого ConfigData
	for i := 0; i < ConfigDataReflectValues.NumField(); i++ {

		// Получаем имя поля из тега
		tagOrFieldNameIfTagEmpty = ConfigDataReflectValues.Type().Field(i).Tag.Get("php")
		if tagOrFieldNameIfTagEmpty == "" { // Если оно пустое
			// Подменяем имя именем поля
			tagOrFieldNameIfTagEmpty = ConfigDataReflectValues.Type().Field(i).Name
		}

		// Прописываем регулярное выражения для получения данных из заголовков curl
		r := regexp.MustCompile(`'` + tagOrFieldNameIfTagEmpty + `: (.*?)'`)
		// Получаем совпадения
		matches := r.FindAllStringSubmatch(string(inputData), -1)
		if len(matches) == 0 { // Если совпадений не нашлось
			// Прописываем регулярное выражения для получения данных из xml тегов внутри body data запроса
			r = regexp.MustCompile(`<` + tagOrFieldNameIfTagEmpty + `>(.*?)<\/` + tagOrFieldNameIfTagEmpty + `>`)
			// Получаем совпадения
			matches = r.FindAllStringSubmatch(string(inputData), -1)
		}
		// Перебираем полученные совпадения (потенциально их может быть больше или меньше 1)
		for _, match := range matches {
			// Заполняем данные из совпадений в соответствующие поля нашей переменной
			reflect.ValueOf(&ConfigData).Elem().Field(i).SetString(match[1])
		}
	}

	// Тестируем вывод
	//fmt.Printf("%+v\n", ConfigData)
	//fmt.Printf("%+v\n", ConfigData.MarshalPHP())

	// Записываем в результирующий файл
	err = ioutil.WriteFile("output.php", []byte(ConfigData.MarshalPHPArray()), 'w')

	// Если записать файл не удалось, выводим ошибку и завершаем работу
	if err != nil {
		log.Fatal("Не удалось записать результирующий файл с конфигурацией в виде PHP массива")
	}
}

// MarshalPHPArray
/*
 * Функция преобразования данных из структуры Config в PHP массив
 * Потенциально можно было обойтись без нее, использовав склейку в строку при заполнении значений, но это было бы чревато ошибками при повторных совпадениях регулярных выражений
 */
func (c *Config) MarshalPHPArray() string {

	// Получаем информацию о конфигурации "c"
	ConfigDataReflectValues := reflect.ValueOf(*c)

	// Объявляем переменные
	var outputString string             // Переменная под строковый вывод
	var tagOrFieldNameIfTagEmpty string // Переменная под имя тега или структуры

	// Итерируем поля из содержимого ConfigData
	for i := 0; i < ConfigDataReflectValues.NumField(); i++ {

		// Получаем имя поля из тега
		tagOrFieldNameIfTagEmpty = ConfigDataReflectValues.Type().Field(i).Tag.Get("php")
		if tagOrFieldNameIfTagEmpty == "" { // Если оно пустое
			// Подменяем имя именем поля
			tagOrFieldNameIfTagEmpty = ConfigDataReflectValues.Type().Field(i).Name
		}

		// Подклеиваем получившиеся данные в результирующую строку
		outputString += "	'" + tagOrFieldNameIfTagEmpty + "' => \"" + reflect.ValueOf(*c).Field(i).String() + "\",\n"
	}

	// Склеиваем получившиеся строки с основным шаблоном
	outputString = "<?php\nreturn [\n" + outputString + "];"

	// Возвращаем результат
	return outputString
}
