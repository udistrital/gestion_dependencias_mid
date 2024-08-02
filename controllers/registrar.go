package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/gestion_dependencias_mid/models"
	"github.com/udistrital/gestion_dependencias_mid/services"
	"github.com/udistrital/gestion_dependencias_mid/helpers"
)

// GestionDependenciasController operations for GestionDependencias
type GestionDependenciasController struct {
	beego.Controller
}

//URLMapping...
func (c *GestionDependenciasController) URLMapping(){
	c.Mapping("RegistrarDependencia", c.RegistrarDependencia)
}

// RegistrarDependencia ...
// @Title RegistrarDependencia
// @Description Registrar dependencia
// @Param	body		body 	{}	true		"body for Registrar Dependencia content"
// @Success 201 {init} 
// @Failure 400 the request contains incorrect syntax
// @router /RegistrarDependencia [post]
func (c *GestionDependenciasController) RegistrarDependencia() {
	defer helpers.ErrorController(c.Controller,"RegistrarDependencia")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "RegistrarDependencia", "err": helpers.ErrorBody, "status": "400"})
	}

	var v models.NuevaDependencia

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if resultado, err := services.RegistrarDependencia(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Dependencia insertada con exito", "Data": resultado}
		} else {
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "RegistrarDependencia", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
