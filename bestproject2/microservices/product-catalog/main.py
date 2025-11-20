from fastapi import FastAPI, HTTPException, Depends, status, Query
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from contextlib import asynccontextmanager
from prometheus_client import make_asgi_app, Counter, Histogram
import uvicorn
import logging
from typing import List, Optional
from pythonjsonlogger import jsonlogger
from elasticsearch import Elasticsearch
from datetime import datetime
import socket
import json

from app.config import settings
from app.database import engine, Base, get_db
from app.routers import products, categories, inventory
from app.middleware.auth import get_current_user
from app.middleware.logging_middleware import LoggingMiddleware

# Elasticsearch client for logging
es_client = None
try:
    es_client = Elasticsearch([settings.ELASTICSEARCH_URL])
    if es_client.ping():
        print("Connected to Elasticsearch successfully")
except Exception as e:
    print(f"Failed to connect to Elasticsearch: {e}. Using console logging only.")

# Custom Elasticsearch Handler
class ElasticsearchHandler(logging.Handler):
    def __init__(self, es_client, index_prefix="product-catalog-logs"):
        super().__init__()
        self.es_client = es_client
        self.index_prefix = index_prefix
        self.hostname = socket.gethostname()
    
    def emit(self, record):
        if self.es_client is None:
            return
        
        try:
            log_entry = {
                "@timestamp": datetime.utcnow().isoformat(),
                "level": record.levelname,
                "logger": record.name,
                "message": record.getMessage(),
                "service": "product-catalog",
                "environment": "production",
                "hostname": self.hostname,
                "module": record.module,
                "function": record.funcName,
                "line": record.lineno
            }
            
            # Add extra fields if present
            if hasattr(record, 'event'):
                log_entry['event'] = record.event
            if hasattr(record, 'method'):
                log_entry['method'] = record.method
            if hasattr(record, 'path'):
                log_entry['path'] = record.path
            if hasattr(record, 'status_code'):
                log_entry['status_code'] = record.status_code
            if hasattr(record, 'duration_seconds'):
                log_entry['duration_seconds'] = record.duration_seconds
            if hasattr(record, 'duration_ms'):
                log_entry['duration_ms'] = record.duration_ms
            if hasattr(record, 'client_ip'):
                log_entry['client_ip'] = record.client_ip
            if hasattr(record, 'user_agent'):
                log_entry['user_agent'] = record.user_agent
            
            index_name = f"{self.index_prefix}-{datetime.utcnow().strftime('%Y.%m.%d')}"
            self.es_client.index(index=index_name, document=log_entry)
        except Exception as e:
            print(f"Failed to send log to Elasticsearch: {e}")

# Configure logging with Elasticsearch and JSON formatting
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

# JSON Formatter for structured logging
class CustomJsonFormatter(jsonlogger.JsonFormatter):
    def add_fields(self, log_record, record, message_dict):
        super(CustomJsonFormatter, self).add_fields(log_record, record, message_dict)
        log_record['timestamp'] = datetime.utcnow().isoformat()
        log_record['service'] = 'product-catalog'
        log_record['environment'] = 'production'
        log_record['hostname'] = socket.gethostname()
        log_record['level'] = record.levelname

# Console Handler with JSON formatting
console_handler = logging.StreamHandler()
console_formatter = CustomJsonFormatter('%(timestamp)s %(level)s %(name)s %(message)s')
console_handler.setFormatter(console_formatter)

# Add handlers
logger.addHandler(console_handler)
if es_client:
    es_handler = ElasticsearchHandler(es_client)
    es_handler.setLevel(logging.INFO)
    logger.addHandler(es_handler)
    logger.info("Elasticsearch logging handler configured successfully")

# Prometheus metrics
REQUEST_COUNT = Counter('http_requests_total', 'Total HTTP Requests', ['method', 'endpoint', 'status'])
REQUEST_LATENCY = Histogram('http_request_duration_seconds', 'HTTP Request Duration', ['method', 'endpoint'])

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    logger.info("Starting Product Catalog Service")
    Base.metadata.create_all(bind=engine)
    yield
    # Shutdown
    logger.info("Shutting down Product Catalog Service")

app = FastAPI(
    title="Product Catalog Service",
    description="E-Commerce Product Catalog and Inventory Management",
    version="1.0.0",
    lifespan=lifespan
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Custom middleware
app.add_middleware(LoggingMiddleware)

# Health check
@app.get("/health")
async def health_check():
    return {
        "status": "healthy",
        "service": "product-catalog",
        "version": "1.0.0"
    }

# Metrics endpoint
metrics_app = make_asgi_app()
app.mount("/metrics", metrics_app)

# Include routers
app.include_router(products.router, prefix="/api/v1/products", tags=["products"])
app.include_router(categories.router, prefix="/api/v1/categories", tags=["categories"])
app.include_router(inventory.router, prefix="/api/v1/inventory", tags=["inventory"])

# Global exception handler
@app.exception_handler(HTTPException)
async def http_exception_handler(request, exc):
    return JSONResponse(
        status_code=exc.status_code,
        content={"success": False, "error": exc.detail}
    )

@app.exception_handler(Exception)
async def general_exception_handler(request, exc):
    logger.error(f"Unhandled exception: {str(exc)}", exc_info=True)
    return JSONResponse(
        status_code=500,
        content={"success": False, "error": "Internal server error"}
    )

if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=settings.PORT,
        reload=settings.DEBUG,
        log_level="info"
    )
