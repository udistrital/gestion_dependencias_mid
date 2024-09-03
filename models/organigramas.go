package models

type Organigramas struct {
	General        []*Organigrama
	Academico      []*Organigrama
	Administrativo []*Organigrama
}

type Organigrama struct {
	Dependencia Dependencia
	Tipo        []string
	Hijos       []*Organigrama
}
