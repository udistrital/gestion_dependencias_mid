package services_test

import (
    "errors"
    "testing"
    "bou.ke/monkey"
    "github.com/udistrital/gestion_dependencias_mid/services"
	"github.com/udistrital/gestion_dependencias_mid/models"
    "github.com/udistrital/utils_oas/request"
    "fmt"
    "strings"
    "encoding/json"
    "strconv"
)

func TestBuscarDependencia(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio TestBuscarDependencia")
    t.Log("//////////////////////////////////")

	t.Run("Caso 1: Todos los datos están completos", func(t *testing.T) {
		transaccion := &models.BusquedaDependencia{
			NombreDependencia: "Dependencia 1",
			TipoDependenciaId: 1,
			FacultadId:        2,
			VicerrectoriaId:   3,
			BusquedaEstado:    &models.Estado{Estado: true},
		}

		dependencias := []models.Dependencia{
			{Id: 1, Nombre: "Dependencia 1"},
			{Id: 2, Nombre: "Dependencia 2"},
		}

		dependenciaTipoDependencia := []models.DependenciaTipoDependencia{
			{DependenciaId: &dependencias[1]}, 
		}

		monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			if strings.Contains(url, "dependencia?query=Nombre:") {
				data, _ := json.Marshal(dependencias)
				_ = json.Unmarshal(data, target)
			} else if strings.Contains(url, "dependencia_tipo_dependencia?query=TipoDependenciaId") {
				data, _ := json.Marshal(dependenciaTipoDependencia)
				_ = json.Unmarshal(data, target)
			} else if strings.Contains(url, "dependencia?query=Id:") {
				id := strings.Split(url, ":")[1]
				idInt, _ := strconv.Atoi(id)
				if idInt == 2 {
					data, _ := json.Marshal([]models.Dependencia{dependencias[1]})
					_ = json.Unmarshal(data, target)
				}
			} else if strings.Contains(url, "dependencia?query=Activo:") {
				data, _ := json.Marshal(dependencias)
				_ = json.Unmarshal(data, target)
			} else {
				return fmt.Errorf("URL no válida")
			}
			return nil
		})
		defer monkey.Unpatch(request.GetJson)

		monkey.Patch(services.CrearRespuestaBusqueda, func(dependencia models.Dependencia) models.RespuestaBusquedaDependencia {
			return models.RespuestaBusquedaDependencia{
				Dependencia: &dependencia, 
			}
		})
		defer monkey.Unpatch(services.CrearRespuestaBusqueda)

		monkey.Patch(services.ExisteDependencia, func(resultadoBusqueda []models.RespuestaBusquedaDependencia, id int) bool {
			for _, res := range resultadoBusqueda {
				if res.Dependencia.Id == id {
					return true
				}
			}
			return false
		})
		defer monkey.Unpatch(services.ExisteDependencia)

		resultadoBusqueda, outputError := services.BuscarDependencia(transaccion)

		if outputError != nil {
            t.Fatalf("No se esperaba un error, pero se obtuvo: %v", outputError)
        }

        if len(resultadoBusqueda) != len(dependencias) {
            t.Errorf("Se esperaban %d resultados, pero se obtuvieron %d", len(dependencias), len(resultadoBusqueda))
        }

        // Verificación adicional de los campos de los resultados
        for i, res := range resultadoBusqueda {
            if res.Dependencia == nil || res.Dependencia.Id != dependencias[i].Id || res.Dependencia.Nombre != dependencias[i].Nombre {
                t.Errorf("El resultado %d no es el esperado. Obtenido: %+v, Esperado: %+v", i, res.Dependencia, dependencias[i])
            }
        }

		t.Log("Test de caso exitoso para BuscarDependencia ejecutado correctamente")
	})
}

func TestCrearRespuestaBusqueda(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio TestCrearRespuestaBusqueda")
    t.Log("//////////////////////////////////")

    t.Run("Caso exitoso", func(t *testing.T) {
        // Generación dinámica de datos de prueba
        dependenciaID := 223
        nombreDependencia := "Dependencia de prueba"

        // Mock para el caso exitoso
        monkey.Patch(request.GetJson, func(url string, target interface{}) error {
            if strings.Contains(url, fmt.Sprintf("dependencia?query=Id:%d", dependenciaID)) {
                *target.(*[]models.Dependencia) = []models.Dependencia{
                    {
                        Id:                      dependenciaID,
                        Nombre:                  nombreDependencia,
                        TelefonoDependencia:     "321654987",
                        CorreoElectronico:       "info@desarrollo.com",
                        Activo:                  true,
                        FechaCreacion:           "2023-01-10T09:00:00Z",
                        FechaModificacion:       "2023-01-10T09:00:00Z",
                        DependenciaTipoDependencia: []*models.DependenciaTipoDependencia{
                            {
                                Activo: true,
                                TipoDependenciaId: &models.TipoDependencia{
                                    Id:     1,
                                    Nombre: "Tipo 1",
                                },
                            },
                        },
                    },
                }
                return nil
            } else if strings.Contains(url, fmt.Sprintf("dependencia_padre?query=HijaId:%d", dependenciaID)) {
                *target.(*[]models.DependenciaPadre) = []models.DependenciaPadre{
                    {
                        Id:              1,
                        PadreId:         &models.Dependencia{Id: 2, Nombre: "Dependencia Padre", Activo: true},
                        HijaId:          &models.Dependencia{Id: dependenciaID},
                        Activo:          true,
                        FechaCreacion:   "2024-01-01T00:00:00Z",
                        FechaModificacion: "2024-01-02T00:00:00Z",
                    },
                }
                return nil
            }
            return errors.New("URL no esperada")
        })
        defer monkey.UnpatchAll()

        // Datos para el caso exitoso
        dependencia := models.Dependencia{
            Id:                      dependenciaID,
            Nombre:                  nombreDependencia,
            TelefonoDependencia:     "321654987",
            CorreoElectronico:       "info@desarrollo.com",
            Activo:                  true,
            DependenciaTipoDependencia: []*models.DependenciaTipoDependencia{},
        }

        resultado := services.CrearRespuestaBusqueda(dependencia)

        // Validaciones para el caso exitoso
        if resultado.Dependencia.Id != dependenciaID {
            t.Errorf("Se esperaba Id %d, pero se obtuvo %d", dependenciaID, resultado.Dependencia.Id)
        }
        if resultado.Dependencia.Nombre != nombreDependencia {
            t.Errorf("Se esperaba Nombre '%s', pero se obtuvo '%s'", nombreDependencia, resultado.Dependencia.Nombre)
        }
        if resultado.DependenciaAsociada == nil || resultado.DependenciaAsociada.Id != 2 {
            t.Errorf("Se esperaba una DependenciaAsociada con Id 2, pero se obtuvo %+v", resultado.DependenciaAsociada)
        }
    })

    t.Run("Caso fallido", func(t *testing.T) {
        // Generación dinámica de datos de prueba
        dependenciaID := 999
        nombreDependencia := "Dependencia inexistente"

        monkey.Patch(request.GetJson, func(url string, target interface{}) error {
            if strings.Contains(url, fmt.Sprintf("dependencia?query=Id:%d", dependenciaID)) {
                *target.(*[]models.Dependencia) = []models.Dependencia{
                    {
                        Id:                      dependenciaID,
                        Nombre:                  nombreDependencia,
                        TelefonoDependencia:     "000000000",
                        CorreoElectronico:       "inexistente@desarrollo.com",
                        Activo:                  true,
                        DependenciaTipoDependencia: []*models.DependenciaTipoDependencia{},
                    },
                }
                return nil
            } else if strings.Contains(url, fmt.Sprintf("dependencia_padre?query=HijaId:%d", dependenciaID)) {
                *target.(*[]models.DependenciaPadre) = []models.DependenciaPadre{}
                return nil
            }
            return errors.New("URL no esperada")
        })
        defer monkey.UnpatchAll()

        dependenciaFallida := models.Dependencia{
            Id:                      dependenciaID,
            Nombre:                  nombreDependencia,
            TelefonoDependencia:     "000000000",
            CorreoElectronico:       "inexistente@desarrollo.com",
            Activo:                  true,
            DependenciaTipoDependencia: []*models.DependenciaTipoDependencia{},
        }

        resultadoFallido := services.CrearRespuestaBusqueda(dependenciaFallida)

        if resultadoFallido.DependenciaAsociada != nil {
            t.Errorf("Se esperaba que no hubiera Dependencia Asociada, pero se obtuvo %+v", resultadoFallido.DependenciaAsociada)
        }

        if resultadoFallido.DependenciaAsociada == nil {
            t.Logf("El test falló porque no se encontró una dependencia hija asociada para el ID %d", dependenciaID)
        }
    })
}

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

