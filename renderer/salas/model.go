package salas

type salasResponse struct {
	Salas []salaTimesheet `json:"salas"`
}

type salaTimesheet struct {
	ID       int        `json:"id"`
	Occupied []timeSlot `json:"occupied"`
}

type timeSlot struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
