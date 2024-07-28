package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/gestion_dependencias_mid/models"
	"github.com/udistrital/gestion_dependencias_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
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
	defer errorhandler.HandlePanic(&c.Controller)
	var v models.NuevaDependencia
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil{
		resultado, err := services.RegistrarDependencia(&v)
		if err == nil {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
		}
	}
	


	c.ServeJSON()
}
