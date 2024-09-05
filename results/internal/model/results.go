package model

type Result struct {
	ID         int  `gorm:"primary_key;auto_increment" json:"id"`
	EventID    int  `gorm:"not null" json:"event_id"`
	P1         int  `json:"p1"`
	P2         int  `json:"p2"`
	P3         int  `json:"p3"`
	P4         int  `json:"p4"`
	P5         int  `json:"p5"`
	P6         int  `json:"p6"`
	P7         int  `json:"p7"`
	P8         int  `json:"p8"`
	P9         int  `json:"p9"`
	P10        int  `json:"p10"`
	P11        int  `json:"p11"`
	P12        int  `json:"p12"`
	P13        int  `json:"p13"`
	P14        int  `json:"p14"`
	P15        int  `json:"p15"`
	P16        int  `json:"p16"`
	P17        int  `json:"p17"`
	P18        int  `json:"p18"`
	P19        int  `json:"p19"`
	P20        int  `json:"p20"`
	PFastLap   *int `json:"p_fast_lap,omitempty"`  // Puede ser nulo si no es un Race
	VSC        *bool `json:"vsc,omitempty"`        // Virtual Safety Car (true si hubo)
	SF         *bool `json:"sf,omitempty"`         // Safety Car (true si hubo)
	DNF        *int  `json:"dnf,omitempty"`        // Número de pilotos que no terminaron
}

//La función TableName() en un modelo de Go (usando GORM) 
//se utiliza para personalizar el nombre de la tabla en la base de datos a la que se asociará el modelo.

// Tabla de results conectada con event
func (Result) TableName() string {
	return "results"
}