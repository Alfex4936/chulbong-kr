import multiprocessing

import uvicorn
from fastapi import FastAPI

from chulbong_fastapi.api.v1.markers import markers
from chulbong_fastapi.core.lifespan import lifespan

app = FastAPI(lifespan=lifespan)

app.include_router(markers.router)

if __name__ == "__main__":

    # Dynamically calculate the number of workers
    num_cpu_cores = multiprocessing.cpu_count()
    workers = 2 * num_cpu_cores + 1

    uvicorn.run(
        app="chulbong_fastapi.main:app",
        host="0.0.0.0",
        port=8000,
        workers=workers,
        limit_concurrency=10000,
        backlog=2048,  # Handle large bursts of traffic
        timeout_keep_alive=5,  # Keep connections alive for 5 seconds
        limit_max_requests=1000,  # Recycle workers after 1000 requests
        timeout_graceful_shutdown=30,  # Allow 30 seconds for graceful shutdown
        proxy_headers=True,  # If behind a reverse proxy
        log_level="info",  # Set logging level to info
    )
