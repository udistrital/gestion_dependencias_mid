package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/gestion_dependencias_mid/helpers"
	"github.com/udistrital/gestion_dependencias_mid/models"
	"github.com/udistrital/gestion_dependencias_mid/services"
)

// GestionDependenciasController operations for GestionDependencias
type GestionDependenciasController struct {
	beego.Controller
}

//URLMapping...
func (c *GestionDependenciasController) URLMapping(){
	c.Mapping("BuscarDependencia", c.BuscarDependencia)
	c.Mapping("EditarDependencia", c.EditarDependencia)
	c.Mapping("Organigramas", c.Organigramas)
}

// BuscarDependencia ...
// @Title BuscarDependencia
// @Description Buscar dependencia
// @Param	body		body 	{}	true		"body for Buscar Dependencia content"
// @Success 201 {init} 
// @Failure 400 the request contains incorrect syntax
// @router /BuscarDependencia [post]
func (c *GestionDependenciasController) BuscarDependencia() {
	defer helpers.ErrorController(c.Controller,"BuscarDependencia")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "BuscarDependencia", "err": helpers.ErrorBody, "status": "400"})
	}

	var v models.BusquedaDependencia

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if resultado, err := services.BuscarDependencia(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Resultado de busqueda", "Data": resultado}
		} else {
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "BuscarDependencia", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}

// EditarDependencia ...
// @Title EditarDependencia
// @Description Editar dependencia
// @Param	body		body 	{}	true		"body for Editar Dependencia content"
// @Success 201 {init} 
// @Failure 400 the request contains incorrect syntax
// @router /EditarDependencia [post]
func (c *GestionDependenciasController) EditarDependencia() {
	defer helpers.ErrorController(c.Controller,"EditarDependencia")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "EditarDependencia", "err": helpers.ErrorBody, "status": "400"})
	}
	var v models.EditarDependencia

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if resultado, err := services.EditarDependencia(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Dependencia editada con exito", "Data": resultado}
		} else {
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "EditarDependencia", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}

// Organigramas ...
// Get ...
// @Title Organigramas
// @Description Organigramas de Dependencias
// @Success 200 {object} models.Organigramas
// @Failure 400 the request contains incorrect syntax
// @router /Organigramas [get]
func (c *GestionDependenciasController) Organigramas() {
	fmt.Println("Entra a organigramas")
	defer helpers.ErrorController(c.Controller,"Organigramas")


	if organigramas, err := services.Organigramas(); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Organigramas", "Data": organigramas}
	} else {
		panic(map[string]interface{}{"funcion": "Organigramas", "err": err, "status": "400"})
	}

	c.ServeJSON()
}