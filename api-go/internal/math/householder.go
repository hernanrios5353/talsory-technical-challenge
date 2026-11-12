// Package math implementa la factorización QR utilizando
// reflexiones de Householder, que es numéricamente estable
// para matrices rectangulares.
package math

import (
"errors"
"math"

"github.com/talsory/coding-challenge/api-go/internal/model"
)

// ErrEmptyMatrix se retorna cuando la matriz no tiene elementos.
var ErrEmptyMatrix = errors.New("la matriz no puede estar vacía")

// QRResult agrupa las matrices Q y R resultantes de la factorización.
type QRResult struct {
Q model.Matrix
R model.Matrix
}

// QR realiza la factorización QR de una matriz A utilizando
// reflexiones de Householder.
func QR(a model.Matrix) (*QRResult, error) {
if len(a) == 0 {
return nil, ErrEmptyMatrix
}

m := len(a)
n := len(a[0])

// Trabajamos sobre una copia para no mutar la matriz original
r := a.Clone()

// Acumulamos las reflexiones en Q, empezando por la identidad
q := identity(m)

// Aplicamos una reflexión de Householder por cada columna
for k := 0; k < n && k < m-1; k++ {
// Extraemos el vector columna desde la diagonal hacia abajo
x := make([]float64, m-k)
for i := k; i < m; i++ {
x[i-k] = r[i][k]
}

// Calculamos el vector de Householder v y el escalar beta
v, beta := householderVector(x)
if beta == 0 {
continue
}

// Aplicamos la reflexión a la submatriz R[k:, k:] (por la izquierda)
applyHouseholderLeft(r, v, beta, k)

// Acumulamos la reflexión en Q (por la derecha)
applyHouseholderRight(q, v, beta, k)
}

// Extraemos Q reducida (m x n) y R reducida (n x n)
qReduced := make(model.Matrix, m)
for i := 0; i < m; i++ {
qReduced[i] = make([]float64, n)
copy(qReduced[i], q[i][:n])
}

rReduced := make(model.Matrix, n)
for i := 0; i < n; i++ {
rReduced[i] = make([]float64, n)
copy(rReduced[i], r[i][:n])
}

// Limpiamos valores muy pequeños en R
cleanSmallValues(rReduced)

return &QRResult{Q: qReduced, R: rReduced}, nil
}

func identity(n int) model.Matrix {
mat := make(model.Matrix, n)
for i := 0; i < n; i++ {
mat[i] = make([]float64, n)
mat[i][i] = 1.0
}
return mat
}

func householderVector(x []float64) ([]float64, float64) {
n := len(x)
if n == 0 {
return nil, 0
}

normX := 0.0
for _, val := range x {
normX += val * val
}
normX = math.Sqrt(normX)

if normX < 1e-12 {
return nil, 0
}

sign := 1.0
if x[0] < 0 {
sign = -1.0
}

v := make([]float64, n)
copy(v, x)
v[0] += sign * normX

vTv := 0.0
for _, val := range v {
vTv += val * val
}
beta := 2.0 / vTv

return v, beta
}

// applyHouseholderLeft aplica (I - beta * v * v^T) * A (por la izquierda)
func applyHouseholderLeft(a model.Matrix, v []float64, beta float64, start int) {
m := len(a)
k := len(v)

for j := start; j < len(a[0]); j++ {
dot := 0.0
for i := 0; i < k; i++ {
dot += a[start+i][j] * v[i]
}
dot *= beta

for i := 0; i < k; i++ {
a[start+i][j] -= v[i] * dot
}
}

// Limpieza numérica en la columna actual
for i := start + 1; i < m; i++ {
a[i][start] = 0.0
}
}

// applyHouseholderRight aplica A * (I - beta * v * v^T) (por la derecha)
func applyHouseholderRight(a model.Matrix, v []float64, beta float64, start int) {
m := len(a)
k := len(v)

for i := 0; i < m; i++ {
dot := 0.0
for j := 0; j < k; j++ {
dot += a[i][start+j] * v[j]
}
dot *= beta

for j := 0; j < k; j++ {
a[i][start+j] -= dot * v[j]
}
}
}

func cleanSmallValues(m model.Matrix) {
const tolerance = 1e-10
for i := range m {
for j := range m[i] {
if math.Abs(m[i][j]) < tolerance {
m[i][j] = 0.0
}
}
}
}
