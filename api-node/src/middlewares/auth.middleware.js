/**
 * Middleware de autenticación JWT.
 * Verifica que el header Authorization contenga un token válido.
 */

const jwt = require('jsonwebtoken');

/**
 * Middleware que protege las rutas.
 * Extrae el token del header `Authorization: Bearer <token>`,
 * lo verifica con JWT_SECRET y, si es válido, continúa.
 */
function authenticateToken(req, res, next) {
  const authHeader = req.headers['authorization'];
  // Formato esperado: "Bearer <token>"
  const token = authHeader && authHeader.split(' ')[1];

  if (!token) {
    return res.status(401).json({ error: 'Token de acceso no proporcionado' });
  }

  const secret = process.env.JWT_SECRET;
  if (!secret) {
    console.error('JWT_SECRET no está definido en las variables de entorno');
    return res.status(500).json({ error: 'Error de configuración del servidor' });
  }

  jwt.verify(token, secret, (err, decoded) => {
    if (err) {
      return res.status(403).json({ error: 'Token inválido o expirado' });
    }
    // Guardamos el payload decodificado en req.user por si se necesita después
    req.user = decoded;
    next();
  });
}

module.exports = { authenticateToken };