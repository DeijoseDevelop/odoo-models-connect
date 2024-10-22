
# odoo-models-connect

Cliente en Go para conectarse a Odoo utilizando XML-RPC. Esta librería permite realizar operaciones CRUD, mapear datos a estructuras Go.

## Tabla de Contenidos

- [odoo-models-connect](#odoo-models-connect)
  - [Tabla de Contenidos](#tabla-de-contenidos)
  - [Características](#características)
  - [Requisitos](#requisitos)
  - [Instalación](#instalación)
  - [Configuración](#configuración)
  - [Uso](#uso)
    - [Inicializar el Cliente](#inicializar-el-cliente)
    - [Operaciones CRUD](#operaciones-crud)
      - [Crear un Registro](#crear-un-registro)
      - [Leer Registros](#leer-registros)
      - [Actualizar un Registro](#actualizar-un-registro)
      - [Eliminar un Registro](#eliminar-un-registro)
    - [Mapeo de Modelos](#mapeo-de-modelos)
      - [Definir una Estructura](#definir-una-estructura)
      - [Mapear Datos a la Estructura](#mapear-datos-a-la-estructura)
    - [Procesamiento en Paralelo con Goroutines](#procesamiento-en-paralelo-con-goroutines)
      - [Ejemplo de Procesamiento en Paralelo](#ejemplo-de-procesamiento-en-paralelo)
  - [Ejemplos](#ejemplos)
  - [Estructura del Proyecto](#estructura-del-proyecto)
  - [Notas Importantes](#notas-importantes)
  - [Contribuciones](#contribuciones)
  - [Licencia](#licencia)

## Características

- Conexión a Odoo mediante XML-RPC.
- Autenticación y manejo de sesiones.
- Operaciones CRUD (Crear, Leer, Actualizar, Eliminar).
- Mapeo de datos a estructuras Go.
- Procesamiento paralelo utilizando goroutines y canales.
- Manejo de errores personalizado.

## Requisitos

- Go 1.13 o superior.
- Acceso a una instancia de Odoo con las credenciales adecuadas.
- Las siguientes dependencias de Go:
  - `github.com/kolo/xmlrpc`
  - `github.com/joho/godotenv`
  - `github.com/mitchellh/mapstructure`

## Instalación

Clona este repositorio o copia los archivos necesarios en tu proyecto.

```bash
git clone github.com/DeijoseDevelop/odoo-models-connect.git
cd odoo-models-connect
```

Instala las dependencias:

```bash
go get github.com/kolo/xmlrpc
go get github.com/joho/godotenv
go get github.com/mitchellh/mapstructure
```

## Configuración

Crea un archivo `.env` en el directorio raíz del proyecto con las siguientes variables:

```
DATABASE=nombre_de_tu_base_de_datos
USERNAME=tu_usuario
PASSWORD=tu_contraseña
URL=http://tu_dominio_odoo.com
```

Asegúrate de reemplazar los valores con tus credenciales y URL de Odoo.

## Uso

### Inicializar el Cliente

Importa los paquetes necesarios y crea una instancia de `OdooClient`:

```go
package main

import (
    "fmt"
    "github.com/DeijoseDevelop/odoo-models-connect" // Importa tu paquete adecuadamente
)

func main() {
    client, err := odoo.NewOdooClient(".env")
    if err != nil {
        fmt.Println("Error al inicializar el cliente Odoo:", err)
        return
    }
    // Ahora puedes usar el cliente para realizar operaciones
}
```

### Operaciones CRUD

#### Crear un Registro

```go
data := map[string]interface{}{
    "name":  "Nuevo Cliente",
    "email": "cliente@example.com",
}
id, err := client.Create("res.partner", data)
if err != nil {
    fmt.Println("Error al crear el registro:", err)
    return
}
fmt.Println("Registro creado con ID:", id)
```

#### Leer Registros

```go
domain := []interface{}{[]interface{}{"is_company", "=", true}}
fields := []string{"id", "name", "email"}
records, err := client.SearchRead("res.partner", domain, fields)
if err != nil {
    fmt.Println("Error al leer los registros:", err)
    return
}
for _, record := range records {
    fmt.Println("Registro:", record)
}
```

#### Actualizar un Registro

```go
updateData := map[string]interface{}{
    "email": "nuevo.email@example.com",
}
success, err := client.Update("res.partner", []int{id}, updateData)
if err != nil || !success {
    fmt.Println("Error al actualizar el registro:", err)
    return
}
fmt.Println("Registro actualizado exitosamente")
```

#### Eliminar un Registro

```go
success, err := client.Delete("res.partner", []int{id})
if err != nil || !success {
    fmt.Println("Error al eliminar el registro:", err)
    return
}
fmt.Println("Registro eliminado exitosamente")
```

### Mapeo de Modelos

Define estructuras que representen tus modelos de Odoo y utiliza la función `MapToStruct` para mapear los datos.

#### Definir una Estructura

```go
type ResPartner struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    IsCompany bool   `json:"is_company"`
}
```

#### Mapear Datos a la Estructura

```go
for _, record := range records {
    var partner ResPartner
    if err := odoo.MapToStruct(record, &partner); err != nil {
        fmt.Println("Error al mapear el registro:", err)
        continue
    }
    fmt.Println("Partner:", partner)
}
```

### Procesamiento en Paralelo con Goroutines

Si necesitas procesar muchos datos, puedes utilizar goroutines y canales para mejorar el rendimiento.

#### Ejemplo de Procesamiento en Paralelo

```go
import (
    "sync"
)

func main() {
    // ... Inicializar cliente ...

    // Obtener IDs de los registros
    domain := []interface{}{}
    ids, err := client.Search("res.partner", domain)
    if err != nil {
        fmt.Println("Error al obtener IDs:", err)
        return
    }

    // Procesar en lotes
    batchSize := 100
    var batches [][]int
    for batchSize < len(ids) {
        ids, batches = ids[batchSize:], append(batches, ids[0:batchSize:batchSize])
    }
    batches = append(batches, ids)

    var wg sync.WaitGroup
    partnersChannel := make(chan ResPartner, len(ids))

    for _, batch := range batches {
        wg.Add(1)
        go func(batch []int) {
            defer wg.Done()
            domain := []interface{}{[]interface{}{"id", "in", batch}}
            fields := []string{"id", "name", "email", "is_company"}
            records, err := client.SearchRead("res.partner", domain, fields)
            if err != nil {
                fmt.Println("Error al obtener registros:", err)
                return
            }
            for _, record := range records {
                var partner ResPartner
                if err := odoo.MapToStruct(record, &partner); err != nil {
                    fmt.Println("Error al mapear el registro:", err)
                    continue
                }
                partnersChannel <- partner
            }
        }(batch)
    }

    go func() {
        wg.Wait()
        close(partnersChannel)
    }()

    for partner := range partnersChannel {
        fmt.Println("Partner:", partner)
    }
}
```

## Ejemplos

Consulta el archivo `example.go` para ver ejemplos completos de cómo utilizar la librería, incluyendo operaciones CRUD y procesamiento paralelo (si copias y pegas, agrega la palabra "odoo" antes de cada componente, ejemplo: odoo.NewOdooClient()).

## Estructura del Proyecto

```
odoo-models-connect/
├── client.go         // Contiene OdooClient y métodos
├── errors.go         // Tipos de error personalizados
├── models.go         // Definiciones de modelos y mapeadores
├── utils.go          // Funciones utilitarias
├── main.go           // Uso de ejemplo
├── .env              // Variables de entorno (no incluir en control de versiones)
└── README.md         // Instrucciones de uso (este archivo)
```

## Notas Importantes

- **Manejo de Errores**: Siempre maneja los errores devueltos por las funciones para evitar comportamientos inesperados.
- **Tipos de Datos**: Verifica los tipos de datos al mapear los resultados. Algunos campos pueden ser `int64`, `float64`, etc.
- **Campos Nulos**: Ten cuidado con los campos que pueden ser `nil` o no estar presentes en los datos retornados.
- **Concurrencia**: Al utilizar goroutines, asegúrate de manejar adecuadamente la sincronización y los recursos compartidos.
- **Variables de Entorno**: No incluyas el archivo `.env` en tu control de versiones ya que contiene información sensible.

## Contribuciones

Las contribuciones son bienvenidas. Si deseas mejorar la librería o agregar nuevas funcionalidades, por favor:

1. Haz un fork del proyecto.
2. Crea una rama para tu nueva funcionalidad (`git checkout -b nueva-funcionalidad`).
3. Realiza tus cambios y haz commit (`git commit -am 'Agregar nueva funcionalidad'`).
4. Envía tus cambios a GitHub (`git push origin nueva-funcionalidad`).
5. Crea un Pull Request.

## Licencia

Este proyecto está bajo la Licencia MIT. Consulta el archivo `LICENSE` para más detalles.
