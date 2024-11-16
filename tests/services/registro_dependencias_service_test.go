package services_test

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/astaxie/beego"
	"github.com/udistrital/gestion_dependencias_mid/models"
	"github.com/udistrital/gestion_dependencias_mid/services"
	"github.com/udistrital/utils_oas/request"
)

// const separadorSlash = "//////////////////////////////////";
// const errorUrlInesperada = "URL no esperada";
// const mensajePanicCapturado = "Panic capturado correctamente: ";
const errorPanicTipoString = "Se esperaba un panic de tipo string, pero se obtuvo: "

// const errorPanicFaltante = "Se esperaba un panic, pero no ocurrió";
const errorOutputEsperadoNulo = "Se esperaba un outputError nulo, pero se obtuvo: "

func TestRegistrarDependencia(t *testing.T) {
	t.Log(separadorSlash)
	t.Log("Inicio TestRegistrarDependencia")
	t.Log(separadorSlash)

	t.Run("Caso 1: Registro exitoso de dependencia", func(t *testing.T) {
		transaccion := &models.NuevaDependencia{
			Dependencia: &models.Dependencia{
				Nombre:              "Dependencia de prueba",
				CorreoElectronico:   "dependencia@prueba.com",
				TelefonoDependencia: "123456789",
			},
			TipoDependenciaId:     []int{1, 2},
			DependenciaAsociadaId: 1,
		}

		monkey.Patch(services.VerificarExistenciaTipo, func(tipoId int) models.TipoDependencia {
			return models.TipoDependencia{Id: tipoId, Nombre: fmt.Sprintf("Tipo %d", tipoId)}
		})
		defer monkey.Unpatch(services.VerificarExistenciaTipo)

		monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			expectedUrl := beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaAsociadaId)
			if url == expectedUrl {
				*target.(*models.Dependencia) = models.Dependencia{Id: transaccion.DependenciaAsociadaId}
				return nil
			}
			return errors.New(errorUrlInesperada)
		})
		defer monkey.Unpatch(request.GetJson)

		monkey.Patch(request.SendJson, func(url, method string, target interface{}, body interface{}) error {
			switch url {
			case beego.AppConfig.String("OikosCrudUrl") + "dependencia":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(10)}
				return nil
			case beego.AppConfig.String("OikosCrudUrl") + "dependencia_tipo_dependencia":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(20)}
				return nil
			case beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(30)}
				return nil
			}
			return errors.New(errorUrlInesperada)
		})
		defer monkey.Unpatch(request.SendJson)

		monkey.Patch(services.RollbackDependenciaTipoDependenciaCreada, func(transaccion *models.Creaciones) map[string]interface{} {
			return nil
		})
		defer monkey.Unpatch(services.RollbackDependenciaTipoDependenciaCreada)

		monkey.Patch(services.RollbackDependenciaCreada, func(transaccion *models.Creaciones) map[string]interface{} {
			return nil
		})
		defer monkey.Unpatch(services.RollbackDependenciaCreada)

		alerta, outputError := services.RegistrarDependencia(transaccion)

		if len(alerta) == 0 || alerta[0] != "Success" {
			t.Errorf("Se esperaba una alerta con 'Success', pero se obtuvo: %v", alerta)
		}

		if outputError != nil {
			t.Errorf("Se esperaba outputError nulo, pero se obtuvo: %v", outputError)
		}

		t.Log("Registro de dependencia ejecutado exitosamente sin errores")
	})

	t.Run("Caso_2: Fallo debido a error en los datos", func(t *testing.T) {
		transaccion := &models.NuevaDependencia{
			Dependencia: &models.Dependencia{
				Nombre:              "Dependencia de prueba",
				CorreoElectronico:   "dependencia@prueba.com",
				TelefonoDependencia: "123456789",
			},
			TipoDependenciaId:     []int{1, 2},
			DependenciaAsociadaId: 1,
		}

		monkey.Patch(services.VerificarExistenciaTipo, func(tipoId int) models.TipoDependencia {
			return models.TipoDependencia{Id: tipoId, Nombre: fmt.Sprintf("Tipo %d", tipoId)}
		})
		defer monkey.Unpatch(services.VerificarExistenciaTipo)

		monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			return errors.New("Error al obtener dependencia asociada")
		})
		defer monkey.Unpatch(request.GetJson)

		monkey.Patch(services.RollbackDependenciaTipoDependenciaCreada, func(transaccion *models.Creaciones) map[string]interface{} {
			return nil
		})
		defer monkey.Unpatch(services.RollbackDependenciaTipoDependenciaCreada)

		monkey.Patch(services.RollbackDependenciaCreada, func(transaccion *models.Creaciones) map[string]interface{} {
			return nil
		})
		defer monkey.Unpatch(services.RollbackDependenciaCreada)

		defer func() {
			if r := recover(); r != nil {
				if outputError, ok := r.(map[string]interface{}); ok {
					if outputError["status"] != "500" {
						t.Errorf("Se esperaba status 500, pero se obtuvo: %v", outputError["status"])
					}
					if outputError["funcion"] != "RegistrarDependencia" {
						t.Errorf("Se esperaba funcion 'RegistrarDependencia', pero se obtuvo: %v", outputError["funcion"])
					}
				} else {
					t.Errorf("Se esperaba un error estructurado, pero se obtuvo: %v", r)
				}
			}
		}()

		alerta, _ := services.RegistrarDependencia(transaccion)

		if len(alerta) > 0 {
			t.Errorf("No se esperaba alerta, pero se obtuvo: %v", alerta)
		}

		t.Log("Test de fallo debido a error en los datos ejecutado correctamente")
	})

}

