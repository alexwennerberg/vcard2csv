package main

import (
	"encoding/csv"
	"github.com/emersion/go-vcard"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// CSV
// ios vcard
// standard vcard

// make sure none of these end in a digit
var asIs = []string{"VERSION", "FN", "NICKNAME", "BDAY", "X-GTALK", "X-PHOENETIC-FIRST-NAME", "X-PHOENETIC-LAST-NAME", "TITLE", "NOTE"}

var repeatedWithType = []string{"TEL", "X-SOCIALPROFILE", "X-ABDATE", "EMAIL", "IMPP", "URL", "ADR"}

func cardToFlatDict(vc vcard.Card) map[string]string {
	out := map[string]string{}
	for _, field := range asIs {
		out[field] = vc.Value(field)
	}
	for _, field := range repeatedWithType {
		fieldObjs := vc[field]
		for _, value := range fieldObjs {
			var localField string
			if len(value.Params.Types()) > 0 {
				localField = strings.ToUpper(field + "_" + strings.Join(value.Params.Types(), "_"))
			} else {
				localField = field
			}
			for {
				if out[localField] != "" {
					if !unicode.IsDigit(rune(localField[len(localField)-1])) {
						localField += "2" // no item should be repeated more than 10 times TBD
					} else {
						count, err := strconv.Atoi(string(localField[len(localField)-1]))
						if err != nil {
							log.Fatal(err)
						}
						localField = localField[:len(localField)-1] + strconv.Itoa(count+1)
					}
				} else {
					break

				}
			}
			out[localField] = value.Value
		}
	}
	return out
}

func main() {
	dec := vcard.NewDecoder(os.Stdin)
	output := csv.NewWriter(os.Stdout)
	headers := map[string]bool{}
	var rows []map[string]string
	for {
		card, err := dec.Decode()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		row := cardToFlatDict(card)
		for key := range row {
			if !headers[key] && row[key] != "" {
				headers[key] = true
			}
		}
		rows = append(rows, row)
	}

	keys := make([]string, len(headers))

	i := 0
	for k := range headers {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	err := output.Write(keys)
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range rows {
		record := make([]string, len(keys))
		for i, k := range keys {
			record[i] = row[k]
		}
		err := output.Write(record)
		if err != nil {
			log.Fatal(err)
		}
	}
	output.Flush()
}
