package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

type record struct {
	Date string `json:"date"`
}

func getLastDate() time.Time {
	content, err := ioutil.ReadFile("./download/coronaparsing/docs/data/12/data.json")
	if err != nil {
		log.Fatal(err)
	}
	var data []record
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(len(data))
	//fmt.Println(data[len(data)-1].Date)
	last, err := time.Parse("060102", data[len(data)-1].Date)
	fmt.Println("getLastDate", last)
	return last
}

func execute(name string) {
	fmt.Println("execute", name)
	cmd := exec.Command(name)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("He")
	/*{
		today := time.Now()
		tomorrow := today.AddDate(0, 0, 1)
		ss := tomorrow.Format("2006-01-02 15:04")
		s := []byte(ss)
		s[11] = '0'
		s[12] = '9'
		s[14] = '0'
		s[15] = '0'
		ss = string(s)
		tommmorowNine, _ := time.Parse("2006-01-02 15:04", ss)
		fmt.Println(tommmorowNine)
		fmt.Println(tommmorowNine.Sub(today))
	}
	return*/
	for {
		// Проверка текущей даты
		today := time.Now()
		fmt.Println(today)
		//fmt.Println(today.Date())
		//fmt.Println(today.Clock())
		// Проверка даты последней записи
		last := getLastDate()
		fmt.Println(last)

		diff := today.Sub(last)
		next := last.Add(time.Hour * 24)
		fmt.Println("diff", diff)
		fmt.Println("next", next)

		// Если нужны новые данные
		//if diff > time.Hour*24*2 {
		if !next.Before(today) {
			// Если нужны данные только за сегодня
			// Запускаем парсинг
			for {
				fmt.Println("today")
				execute(`.\service\win\download.bat`)
				// Если данных нет - повтор через час
				if last != getLastDate() {
					break
				}
				fmt.Println("sleep 1 hour")
				time.Sleep(time.Hour)
			}
			// Если есть - ждем до завтра до 9 часов
		} else {
			// Если данные нужны не только за сегодня
			fmt.Println("a few days")
			// Повторяем процедуру пока не останутся данные только за сегодня
			for next.Before(today) {
				// Запускаем парсинг
				execute(`.\service\win\download.bat`)
				next2 := getLastDate().Add(time.Hour * 24)
				if next2 == next {
					break
				}
				next = next2
				// ждем 30 секунд
				fmt.Println("sleep 30 sec")
				time.Sleep(30 * time.Second)
			}
		}
		// Если новые данные появились делаем git commit
		/*last2 := getLastDate()
		if last2 != last {
			execute(`.\service\win\git.bat`)
		}*/

		// ждем до завтра до 9 часов
		tomorrow := today.AddDate(0, 0, 1)
		ss := tomorrow.Format("2006-01-02 15:04")
		s := []byte(ss)
		s[11] = '0'
		s[12] = '9'
		s[14] = '0'
		s[15] = '0'
		ss = string(s)
		tommmorowNine, err := time.Parse("2006-01-02 15:04", ss)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("sleep", tommmorowNine)
		time.Sleep(tommmorowNine.Sub(today))
	}
}
