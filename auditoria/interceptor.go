package auditoria

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/udistrital/utils_oas/request"
)

func InterceptMidRequest(ctx *context.Context) {
	fmt.Println("Entra a interceptor");
	end_point := ctx.Request.URL.String()
	if end_point != "/" {
		defer func() {
			//Catch
			if r := recover(); r != nil {
				errMsg := fmt.Sprintf("Panic interceptado: %v", r)
				CreateLog(ctx, 400, "Error interno", errMsg) // Llama al sistema de logs
			}
		}()
		// try
		request.SetHeader(ctx.Request.Header["Authorization"][0])
	}

}

func InitInterceptor() {
	fmt.Println("Entra al InitInterceptor")
	beego.InsertFilter("*", beego.BeforeExec, InterceptMidRequest, false)
}
