package odoo

import (
    "fmt"
    "reflect"

    "github.com/mitchellh/mapstructure"
)

// Definición de modelos

type ResPartner struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    IsCompany bool   `json:"is_company"`
    // Agrega otros campos necesarios
}

type ProductProduct struct {
    ID           int     `json:"id"`
    Name         string  `json:"name"`
    Type         string  `json:"type"`
    SaleOk       bool    `json:"sale_ok"`
    PurchaseOk   bool    `json:"purchase_ok"`
    ListPrice    float64 `json:"list_price"`
    StandardPrice float64 `json:"standard_price"`
    // Agrega otros campos necesarios
}

type AccountMove struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    PartnerID   int     `json:"partner_id"`
    InvoiceDate string  `json:"invoice_date"`
    AmountTotal float64 `json:"amount_total"`
    State       string  `json:"state"`
    // Agrega otros campos necesarios
}

// Función genérica para mapear datos a estructuras
func MapToStruct(data map[string]interface{}, result interface{}) error {
    decoderConfig := &mapstructure.DecoderConfig{
        TagName:          "json",
        WeaklyTypedInput: true,
        Result:           result,
    }
    decoder, err := mapstructure.NewDecoder(decoderConfig)
    if err != nil {
        return err
    }
    return decoder.Decode(data)
}

// Funciones específicas pueden agregarse si se requiere lógica adicional

// Ejemplo de función para mostrar información de cualquier modelo
func PrintModel(model interface{}) {
    v := reflect.ValueOf(model)
    typeOfS := v.Type()
    fmt.Println("-----")
    for i := 0; i < v.NumField(); i++ {
        fmt.Printf("%s: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
    }
    fmt.Println("-----")
}
