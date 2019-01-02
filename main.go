package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type keys []string

var paramKeys keys

func (i *keys) String() string {
	return fmt.Sprint(*i)
}


func (i *keys) Set(value string) error {
	*i = append(*i, value)
	return nil
}
func main() {
	var filename = flag.String("f", "log.txt", "log file")
	var timefilter = flag.String("t", "", "time")
	var isListKeys = flag.Bool("lk", false, "list keys")
	flag.Var(&paramKeys, "key",  "key to print")

	flag.Parse()
	var scanner *bufio.Scanner
	if (filename != nil) {
		f, err := os.Open(*filename)
		if err != nil {
			log.Fatal(err)
			return
		}
		scanner = bufio.NewScanner(f)

	} else {
		return
	}
	lc := 0

	for scanner.Scan() {
		lc++
		r := strings.Split(scanner.Text(), " ")
		if 14 != len(r) {
			continue
		}

		if strings.HasPrefix(r[4], *timefilter) == true {
			fmt.Println(scanner.Text())
			params := make(map[string]string)

			path := r[10][8:(len(r[10]) - 1)]
			ps := strings.Split(path, "&")
			if len(ps) < 1 {
				continue
			}
			for _, v := range ps {
				kv := strings.Split(v, "=")
				if *isListKeys {
					fmt.Print(kv[0]+",")
				}
				params[kv[0]] = kv[1]
			}
			fmt.Println()

			params["time"] = r[4]
			params["hits"] = strings.Split(r[11], "=")[1]
			params["qtime"] = strings.Split(r[12], "=")[1]
			if 0 < len(paramKeys) {
				for _, v := range paramKeys {
					fmt.Printf("%v=%v\n", v, params[v])
				}
			}
		}

	}

}
