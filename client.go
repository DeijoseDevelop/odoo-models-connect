package odoo

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
    "github.com/kolo/xmlrpc"
)

type OdooClient struct {
    Database string
    Username string
    Password string
    URL      string
    UID      int
    common   *xmlrpc.Client
    models   *xmlrpc.Client
}

// Inicializa el cliente Odoo y autentica
func NewOdooClient(envPath string) (*OdooClient, error) {
    // Cargar variables de entorno
    if err := godotenv.Load(envPath); err != nil {
        return nil, fmt.Errorf("Error al cargar el archivo .env: %v", err)
    }

    client := &OdooClient{
        Database: os.Getenv("DATABASE"),
        Username: os.Getenv("USERNAME"),
        Password: os.Getenv("PASSWORD"),
        URL:      os.Getenv("URL"),
    }

    if err := client.authenticate(); err != nil {
        return nil, err
    }

    if err := client.setupModels(); err != nil {
        return nil, err
    }

    return client, nil
}

func (c *OdooClient) authenticate() error {
    commonURL := fmt.Sprintf("%s/xmlrpc/2/common", c.URL)
    client, err := xmlrpc.NewClient(commonURL, nil)
    if err != nil {
        return fmt.Errorf("Error al crear cliente común: %v", err)
    }
    c.common = client

    var uid int
    err = c.common.Call("authenticate", []interface{}{c.Database, c.Username, c.Password, nil}, &uid)
    if err != nil {
        return fmt.Errorf("Error en la autenticación: %v", err)
    }
    c.UID = uid
    return nil
}

func (c *OdooClient) setupModels() error {
    modelsURL := fmt.Sprintf("%s/xmlrpc/2/object", c.URL)
    client, err := xmlrpc.NewClient(modelsURL, nil)
    if err != nil {
        return fmt.Errorf("Error al crear cliente de modelos: %v", err)
    }
    c.models = client
    return nil
}

// Método genérico para ejecutar llamadas XML-RPC con manejo de errores
func (c *OdooClient) execute(model, method string, args []interface{}, kwargs map[string]interface{}, result interface{}) error {
    var fullArgs []interface{}
    fullArgs = append(fullArgs, c.Database, c.UID, c.Password, model, method, args)
    if kwargs != nil {
        fullArgs = append(fullArgs, kwargs)
    }

    err := c.models.Call("execute_kw", fullArgs, result)
    if err != nil {
        if fault, ok := err.(*xmlrpc.FaultError); ok {
            switch fault.Code {
            case 3:
                return &AccessDeniedError{Message: "Acceso Denegado"}
            default:
                return fmt.Errorf("Error XML-RPC: Código %d, Mensaje: %s", fault.Code, fault.String)
            }
        }
        return err
    }
    return nil
}

// Implementación de operaciones CRUD utilizando el método execute
func (c *OdooClient) SearchRead(model string, domain []interface{}, fields []string) ([]map[string]interface{}, error) {
    var records []map[string]interface{}
    kwargs := map[string]interface{}{
        "fields": fields,
    }
    args := []interface{}{domain}
    if err := c.execute(model, "search_read", args, kwargs, &records); err != nil {
        return nil, err
    }
    return records, nil
}

func (c *OdooClient) Create(model string, data map[string]interface{}) (int, error) {
    var id int
    args := []interface{}{data}
    if err := c.execute(model, "create", args, nil, &id); err != nil {
        return 0, err
    }
    return id, nil
}

func (c *OdooClient) Update(model string, ids []int, data map[string]interface{}) (bool, error) {
    var success bool
    args := []interface{}{ids, data}
    if err := c.execute(model, "write", args, nil, &success); err != nil {
        return false, err
    }
    return success, nil
}

func (c *OdooClient) Delete(model string, ids []int) (bool, error) {
    var success bool
    args := []interface{}{ids}
    if err := c.execute(model, "unlink", args, nil, &success); err != nil {
        return false, err
    }
    return success, nil
}

// Método para obtener IDs utilizando execute
func (c *OdooClient) Search(model string, domain []interface{}) ([]int, error) {
    var ids []int
    args := []interface{}{domain}
    if err := c.execute(model, "search", args, nil, &ids); err != nil {
        return nil, err
    }
    return ids, nil
}