func TestEditarDependencia(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio TestEditarDependencia")
    t.Log("//////////////////////////////////")

    transaccion := &models.EditarDependencia{
        DependenciaId:        1,
        DependenciaAsociadaId: 2,
        Nombre:               "Dependencia Editada",
        CorreoElectronico:    "test@example.com",
        TelefonoDependencia:  "123456789",
        TipoDependenciaId:    []int{1, 2},
    }

    t.Run("Caso 1: La dependencia editada", func(t *testing.T) {

        monkey.Patch(request.GetJson, func(url string, target interface{}) error {
            fmt.Println("URL recibida en el mock:", url)

            if strings.Contains(url, "dependencia/1") {
                dependencia := models.Dependencia{
                    Id:                      1,
                    Nombre:                  "Dependencia Editada",
                    TelefonoDependencia:     "123456789",
                    CorreoElectronico:       "test@example.com",
                    Activo:                  true,
                    FechaCreacion:           "2023-01-10T09:00:00Z",
                    FechaModificacion:       "2023-01-10T09:00:00Z",
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
                *target.(*models.Dependencia) = dependencia
                return nil
            } else if strings.Contains(url, "dependencia_tipo_dependencia?query=dependenciaId:1") {
                *target.(*[]models.DependenciaTipoDependencia) = []models.DependenciaTipoDependencia{
                    {
                        Id: 1,
                        TipoDependenciaId: &models.TipoDependencia{
                            Id:     1,
                            Nombre: "Tipo 1",
                        },
                        DependenciaId: &models.Dependencia{
                            Id: 1,
                            Nombre: "Dependencia Editada",
                            Activo: true,
                        },
                        Activo:          true,
                        FechaCreacion:   "2023-01-10T09:00:00Z",
                        FechaModificacion: "2023-01-10T09:00:00Z",
                    },
                    {
                        Id: 2,
                        TipoDependenciaId: &models.TipoDependencia{
                            Id:     2,
                            Nombre: "Tipo 2",
                        },
                        DependenciaId: &models.Dependencia{
                            Id: 1,
                            Nombre: "Dependencia Editada",
                            Activo: true,
                        },
                        Activo:          true,
                        FechaCreacion:   "2023-01-10T09:00:00Z",
                        FechaModificacion: "2023-01-10T09:00:00Z",
                    },
                }
                return nil
            } else if strings.Contains(url, "dependencia_padre?query=HijaId:1") {
                *target.(*[]models.DependenciaPadre) = []models.DependenciaPadre{
                    {
                        Id:              1,
                        PadreId:         &models.Dependencia{Id: 2, Nombre: "Dependencia Padre", Activo: true},
                        HijaId:          &models.Dependencia{Id: 1, Nombre: "Dependencia Editada", Activo: true},
                        Activo:          true,
                        FechaCreacion:   "2024-01-01T00:00:00Z",
                        FechaModificacion: "2024-01-02T00:00:00Z",
                    },
                }
                return nil
            } else if strings.Contains(url, "dependencia/2") {
                dependenciaAsociada := models.Dependencia{
                    Id:                      2,
                    Nombre:                  "Dependencia Asociada",
                    TelefonoDependencia:     "987654321",
                    CorreoElectronico:       "associate@example.com",
                    Activo:                  true,
                    FechaCreacion:           "2023-01-15T09:00:00Z",
                    FechaModificacion:       "2023-01-15T09:00:00Z",
                }
                *target.(*models.Dependencia) = dependenciaAsociada
                return nil
            }

            return errors.New("URL no esperada")
        })

        monkey.Patch(request.SendJson, func(url, method string, response interface{}, data interface{}) error {
            fmt.Println("URL recibida en el mock de SendJson:", url)

            if strings.Contains(url, "dependencia/1") && method == "PUT" {
                respuestaDependencia := map[string]interface{}{
                    "success": true,
                    "message": "Dependencia actualizada correctamente",
                }

                if responseMap, ok := response.(*map[string]interface{}); ok {
                    *responseMap = respuestaDependencia
                }

                return nil 
            } else if strings.Contains(url, "dependencia_padre/1") && method == "PUT" {
                respuestaPadre := map[string]interface{}{
                    "Id": 1,
                    "success": true,
                    "message": "Dependencia padre actualizada correctamente",
                }

                if responseMap, ok := response.(*map[string]interface{}); ok {
                    *responseMap = respuestaPadre
                }

                return nil 
            }

            return errors.New("la URL o el método no son los esperados")
        })


        defer monkey.UnpatchAll()

        alerta, outputError := services.EditarDependencia(transaccion)

        if len(alerta) == 0 || alerta[0] != "Success" {
            t.Errorf("Se esperaba alerta de éxito, pero se obtuvo: %v", alerta)
        }
        if outputError != nil {
            t.Errorf("Se esperaba ningún error, pero se obtuvo: %v", outputError)
        }
    })
}

func TestNuevoTipoDependencia(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio TestNuevoTipoDependencia")
    t.Log("//////////////////////////////////")

    t.Run("Caso 1: Creación exitosa de nuevo tipo de dependencia", func(t *testing.T) {
        tipo := 1
        dependenciaId := models.Dependencia{
            Id:                  1,
            Nombre:              "Dependencia Prueba",
            TelefonoDependencia: "123456789",
            CorreoElectronico:   "prueba@example.com",
            Activo:              true,
            FechaCreacion:       "2023-01-10T09:00:00Z",
            FechaModificacion:   "2023-01-10T09:00:00Z",
        }
        tiposRegistrados := []int{2, 3}
        dependenciaOriginal := models.Dependencia{
            Id:                  1,
            Nombre:              "Dependencia Original",
            TelefonoDependencia: "123456789",
            CorreoElectronico:   "original@example.com",
            Activo:              true,
        }
    
        tipoDependencia := models.TipoDependencia{
            Id:     1,
            Nombre: "Tipo 1",
        }
    
        monkey.Patch(request.GetJson, func(url string, target interface{}) error {
            if strings.Contains(url, "tipo_dependencia/1") {
                *target.(*models.TipoDependencia) = tipoDependencia
                return nil
            }
            return errors.New("URL no esperada")
        })
        defer monkey.UnpatchAll()
    
        monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
            if strings.Contains(url, "dependencia_tipo_dependencia") && method == "POST" {
                *target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(4)} // Cambiar a float64
                return nil
            }
            return errors.New("Error en el envío de JSON")
        })
        defer monkey.UnpatchAll()
    
        monkey.Patch(services.RollbackDependenciaTipoDependencia, func(dependencia models.Dependencia, tiposRegistrados *[]int) {
            t.Errorf("No se esperaba que se llamara a RollbackDependenciaTipoDependencia")
        })
        defer monkey.UnpatchAll()
    
        defer func() {
            if r := recover(); r != nil {
                t.Errorf("Se esperaba que no se generara un panic, pero se obtuvo: %v", r)
            }
        }()
    
        services.NuevoTipoDependencia(tipo, dependenciaId, &tiposRegistrados, dependenciaOriginal)
    
        tipoAgregado := false
        for _, v := range tiposRegistrados {
            if v == 4 {
                tipoAgregado = true
                break
            }
        }
        if !tipoAgregado {
            t.Errorf("Se esperaba que el tipo 4 se agregara a la lista de tipos registrados")
        }
    
        t.Log("Creación de nuevo tipo de dependencia ejecutada exitosamente sin errores")
    })
    

    t.Run("Caso 2: Error en la creación debido a un fallo en la obtención del tipo de dependencia", func(t *testing.T) {
        tipo := 1
        dependenciaId := models.Dependencia{
            Id:                  1,
            Nombre:              "Dependencia Prueba",
            TelefonoDependencia: "123456789",
            CorreoElectronico:   "prueba@example.com",
            Activo:              true,
        }
        tiposRegistrados := []int{2, 3}
        dependenciaOriginal := models.Dependencia{
            Id:                  1,
            Nombre:              "Dependencia Original",
            TelefonoDependencia: "123456789",
            CorreoElectronico:   "original@example.com",
            Activo:              true,
        }

        monkey.Patch(request.GetJson, func(url string, target interface{}) error {
            return errors.New("error al intentar obtener JSON")
        })
        defer monkey.UnpatchAll()

        rollbackCalled := false
        monkey.Patch(services.RollbackDependenciaTipoDependencia, func(dependencia models.Dependencia, tiposRegistrados *[]int) {
            rollbackCalled = true
        })
        defer monkey.UnpatchAll()

        defer func() {
            if r := recover(); r != nil {
                errorMessage, ok := r.(string)
                if ok {
                    t.Logf("Panic capturado correctamente: %v", errorMessage)
                    if !strings.Contains(errorMessage, "error al intentar obtener JSON") {
                        t.Errorf("Se esperaba un mensaje de error relacionado con la obtención de JSON, pero se obtuvo: %v", errorMessage)
                    }
                } else {
                    t.Errorf("Se esperaba un panic de tipo string, pero se obtuvo: %v", r)
                }
            } else {
                t.Errorf("Se esperaba un panic, pero no ocurrió")
            }
        }()

        services.NuevoTipoDependencia(tipo, dependenciaId, &tiposRegistrados, dependenciaOriginal)

        if !rollbackCalled {
            t.Errorf("Se esperaba que se llamara a RollbackDependenciaTipoDependencia")
        }
    })
}

