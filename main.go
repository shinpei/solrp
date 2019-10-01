package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
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
	var passthrough = flag.Bool("passthrough", false, "pass this to pass through every parameter")
	var tabmode = flag.Bool("tab", false, "tab mode")

	var qtimeFilt = flag.Int("qf", 0, "qtime filter")
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

	timeFiltStarted := false
	timeFiltDone := false
	var joiner string
	if (*tabmode) {
		joiner = "\t"
	} else {
		joiner = ","
	}
	for scanner.Scan() {

		if timeFiltDone {
			break
		}
		text := scanner.Text()
		lc++
		r := strings.Split(text, " ")
		if 14 != len(r) {
			continue
		}

		if (r[9] != "path=/select") {
			continue
		}
		if strings.HasPrefix(r[4], *timefilter) == true {
			timeFiltStarted = true
			params := make(map[string]interface{})

			path := r[10][8:(len(r[10]) - 1)]
			ps := strings.Split(path, "&")
			if len(ps) < 1 {
				continue
			}
			if *passthrough {

				fmt.Println(text)
				continue
			}

			for _, v := range ps {
				kv := strings.Split(v, "=")
				if *isListKeys {
					fmt.Print(kv[0]+ joiner)
				}
				switch kv[0] {
				case "fq":
					_, ok := params["fq"]
					if !ok {
						params["fq"] = new([]string)
					}
					arr := params["fq"].(*[]string)
					*arr = append(*arr, kv[1])
					params["fq"] = arr
				default:
					params[kv[0]] = kv[1]
				}

			}
			if *isListKeys {
				fmt.Println("")
			}
			params["time"] = r[4]
			params["hits"] = strings.Split(r[11], "=")[1]
			params["status"] = strings.Split(r[12], "=")[1]
			params["qtime"] = strings.Split(r[13], "=")[1]
			if 0 < len(paramKeys) {
				arr := new([]string)
				qtime, _:= params["qtime"].(string)
				if i, _ := strconv.Atoi(qtime); i > *qtimeFilt {
					for _, v := range paramKeys {
						s := fmt.Sprintf("%v=%v", v, params[v])
						*arr = append(*arr, s)
					}
					fmt.Println(strings.Join(*arr, joiner))
				}
			}
		} else {
			if timeFiltStarted {
				timeFiltDone = true
			}
		}

	}

}
