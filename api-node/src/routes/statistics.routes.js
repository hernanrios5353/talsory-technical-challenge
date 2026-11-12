/**
 * Rutas de estadísticas.
 * Define el endpoint POST /api/v1/statistics protegido por JWT.
 */

const express = require('express');
const router = express.Router();
const { getStatistics } = require('../controllers/statistics.controller');
const { authenticateToken } = require('../middlewares/auth.middleware');

// POST /api/v1/statistics
// Protegido por JWT: requiere header Authorization: Bearer <token>
router.post('/statistics', authenticateToken, getStatistics);

module.exports = router;