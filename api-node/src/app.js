/**
 * Punto de entrada de la API de Node.js.
 * Configura Express, middlewares globales, rutas y health check.
 */

const express = require('express');
const statisticsRoutes = require('./routes/statistics.routes');

const app = express();

// Middleware para parsear JSON en el body
app.use(express.json());

// Middleware de logs simple (opcional, pero útil para depurar)
app.use((req, res, next) => {
  console.log(`${new Date().toISOString()} - ${req.method} ${req.path}`);
  next();
});

// Health check (crítico para Docker y despliegue cloud)
app.get('/health', (req, res) => {
  res.status(200).json({ status: 'ok', service: 'api-node' });
});

// Rutas de la API
app.use('/api/v1', statisticsRoutes);

// Middleware para rutas no encontradas (404)
app.use((req, res) => {
  res.status(404).json({ error: 'Ruta no encontrada' });
});

// Middleware de manejo de errores global
app.use((err, req, res, next) => {
  console.error('Error no manejado:', err);
  res.status(500).json({ error: 'Error interno del servidor' });
});

// Arrancar el servidor
const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`🚀 API Node.js escuchando en el puerto ${PORT}`);
});

module.exports = app; // Exportamos para pruebas con supertest