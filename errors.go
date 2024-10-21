package odoo

import "fmt"

// Error personalizado para acceso denegado
type AccessDeniedError struct {
    Message string
}

func (e *AccessDeniedError) Error() string {
    return e.Message
}

// Error cuando un objeto no existe
type ObjectDoesNotExistError struct {
    ID int
}

func (e *ObjectDoesNotExistError) Error() string {
    return fmt.Sprintf("El objeto con ID %d no existe", e.ID)
}
