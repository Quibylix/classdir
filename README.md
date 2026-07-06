ClassDir is a web application designed to help educators manage their classes more dynamically. Teachers can control presentations from any device, annotate content in real-time, and engage students through interactive participation techniques.

## Quick Start

1. **Prerequisites**: [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
2. **Clone** the repository
3. **Create environment file**:
   ```sh
   cp .env.example .env
   ```
4. **Fill in the required variables** in `.env`:

   | Variable | Description |
   |---|---|
   | `DB_PASSWORD` | PostgreSQL database password |
   | `ADMIN_PASSWORD` | Password used to log into the app |
   | `JWT_SECRET` | Secret key for signing authentication tokens |
   | `WS_ORIGIN` | Allowed WebSocket origin (e.g. `http://localhost:3000` or `*`) |

5. **Start the app**:
   ```sh
   docker compose up
   ```
6. Open **http://localhost:3000** and log in with the password set in `ADMIN_PASSWORD`

### Architecture

See [ARCHITECTURE.md](./ARCHITECTURE.md) for system design and domain rules.
