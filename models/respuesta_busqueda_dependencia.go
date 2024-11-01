package models

type RespuestaBusquedaDependencia struct {
	Dependencia				*Dependencia
    DependenciaAsociada		*Dependencia
	TipoDependencia			*[]TipoDependencia
    Estado					bool
}