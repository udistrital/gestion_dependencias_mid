package models

import "time"

type TipoDependencia struct {
	Id							int
	Nombre						string
	Descripcion			 		string
	CodigoAbreviacion 			string
	Activo						bool
	FechaCreacion				time.Time 
	FechaModificacion			time.Time 
}