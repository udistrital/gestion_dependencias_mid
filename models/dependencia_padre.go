package models

type DependenciaPadre struct {
	Id					int
	PadreId				*Dependencia
	HijaId				*Dependencia
	Activo				bool
	FechaCreacion		string
	FechaModificacion	string
}