package models

type NuevaDependencia struct {
	Dependencia 			*Dependencia
	TipoDependenciaId		[]int
	DependenciaAsociadaId	int
}