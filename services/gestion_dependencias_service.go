package services

import (
	"strconv"
	"net/url"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/gestion_dependencias_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func BuscarDependencia(transaccion *models.BusquedaDependencia) (resultadoBusqueda []models.RespuestaBusquedaDependencia, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "BuscarDependencia", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	if transaccion.NombreDependencia != "" {
		var dependenciasxNombre []models.Dependencia
		nombreDependencia := url.QueryEscape(transaccion.NombreDependencia)
		url := beego.AppConfig.String("OikosCrudUrl") + "dependencia?query=Nombre:" + nombreDependencia
		if err := request.GetJson(url, &dependenciasxNombre); err != nil {
			logs.Error(err)
			panic(err.Error())
		}

		for _, dependencia := range dependenciasxNombre {
			resultado := crearRespuestaBusqueda(dependencia)
			if !ExisteDependencia(resultadoBusqueda, resultado.Dependencia.Id) {
				resultadoBusqueda = append(resultadoBusqueda, resultado)
			}
		}
	}

	if transaccion.TipoDependenciaId != 0 {
		var dependenciaTipoDependencia []models.DependenciaTipoDependencia
		url := beego.AppConfig.String("OikosCrudUrl") + "dependencia_tipo_dependencia?query=TipoDependenciaId.Id:" + strconv.Itoa(transaccion.TipoDependenciaId) + "&limit=-1"
		if err := request.GetJson(url, &dependenciaTipoDependencia); err != nil {
			logs.Error(err)
			panic(err.Error())
		}

		for _, dt := range dependenciaTipoDependencia {
			resultado := crearRespuestaBusqueda(*dt.DependenciaId)
			if !ExisteDependencia(resultadoBusqueda, resultado.Dependencia.Id) {
				resultadoBusqueda = append(resultadoBusqueda, resultado)
			}
		}
	}

	if transaccion.FacultadId != 0{
		var dependenciasxNombre []models.Dependencia
		url := beego.AppConfig.String("OikosCrudUrl") + "dependencia?query=Id:" + strconv.Itoa(transaccion.FacultadId)
		if err := request.GetJson(url, &dependenciasxNombre); err != nil {
			logs.Error(err)
			panic(err.Error())
		}

		for _, dependencia := range dependenciasxNombre {
			resultado := crearRespuestaBusqueda(dependencia)
			if !ExisteDependencia(resultadoBusqueda, resultado.Dependencia.Id) {
				resultadoBusqueda = append(resultadoBusqueda, resultado)
			}
		}
	}

	if transaccion.VicerrectoriaId != 0{
		var dependenciasxNombre []models.Dependencia
		url := beego.AppConfig.String("OikosCrudUrl") + "dependencia?query=Id:" + strconv.Itoa(transaccion.VicerrectoriaId)
		if err := request.GetJson(url, &dependenciasxNombre); err != nil {
			logs.Error(err)
			panic(err.Error())
		}

		for _, dependencia := range dependenciasxNombre {
			resultado := crearRespuestaBusqueda(dependencia)
			if !ExisteDependencia(resultadoBusqueda, resultado.Dependencia.Id) {
				resultadoBusqueda = append(resultadoBusqueda, resultado)
			}
		}
	}

	if transaccion.BusquedaEstado != nil{
		var dependenciasxNombre []models.Dependencia
		url := beego.AppConfig.String("OikosCrudUrl") + "dependencia?query=Activo:" + strconv.FormatBool(transaccion.BusquedaEstado.Estado) + "&limit=-1"
		if err := request.GetJson(url, &dependenciasxNombre); err != nil {
			logs.Error(err)
			panic(err.Error())
		}

		for _, dependencia := range dependenciasxNombre {
			resultado := crearRespuestaBusqueda(dependencia)
			if !ExisteDependencia(resultadoBusqueda, resultado.Dependencia.Id) {
				resultadoBusqueda = append(resultadoBusqueda, resultado)
			}
		}
	}

	return resultadoBusqueda, outputError
}

func crearRespuestaBusqueda(dependencia models.Dependencia) models.RespuestaBusquedaDependencia {
	var resultado models.RespuestaBusquedaDependencia
	resultado.Dependencia = &dependencia
	resultado.Estado = dependencia.Activo
	if len(dependencia.DependenciaTipoDependencia) == 0{
		var dependenciaAux []models.Dependencia
		url := beego.AppConfig.String("OikosCrudUrl") + "dependencia?query=Id:" +  strconv.Itoa(dependencia.Id) 
		if err := request.GetJson(url, &dependenciaAux); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
		dependencia = dependenciaAux[0]
	}
	tiposDependencia := make([]models.TipoDependencia, 0)
	for _, dt := range dependencia.DependenciaTipoDependencia {
		tipoDependencia := models.TipoDependencia{
			Id:     dt.TipoDependenciaId.Id,
			Nombre: dt.TipoDependenciaId.Nombre,
		}
		tiposDependencia = append(tiposDependencia, tipoDependencia)
	}
	resultado.TipoDependencia = &tiposDependencia

	var dependenciaPadre []models.DependenciaPadre
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre?query=HijaId:" + strconv.Itoa(dependencia.Id)
	if err := request.GetJson(url, &dependenciaPadre); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	if len(dependenciaPadre) > 0 {
		resultado.DependenciaAsociada = dependenciaPadre[0].PadreId
	}

	return resultado
}

func ExisteDependencia(resultadoBusqueda []models.RespuestaBusquedaDependencia, dependenciaId int) bool {
	for _, r := range resultadoBusqueda {
		if r.Dependencia != nil && r.Dependencia.Id == dependenciaId {
			return true
		}
	}
	return false
}
