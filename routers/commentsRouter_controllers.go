package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

const controllerPathOne = "github.com/udistrital/gestion_dependencias_mid/controllers:GestionDependenciasController"
const controllerPathTwo = "github.com/udistrital/gestion_dependencias_mid/controllers:RegistroDependenciasController"

func init() {

    beego.GlobalControllerRouter[controllerPathOne] = append(beego.GlobalControllerRouter[controllerPathOne],
        beego.ControllerComments{
            Method: "BuscarDependencia",
            Router: "/BuscarDependencia",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter[controllerPathOne] = append(beego.GlobalControllerRouter[controllerPathOne],
        beego.ControllerComments{
            Method: "EditarDependencia",
            Router: "/EditarDependencia",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter[controllerPathOne] = append(beego.GlobalControllerRouter[controllerPathOne],
        beego.ControllerComments{
            Method: "Organigramas",
            Router: "/Organigramas",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter[controllerPathTwo] = append(beego.GlobalControllerRouter[controllerPathTwo],
        beego.ControllerComments{
            Method: "RegistrarDependencia",
            Router: "/RegistrarDependencia",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
