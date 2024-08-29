package services

import (
	"net/url"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/gestion_dependencias_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
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
			resultado := CrearRespuestaBusqueda(dependencia)
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
			resultado := CrearRespuestaBusqueda(*dt.DependenciaId)
			if !ExisteDependencia(resultadoBusqueda, resultado.Dependencia.Id) {
				resultadoBusqueda = append(resultadoBusqueda, resultado)
			}
		}
	}

	if transaccion.FacultadId != 0 {
		var dependenciasxNombre []models.Dependencia
		url := beego.AppConfig.String("OikosCrudUrl") + "dependencia?query=Id:" + strconv.Itoa(transaccion.FacultadId)
		if err := request.GetJson(url, &dependenciasxNombre); err != nil {
			logs.Error(err)
			panic(err.Error())
		}

		for _, dependencia := range dependenciasxNombre {
			resultado := CrearRespuestaBusqueda(dependencia)
			if !ExisteDependencia(resultadoBusqueda, resultado.Dependencia.Id) {
				resultadoBusqueda = append(resultadoBusqueda, resultado)
			}
		}
	}

	if transaccion.VicerrectoriaId != 0 {
		var dependenciasxNombre []models.Dependencia
		url := beego.AppConfig.String("OikosCrudUrl") + "dependencia?query=Id:" + strconv.Itoa(transaccion.VicerrectoriaId)
		if err := request.GetJson(url, &dependenciasxNombre); err != nil {
			logs.Error(err)
			panic(err.Error())
		}

		for _, dependencia := range dependenciasxNombre {
			resultado := CrearRespuestaBusqueda(dependencia)
			if !ExisteDependencia(resultadoBusqueda, resultado.Dependencia.Id) {
				resultadoBusqueda = append(resultadoBusqueda, resultado)
			}
		}
	}

	if transaccion.BusquedaEstado != nil {
		var dependenciasxNombre []models.Dependencia
		url := beego.AppConfig.String("OikosCrudUrl") + "dependencia?query=Activo:" + strconv.FormatBool(transaccion.BusquedaEstado.Estado) + "&limit=-1"
		if err := request.GetJson(url, &dependenciasxNombre); err != nil {
			logs.Error(err)
			panic(err.Error())
		}

		for _, dependencia := range dependenciasxNombre {
			resultado := CrearRespuestaBusqueda(dependencia)
			if !ExisteDependencia(resultadoBusqueda, resultado.Dependencia.Id) {
				resultadoBusqueda = append(resultadoBusqueda, resultado)
			}
		}
	}

	return resultadoBusqueda, outputError
}

func CrearRespuestaBusqueda(dependencia models.Dependencia) models.RespuestaBusquedaDependencia {
	var resultado models.RespuestaBusquedaDependencia
	resultado.Dependencia = &dependencia
	resultado.Estado = dependencia.Activo
	if len(dependencia.DependenciaTipoDependencia) == 0 {
		var dependenciaAux []models.Dependencia
		url := beego.AppConfig.String("OikosCrudUrl") + "dependencia?query=Id:" + strconv.Itoa(dependencia.Id)
		if err := request.GetJson(url, &dependenciaAux); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
		dependencia = dependenciaAux[0]
	}
	tiposDependencia := make([]models.TipoDependencia, 0)
	for _, dt := range dependencia.DependenciaTipoDependencia {
		if dt.Activo {
			tipoDependencia := models.TipoDependencia{
				Id:     dt.TipoDependenciaId.Id,
				Nombre: dt.TipoDependenciaId.Nombre,
			}
			tiposDependencia = append(tiposDependencia, tipoDependencia)
		}
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

func EditarDependencia(transaccion *models.EditarDependencia) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "BuscarDependencia", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	alerta = append(alerta, "Success")

	var dependencia models.Dependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaId)
	if err := request.GetJson(url, &dependencia); err != nil || dependencia.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	dependencia.Nombre = transaccion.Nombre
	dependencia.CorreoElectronico = transaccion.CorreoElectronico
	dependencia.TelefonoDependencia = transaccion.TelefonoDependencia
	dependencia.FechaModificacion = time_bogota.TiempoBogotaFormato()

	var err error
	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaId)
	var respuestaDependencia map[string]interface{}
	if err = request.SendJson(url, "PUT", &respuestaDependencia, dependencia); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	var (
		dependenciaTipoDependencia []models.DependenciaTipoDependencia
		tiposActualesFalse         []int
		tiposMap                   = map[int]bool{}
	)

	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia_tipo_dependencia?query=dependenciaId:" + strconv.Itoa(transaccion.DependenciaId)
	if err := request.GetJson(url, &dependenciaTipoDependencia); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	for _, dependenciaTipo := range dependenciaTipoDependencia {
		tipoID := dependenciaTipo.TipoDependenciaId.Id
		tiposMap[tipoID] = false
		if !dependenciaTipo.Activo {
			tiposActualesFalse = append(tiposActualesFalse, tipoID)
		}
	}

	for _, tipo := range transaccion.TipoDependenciaId {
		if _, exists := tiposMap[tipo]; exists {
			tiposMap[tipo] = true
		} else {
			nuevoTipoDependencia(tipo, dependencia)
		}
	}

	for tipo, activo := range tiposMap {
		if !activo {
			actualizarDependenciaTipoDependencia(tipo, false, transaccion.DependenciaId)
		} else {
			if contiene(tiposActualesFalse, tipo) {
				actualizarDependenciaTipoDependencia(tipo, true, transaccion.DependenciaId)
			}
		}
	}

	var depedencia_padre []models.DependenciaPadre
	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre?query=HijaId:" + strconv.Itoa(transaccion.DependenciaId) + ",Activo:true"
	if err := request.GetJson(url, &depedencia_padre); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	if transaccion.DependenciaAsociadaId != depedencia_padre[0].HijaId.Id {

		var depedencia_padre_nueva models.Dependencia
		url := beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaAsociadaId)
		if err := request.GetJson(url, &depedencia_padre_nueva); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
		depedencia_padre[0].PadreId = &depedencia_padre_nueva
		depedencia_padre[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
		url = beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre/" + strconv.Itoa(depedencia_padre[0].Id)
		var res map[string]interface{}
		if err := request.SendJson(url, "PUT", &res, depedencia_padre[0]); err != nil {
			logs.Error(err)
			panic(err.Error())
		}

	}

	return alerta, outputError
}

