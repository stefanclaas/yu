package main

import (
	"bufio"
	"fmt"
	"mime"
	"os"
	"strings"
	"log"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("Usage: yu < infile > outfile")
		os.Exit(1)
	}

	body := ""
	subject := []string{}
	finalText := ""
	scanner := bufio.NewScanner(os.Stdin)
	str := []string{}
	for scanner.Scan() {
		str = strings.Split(scanner.Text(), ":")

		if len(str) > 1 {
			if str[0] == "To" {
				finalText += scanner.Text()
				finalText += "\n"
			} else if str[0] == "Subject" {
				subject = str
			} else {
				if len(subject) > 0 {
					body += scanner.Text() + "\n"
				} else {
					finalText += scanner.Text() + "\n"
				}
			}
		} else {
			body += scanner.Text() + "\n"
		}
	}
	
	finalText += subject[0] + ": "
	subEncoded := mime.BEncoding.Encode("UTF-8", strings.Join(subject[1:], ":"))
	if len(strings.Split(subEncoded, " ")) > 1 && strings.Join(subject[1:], ":") != subEncoded {
		for i, se := range strings.Split(subEncoded, " ") {
			if i > 0 {
				finalText += "  " + se
			} else {
				finalText += se
			}
			finalText += "\n"
		}
	} else {
		finalText += subEncoded + "\n"
	}

	finalText += "MIME-Version: 1.0\nContent-Type: text/plain; charset=UTF-8\nContent-Transfer-Encoding: 8bit\n"
	finalText += body
	fmt.Println(finalText)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

