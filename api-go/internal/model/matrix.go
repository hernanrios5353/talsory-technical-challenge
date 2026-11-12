// Package model contiene los tipos de datos y validaciones
// relacionados con las matrices utilizadas en la factorización QR.
package model

import (
	"errors"
	"fmt"
)

// Matrix representa una matriz rectangular de números reales.
// Se utiliza [][]float64 porque es el tipo natural para JSON
// y para las operaciones matemáticas que realizaremos.
type Matrix [][]float64

// Errores de validación predefinidos para facilitar las pruebas
// y el manejo de errores en los handlers.
var (
	ErrEmptyMatrix        = errors.New("la matriz no puede estar vacía")
	ErrIrregularMatrix    = errors.New("todas las filas deben tener la misma cantidad de columnas")
	ErrInvalidDimensions  = errors.New("la matriz debe tener al menos 1 fila y 1 columna")
	ErrMoreColumnsThanRows = errors.New("para la factorización QR se requiere que filas >= columnas")
)

// Validate verifica que la matriz cumpla con las condiciones
// necesarias para la factorización QR:
//   - No está vacía.
//   - Todas las filas tienen la misma longitud (matriz rectangular).
//   - Tiene al menos 1 fila y 1 columna.
//   - Tiene más filas que columnas (o iguales), condición necesaria
//     para la factorización QR estándar (Householder).
func (m Matrix) Validate() error {
	rows := len(m)
	if rows == 0 {
		return ErrEmptyMatrix
	}

	cols := len(m[0])
	if cols == 0 {
		return ErrInvalidDimensions
	}

	// Verificamos que todas las filas tengan la misma longitud
	for i, row := range m {
		if len(row) != cols {
			return fmt.Errorf("%w: fila %d tiene %d columnas, se esperaban %d",
				ErrIrregularMatrix, i, len(row), cols)
		}
	}

	// Para QR estándar (Householder) necesitamos filas >= columnas
	if rows < cols {
		return fmt.Errorf("%w: filas=%d, columnas=%d",
			ErrMoreColumnsThanRows, rows, cols)
	}

	return nil
}

// Rows devuelve el número de filas de la matriz.
func (m Matrix) Rows() int {
	return len(m)
}

// Cols devuelve el número de columnas de la matriz.
// Asume que la matriz ya fue validada (todas las filas iguales).
func (m Matrix) Cols() int {
	if len(m) == 0 {
		return 0
	}
	return len(m[0])
}

// Clone devuelve una copia profunda de la matriz.
// Útil para no modificar la matriz original durante los cálculos.
func (m Matrix) Clone() Matrix {
	clone := make(Matrix, len(m))
	for i, row := range m {
		clone[i] = make([]float64, len(row))
		copy(clone[i], row)
	}
	return clone
}