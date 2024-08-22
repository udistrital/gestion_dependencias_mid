package models


type Dependencia struct {
	Id							int
	Nombre						string
	TelefonoDependencia 		string
	CorreoElectronico 			string
	Activo						bool
	FechaCreacion				string
	FechaModificacion			string
	DependenciaTipoDependencia	[]*DependenciaTipoDependencia
}