func TestActualizarDependenciaTipoDependencia(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio TestActualizarDependenciaTipoDependencia")
    t.Log("//////////////////////////////////")

    t.Run("Caso 1: Actualización exitosa de tipo de dependencia", func(t *testing.T) {
        tipo := 1
        activo := true
        dependenciaId := 1

        tiposOriginales := []models.DependenciaTipoDependencia{}
        dependenciaOriginal := models.Dependencia{
            Id:                  1,
            Nombre:              "Dependencia Original",
            TelefonoDependencia: "123456789",
            CorreoElectronico:   "original@example.com",
            Activo:              true,
            FechaCreacion:       "2023-01-10T09:00:00Z",
            FechaModificacion:   "2023-01-10T09:00:00Z",
        }
        tiposRegistrados := []int{1, 2}

        dependenciaTipoDependenciaActual := []models.DependenciaTipoDependencia{
            {
                Id:     1,
                Activo: false,
                TipoDependenciaId: &models.TipoDependencia{
                    Id:     1,
                    Nombre: "Tipo 1",
                },
            },
        }

        monkey.Patch(request.GetJson, func(url string, target interface{}) error {
            if strings.Contains(url, "dependencia_tipo_dependencia?query=dependenciaId:1,tipoDependenciaId:1") {
                *target.(*[]models.DependenciaTipoDependencia) = dependenciaTipoDependenciaActual
                return nil
            }
            return errors.New("URL no esperada")
        })
        defer monkey.UnpatchAll()

        monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
            if strings.Contains(url, "dependencia_tipo_dependencia/1") && method == "PUT" {
                *target.(*map[string]interface{}) = map[string]interface{}{"Id": 1}
                return nil
            }
            return errors.New("Error en el envío de JSON")
        })
        defer monkey.UnpatchAll()

        monkey.Patch(services.RollbackActualizacionTipoDependencia, func(dependencia models.Dependencia, tiposRegistrados *[]int, tiposOriginales *[]models.DependenciaTipoDependencia) {
            t.Errorf("No se esperaba que se llamara a RollbackActualizacionTipoDependencia")
        })
        defer monkey.UnpatchAll()

        defer func() {
            if r := recover(); r != nil {
                t.Errorf("Se esperaba que no se generara un panic, pero se obtuvo: %v", r)
            }
        }()

        services.ActualizarDependenciaTipoDependencia(tipo, activo, dependenciaId, &tiposOriginales, dependenciaOriginal, &tiposRegistrados)

        t.Log("Actualización de tipo de dependencia ejecutada exitosamente sin errores")
    })

    t.Run("Caso 2: Error en la actualización por fallo en la obtención del tipo de dependencia", func(t *testing.T) {
        tipo := 1
        activo := true
        dependenciaId := 1

        tiposOriginales := []models.DependenciaTipoDependencia{}
        dependenciaOriginal := models.Dependencia{
            Id:                  1,
            Nombre:              "Dependencia Original",
            TelefonoDependencia: "123456789",
            CorreoElectronico:   "original@example.com",
            Activo:              true,
        }
        tiposRegistrados := []int{1, 2}

        monkey.Patch(request.GetJson, func(url string, target interface{}) error {
            return errors.New("error al intentar obtener JSON")
        })
        defer monkey.UnpatchAll()

        rollbackCalled := false
        monkey.Patch(services.RollbackActualizacionTipoDependencia, func(dependencia models.Dependencia, tiposRegistrados *[]int, tiposOriginales *[]models.DependenciaTipoDependencia) {
            rollbackCalled = true
        })
        defer monkey.UnpatchAll()

        defer func() {
            if r := recover(); r != nil {
                errorMessage, ok := r.(string)
                if ok {
                    t.Logf("Panic capturado correctamente: %v", errorMessage)
                    if !strings.Contains(errorMessage, "error al intentar obtener JSON") {
                        t.Errorf("Se esperaba un mensaje de error relacionado con la obtención de JSON, pero se obtuvo: %v", errorMessage)
                    }
                } else {
                    t.Errorf("Se esperaba un panic de tipo string, pero se obtuvo: %v", r)
                }
            } else {
                t.Errorf("Se esperaba un panic, pero no ocurrió")
            }
        }()

        services.ActualizarDependenciaTipoDependencia(tipo, activo, dependenciaId, &tiposOriginales, dependenciaOriginal, &tiposRegistrados)

        if !rollbackCalled {
            t.Errorf("Se esperaba que se llamara a RollbackActualizacionTipoDependencia")
        }
    })
}

