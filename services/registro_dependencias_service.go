package services

import (
	"fmt"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/udistrital/gestion_dependencias_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/astaxie/beego/logs"
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
	var tipoDependencia models.TipoDependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_dependencia/" + strconv.Itoa(transaccion.TipoDependenciaId)
	if err := request.GetJson(url,&tipoDependencia); err != nil || tipoDependencia.Id == 0{
		logs.Error(err)
		panic(err.Error())
	}
	var dependenciaAsociada models.Dependencia
	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaAsociadaId)
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
	var dependenciaTipoDependenciaNueva models.DependenciaTipoDependencia
	dependenciaTipoDependenciaNueva.TipoDependenciaId = &tipoDependencia
	dependenciaTipoDependenciaNueva.DependenciaId = &dependenciaNueva
	dependenciaTipoDependenciaNueva.Activo = true
	dependenciaTipoDependenciaNueva.FechaCreacion = time_bogota.TiempoBogotaFormato()
	dependenciaTipoDependenciaNueva.FechaModificacion = time_bogota.TiempoBogotaFormato()
	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia_tipo_dependencia"
	var resDependenciaTipoDependenciaRegistrada map[string]interface{}
	if err := request.SendJson(url,"POST",&resDependenciaTipoDependenciaRegistrada,dependenciaTipoDependenciaNueva); err != nil{
		rollbackDependenciaCreada(&creaciones)
		logs.Error(err)
		panic(err.Error())
	}
	fmt.Println("FUNCIONO DEPENDENCIA TIPO DEPENDENCIA")
	fmt.Println(resDependenciaTipoDependenciaRegistrada["Id"])
	dependenciaTipoDependenciaNueva.Id = int(resDependenciaTipoDependenciaRegistrada["Id"].(float64))
	creaciones.DependenciaTipoDependenciaId = int(resDependenciaTipoDependenciaRegistrada["Id"].(float64))
	var depedencia_padre models.DependenciaPadre
	depedencia_padre.PadreId = &dependenciaAsociada
	depedencia_padre.HijaId = &dependenciaNueva
	depedencia_padre.Activo = true
	depedencia_padre.FechaCreacion = time_bogota.TiempoBogotaFormato()
	depedencia_padre.FechaModificacion = time_bogota.TiempoBogotaFormato()
	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre"
	var resDependenciaPadre map[string]interface{}
	if err := request.SendJson(url,"POST",&resDependenciaPadre,depedencia_padre);err != nil{
		rollbackDependenciaTipoDependenciaCreada(&creaciones)
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
// 	creaciones.DependenciaId = 282
// 	creaciones.DependenciaTipoDependenciaId = 433
// 	creaciones.DependenciaPadreId = 225
// 	rollbackDependenciaPadreCreada(&creaciones)
// 	alerta = append(alerta, "Success")
// 	return alerta, outputError
// }

func rollbackDependenciaCreada(transaccion *models.Creaciones) (outputError map[string]interface{}) {
	var respuesta map[string]interface{}
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaId)
	if err := request.SendJson(url,"DELETE",&respuesta,nil); err != nil{
		panic("Rollback de dependencia" + err.Error())
	}
	return nil
}


func rollbackDependenciaTipoDependenciaCreada(transaccion *models.Creaciones) (outputError map[string]interface{}) {
	var respuesta map[string]interface{}
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia_tipo_dependencia/" + strconv.Itoa(transaccion.DependenciaTipoDependenciaId)
	if err := request.SendJson(url,"DELETE",&respuesta,nil); err != nil{
		panic("Rollback de dependencia tipo dependencia" + err.Error())
	}
	rollbackDependenciaCreada(transaccion)

	return nil
}

func rollbackDependenciaPadreCreada(transaccion *models.Creaciones) (outputError map[string]interface{}) {
	var respuesta map[string]interface{}
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre/" + strconv.Itoa(transaccion.DependenciaPadreId)
	if err := request.SendJson(url,"DELETE",&respuesta,nil); err != nil{
		panic("Rollback de dependencia padre" + err.Error())
	}
	rollbackDependenciaTipoDependenciaCreada(transaccion)
	return nil
}