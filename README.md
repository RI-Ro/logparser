Простое приложение на GOLANG для приема логов от syslog и записи в базу данных (postgresql) или в файл, если БД не доступна.

Для запуска:
go run main.go -configPath="./config/config.yaml" -fileToSave="/tmp/log2"

go buld -o main main.go

./main -configPath="./config/config.yaml" -fileToSave="/tmp/log2"
