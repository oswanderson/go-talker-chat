// Custom Logger ...

package ctmlogger

import (
	"log"
	"os"
	"sync"
)

type customLogger struct {
	fileName string
	*log.Logger
}

var ctmLogger *customLogger
var once sync.Once

func GetInstance() *customLogger {
	once.Do(func() {
		instance, _ := createLogger("custom_logger.log")
		ctmLogger = instance
	})
	return ctmLogger
}

func createLogger(fileName string) (*customLogger, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)

	if err != nil {
		return nil, err
	}
	return &customLogger{
			fileName: fileName,
			Logger:   log.New(file, "LOGGER >>>> ", log.Ldate|log.Ltime|log.Lshortfile)},
		nil
}
