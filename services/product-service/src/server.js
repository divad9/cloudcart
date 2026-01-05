/**
 * Product Service Entry Point
 */
const app = require('./app');
const { sequelize, testConnection } = require('./config/database');
const Product = require('./models/Product');
const Category = require('./models/Category');

const PORT = process.env.PORT || 5002;

// Start server
const startServer = async () => {
  try {
    // Test database connection
    await testConnection();

    // Sync database models
    await sequelize.sync({ alter: true });
    console.log('✅ Database models synchronized');

    // Start listening
    app.listen(PORT, '0.0.0.0', () => {
      console.log(`🚀 Product Service running on port ${PORT}`);
      console.log(`📍 http://localhost:${PORT}/health`);
    });

  } catch (error) {
    console.error('❌ Failed to start server:', error);
    process.exit(1);
  }
};

startServer();