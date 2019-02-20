package main

import (
	"regexp"
	"strconv"
	"strings"
)

type CafeMenu struct {
	Menu []CafeMenuDia `json:"menu"`
	Mes  string        `json:"mes"`
	From int           `json:"desde"`
	To   int           `json:"hasta"`
}

type CafeMenuDia struct {
	PrimerPlato  []string `json:"primer"`
	SegundoPlato []string `json:"segundo"`
}

var monthre = regexp.MustCompile(`MES: ([A-Z\-]+)DEL `)
var fromre = regexp.MustCompile(`DEL D[ÍI]A: (\d{1,2})`)
var tore = regexp.MustCompile(`AL D[ÍI]A: (\d{1,2})`)

func getMetadata(data string) (mes string, desde int, hasta int) {
	month := monthre.FindStringSubmatch(data)
	from := fromre.FindStringSubmatch(data)
	to := tore.FindStringSubmatch(data)

	if len(month) > 0 && len(from) > 0 && len(to) > 0 {
		mes = month[1]
		desde, _ = strconv.Atoi(from[1])
		hasta, _ = strconv.Atoi(to[1])
		return
	}

	return "", 0, 0
}

var platore = regexp.MustCompile(`([A-Z\-ÁÉÍÓÚÑ]+\s{1,2})*([A-Z\-ÁÉÍÓÚÑ]{3,})`)

func parsePdfTable(tables []pdfTableResponse) *CafeMenu {
	segundoPlato := 3

	menu := &CafeMenu{
		Menu: make([]CafeMenuDia, 5),
	}

	for i := 0; i < 5; i++ {
		menu.Menu[i].PrimerPlato = []string{}
		menu.Menu[i].SegundoPlato = []string{}
	}

	menu.Mes, menu.From, menu.To = getMetadata(strings.Join(tables[0].Data[0], ""))

	pdf := tables[1].Data
	platos := [5][]string{}
	indices := [5]int{}

	for row := range pdf {
		for col := range pdf[row] {
			text := pdf[row][col]
			text = strings.Replace(text, "\n", "", -1)
			text = strings.ToUpper(text)

			nombres := platore.FindAllString(text, -1)

			for day, nombre := range nombres {
				day = day + col

				day -= len(nombres) - 1

				if nombre == "KCAL" || nombre == "PESCADO" {
					// pop
					if len(platos[day]) > 0 {
						if indices[day] >= segundoPlato {
							menu.Menu[day].SegundoPlato = append(menu.Menu[day].SegundoPlato, strings.Join(platos[day], " "))
						} else {
							menu.Menu[day].PrimerPlato = append(menu.Menu[day].PrimerPlato, strings.Join(platos[day], " "))
						}
						platos[day] = []string{}
						indices[day]++
					}
				} else {
					platos[day] = append(platos[day], nombre)
				}
			}
		}
	}

	return menu
}
