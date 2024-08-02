package models


type Dependencia struct {
	Id							int
	Nombre						string
	TelefonoDependencia 		string
	CorreoElectronico 			string
	DependenciaTipoDependencia	[]*DependenciaTipoDependencia
	Activo						bool
	FechaCreacion				string
	FechaModificacion			string
}