func TestContiene(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio TestContiene")
    t.Log("//////////////////////////////////")

    t.Run("Caso 1: Valor presente en el slice", func(t *testing.T) {
        slice := []int{1, 2, 3, 4, 5}
        valor := 3
        if !services.Contiene(slice, valor) {
            t.Errorf("Se esperaba que el valor %d estuviera presente en el slice, pero no se encontró", valor)
        }
    })

    t.Run("Caso 2: Valor no presente en el slice", func(t *testing.T) {
        slice := []int{1, 2, 3, 4, 5}
        valor := 6
        if services.Contiene(slice, valor) {
            t.Errorf("Se esperaba que el valor %d no estuviera presente en el slice, pero se encontró", valor)
        }
    })

}

func TestRollbackActualizacionTipoDependencia(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio TestRollbackActualizacionTipoDependencia")
    t.Log("//////////////////////////////////")

    t.Run("Caso 1: Rollback exitoso de actualización de tipos de dependencia", func(t *testing.T) {
        dependencia := models.Dependencia{
            Id:                      1,
            Nombre:                  "Dependencia Editada",
            TelefonoDependencia:     "123456789",
            CorreoElectronico:       "test@example.com",
            Activo:                  true,
        }

        tiposRegistrados := []int{1, 2, 3}
        tiposOriginales := []models.DependenciaTipoDependencia{
            {
                Id:     1,
                Activo: true,
                TipoDependenciaId: &models.TipoDependencia{
                    Id:     1,
                    Nombre: "Tipo 1",
                },
            },
            {
                Id:     2,
                Activo: true,
                TipoDependenciaId: &models.TipoDependencia{
                    Id:     2,
                    Nombre: "Tipo 2",
                },
            },
        }

        monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
            if strings.Contains(url, "dependencia_tipo_dependencia/") {
                return nil
            }
            return errors.New("URL no esperada")
        })
        defer monkey.UnpatchAll()

        monkey.Patch(services.RollbackDependenciaTipoDependencia, func(dependencia models.Dependencia, tiposRegistrados *[]int) {
        })
        defer monkey.UnpatchAll()

        defer func() {
            if r := recover(); r != nil {
                t.Errorf("Se esperaba que no se generara un panic, pero se obtuvo: %v", r)
            }
        }()

        services.RollbackActualizacionTipoDependencia(dependencia, &tiposRegistrados, &tiposOriginales)

        t.Log("Rollback de actualización de tipos de dependencia ejecutado exitosamente sin errores")
    })

    t.Run("Caso 2: Error en el rollback de actualización por fallo en el envío de JSON", func(t *testing.T) {
        dependencia := models.Dependencia{
            Id:                      0,
            Nombre:                  "",
            TelefonoDependencia:     "",
            CorreoElectronico:       "",
            Activo:                  false,
        }

        tiposRegistrados := []int{1}
        tiposOriginales := []models.DependenciaTipoDependencia{
            {
                Id:     1,
                Activo: true,
                TipoDependenciaId: &models.TipoDependencia{
                    Id:     1,
                    Nombre: "Tipo 1",
                },
            },
        }

        monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
            return errors.New("error al intentar enviar JSON")
        })
        defer monkey.UnpatchAll()

        monkey.Patch(services.RollbackDependenciaTipoDependencia, func(dependencia models.Dependencia, tiposRegistrados *[]int) {
        })
        defer monkey.UnpatchAll()

        defer func() {
            if r := recover(); r != nil {
                errorMessage, ok := r.(string)
                if ok {
                    t.Logf("Panic capturado correctamente: %v", errorMessage)
                    if !strings.Contains(errorMessage, "error al intentar enviar JSON") {
                        t.Errorf("Se esperaba un mensaje de error relacionado con el envío de JSON, pero se obtuvo: %v", errorMessage)
                    }
                } else {
                    t.Errorf("Se esperaba un panic de tipo string, pero se obtuvo: %v", r)
                }
            } else {
                t.Errorf("Se esperaba un panic, pero no ocurrió")
            }
        }()

        services.RollbackActualizacionTipoDependencia(dependencia, &tiposRegistrados, &tiposOriginales)
    })
}

