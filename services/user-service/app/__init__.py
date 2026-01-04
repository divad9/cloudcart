"""
Flask application factory
"""
from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from flask_jwt_extended import JWTManager
from flask_cors import CORS
from config import Config

# Initialize extensions
db = SQLAlchemy()
jwt = JWTManager()


def create_app(config_class=Config):
    """Create and configure Flask application"""
    
    app = Flask(__name__)
    app.config.from_object(config_class)
    
    # Initialize extensions
    db.init_app(app)
    jwt.init_app(app)
    CORS(app)
    
    # Register blueprints
    from app.routes.auth import auth_bp
    app.register_blueprint(auth_bp)
    
    # Health check endpoint
    @app.route('/health')
    def health():
        return {'status': 'healthy', 'service': 'user-service'}, 200
    
    # Create tables
    with app.app_context():
        db.create_all()
    
    return app