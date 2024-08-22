package models


type EditarDependencia struct {
	DependenciaId			int
	Nombre					string
	TelefonoDependencia 	string
	CorreoElectronico 		string
	DependenciaAsociadaId	int
	TipoDependenciaId		[]int
}