func TestRollbackDependenciaTipoDependencia(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio TestRollbackDependenciaTipoDependencia")
    t.Log("//////////////////////////////////")

    t.Run("Caso 1: Rollback exitoso de tipos de dependencia", func(t *testing.T) {
        dependencia := models.Dependencia{
            Id:                      1,
            Nombre:                  "Dependencia Editada",
            TelefonoDependencia:     "123456789",
            CorreoElectronico:       "test@example.com",
            Activo:                  true,
            FechaCreacion:           "2023-01-10T09:00:00Z",
            FechaModificacion:       "2023-01-10T09:00:00Z",
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

        tiposRegistrados := []int{1, 2, 3}

        monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
            if strings.Contains(url, "dependencia_tipo_dependencia/") {
                return nil 
            }
            return errors.New("URL no esperada")
        })
        defer monkey.UnpatchAll()

        monkey.Patch(services.RollbackDependenciaOriginal, func(dependencia models.Dependencia) map[string]interface{} {
            return nil 
        })
        defer monkey.UnpatchAll()

        defer func() {
            if r := recover(); r != nil {
                t.Errorf("Se esperaba que no se generara un panic, pero se obtuvo: %v", r)
            }
        }()
        
        services.RollbackDependenciaTipoDependencia(dependencia, &tiposRegistrados)

        t.Log("Rollback de tipos de dependencia ejecutado exitosamente sin errores")
        
    })

    t.Run("Caso 2: Error en el rollback por datos incompletos", func(t *testing.T) {
        dependencia := models.Dependencia{
            Id:                      0, 
            Nombre:                  "", 
            TelefonoDependencia:     "", 
            CorreoElectronico:       "", 
            Activo:                  false, 
            DependenciaTipoDependencia: nil, 
        }
    
        tiposRegistrados := []int{1} 
    
        monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
            return errors.New("error al intentar enviar JSON")
        })
        defer monkey.UnpatchAll()
    
        monkey.Patch(services.RollbackDependenciaOriginal, func(dependencia models.Dependencia) map[string]interface{} {
            return nil 
        })
        defer monkey.UnpatchAll()
    
        // Captura de panic
        defer func() {
            if r := recover(); r != nil {
                errorMessage, ok := r.(string)
                if ok {
                    t.Logf("Panic capturado correctamente: %v", errorMessage)
                    if !strings.Contains(errorMessage, "error al intentar enviar JSON") {
                        t.Errorf("Se esperaba un mensaje de error relacionado con el envío de JSON, pero se obtuvo: %v", errorMessage)
                    }
                } else {
                    t.Errorf("Se esperaba un panic de tipo string, pero se obtuvo: %v", r)
                }
            } else {
                t.Errorf("Se esperaba un panic, pero no ocurrió")
            }
        }()
        services.RollbackDependenciaTipoDependencia(dependencia, &tiposRegistrados)
    }) 
}

func TestRollbackDependenciaOriginal(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio TestRollbackDependenciaOriginal")
    t.Log("//////////////////////////////////")

    t.Run("Caso 1: Rollback exitoso", func(t *testing.T) {
        dependencia := models.Dependencia{
            Id:                      1,
            Nombre:                  "Dependencia Editada",
            TelefonoDependencia:     "123456789",
            CorreoElectronico:       "test@example.com",
            Activo:                  true,
            FechaCreacion:           "2023-01-10T09:00:00Z",
            FechaModificacion:       "2023-01-10T09:00:00Z",
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

        monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
            if strings.Contains(url, fmt.Sprintf("dependencia/%d", dependencia.Id)) {
                return nil
            }
            return errors.New("URL no esperada")
        })
        defer monkey.UnpatchAll()

        outputError := services.RollbackDependenciaOriginal(dependencia)

        if outputError != nil {
            t.Errorf("Se esperaba un error nil, pero se obtuvo %+v", outputError)
        } else {
            t.Log("Rollback ejecutado exitosamente sin errores")
        }
    })

    t.Run("Caso 2: Error al enviar la solicitud", func(t *testing.T) {
        dependencia := models.Dependencia{
            Id:                  1,
            Nombre:              "Dependencia Editada",
            TelefonoDependencia: "123456789",
            Activo:                  true,
            FechaCreacion:           "2023-01-10T09:00:00Z",
            FechaModificacion:       "2023-01-10T09:00:00Z",
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
    
        monkey.Patch(request.SendJson, func(url string, method string, target interface{}, body interface{}) error {
            if strings.Contains(url, fmt.Sprintf("dependencia/%d", dependencia.Id)) {
                return errors.New("error al enviar la solicitud: campo 'CorreoElectronico' faltante")
            }
            return errors.New("URL no esperada")
        })
        defer monkey.UnpatchAll()
    
        defer func() {
            if r := recover(); r != nil {
                expectedPanic := "error al enviar la solicitud: campo 'CorreoElectronico' faltante"
                if r.(string) == expectedPanic {
                    t.Logf("Panic capturado correctamente: %v", r)
                } else {
                    t.Logf("Se esperaba panic con mensaje: %v, pero se obtuvo: %v", expectedPanic, r)
                }
            } else {
                t.Logf("Se esperaba un panic, pero no ocurrió")
            }
        }()
        _ = services.RollbackDependenciaOriginal(dependencia)
    })
}

