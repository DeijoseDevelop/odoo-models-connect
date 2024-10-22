package odoo

import (
	"fmt"
	"sync"
)

func main() {
    client, err := NewOdooClient(".env")
    if err != nil {
        fmt.Println("Error al inicializar el cliente Odoo:", err)
        return
    }

    // Ejemplo de uso con goroutines para procesar muchos datos
    // Vamos a obtener muchos partners y procesarlos en paralelo

    // Paso 1: Obtener IDs de partners (por lotes si es necesario)
    domain := []interface{}{} // Sin filtros para obtener todos
    ids, err := client.Search("res.partner", domain)
    if err != nil {
        fmt.Println("Error al obtener IDs de partners:", err)
        return
    }

    // Para este ejemplo, limitamos a los primeros 1000 IDs
    if len(ids) > 1000 {
        ids = ids[:1000]
    }

    // Paso 2: Dividir los IDs en lotes para procesar con goroutines
    batchSize := 100
    var batches [][]int
    for batchSize < len(ids) {
        ids, batches = ids[batchSize:], append(batches, ids[0:batchSize:batchSize])
    }
    batches = append(batches, ids)

    // Paso 3: Procesar cada lote en una goroutine
    var wg sync.WaitGroup
    partnersChannel := make(chan ResPartner, len(ids))

    for _, batch := range batches {
        wg.Add(1)
        go func(batch []int) {
            defer wg.Done()
            // Obtener registros del lote
            domain := []interface{}{[]interface{}{"id", "in", batch}}
            fields := []string{"id", "name", "email", "is_company"}
            records, err := client.SearchRead("res.partner", domain, fields)
            if err != nil {
                fmt.Println("Error al obtener partners:", err)
                return
            }
            // Mapear y enviar a través del canal
            for _, record := range records {
                var partner ResPartner
                if err := MapToStruct(record, &partner); err != nil {
                    fmt.Println("Error al mapear partner:", err)
                    continue
                }
                partnersChannel <- partner
            }
        }(batch)
    }

    // Paso 4: Cerrar el canal cuando todas las goroutines terminen
    go func() {
        wg.Wait()
        close(partnersChannel)
    }()

    // Paso 5: Procesar los partners recibidos
    for partner := range partnersChannel {
        // Aquí puedes procesar cada partner como desees
        PrintModel(partner)
    }

    // Ejemplo similar para otro modelo, como product.product
    // Puedes repetir los pasos anteriores para procesar productos u otros datos
}
