import os
import pytest
from testcontainers.postgres import PostgresContainer
from testcontainers.redis import RedisContainer

@pytest.mark.integration
def test_postgres_and_redis_containers_start():
    with PostgresContainer("postgres:15-alpine") as pg, RedisContainer("redis:7-alpine") as rc:
        os.environ["DATABASE_URL"] = pg.get_connection_url()
        os.environ["REDIS_URL"] = f"redis://{rc.get_container_host_ip()}:{rc.get_exposed_port(6379)}/0"
        # Containers started successfully; in a real test, initialize app and perform DB ops
        assert pg.get_connection_url().startswith("postgresql")
        assert rc.get_exposed_port(6379)
