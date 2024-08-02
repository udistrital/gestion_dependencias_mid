package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/gestion_dependencias_mid/controllers:GestionDependenciasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/gestion_dependencias_mid/controllers:GestionDependenciasController"],
        beego.ControllerComments{
            Method: "RegistrarDependencia",
            Router: "/RegistrarDependencia",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
