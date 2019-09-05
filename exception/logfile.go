package exception

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	Info *log.Logger
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Info = log.New(os.Stdout, "Info:", (log.LstdFlags | log.Lshortfile))
}

/* A function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func LogFile(message string) {
	//filename := time.Now().Format("20060102150405")
	filename := time.Now().Format("2006010215")

	filepath := []string{"./log/", filename}

	f, err := os.OpenFile(strings.Join(filepath, ""), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	fmt.Fprintf(f, "%s %s", time.Now().Format("20060102150405"), message)
}
