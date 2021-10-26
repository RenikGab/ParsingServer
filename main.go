package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type record struct {
	Date string `json:"date"`
}

var Log *log.Logger

func getLastDate() time.Time {
	content, err := ioutil.ReadFile("download/coronaParsing/docs/data/12/data.json")

	if err != nil {
		log.Fatal(err)
	}
	var data []record
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Fatal(err)
	}
	//Log.Println(len(data))
	//Log.Println(data[len(data)-1].Date)
	last, err := time.Parse("060102", data[len(data)-1].Date)
	Log.Println("getLastDate", last)
	return last
}

func execute(name string) {
	Log.Println("execute", name)
	cmd := exec.Command(name)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	logFile, _ := os.OpenFile("logFile", os.O_RDWR|os.O_CREATE, 0666)
	multi := io.MultiWriter(logFile, os.Stdout)
	Log = log.New(multi, "INFO: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	Log.Println("Start logging...")
	Log.Println(runtime.GOOS)
	downloadScript := `service/linux/download.sh`
	gitScript := `service/linux/git.sh`
	if runtime.GOOS == "windows" {
		downloadScript = `.\service\win\download.bat`
		gitScript = `.\service\win\git.bat`
	}

	defer logFile.Close()

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
		Log.Println(tommmorowNine)
		Log.Println(tommmorowNine.Sub(today))
	}
	return*/
	for {
		// Проверка текущей даты
		today := time.Now()
		Log.Println(today)
		//Log.Println(today.Date())
		//Log.Println(today.Clock())
		// Проверка даты последней записи
		last := getLastDate()
		Log.Println(last)

		diff := today.Sub(last)
		next := last.Add(time.Hour * 24)
		Log.Println("diff", diff)
		Log.Println("next", next)

		// Если нужны новые данные
		//if diff > time.Hour*24*2 {
		if !next.Before(today) {
			// Если нужны данные только за сегодня
			// Запускаем парсинг
			for {
				Log.Println("today")
				execute(downloadScript)
				// Если данных нет - повтор через час
				if last != getLastDate() {
					break
				}
				Log.Println("sleep 1 hour")
				time.Sleep(time.Hour)
			}
			// Если есть - ждем до завтра до 9 часов
		} else {
			// Если данные нужны не только за сегодня
			Log.Println("a few days")
			// Повторяем процедуру пока не останутся данные только за сегодня
			for next.Before(today) {
				// Запускаем парсинг
				execute(downloadScript)
				next2 := getLastDate().Add(time.Hour * 24)
				if next2 == next {
					break
				}
				next = next2
				// ждем 30 секунд
				Log.Println("sleep 30 sec")
				time.Sleep(30 * time.Second)
			}
		}
		// Если новые данные появились делаем git commit
		last2 := getLastDate()
		if last2 != last {
			execute(gitScript)
		}

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
		Log.Println("sleep", tommmorowNine)
		time.Sleep(tommmorowNine.Sub(today))
	}
}
