package models

type BusquedaDependencia struct {
	NombreDependencia	string
    TipoDependenciaId	int
    FacultadId			int
    VicerrectoriaId		int
    BusquedaEstado		*Estado
} 