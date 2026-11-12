/**
 * Controlador de estadísticas.
 * Recibe las matrices Q y R de la API de Go, calcula las estadísticas
 * y devuelve la respuesta consolidada.
 */

const { calculateStats } = require('../services/statistics.service');

/**
 * Maneja POST /api/v1/statistics
 * Body esperado: { q: number[][], r: number[][] }
 */
function getStatistics(req, res) {
  const { q, r } = req.body;

  // Validación básica de entrada
  if (!q || !r) {
    return res.status(400).json({
      error: "Los campos 'q' y 'r' son obligatorios",
    });
  }

  if (!Array.isArray(q) || !Array.isArray(r)) {
    return res.status(400).json({
      error: "Los campos 'q' y 'r' deben ser arrays bidimensionales",
    });
  }

  try {
    // Calcular estadísticas individuales
    const statsQ = calculateStats(q);
    const statsR = calculateStats(r);

    // Calcular estadísticas globales (sobre todos los valores de Q y R juntos)
    // Para max, min, average y sum, concatenamos y calculamos.
    const combined = q.concat(r);
    const statsGlobal = calculateStats(combined);

    // CORRECCIÓN: Para la propiedad isDiagonal, verificamos si Q o R lo son.
    // No tiene sentido matemático verificarlo sobre la concatenación.
    statsGlobal.isDiagonal = statsQ.isDiagonal || statsR.isDiagonal;

    return res.status(200).json({
      q: statsQ,
      r: statsR,
      global: statsGlobal,
    });
  } catch (err) {
    console.error('Error calculando estadísticas:', err);
    return res.status(500).json({ error: 'Error interno calculando estadísticas' });
  }
}

module.exports = { getStatistics };