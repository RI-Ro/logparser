package main

import (
    "fmt"
    "flag"
    "log"
    "gopkg.in/mcuadros/go-syslog.v2"

    localConfig "logparser/config"
    localDB "logparser/db"
    )


func main() {
    var configPath string;
    var fileIfNotWorkDB string;
    flag.StringVar(&configPath, "configPath", "/etc/logparser/config.yaml", "Путь по которому размещен файл конфигурации")
    flag.StringVar(&fileIfNotWorkDB, "fileToSave", "/tmp/logparser.log", "Файл в который происходит запись логов, если DB не доступна")
    flag.Parse()

    config := localConfig.CreateConfig(fmt.Sprintf(configPath))
    psqlConnect := localConfig.CreateConnectString(config)

    // Create a channel to receive parsed syslog messages
    logChannel := make(syslog.LogPartsChannel)

    // Create a handler that sends messages to the channel
    handler := syslog.NewChannelHandler(logChannel)

    // Create a new syslog server
    server := syslog.NewServer()

    // Set the desired syslog format (e.g., RFC5424)
//  server.SetFormat(syslog.RFC5424)
    server.SetFormat(syslog.RFC3164)

    // Set the handler for the server
    server.SetHandler(handler)

    // Listen for UDP syslog messages on all interfaces and port 514
    err := server.ListenUDP(fmt.Sprintf("%v:%v", config.Server.Ipaddress, config.Server.Port))
    if err != nil {
        log.Fatalf("Error listening on UDP: %v", err)
    }

    // Start the syslog server
    server.Boot()

    // Goroutine to process incoming log messages from the channel
    go func(channel syslog.LogPartsChannel) {
        for logParts := range channel {
            // Сохраняем результаты syslog в БД или файл, если БД не доступна
            localDB.SaveToDB(psqlConnect, fileIfNotWorkDB ,logParts)
            // In a real application, you would process or store these logs
            }
    }(logChannel)

    // Wait for the server to gracefully shut down (e.g., on signal)
    // This will keep the main goroutine alive and the server running
    server.Wait()
}


