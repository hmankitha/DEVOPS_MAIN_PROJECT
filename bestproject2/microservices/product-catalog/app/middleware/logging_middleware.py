from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
import logging
import time
import json

logger = logging.getLogger(__name__)

class LoggingMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        start_time = time.time()
        
        # Prepare request log data
        request_data = {
            'event': 'request_received',
            'method': request.method,
            'path': request.url.path,
            'query_params': str(request.query_params),
            'client_ip': request.client.host if request.client else None,
            'user_agent': request.headers.get('user-agent', 'unknown')
        }
        
        # Log request with structured data
        logger.info("HTTP Request", extra=request_data)
        
        response = await call_next(request)
        
        # Calculate process time
        process_time = time.time() - start_time
        
        # Prepare response log data
        response_data = {
            'event': 'request_completed',
            'method': request.method,
            'path': request.url.path,
            'status_code': response.status_code,
            'duration_seconds': round(process_time, 4),
            'duration_ms': round(process_time * 1000, 2)
        }
        
        # Log response with structured data
        if response.status_code >= 400:
            logger.warning("HTTP Request Failed", extra=response_data)
        else:
            logger.info("HTTP Request Success", extra=response_data)
        
        # Add process time to response headers
        response.headers["X-Process-Time"] = str(process_time)
        
        return response
