package db

import (
    "database/sql"
    "fmt"
    "os"
    "log"
    "strings"

    _ "github.com/lib/pq"
)

type Message struct {
	ip interface{}
	content interface{}
	facility interface{}
	hostname interface{}
	priority interface{}
	severity interface{}
	tag interface{}
	timestamp interface{}
}

func ParseLog(logParts map[string]interface{}) (message Message) {
	message.ip = logParts["client"]
	message.content = logParts["content"]
	message.facility = logParts["facility"]
	message.hostname = logParts["hostname"]
	message.priority = logParts["priority"]
	message.severity = logParts["severity"]
	message.tag = logParts["tag"]
	message.timestamp = logParts["timestamp"]
	return message
}

func SaveToDB(psqlConnect string, fileIfNotWorkDB string, logParts map[string]interface{}, severity int) {
    db, err := sql.Open("postgres", psqlConnect)
    if err != nil {
        panic(err)
    }
    defer db.Close()
  
    msg := ParseLog(logParts)

    ip := fmt.Sprintf("%v", msg.ip)
	
	indexDelete := strings.Index(ip, ":")

	// Если двоеточие найдено (indexDelete != -1)
	if indexDelete != -1 {
		// Удаляем все после двоеточия (включая двоеточие)
		// Оставляем только IP-адрес
		// Приходит 127.0.0.1.34583 оставляем только 127.0.0.1
		ip = ip[:indexDelete]
	} else {
		// Двоеточия не было
		ip = ip
	}

    err = db.Ping()
    if err != nil {
//        panic(err)
		fmt.Println("В случае ошибки подключения к БД производим запись в файл.\n")
		SaveToFile(fileIfNotWorkDB, msg, severity, ip)
    } else {
    	fmt.Println("Successfully connected!")
    	sqlStatement := `
		INSERT INTO logparser_logs (timestamp, remotetimestamp, hostname_id, facility_id, severity_id, content, tag)
		VALUES (NOW(), $1, (SELECT id FROM hostname WHERE hostname LIKE $2 AND ip LIKE $3 LIMIT 1), $4, $5, $6, $7) `



		if (severity >= msg.severity.(int)) {

			addHostname := `INSERT INTO hostname (ip, hostname) VALUES ($1, $2) ON CONFLICT (ip, hostname) DO NOTHING`
			_, err = db.Exec(addHostname, ip, msg.hostname)
			if err != nil {
				fmt.Printf("Ошибка записи в БД! %v\n", err)
				}

			_, err = db.Exec(sqlStatement, msg.timestamp, msg.hostname, ip, 
									msg.facility, msg.severity, 
									msg.content, msg.tag)
			if err != nil {
				fmt.Println("Ошибка записи в БД!\n")
				SaveToFile(fileIfNotWorkDB, msg, severity, ip)
			} else {
				fmt.Println("Запись в БД прошла успешно!\n")	
			}
		}
	}
}

func SaveToFile(fileIfNotWorkDB string, msg Message, severity int, ip string){
	
	file, err := os.OpenFile(fileIfNotWorkDB, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // Важно закрыть файл после использования

	var textToSave string;

	textToSave = fmt.Sprintf("%v ip:%v content:%v faciliti:%v hostname:%v severity:%v tag:%v\n", 
								msg.timestamp, ip, 
								msg.content, msg.facility, msg.hostname, 
								msg.severity, msg.tag)

	// Запись в файл
	if (severity >= msg.severity.(int)) {
		_, err = file.WriteString(textToSave)
		if err != nil {
	//		fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}
}