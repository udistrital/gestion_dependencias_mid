package services_test

import (
    "testing"
    "net/http"
    "strconv"
    "github.com/jarcoal/httpmock"
    "github.com/udistrital/gestion_dependencias_mid/services"
	"github.com/udistrital/gestion_dependencias_mid/models"
    "github.com/astaxie/beego"
)
//ver que responda como se espera, como datos de prueba y que falle como se espera

func TestExisteDependencia(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio TestExisteDependencia")
    t.Log("//////////////////////////////////")

    t.Run("Caso 1: La dependencia existe", func(t *testing.T) {
        dependencias := []models.RespuestaBusquedaDependencia{
            {Dependencia: &models.Dependencia{Id: 1}},
        }
        if result := services.ExisteDependencia(dependencias, 1); result {
			t.Log("La dependencia con ID 1 existe")
        } else {
            t.Errorf("No existe la dependencia con ID 1")
        }
    })
}

func TestCrearRespuestaBusqueda(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio TestCrearRespuestaBusqueda")
    t.Log("//////////////////////////////////")

    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    t.Run("Caso 1: La dependencia tiene los datos son completos", func(t *testing.T) {
        dependencia := models.Dependencia{
            Id:     1,
            Activo: true,
            DependenciaTipoDependencia: []*models.DependenciaTipoDependencia{
                {
                    Activo: true,
                    TipoDependenciaId: &models.TipoDependencia{
                        Id:     1,
                        Nombre: "Tipo 1",
                    },
                },
            },
        }

        urlDependencia := beego.AppConfig.String("OikosCrudUrl") + "dependencia?query=Id:" + strconv.Itoa(dependencia.Id)
        httpmock.RegisterResponder("GET", urlDependencia,
            httpmock.NewJsonResponderOrPanic(http.StatusOK, []models.Dependencia{dependencia}),
        )

        urlDependenciaPadre := beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre?query=HijaId:" + strconv.Itoa(dependencia.Id)
        httpmock.RegisterResponder("GET", urlDependenciaPadre,
            httpmock.NewJsonResponderOrPanic(http.StatusOK, []models.DependenciaPadre{
                {
                    Id:     1,
                    PadreId: &dependencia,
                    Activo: true,
                },
            }),
        )
        t.Log("111111111111111111")
        resultado := services.CrearRespuestaBusqueda(dependencia)
        t.Log("22222222222222222222222")
        /*if resultado.Estado != true {
            t.Errorf("Se esperaba que el estado fuera true, pero se obtuvo %v", resultado.Estado)
        } 

        if resultado.TipoDependencia == nil || len(*resultado.TipoDependencia) == 0 {
            t.Errorf("Se esperaba que el tipo de dependencia no fuera nil o vacío")
        } 

        if resultado.DependenciaAsociada == nil || resultado.DependenciaAsociada.Id != 1 {
            t.Errorf("Se esperaba que la dependencia asociada tuviera el ID 1, pero se obtuvo %v", resultado.DependenciaAsociada)
        } */

        if resultado.Dependencia == nil || resultado.Dependencia.Id != 1 {
            t.Errorf("Se esperaba que el ID de la dependencia fuera 1, pero se obtuvo %v", resultado.Dependencia)
        } else {
            t.Log("Se encontro el ID de la dependencia con sus datos completos")
        }
        t.Log("termino caso 1")
    })
    t.Log("empieza caso 2")
    t.Run("Caso 2: Falta la dependencia o el tipo de dependencia", func(t *testing.T) {
        dependencia := models.Dependencia{
            Id:     2,
            Activo: true,
        }

        urlDependencia := beego.AppConfig.String("OikosCrudUrl") + "dependencia?query=Id:" + strconv.Itoa(dependencia.Id)
        httpmock.RegisterResponder("GET", urlDependencia,
            httpmock.NewJsonResponderOrPanic(http.StatusOK, []models.Dependencia{}), // Sin dependencias
        )

        urlDependenciaPadre := beego.AppConfig.String("OikosCrudUrl") + "dependencia_padre?query=HijaId:" + strconv.Itoa(dependencia.Id)
        httpmock.RegisterResponder("GET", urlDependenciaPadre,
            httpmock.NewJsonResponderOrPanic(http.StatusOK, []models.DependenciaPadre{}), // Sin dependencia padre
        )

        defer func() {
            if r := recover(); r == nil {
                t.Errorf("Se esperaba que la función panic cuando los datos no son completos, pero no ocurrió")
            }
        }()

        services.CrearRespuestaBusqueda(dependencia)
    })
}
