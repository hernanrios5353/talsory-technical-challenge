package tests

import (
"math"
"testing"

qrmath "github.com/talsory/coding-challenge/api-go/internal/math"
"github.com/talsory/coding-challenge/api-go/internal/model"
)

const tolerance = 1e-9

func almostEqual(a, b, tol float64) bool {
return math.Abs(a-b) < tol
}

func assertMatrixAlmostEqual(t *testing.T, got, want model.Matrix, tol float64, msg string) {
t.Helper()
if len(got) != len(want) {
t.Fatalf("%s: número de filas distinto: got=%d, want=%d", msg, len(got), len(want))
}
for i := range got {
if len(got[i]) != len(want[i]) {
t.Fatalf("%s: fila %d tiene %d columnas, se esperaban %d", msg, i, len(got[i]), len(want[i]))
}
for j := range got[i] {
if !almostEqual(got[i][j], want[i][j], tol) {
t.Errorf("%s: diferencia en [%d][%d]: got=%.12f, want=%.12f, diff=%.12e",
msg, i, j, got[i][j], want[i][j], math.Abs(got[i][j]-want[i][j]))
}
}
}
}

func multiplyMatrices(a, b model.Matrix) model.Matrix {
m := len(a)
n := len(b[0])
k := len(b)
result := make(model.Matrix, m)
for i := 0; i < m; i++ {
result[i] = make([]float64, n)
for j := 0; j < n; j++ {
sum := 0.0
for l := 0; l < k; l++ {
sum += a[i][l] * b[l][j]
}
result[i][j] = sum
}
}
return result
}

func transpose(a model.Matrix) model.Matrix {
m := len(a)
n := len(a[0])
result := make(model.Matrix, n)
for i := 0; i < n; i++ {
result[i] = make([]float64, m)
for j := 0; j < m; j++ {
result[i][j] = a[j][i]
}
}
return result
}

func TestQRReconstructsA(t *testing.T) {
cases := []struct {
name string
a    model.Matrix
}{
{"3x2", model.Matrix{{1, 2}, {3, 4}, {5, 6}}},
{"4x2", model.Matrix{{1, 2}, {3, 4}, {5, 6}, {7, 8}}},
{"5x3", model.Matrix{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}}},
{"3x3", model.Matrix{{2, 1, 0}, {1, 3, 1}, {0, 1, 2}}},
{"cuadrada identidad", model.Matrix{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}},
}

for _, tc := range cases {
t.Run(tc.name, func(t *testing.T) {
result, err := qrmath.QR(tc.a)
if err != nil {
t.Fatalf("QR devolvió error: %v", err)
}
reconstructed := multiplyMatrices(result.Q, result.R)
assertMatrixAlmostEqual(t, reconstructed, tc.a, tolerance, "Q * R no reconstruye A")
})
}
}

func TestQOrthogonal(t *testing.T) {
cases := []struct {
name string
a    model.Matrix
}{
{"3x2", model.Matrix{{1, 2}, {3, 4}, {5, 6}}},
{"4x2", model.Matrix{{1, 2}, {3, 4}, {5, 6}, {7, 8}}},
{"5x3", model.Matrix{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}}},
}

for _, tc := range cases {
t.Run(tc.name, func(t *testing.T) {
result, err := qrmath.QR(tc.a)
if err != nil {
t.Fatalf("QR devolvió error: %v", err)
}
qT := transpose(result.Q)
qtq := multiplyMatrices(qT, result.Q)
n := len(result.Q[0])
expected := make(model.Matrix, n)
for i := 0; i < n; i++ {
expected[i] = make([]float64, n)
expected[i][i] = 1.0
}
assertMatrixAlmostEqual(t, qtq, expected, tolerance, "Q^T * Q no es la identidad (Q no es ortogonal)")
})
}
}

func TestRTriangularSuperior(t *testing.T) {
cases := []struct {
name string
a    model.Matrix
}{
{"3x2", model.Matrix{{1, 2}, {3, 4}, {5, 6}}},
{"5x3", model.Matrix{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}}},
{"3x3", model.Matrix{{2, 1, 0}, {1, 3, 1}, {0, 1, 2}}},
}

for _, tc := range cases {
t.Run(tc.name, func(t *testing.T) {
result, err := qrmath.QR(tc.a)
if err != nil {
t.Fatalf("QR devolvió error: %v", err)
}
for i := range result.R {
for j := 0; j < i && j < len(result.R[i]); j++ {
if !almostEqual(result.R[i][j], 0, tolerance) {
t.Errorf("R[%d][%d] = %.12e, se esperaba 0 (R no es triangular superior)", i, j, result.R[i][j])
}
}
}
})
}
}

func TestQRDimensions(t *testing.T) {
cases := []struct {
name     string
a        model.Matrix
wantQRows int
wantQCols int
wantRRows int
wantRCols int
}{
{"3x2", model.Matrix{{1, 2}, {3, 4}, {5, 6}}, 3, 2, 2, 2},
{"4x2", model.Matrix{{1, 2}, {3, 4}, {5, 6}, {7, 8}}, 4, 2, 2, 2},
{"5x3", model.Matrix{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}}, 5, 3, 3, 3},
}

for _, tc := range cases {
t.Run(tc.name, func(t *testing.T) {
result, err := qrmath.QR(tc.a)
if err != nil {
t.Fatalf("QR devolvió error: %v", err)
}
if len(result.Q) != tc.wantQRows {
t.Errorf("Q tiene %d filas, se esperaban %d", len(result.Q), tc.wantQRows)
}
if len(result.Q[0]) != tc.wantQCols {
t.Errorf("Q tiene %d columnas, se esperaban %d", len(result.Q[0]), tc.wantQCols)
}
if len(result.R) != tc.wantRRows {
t.Errorf("R tiene %d filas, se esperaban %d", len(result.R), tc.wantRRows)
}
if len(result.R[0]) != tc.wantRCols {
t.Errorf("R tiene %d columnas, se esperaban %d", len(result.R[0]), tc.wantRCols)
}
})
}
}

func TestQRErrors(t *testing.T) {
t.Run("matriz vacía", func(t *testing.T) {
_, err := qrmath.QR(model.Matrix{})
if err == nil {
t.Error("se esperaba error para matriz vacía")
}
})
}

func TestQRColumnOfZeros(t *testing.T) {
a := model.Matrix{
{0, 1},
{0, 2},
{0, 3},
}
result, err := qrmath.QR(a)
if err != nil {
t.Fatalf("QR devolvió error: %v", err)
}
reconstructed := multiplyMatrices(result.Q, result.R)
assertMatrixAlmostEqual(t, reconstructed, a, tolerance, "Q * R no reconstruye A con columna de ceros")
}