func TestVerificarExistenciaTipo(t *testing.T) {
	t.Log(separadorSlash)
	t.Log("Inicio TestVerificarExistenciaTipo")
	t.Log(separadorSlash)

	t.Run("Caso 1: Tipo de dependencia existe", func(t *testing.T) {
		tipoId := 1
		expectedTipo := models.TipoDependencia{
			Id:                tipoId,
			Nombre:            "Tipo Ejemplo",
			Descripcion:       "Descripción de ejemplo",
			CodigoAbreviacion: "TE",
			Activo:            true,
			FechaCreacion:     time.Date(2022, 10, 22, 10, 0, 0, 0, time.UTC),
			FechaModificacion: time.Date(2023, 10, 22, 10, 0, 0, 0, time.UTC),
		}

		monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			if url == beego.AppConfig.String("OikosCrudUrl")+"tipo_dependencia/"+strconv.Itoa(tipoId) {
				tipoDependenciaPtr := target.(*models.TipoDependencia)
				*tipoDependenciaPtr = expectedTipo
				return nil
			}
			return errors.New(errorUrlInesperada)
		})
		defer monkey.UnpatchAll()

		tipoDependencia := services.VerificarExistenciaTipo(tipoId)

		if tipoDependencia != expectedTipo {
			t.Errorf("Se esperaba %v, pero se obtuvo %v", expectedTipo, tipoDependencia)
		}

		t.Log("Tipo de dependencia verificado exitosamente")
	})

	t.Run("Caso 2: Tipo de dependencia no existe", func(t *testing.T) {
		tipoId := 2
		monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			return errors.New("Tipo de dependencia no encontrado")
		})
		defer monkey.UnpatchAll()

		defer func() {
			if r := recover(); r != nil {
				errorMessage, ok := r.(string)
				if ok {
					t.Logf(mensajePanicCapturado, errorMessage)
					if !strings.Contains(errorMessage, "Tipo de dependencia no encontrado") {
						t.Errorf("Se esperaba un mensaje de error relacionado con tipo no encontrado, pero se obtuvo: %v", errorMessage)
					}
				} else {
					t.Errorf(errorPanicTipoString, r)
				}
			} else {
				t.Errorf(errorPanicFaltante)
			}
		}()

		services.VerificarExistenciaTipo(tipoId)
	})
}

func TestRollbackDependenciaCreada(t *testing.T) {
	t.Log(separadorSlash)
	t.Log("Inicio TestRollbackDependenciaCreada")
	t.Log(separadorSlash)

	t.Run("Caso 1: Rollback exitoso de dependencia creada", func(t *testing.T) {
		transaccion := &models.Creaciones{
			DependenciaId:                1,
			DependenciaTipoDependenciaId: []int{1, 2},
			DependenciaPadreId:           3,
		}

		monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
			if strings.Contains(url, "dependencia/1") && method == "DELETE" {
				return nil
			}
			return errors.New(errorUrlInesperada)
		})
		defer monkey.UnpatchAll()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Se esperaba que no se generara un panic, pero se obtuvo: %v", r)
			}
		}()

		outputError := services.RollbackDependenciaCreada(transaccion)

		if outputError != nil {
			t.Errorf(errorOutputEsperadoNulo, outputError)
		}

		t.Log("Rollback de dependencia ejecutado exitosamente sin errores")
	})

	t.Run("Caso 2: Error en el rollback por datos mal formados", func(t *testing.T) {
		transaccion := &models.Creaciones{
			DependenciaId: 0,
		}

		monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
			return errors.New("Error en el rollback por datos mal formados")
		})

		defer monkey.UnpatchAll()

		defer func() {
			if r := recover(); r != nil {
				errorMessage, ok := r.(string)
				if ok {
					t.Logf(mensajePanicCapturado, errorMessage)
					if !strings.Contains(errorMessage, "Rollback de dependencia") {
						t.Errorf("Se esperaba un mensaje de error relacionado con el rollback, pero se obtuvo: %v", errorMessage)
					}
				} else {
					t.Errorf(errorPanicTipoString, r)
				}
			} else {
				t.Errorf(errorPanicFaltante)
			}
		}()

		services.RollbackDependenciaCreada(transaccion)
	})
}

