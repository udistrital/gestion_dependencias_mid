package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/xray"
)

const (
	JSON_error          string = "Error en el archivo JSON"
	ErrorParametros     string = "Error en los parametros de ingreso"
	ErrorBody           string = "Cuerpo de la peticion invalido"
	AppJson             string = "application/json"
	Calibri             string = "Calibri"
	CalibriBold         string = "Calibri-Bold"
	MinionProBoldCn     string = "MinionPro-BoldCn"
	MinionProMediumCn   string = "MinionPro-MediumCn"
	MinionProBoldItalic string = "MinionProBoldItalic"
)




// Valida que el body recibido en la petición tenga contenido válido
func ValidarBody(body []byte) (valid bool, err error) {
	var test interface{}
	if err = json.Unmarshal(body, &test); err != nil {
		return false, err
	} else {
		content := fmt.Sprintf("%v", test)
		fmt.Println(content)
		switch content {
		case "map[]", "[map[]]": // body vacio
			return false, nil
		}
	}
	return true, nil
}

// Manejo único de errores para controladores sin repetir código
func ErrorController(c beego.Controller, controller string) {
	if err := recover(); err != nil {
		logs.Error(err)
		localError := err.(map[string]interface{})
		c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + controller + "/" + (localError["funcion"]).(string))
		c.Data["data"] = (localError["err"])
		xray.EndSegmentErr(http.StatusBadRequest, localError["err"])
		if status, ok := localError["status"]; ok {
			c.Abort(status.(string))
		} else {
			c.Abort("500")
		}
	}
}