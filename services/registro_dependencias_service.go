package services

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/gestion_dependencias_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

func RegistrarDependencia(transaccion *models.NuevaDependencia) (alerta []string, outputError map[string]interface{}){
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "RegistrarDependencia", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	alerta = append(alerta, "Success")

	var tiposDependencia []*models.TipoDependencia
	for _, tipo := range transaccion.TipoDependenciaId {
		var tipoDependencia models.TipoDependencia
		tipoDependencia = VerificarExistenciaTipo(tipo)
		tiposDependencia = append(tiposDependencia, &tipoDependencia)
	}

	var dependenciaAsociada models.Dependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaAsociadaId)
	if err := request.GetJson(url,&dependenciaAsociada); err != nil || dependenciaAsociada.Id == 0{
		logs.Error(err)
		panic(err.Error())
	}

	var creaciones models.Creaciones
	var dependenciaNueva models.Dependencia
	dependenciaNueva.Nombre = transaccion.Dependencia.Nombre
	dependenciaNueva.CorreoElectronico = transaccion.Dependencia.CorreoElectronico
	dependenciaNueva.TelefonoDependencia = transaccion.Dependencia.TelefonoDependencia
	dependenciaNueva.Activo = true
	dependenciaNueva.FechaCreacion = time_bogota.TiempoBogotaFormato()
	dependenciaNueva.FechaModificacion = time_bogota.TiempoBogotaFormato()
	var resDependenciaRegistrada map[string]interface{}
	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia"
	if err := request.SendJson(url,"POST",&resDependenciaRegistrada,dependenciaNueva); err != nil{
		logs.Error(err)
		panic(err.Error())
	}
	fmt.Println("FUNCIONO REGISTRO DE DEPENDENCIA")
	fmt.Println(resDependenciaRegistrada["Id"])
	dependenciaNueva.Id = int(resDependenciaRegistrada["Id"].(float64))
	creaciones.DependenciaId = 	int(resDependenciaRegistrada["Id"].(float64))
	for _, tipoDependencia := range tiposDependencia{
		var dependenciaTipoDependenciaNueva models.DependenciaTipoDependencia
		dependenciaTipoDependenciaNueva.TipoDependenciaId = tipoDependencia
		dependenciaTipoDependenciaNueva.DependenciaId = &dependenciaNueva
		dependenciaTipoDependenciaNueva.Activo = true
		dependenciaTipoDependenciaNueva.FechaCreacion = time_bogota.TiempoBogotaFormato()
		dependenciaTipoDependenciaNueva.FechaModificacion = time_bogota.TiempoBogotaFormato()
		url = beego.AppConfig.String("OikosCrudUrl") + "dependencia_tipo_dependencia"
		var resDependenciaTipoDependenciaRegistrada map[string]interface{}
		if err := request.SendJson(url,"POST",&resDependenciaTipoDependenciaRegistrada,dependenciaTipoDependenciaNueva); err != nil{
			if (len(creaciones.DependenciaTipoDependenciaId) >= 1){
				RollbackDependenciaTipoDependenciaCreada(&creaciones)
			}else{
				RollbackDependenciaCreada(&creaciones)
			}
			logs.Error(err)
			panic(err.Error())
		}
		creaciones.DependenciaTipoDependenciaId = append(creaciones.DependenciaTipoDependenciaId, int(resDependenciaTipoDependenciaRegistrada["Id"].(float64)))
	}
	fmt.Println("FUNCIONO DEPENDENCIA TIPO DEPENDENCIA")
	fmt.Println(creaciones.DependenciaTipoDependenciaId)
	var depedencia_padre models.DependenciaPadre
	depedencia_padre.PadreId = &dependenciaAsociada
	depedencia_padre.HijaId = &dependenciaNueva
	depedencia_padre.Activo = true
	depedencia_padre.FechaCreacion = time_bogota.TiempoBogotaFormato()
	depedencia_padre.FechaModificacion = time_bogota.TiempoBogotaFormato()
	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre"
	var resDependenciaPadre map[string]interface{}
	if err := request.SendJson(url,"POST",&resDependenciaPadre,depedencia_padre);err != nil{
		RollbackDependenciaTipoDependenciaCreada(&creaciones)
		logs.Error(err)
		panic(err.Error())
	}
	fmt.Println("FUNCIONO DEPENDENCIA PADRE")
	fmt.Println(resDependenciaPadre["Id"])

	return alerta, outputError
}

// func RegistrarDependencia(transaccion *models.NuevaDependencia) (alerta []string, outputError map[string]interface{}){
// 	defer func() {
// 		if err := recover(); err != nil {
// 			outputError = map[string]interface{}{"funcion": "RegistrarDependencia", "err": err, "status": "500"}
// 			panic(outputError)
// 		}
// 	}()
// 	var creaciones models.Creaciones
// 	creaciones.DependenciaId = 298
// 	creaciones.DependenciaTipoDependenciaId = append(creaciones.DependenciaTipoDependenciaId, 459)
// 	creaciones.DependenciaTipoDependenciaId = append(creaciones.DependenciaTipoDependenciaId, 460)
// 	creaciones.DependenciaTipoDependenciaId = append(creaciones.DependenciaTipoDependenciaId, 461)	
// 	creaciones.DependenciaPadreId = 246
// 	RollbackDependenciaPadreCreada(&creaciones)
// 	alerta = append(alerta, "Success")
// 	return alerta, outputError
// }

func VerificarExistenciaTipo(tipo int) (tipoDependencia models.TipoDependencia){
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_dependencia/" + strconv.Itoa(tipo)
	if err := request.GetJson(url,&tipoDependencia); err != nil || tipoDependencia.Id == 0{
		logs.Error(err)
		panic(err.Error())
	}
	return tipoDependencia
}

func RollbackDependenciaCreada(transaccion *models.Creaciones) (outputError map[string]interface{}) {
	var respuesta map[string]interface{}
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaId)
	if err := request.SendJson(url,"DELETE",&respuesta,nil); err != nil{
		panic("Rollback de dependencia" + err.Error())
	}
	return nil
}


func RollbackDependenciaTipoDependenciaCreada(transaccion *models.Creaciones) (outputError map[string]interface{}) {
	fmt.Println(transaccion.DependenciaTipoDependenciaId)
	for _, tipo := range transaccion.DependenciaTipoDependenciaId{
		var respuesta map[string]interface{}
		url := beego.AppConfig.String("OikosCrudUrl") + "dependencia_tipo_dependencia/" + strconv.Itoa(tipo)
		if err := request.SendJson(url,"DELETE",&respuesta,nil); err != nil{
			panic("Rollback de dependencia tipo dependencia" + err.Error())
		}
	}
	RollbackDependenciaCreada(transaccion)

	return nil
}

func RollbackDependenciaPadreCreada(transaccion *models.Creaciones) (outputError map[string]interface{}) {
	var respuesta map[string]interface{}
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre/" + strconv.Itoa(transaccion.DependenciaPadreId)
	if err := request.SendJson(url,"DELETE",&respuesta,nil); err != nil{
		panic("Rollback de dependencia padre" + err.Error())
	}
	RollbackDependenciaTipoDependenciaCreada(transaccion)
	return nil
}