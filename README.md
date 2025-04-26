# How to launch

Postgres
```bash
docker compose -f docker/pg.yml up
```

App
```bash
docker compose -f docker/app.yml up --build app
```