/**
 * Servicio de estadísticas.
 * Contiene funciones puras para calcular métricas sobre matrices.
 * No depende de Express ni de HTTP, lo que facilita las pruebas unitarias.
 */

const DIAGONAL_TOLERANCE = 1e-9;

/**
 * Aplana una matriz (array de arrays) en un solo array unidimensional.
 * @param {number[][]} matrix - Matriz de entrada.
 * @returns {number[]} Array plano con todos los elementos.
 */
function flattenMatrix(matrix) {
  return matrix.reduce((acc, row) => acc.concat(row), []);
}

/**
 * Calcula las estadísticas básicas de una matriz:
 * máximo, mínimo, promedio y suma total.
 * @param {number[][]} matrix - Matriz de entrada.
 * @returns {{ max: number, min: number, average: number, sum: number }}
 */
function calculateBasicStats(matrix) {
  const flat = flattenMatrix(matrix);

  if (flat.length === 0) {
    return { max: 0, min: 0, average: 0, sum: 0 };
  }

  let max = flat[0];
  let min = flat[0];
  let sum = 0;

  for (const value of flat) {
    if (value > max) max = value;
    if (value < min) min = value;
    sum += value;
  }

  const average = sum / flat.length;

  return { max, min, average, sum };
}

/**
 * Verifica si una matriz es diagonal.
 * Una matriz es diagonal si es cuadrada y todos los elementos
 * fuera de la diagonal principal son cero (con tolerancia numérica).
 * @param {number[][]} matrix - Matriz de entrada.
 * @returns {boolean} true si es diagonal, false en caso contrario.
 */
function isDiagonal(matrix) {
  const rows = matrix.length;
  if (rows === 0) return false;

  // Una matriz diagonal debe ser cuadrada
  for (const row of matrix) {
    if (row.length !== rows) return false;
  }

  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < rows; j++) {
      if (i !== j && Math.abs(matrix[i][j]) > DIAGONAL_TOLERANCE) {
        return false;
      }
    }
  }

  return true;
}

/**
 * Calcula todas las estadísticas de una matriz, incluyendo
 * la verificación de si es diagonal.
 * @param {number[][]} matrix - Matriz de entrada.
 * @returns {{
 *   max: number,
 *   min: number,
 *   average: number,
 *   sum: number,
 *   isDiagonal: boolean
 * }}
 */
function calculateStats(matrix) {
  const basic = calculateBasicStats(matrix);
  const diagonal = isDiagonal(matrix);

  return {
    ...basic,
    isDiagonal: diagonal,
  };
}

module.exports = {
  calculateStats,
  calculateBasicStats,
  isDiagonal,
  flattenMatrix,
};