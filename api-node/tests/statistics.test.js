/**
 * Pruebas unitarias para el servicio de estadísticas.
 */

const {
  calculateStats,
  calculateBasicStats,
  isDiagonal,
  flattenMatrix,
} = require('../src/services/statistics.service');

describe('flattenMatrix', () => {
  test('aplana una matriz 2x2 correctamente', () => {
    expect(flattenMatrix([[1, 2], [3, 4]])).toEqual([1, 2, 3, 4]);
  });

  test('aplana una matriz rectangular 3x2', () => {
    expect(flattenMatrix([[1, 2], [3, 4], [5, 6]])).toEqual([1, 2, 3, 4, 5, 6]);
  });

  test('matriz vacía devuelve array vacío', () => {
    expect(flattenMatrix([])).toEqual([]);
  });
});

describe('calculateBasicStats', () => {
  test('calcula max, min, average y sum de una matriz 2x2', () => {
    const result = calculateBasicStats([[1, 2], [3, 4]]);
    expect(result.max).toBe(4);
    expect(result.min).toBe(1);
    expect(result.sum).toBe(10);
    expect(result.average).toBe(2.5);
  });

  test('maneja números negativos', () => {
    const result = calculateBasicStats([[-5, 0], [5, 10]]);
    expect(result.max).toBe(10);
    expect(result.min).toBe(-5);
    expect(result.sum).toBe(10);
    expect(result.average).toBe(2.5);
  });

  test('matriz vacía devuelve ceros', () => {
    const result = calculateBasicStats([]);
    expect(result).toEqual({ max: 0, min: 0, average: 0, sum: 0 });
  });
});

describe('isDiagonal', () => {
  test('matriz diagonal verdadera (cuadrada con ceros fuera)', () => {
    expect(isDiagonal([[1, 0], [0, 2]])).toBe(true);
  });

  test('matriz diagonal 3x3 verdadera', () => {
    expect(isDiagonal([[1, 0, 0], [0, 5, 0], [0, 0, 9]])).toBe(true);
  });

  test('matriz no diagonal (tiene valores fuera)', () => {
    expect(isDiagonal([[1, 2], [0, 3]])).toBe(false);
  });

  test('matriz no cuadrada no es diagonal', () => {
    expect(isDiagonal([[1, 0], [0, 2], [0, 0]])).toBe(false);
  });

  test('tolerancia numérica: valores muy pequeños se consideran cero', () => {
    expect(isDiagonal([[1, 1e-12], [0, 2]])).toBe(true);
  });
});

describe('calculateStats', () => {
  test('combina todas las estadísticas correctamente', () => {
    const result = calculateStats([[1, 0], [0, 2]]);
    expect(result.max).toBe(2);
    expect(result.min).toBe(0);
    expect(result.sum).toBe(3);
    expect(result.average).toBe(0.75);
    expect(result.isDiagonal).toBe(true);
  });

  test('matriz no diagonal', () => {
    const result = calculateStats([[1, 2], [3, 4]]);
    expect(result.isDiagonal).toBe(false);
  });
});