func TestOrganigramas(t *testing.T) {
    t.Log("//////////////////////////////////")
	t.Log("Inicio TestOrganigramas")
	t.Log("//////////////////////////////////")

	t.Run("Caso 1: Filtrar correctamente los hijos de la Rectoría", func(t *testing.T) {
        dependencias := []models.Dependencia{
            {Id: 1, Nombre: "Dependencia 1", DependenciaTipoDependencia: []*models.DependenciaTipoDependencia{
                {Activo: true, TipoDependenciaId: &models.TipoDependencia{Id: 1, Nombre: "Tipo 1"}},
            }},
            {Id: 2, Nombre: "Dependencia 2", DependenciaTipoDependencia: []*models.DependenciaTipoDependencia{
                {Activo: true, TipoDependenciaId: &models.TipoDependencia{Id: 2, Nombre: "Tipo 2"}},
            }},
            {Id: 3, Nombre: "Dependencia 3", DependenciaTipoDependencia: []*models.DependenciaTipoDependencia{
                {Activo: false, TipoDependenciaId: &models.TipoDependencia{Id: 1, Nombre: "Tipo 1"}},
            }},
        }

        dependenciasPadre := []models.DependenciaPadre{
            {PadreId: &models.Dependencia{Id: 1}, HijaId: &models.Dependencia{Id: 2}},
            {PadreId: &models.Dependencia{Id: 1}, HijaId: &models.Dependencia{Id: 3}},
        }

        monkey.Patch(request.GetJson, func(url string, target interface{}) error {
            if strings.Contains(url, "dependencia?limit=-1") {
                data, _ := json.Marshal(dependencias)
                _ = json.Unmarshal(data, target)
            } else if strings.Contains(url, "dependencia_padre?limit=-1") {
                data, _ := json.Marshal(dependenciasPadre)
                _ = json.Unmarshal(data, target)
            } else {
                return errors.New("URL no válida")
            }
            return nil
        })
        defer monkey.Unpatch(request.GetJson)

        monkey.Patch(services.CopiarOrganigrama, func(organigrama []*models.Organigrama) []*models.Organigrama {
            return organigrama 
        })
        defer monkey.Unpatch(services.CopiarOrganigrama)

        monkey.Patch(services.FiltrarOrganigrama, func(organigrama []*models.Organigrama, dependenciasPadre []models.DependenciaPadre) []*models.Organigrama {
            return organigrama 
        })
        defer monkey.Unpatch(services.FiltrarOrganigrama)

        monkey.Patch(services.PodarOrganigramaAcademico, func(organigrama []*models.Organigrama) []*models.Organigrama {
            return organigrama
        })
        defer monkey.Unpatch(services.PodarOrganigramaAcademico)

        monkey.Patch(services.PodarOrganigramaAdministrativo, func(organigrama []*models.Organigrama) []*models.Organigrama {
            return organigrama 
        })
        defer monkey.Unpatch(services.PodarOrganigramaAdministrativo)

        organigramas, outputError := services.Organigramas()

        if outputError != nil {
            t.Errorf("No se esperaba un error, pero se obtuvo: %v", outputError)
        }

        if len(organigramas.General) != 1 {
            t.Errorf("Se esperaba 1 raíz en el organigrama general, pero se obtuvieron: %d", len(organigramas.General))
        }
        if len(organigramas.General[0].Hijos) != 2 {
            t.Errorf("Se esperaban 2 hijos en la raíz del organigrama general, pero se obtuvieron: %d", len(organigramas.General[0].Hijos))
        }
        
        t.Log("Test de caso exitoso para Organigramas ejecutado correctamente")

    })

    t.Run("Caso 2: Datos incompletos para dependencias", func(t *testing.T) {
        dependencias := []models.Dependencia{
            {Id: 1, Nombre: "Dependencia 1", DependenciaTipoDependencia: []*models.DependenciaTipoDependencia{
                {Activo: true, TipoDependenciaId: &models.TipoDependencia{Id: 1, Nombre: "Tipo 1"}},
            }},
            {Id: 2, Nombre: "Dependencia 2", DependenciaTipoDependencia: []*models.DependenciaTipoDependencia{
                {Activo: true, TipoDependenciaId: &models.TipoDependencia{Id: 2, Nombre: "Tipo 2"}},
            }},
            {Id: 3, Nombre: "", DependenciaTipoDependencia: []*models.DependenciaTipoDependencia{
                {Activo: false, TipoDependenciaId: nil},
            }},
        }

        dependenciasPadre := []models.DependenciaPadre{
            {PadreId: &models.Dependencia{Id: 1}, HijaId: &models.Dependencia{Id: 2}},
            {PadreId: nil, HijaId: &models.Dependencia{Id: 3}},
        }

        monkey.Patch(request.GetJson, func(url string, target interface{}) error {
            if strings.Contains(url, "dependencia?limit=-1") {
                data, _ := json.Marshal(dependencias)
                _ = json.Unmarshal(data, target)
            } else if strings.Contains(url, "dependencia_padre?limit=-1") {
                data, _ := json.Marshal(dependenciasPadre)
                _ = json.Unmarshal(data, target)
            } else {
                return errors.New("URL no válida")
            }
            return nil
        })
        defer monkey.Unpatch(request.GetJson)

        monkey.Patch(services.CopiarOrganigrama, func(organigrama []*models.Organigrama) []*models.Organigrama {
            return organigrama
        })
        defer monkey.Unpatch(services.CopiarOrganigrama)

        monkey.Patch(services.FiltrarOrganigrama, func(organigrama []*models.Organigrama, dependenciasPadre []models.DependenciaPadre) []*models.Organigrama {
            return organigrama
        })
        defer monkey.Unpatch(services.FiltrarOrganigrama)

        monkey.Patch(services.PodarOrganigramaAcademico, func(organigrama []*models.Organigrama) []*models.Organigrama {
            return organigrama
        })
        defer monkey.Unpatch(services.PodarOrganigramaAcademico)

        monkey.Patch(services.PodarOrganigramaAdministrativo, func(organigrama []*models.Organigrama) []*models.Organigrama {
            return organigrama
        })
        defer monkey.Unpatch(services.PodarOrganigramaAdministrativo)

        defer func() {
            if r := recover(); r != nil {
                t.Logf("Se capturó el pánico esperado debido a los datos incompletos: %v", r)
            }
        }()

        organigramas, outputError := services.Organigramas()

        if outputError == nil {
            t.Errorf("Se esperaba un error debido a los datos incompletos, pero no se obtuvo ninguno")
        } else {
            t.Logf("Se obtuvo el error esperado: %v", outputError)
        }

        if organigramas.General != nil && len(organigramas.General) > 0 {
            t.Errorf("No se esperaba una estructura válida de organigramas, pero se obtuvo: %v", organigramas.General)
        }

        t.Log("Test de caso fallido para Organigramas ejecutado correctamente")
    })
    
}

