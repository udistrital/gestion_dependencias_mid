package models


type DependenciaTipoDependencia struct {
	Id							int
	TipoDependenciaId			*TipoDependencia
	DependenciaId				*Dependencia
	Activo						bool
	FechaCreacion				string
	FechaModificacion			string
}