package main

import "github.com/CoreDumped-ETSISI/uni-services/emt/api"

type UniversityStops struct {
	SentidoSierra []api.Bus `json:"sierra"`
	SentidoConde  []api.Bus `json:"conde"`
}
