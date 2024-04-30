package utilities_test

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

/* -------------------------------------------------------------------------- */
/*                                 AUXILIARES                                 */
/* -------------------------------------------------------------------------- */
var Resultados strings.Builder

func AddText(text string) {
	Resultados.WriteString(text + "\n")
}

func CreateFile(name string) error {
	//Ensure the directory exists
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err CreateFile dir==", err)
		return err
	}

	// Create file
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Err CreateFile create==", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

// Funtion to open bin file in read/write mode
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile==", err)
		return nil, err
	}
	return file, nil
}

// Function to Write an object in a bin file
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==", err)
		return err
	}
	return nil
}

// Function to Read an object from a bin file
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}

func ConvertToZeros(filename string, start int64, end int64) error {
	// Abrir el archivo en modo lectura y escritura
	file, err := OpenFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Obtener el tamaño del archivo
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	// Verificar que el rango especificado esté dentro del tamaño del archivo
	if start >= fileSize || end >= fileSize || start > end {
		return fmt.Errorf("rango fuera de los límites del archivo")
	}

	// Crear un slice de bytes lleno de ceros
	zeroBytes := make([]byte, end-start+1)

	// Escribir los ceros en el rango definido del archivo
	_, err = file.WriteAt(zeroBytes, start)
	if err != nil {
		return err
	}

	return nil
}

func CalcularPorcentaje(tamanoParticion int64, tamanoDisco int64) int64 {
	return (int64(tamanoParticion) * 100 / int64(tamanoDisco))
}

func LimpiarCerosBinarios(bytes []byte) []byte {
	// Crear un slice para almacenar los bytes sin ceros
	resultado := []byte{}

	// Recorrer el slice de bytes
	for _, b := range bytes {
		// Si el byte no es cero, agregarlo al resultado
		if b != 0 {
			resultado = append(resultado, b)
		}
	}

	return resultado
}

func EsNumero(caracter string) bool {
	runeValue := []rune(caracter)
	if len(runeValue) != 1 {
		// No es un solo caracter
		return false
	}
	return unicode.IsDigit(runeValue[0])
}