func TestRollbackDependenciaTipoDependenciaCreada(t *testing.T) {
	t.Log(separadorSlash)
	t.Log("Inicio TestRollbackDependenciaTipoDependenciaCreada")
	t.Log(separadorSlash)

	t.Run("Caso 1: Rollback exitoso de tipo de dependencia", func(t *testing.T) {
		transaccion := &models.Creaciones{
			DependenciaTipoDependenciaId: []int{1, 2},
		}

		monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
			if method == "DELETE" && (strings.Contains(url, "dependencia_tipo_dependencia/1") || strings.Contains(url, "dependencia_tipo_dependencia/2")) {
				return nil
			} else if strings.Contains(url, "dependencia/") {
				return nil
			}
			return errors.New(errorUrlInesperada) // Para URLs no reconocidas
		})
		defer monkey.UnpatchAll()

		outputError := services.RollbackDependenciaTipoDependenciaCreada(transaccion)

		if outputError != nil {
			t.Errorf(errorOutputEsperadoNulo, outputError)
		}

		t.Log("Rollback de tipo de dependencia ejecutado exitosamente sin errores")
	})

	t.Run("Caso 2: Error en el rollback de tipo de dependencia", func(t *testing.T) {
		transaccion := &models.Creaciones{
			DependenciaTipoDependenciaId: []int{1},
		}

		monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
			if method == "DELETE" && strings.Contains(url, "dependencia_tipo_dependencia/1") {
				return errors.New("Error al eliminar tipo de dependencia")
			}
			return errors.New(errorUrlInesperada)
		})
		defer monkey.UnpatchAll()

		defer func() {
			if r := recover(); r != nil {
				errorMessage, ok := r.(string)
				if ok {
					t.Logf(mensajePanicCapturado, errorMessage)
					if !strings.Contains(errorMessage, "Rollback de dependencia tipo dependencia") {
						t.Errorf("Se esperaba un mensaje de error relacionado con rollback, pero se obtuvo: %v", errorMessage)
					}
				} else {
					t.Errorf(errorPanicTipoString, r)
				}
			} else {
				t.Errorf(errorPanicFaltante)
			}
		}()

		services.RollbackDependenciaTipoDependenciaCreada(transaccion)
	})
}

func TestRollbackDependenciaPadreCreada(t *testing.T) {
	t.Log(separadorSlash)
	t.Log("Inicio TestRollbackDependenciaPadreCreada")
	t.Log(separadorSlash)

	t.Run("Caso 1: Rollback exitoso de dependencia padre", func(t *testing.T) {
		transaccion := &models.Creaciones{
			DependenciaPadreId:           1,
			DependenciaTipoDependenciaId: []int{2, 3},
		}

		monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
			if method == "DELETE" && url == beego.AppConfig.String("OikosCrudUrl")+"dependencia_padre/1" {
				return nil
			}
			return errors.New(errorUrlInesperada)
		})
		defer monkey.Unpatch(request.SendJson)

		monkey.Patch(services.RollbackDependenciaTipoDependenciaCreada, func(transaccion *models.Creaciones) map[string]interface{} {
			return nil
		})
		defer monkey.Unpatch(services.RollbackDependenciaTipoDependenciaCreada)

		outputError := services.RollbackDependenciaPadreCreada(transaccion)

		if outputError != nil {
			t.Errorf(errorOutputEsperadoNulo, outputError)
		}

		t.Log("Rollback de dependencia padre ejecutado exitosamente sin errores")
	})

	t.Run("Caso 2: Error en el rollback de dependencia padre", func(t *testing.T) {
		transaccion := &models.Creaciones{
			DependenciaPadreId: 1,
		}

		monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
			if method == "DELETE" && url == beego.AppConfig.String("OikosCrudUrl")+"dependencia_padre/1" {
				return errors.New("Error al eliminar dependencia padre")
			}
			return errors.New(errorUrlInesperada)
		})
		defer monkey.Unpatch(request.SendJson)

		monkey.Patch(services.RollbackDependenciaTipoDependenciaCreada, func(transaccion *models.Creaciones) map[string]interface{} {
			return nil
		})
		defer monkey.Unpatch(services.RollbackDependenciaTipoDependenciaCreada)

		defer func() {
			if r := recover(); r != nil {
				errorMessage, ok := r.(string)
				if ok {
					t.Logf(mensajePanicCapturado, errorMessage)
					if !strings.Contains(errorMessage, "Rollback de dependencia padre") {
						t.Errorf("Se esperaba un mensaje de error relacionado con rollback, pero se obtuvo: %v", errorMessage)
					}
				} else {
					t.Errorf(errorPanicTipoString, r)
				}
			} else {
				t.Errorf(errorPanicFaltante)
			}
		}()

		services.RollbackDependenciaPadreCreada(transaccion)
	})
}
