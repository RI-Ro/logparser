package db

//CREATE table logparser_logs 
//(id SERIAL PRIMARY KEY, timestamp timestamptz default now(), remotetimestamp timestamptz, client text, 
//content text, facility int, hostname text, priority int, severity int, tag text);

//create index ON logparser_logs (client);
//create index ON logparser_logs (facility);
//create index ON logparser_logs (hostname);
//create index ON logparser_logs (priority);
//create index ON logparser_logs (severity);
//create index ON logparser_logs (remotetimestamp);
//create index ON logparser_logs (timestamp);

import (
    "database/sql"
    "fmt"
    "os"
    "log"
    "strings"

    _ "github.com/lib/pq"
)

type Message struct {
	client interface{}
	content interface{}
	facility interface{}
	hostname interface{}
	priority interface{}
	severity interface{}
	tag interface{}
	timestamp interface{}
}

func ParseLog(logParts map[string]interface{}) (message Message) {
	message.client = logParts["client"]
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

    client := fmt.Strintf("%v", msg.client)
	
	indexDelete := strings.Index(client, ":")

	// Если двоеточие найдено (indexDelete != -1)
	if indexDelete != -1 {
		// Удаляем все после двоеточия (включая двоеточие)
		// Оставляем только IP-адрес
		// Приходит 127.0.0.1.34583 оставляем только 127.0.0.1
		client = client[:indexDelete]
	} else {
		// Двоеточия не было
		client = client
	}

    err = db.Ping()
    if err != nil {
//        panic(err)
		fmt.Println("В случае ошибки подключения к БД производим запись в файл.\n")
		SaveToFile(fileIfNotWorkDB, msg, severity, client)
    } else {
    	fmt.Println("Successfully connected!")
    	sqlStatement := `
		INSERT INTO logparser_logs (timestamp, remotetimestamp, client, content, facility, hostname, priority, severity, tag)
		VALUES (NOW(), $1, $2, $3, $4, $5, $6, $7, $8)`



		if (severity >= msg.severity.(int)) {
			_, err = db.Exec(sqlStatement, msg.timestamp, client, 
									msg.content, msg.facility, msg.hostname, 
									msg.priority, msg.severity, msg.tag)
			if err != nil {
				fmt.Println("Ошибка записи в БД!\n")
				SaveToFile(fileIfNotWorkDB, msg, severity, client)
			addClient := `INSERT INTO clients (client) VALUES ($1) ON CONFLICT (client) DO NOTHING`
			_, _ = db.Exec(addClient, client)
			} else {
				fmt.Println("Запись в БД прошла успешно!\n")	
			}
		}
	}
}

func SaveToFile(fileIfNotWorkDB string, msg Message, severity int, client string){
	
	file, err := os.OpenFile(fileIfNotWorkDB, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // Важно закрыть файл после использования

	var textToSave string;

	textToSave = fmt.Sprintf("%v client:%v content:%v faciliti:%v hostname:%v priority:%v severity:%v tag:%v\n", 
								msg.timestamp, client, 
								msg.content, msg.facility, msg.hostname, 
								msg.priority, msg.severity, msg.tag)

	// Запись в файл
	if (severity >= msg.severity.(int)) {
		_, err = file.WriteString(textToSave)
		if err != nil {
	//		fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}
}