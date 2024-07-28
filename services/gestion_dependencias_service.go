package services

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/udistrital/gestion_dependencias_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/time_bogota"
)

func RegistrarDependencia(transaccion *models.NuevaDependencia) (interface{}, error){

	var tipoDependencia models.TipoDependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_dependencia/" + strconv.Itoa(transaccion.TipoDependenciaId)
	if err := request.GetJson(url,&tipoDependencia); err == nil{
		var dependenciaAsociada models.Dependencia
		url = beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaAsociadaId)
		if err := request.GetJson(url,&dependenciaAsociada); err == nil{
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
			if err := request.SendJson(url,"POST",&resDependenciaRegistrada,dependenciaNueva); err == nil{
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
				if err := request.SendJson(url,"POST",&resDependenciaTipoDependenciaRegistrada,dependenciaTipoDependenciaNueva); err == nil{
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
					if err := request.SendJson(url,"POST",&resDependenciaPadre,depedencia_padre);err == nil{
						fmt.Println("FUNCIONO DEPENDENCIA PADRE")
						fmt.Println(resDependenciaPadre["Id"])
						depedencia_padre.Id = int(resDependenciaPadre["Id"].(float64))
						creaciones.DependenciaPadreId = int(resDependenciaPadre["Id"].(float64))
						return dependenciaNueva, nil
					}else{
						logs.Error(err)
						rollbackDependenciaPadreCreada(&creaciones)
						return nil, errors.New(err.Error())
					}			
				}else{
					logs.Error(err)
					rollbackDependenciaTipoDependenciaCreada(&creaciones)
					return nil, errors.New(err.Error())
				}
			}else{
				logs.Error(err)
				return nil, errors.New(err.Error())
			}
		}else{
			logs.Error(err)
			return nil, errors.New("no se encontró la dependencia asociada")
		}
	}else{
		logs.Error(err)
		return nil, errors.New("No se encontró el tipo de dependencia")
	}
}


// func RegistrarDependencia(transaccion *models.NuevaDependencia) (interface{}, error){
// 	var creaciones models.Creaciones
// 	creaciones.DependenciaId = 278
// 	creaciones.DependenciaTipoDependenciaId = 429
// 	creaciones.DependenciaPadreId = 221
// 	rollbackDependenciaPadreCreada(&creaciones)
// 	return nil, errors.New("Se borraron los registros correctamente")
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
	}else{
		rollbackDependenciaCreada(transaccion)
	}
	return nil
}

func rollbackDependenciaPadreCreada(transaccion *models.Creaciones) (outputError map[string]interface{}) {
	var respuesta map[string]interface{}
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre/" + strconv.Itoa(transaccion.DependenciaPadreId)
	if err := request.SendJson(url,"DELETE",&respuesta,nil); err != nil{
		panic("Rollback de dependencia padre" + err.Error())
	}else{
		rollbackDependenciaTipoDependenciaCreada(transaccion)
	}
	return nil
}