func TestCopiarOrganigrama(t *testing.T) {
	t.Log("//////////////////////////////////")
	t.Log("Inicio TestCopiarOrganigrama")
	t.Log("//////////////////////////////////")

	dependencia1 := models.Dependencia{Id: 1, Nombre: "Dependencia 1"}
	dependencia2 := models.Dependencia{Id: 2, Nombre: "Dependencia 2"}
	dependencia3 := models.Dependencia{Id: 3, Nombre: "Dependencia 3"}

	organigramaOriginal := []*models.Organigrama{
		{
			Dependencia: dependencia1,
			Tipo:        []string{"Tipo A"}, 
			Hijos: []*models.Organigrama{
				{
					Dependencia: dependencia2,
					Tipo:        []string{"Tipo B"}, 
					Hijos:       []*models.Organigrama{},
				},
			},
		},
		{
			Dependencia: dependencia3,
			Tipo:        []string{"Tipo C"}, 
			Hijos:       []*models.Organigrama{},
		},
	}

	t.Run("Caso 1: Copiar un organigrama con múltiples niveles", func(t *testing.T) {
		resultado := services.CopiarOrganigrama(organigramaOriginal)

		if len(resultado) != len(organigramaOriginal) {
			t.Errorf("Se esperaban %d elementos, pero se obtuvieron: %d", len(organigramaOriginal), len(resultado))
		}

		for i, org := range resultado {
			if org.Dependencia.Id != organigramaOriginal[i].Dependencia.Id { 
				t.Errorf("Se esperaban %v, pero se obtuvieron %v en el índice %d", organigramaOriginal[i].Dependencia.Id, org.Dependencia.Id, i)
			}
			
			if len(org.Tipo) != len(organigramaOriginal[i].Tipo) {
				t.Errorf("El tamaño de los tipos no coincide en el índice %d", i)
			} else {
				for j, tipo := range org.Tipo {
					if tipo != organigramaOriginal[i].Tipo[j] {
						t.Errorf("Se esperaban el tipo %s, pero se obtuvieron %s en el índice %d", organigramaOriginal[i].Tipo[j], tipo, i)
					}
				}
			}
		}
	})

	t.Run("Caso 2: Copiar un organigrama vacío", func(t *testing.T) {
		organigramaVacio := []*models.Organigrama{}

		resultado := services.CopiarOrganigrama(organigramaVacio)

		if len(resultado) != 0 {
			t.Errorf("Se esperaba un organigrama vacío, pero se obtuvieron: %d", len(resultado))
		}
	})
}

func TestFiltrarOrganigrama(t *testing.T) {
	t.Log("//////////////////////////////////")
	t.Log("Inicio TestFiltrarOrganigrama")
	t.Log("//////////////////////////////////")

	dependencia1 := &models.Dependencia{Id: 1, Nombre: "Dependencia 1"}
	dependencia2 := &models.Dependencia{Id: 2, Nombre: "Dependencia 2"}
	dependencia3 := &models.Dependencia{Id: 3, Nombre: "Dependencia 3"}
	dependencia4 := &models.Dependencia{Id: 4, Nombre: "Dependencia 4"}

	dependencias_padre := []models.DependenciaPadre{
		{HijaId: dependencia1},
		{HijaId: dependencia2},
	}

	organigramaCumpleCriterios := []*models.Organigrama{
		{
			Dependencia: *dependencia1,
			Hijos: []*models.Organigrama{
				{
					Dependencia: *dependencia3,
					Hijos: []*models.Organigrama{},
				},
			},
		},
		{
			Dependencia: *dependencia2,
			Hijos: []*models.Organigrama{},
		},
		{
			Dependencia: *dependencia4,
			Hijos: []*models.Organigrama{
				{
					Dependencia: *dependencia1,
					Hijos: []*models.Organigrama{},
				},
			},
		},
	}

	t.Run("Caso 1: Filtrar organigrama con hijos y/o dependencias padre", func(t *testing.T) {
		monkey.Patch(services.TienePadre, func(nodo *models.Organigrama, dependencias_padre []models.DependenciaPadre) bool {
			return true 
		})
		defer monkey.UnpatchAll()

		resultado := services.FiltrarOrganigrama(organigramaCumpleCriterios, dependencias_padre)

		if len(resultado) != 3 {
			t.Errorf("Se esperaban 3 elementos, pero se obtuvieron: %d", len(resultado))
		}
	})

	organigramaSinCriterios := []*models.Organigrama{
		{
			Dependencia: *dependencia4,
			Hijos:       []*models.Organigrama{},
		},
	}

	t.Run("Caso 2: Filtrar organigrama sin hijos y sin dependencias padre", func(t *testing.T) {
		monkey.Patch(services.TienePadre, func(nodo *models.Organigrama, dependencias_padre []models.DependenciaPadre) bool {
			return false 
		})
		defer monkey.UnpatchAll()

		resultado := services.FiltrarOrganigrama(organigramaSinCriterios, dependencias_padre)

		if len(resultado) != 0 {
			t.Errorf("Se esperaba un organigrama vacío, pero se obtuvieron: %d", len(resultado))
		}
	})
}

