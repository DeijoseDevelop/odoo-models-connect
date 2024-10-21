package odoo

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

// Función para convertir una imagen a base64
func ConvertImageToBase64(imagePath string) (string, error) {
    file, err := os.Open(imagePath)
    if err != nil {
        return "", fmt.Errorf("Error al abrir la imagen: %v", err)
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return "", fmt.Errorf("Error al decodificar la imagen: %v", err)
    }

    buf := new(bytes.Buffer)
    if err := jpeg.Encode(buf, img, nil); err != nil {
        return "", fmt.Errorf("Error al codificar la imagen: %v", err)
    }

    return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