func nuevoTipoDependencia(tipo int, dependenciaId models.Dependencia) {
	var tipoDependencia models.TipoDependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_dependencia/" + strconv.Itoa(tipo)
	if err := request.GetJson(url, &tipoDependencia); err != nil || tipoDependencia.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	nuevaDependenciaTipoDependencia := models.DependenciaTipoDependencia{
		TipoDependenciaId: &tipoDependencia,
		DependenciaId:     &dependenciaId,
		Activo:            true,
		FechaCreacion:     time_bogota.TiempoBogotaFormato(),
		FechaModificacion: time_bogota.TiempoBogotaFormato(),
	}

	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia_tipo_dependencia"
	var res map[string]interface{}
	if err := request.SendJson(url, "POST", &res, nuevaDependenciaTipoDependencia); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

}

func actualizarDependenciaTipoDependencia(tipo int, activo bool, dependenciaId int) {
	var dependenciaTipoDependenciaActual []models.DependenciaTipoDependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia_tipo_dependencia?query=dependenciaId:" + strconv.Itoa(dependenciaId) + ",tipoDependenciaId:" + strconv.Itoa(tipo)
	if err := request.GetJson(url, &dependenciaTipoDependenciaActual); err != nil || len(dependenciaTipoDependenciaActual) == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	dependenciaTipoDependenciaActual[0].Activo = activo
	dependenciaTipoDependenciaActual[0].FechaModificacion = time_bogota.TiempoBogotaFormato()

	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia_tipo_dependencia/" + strconv.Itoa(dependenciaTipoDependenciaActual[0].Id)
	var res map[string]interface{}
	if err := request.SendJson(url, "PUT", &res, dependenciaTipoDependenciaActual[0]); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
}

func contiene(slice []int, valor int) bool {
	for _, item := range slice {
		if item == valor {
			return true
		}
	}
	return false
}

func Organigramas() (organigramas models.Organigramas, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "BuscarDependencia", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var dependencias []models.Dependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia?limit=-1"
	if err := request.GetJson(url, &dependencias); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	var dependencias_padre []models.DependenciaPadre
	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre?limit=-1"
	if err := request.GetJson(url, &dependencias_padre); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	organigramaMap := make(map[int]*models.Organigrama)

	for _, dep := range dependencias {
		organigramaMap[dep.Id] = &models.Organigrama{Dependencia: dep}
	}

	esHijo := make(map[int]bool)

	//Construir el arbol
	for _, dep_padre := range dependencias_padre {
		padre := organigramaMap[dep_padre.PadreId.Id]
		hija := organigramaMap[dep_padre.HijaId.Id]
		padre.Hijos = append(padre.Hijos, hija)
		esHijo[dep_padre.HijaId.Id] = true
	}

	// Encontrar las raíces del árbol (nodos que no son hijos de nadie)
	var raiz []*models.Organigrama
	for id, org := range organigramaMap {
		if !esHijo[id] {
			raiz = append(raiz, org)
		}
	}

	organigramas.General = raiz

	/*var dependencias []models.Dependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia?limit=-1"
	if err := request.GetJson(url, &dependencias); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	var dependencias_struct []*models.Organigrama
	for _, dependencia := range dependencias {
		dependencia_items := &models.Organigrama{
			Dependencia: dependencia.Nombre,
		}
		var tiposDependencia []string
		for _, tipos := range dependencia.DependenciaTipoDependencia {
			if tipos.Activo {
				tiposDependencia = append(tiposDependencia, tipos.TipoDependenciaId.Nombre)
			}
		}
		dependencia_items.Tipo = tiposDependencia
		dependencias_struct = append(dependencias_struct, dependencia_items)
	}


	var dependencias_padre []models.DependenciaPadre
	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre?limit=-1"
	if err := request.GetJson(url, &dependencias_padre); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	dependenciasPasadas := make(map[string]bool)
	for i := 0; i < len(dependencias_struct); i++ {
		dependencia := dependencias_struct[i]
		for j := 0; j < len(dependencias_padre); j++ {
			arbol := dependencias_padre[j]
			if dependencia.Dependencia == arbol.PadreId.Nombre {
				for k := 0; k < len(dependencias_struct); k++ {
					dependenciaHija := dependencias_struct[k]
					if dependenciaHija.Dependencia == arbol.HijaId.Nombre {
						dependencia.Hijos = append(dependencia.Hijos, dependenciaHija)
						if dependenciasPasadas[dependenciaHija.Dependencia] {
							dependencias_struct = append(dependencias_struct[:k], dependencias_struct[k+1:]...)
							k--
						}
					}
				}
			}
		}
		dependenciasPasadas[dependencia.Dependencia] = true
	}
	organigramas.General = dependencias_struct*/

	return organigramas, outputError
}