func TestTienePadre(t *testing.T) {
	t.Log("//////////////////////////////////")
	t.Log("Inicio TestTienePadre")
	t.Log("//////////////////////////////////")

	dependencia1 := &models.Dependencia{Id: 1, Nombre: "Dependencia 1"}
	dependencia2 := &models.Dependencia{Id: 2, Nombre: "Dependencia 2"}
	dependencia3 := &models.Dependencia{Id: 3, Nombre: "Dependencia 3"}

	t.Run("Caso 1: Nodo tiene padre", func(t *testing.T) {
		nodo := &models.Organigrama{
			Dependencia: *dependencia1,
		}

		dependencias_padre := []models.DependenciaPadre{
			{HijaId: dependencia1},
			{HijaId: dependencia2},
		}

		resultado := services.TienePadre(nodo, dependencias_padre)

		if !resultado {
			t.Errorf("Se esperaba que el nodo tuviera padre, pero no lo tiene")
		}
	})

	t.Run("Caso 2: Nodo no tiene padre", func(t *testing.T) {
		nodo := &models.Organigrama{
			Dependencia: *dependencia3,
		}

		dependencias_padre := []models.DependenciaPadre{
			{HijaId: dependencia1},
			{HijaId: dependencia2},
		}

		resultado := services.TienePadre(nodo, dependencias_padre)

		if resultado {
			t.Errorf("Se esperaba que el nodo no tuviera padre, pero lo tiene")
		}
	})
}

func TestPodarOrganigramaAcademico(t *testing.T) {
	t.Log("//////////////////////////////////")
	t.Log("Inicio TestPodarOrganigramaAcademico")
	t.Log("//////////////////////////////////")

	t.Run("Caso 1: Filtrar correctamente los hijos de la Rectoría", func(t *testing.T) {
		organigrama := []*models.Organigrama{
			{
				Dependencia: models.Dependencia{Nombre: "RECTORIA"},
				Hijos: []*models.Organigrama{
					{
						Dependencia: models.Dependencia{Nombre: "VICERRECTORIA ACADEMICA"},
					},
					{
						Dependencia: models.Dependencia{Nombre: "VICERRECTORIA ADMINISTRATIVA"},
					},
				},
			},
			{
				Dependencia: models.Dependencia{Nombre: "OTRA DEPENDENCIA"},
				Hijos: []*models.Organigrama{
					{
						Dependencia: models.Dependencia{Nombre: "VICERRECTORIA ACADEMICA"},
					},
				},
			},
		}

		resultado := services.PodarOrganigramaAcademico(organigrama)

		if len(resultado[0].Hijos) != 1 {
			t.Errorf("Se esperaban 1 hijo después de podar, pero se obtuvieron: %d", len(resultado[0].Hijos))
		}

		if resultado[0].Hijos[0].Dependencia.Nombre != "VICERRECTORIA ACADEMICA" {
			t.Errorf("Se esperaba que el hijo fuera 'VICERRECTORIA ACADEMICA', pero se obtuvo: %s", resultado[0].Hijos[0].Dependencia.Nombre)
		}

		if len(resultado[1].Hijos) != 1 {
			t.Errorf("Se esperaban 1 hijo en 'OTRA DEPENDENCIA', pero se obtuvieron: %d", len(resultado[1].Hijos))
		}
	})

	t.Run("Caso 2: Sin dependencias de tipo 'RECTORIA'", func(t *testing.T) {
		organigrama := []*models.Organigrama{
			{
				Dependencia: models.Dependencia{Nombre: "OTRA DEPENDENCIA"},
				Hijos: []*models.Organigrama{
					{
						Dependencia: models.Dependencia{Nombre: "VICERRECTORIA ACADEMICA"},
					},
				},
			},
		}

		resultado := services.PodarOrganigramaAcademico(organigrama)

		if len(resultado) != 1 || len(resultado[0].Hijos) != 1 {
			t.Errorf("Se esperaba que el organigrama no se modificara, pero se obtuvieron: %d elementos", len(resultado))
		}
	})
}

func TestPodarOrganigramaAdministrativo(t *testing.T) {
	t.Log("//////////////////////////////////")
	t.Log("Inicio TestPodarOrganigramaAdministrativo")
	t.Log("//////////////////////////////////")

	t.Run("Caso 1: Filtrar correctamente los hijos de la Rectoría", func(t *testing.T) {
		organigrama := []*models.Organigrama{
			{
				Dependencia: models.Dependencia{Nombre: "RECTORIA"},
				Hijos: []*models.Organigrama{
					{
						Dependencia: models.Dependencia{Nombre: "VICERRECTORIA ACADEMICA"},
					},
					{
						Dependencia: models.Dependencia{Nombre: "VICERRECTORIA ADMINISTRATIVA"},
					},
				},
			},
			{
				Dependencia: models.Dependencia{Nombre: "OTRA DEPENDENCIA"},
				Hijos: []*models.Organigrama{
					{
						Dependencia: models.Dependencia{Nombre: "VICERRECTORIA ACADEMICA"},
					},
				},
			},
		}

		resultado := services.PodarOrganigramaAdministrativo(organigrama)

		if len(resultado[0].Hijos) != 1 {
			t.Errorf("Se esperaban 1 hijo después de podar, pero se obtuvieron: %d", len(resultado[0].Hijos))
		}

		if resultado[0].Hijos[0].Dependencia.Nombre != "VICERRECTORIA ADMINISTRATIVA" {
			t.Errorf("Se esperaba que el hijo fuera 'VICERRECTORIA ADMINISTRATIVA', pero se obtuvo: %s", resultado[0].Hijos[0].Dependencia.Nombre)
		}

		if len(resultado[1].Hijos) != 1 {
			t.Errorf("Se esperaban 1 hijo en 'OTRA DEPENDENCIA', pero se obtuvieron: %d", len(resultado[1].Hijos))
		}
	})

	t.Run("Caso 2: Sin dependencias de tipo 'RECTORIA'", func(t *testing.T) {
		organigrama := []*models.Organigrama{
			{
				Dependencia: models.Dependencia{Nombre: "OTRA DEPENDENCIA"},
				Hijos: []*models.Organigrama{
					{
						Dependencia: models.Dependencia{Nombre: "VICERRECTORIA ACADEMICA"},
					},
				},
			},
		}

		resultado := services.PodarOrganigramaAdministrativo(organigrama)

		if len(resultado) != 1 || len(resultado[0].Hijos) != 1 {
			t.Errorf("Se esperaba que el organigrama no se modificara, pero se obtuvieron: %d elementos", len(resultado))
		}
	})
}
