# Sports Events Calendar - Go REST API
## Table of Contents
- [Information](#information)
- [Installation](#installation)
- [Usage](#usage)
- [ERD](#erd)
- [API Documentation](#api-documentation)

## Information
> [!NOTE]
> **Requirements:**
> - `go>=1.21`
> - PostgreSQL database
> - [goose](https://github.com/pressly/goose) for running migrations

## Installation
```bash
git clone https://github.com/MedrekIT/sca-recrutation-task.git
cd sca-recrutation-task
go mod tidy
```

Initialize your database
```
psql -U postgres
CREATE DATABASE sca_db;
\c sca_db
ALTER USER postgres PASSWORD 'postgres';
\q
```

Set up environment variables in a `.env` file:
```env
SERVER_PORT=8080
DB_URL=postgres://user:password@localhost:5432/sca_db?sslmode=disable
```

Run database migrations in the server directory:
```bash
goose postgres "postgres://postgres:postgres@localhost:5432/sca_db" up
```

## Usage
### Run the server
```bash
go run ./server/cmd/main.go
```

The server will start listening on the configured port (default: `8080`).

### Frontend
Open `index.html` in a browser to access the sports events calendar.
Navigate to `createEvent.html` to add new events through the UI.

> [!IMPORTANT]
> The frontend expects the API to be running at `http://localhost:8080`.
> I used `python -m http.server 3000` in my terminal to make it work.
> Make sure CORS is configured to allow your frontend origin.

---

## ERD

![Entity Relationship Diagram](./assets/sca_db_erd.svg)

---

# API Documentation
Base development URL: `http://localhost:8080`

## Sports

### `GET /api/sports`
**Description:** Returns a list of all sport names registered in the system

**Response (`200 OK`):**
```json
["Football", "Basketball", "Tennis"]
```

**Possible errors:**
- `500 Internal Server Error`:
```json
{
    "code": "DATABASE_ERROR",
    "message": "Something went wrong"
}
```

---

### `POST /api/sports`
**Description:** Adds a new sport to the system

**Request Body:**
```json
{
    "name": "Football"
}
```

**Response (`201 Created`):**
```json
{
    "id": 1
}
```

**Possible errors:**
- `400 Bad Request`: invalid or missing request body
- `500 Internal Server Error`: failure during database access

---

## Countries

### `POST /api/countries`
**Description:** Adds a new country to the system

**Request Body:**
```json
{
    "country_code": "POL",
    "name": "Poland"
}
```

**Response (`201 Created`):**
```json
{
    "message": "Country added successfully"
}
```

**Possible errors:**
- `400 Bad Request`: invalid JSON, missing fields, or `country_code` longer than 3 characters
- `500 Internal Server Error`: failure during database access

---

## Competitions

### `GET /api/sports/{sportName}/competitions`
**Description:** Returns a list of all competition editions for a given sport

**Path Parameters:**
- `sportName` (string) - Name of the sport

**Response (`200 OK`):**
```json
[
    {
        "edition_id": 1,
        "label": "Premier League 2024/25"
    },
    ...
]
```

**Possible errors:**
- `400 Bad Request`: missing sport name
- `500 Internal Server Error`: failure during database access

---

### `POST /api/competitions`
**Description:** Creates a new competition along with its first edition in the system

**Request Body:**
```json
{
    "competition_name": "Premier League",
    "sport_id": 1,
    "season": "2024/25"
}
```

**Response (`201 Created`):**
```json
{
    "id": 1
}
```

> [!NOTE]
> The returned `id` is the **edition ID**, not the competition ID. Use it when referencing stages.

**Possible errors:**
- `400 Bad Request`: invalid JSON or non-existent `sport_id`
- `500 Internal Server Error`: failure during database access

---

## Stages

### `GET /api/competitions/{editionID}/stages`
**Description:** Returns a list of all stages for a given competition edition

**Path Parameters:**
- `editionID` (integer) - ID of the competition edition

**Response (`200 OK`):**
```json
[
    {
        "stage_id": 1,
        "name": "Group Stage"
    },
    ...
]
```

**Possible errors:**
- `400 Bad Request`: missing or invalid edition ID
- `500 Internal Server Error`: failure during database access

---

### `POST /api/competitions/{editionID}/stages`
**Description:** Adds a new stage to a competition edition

**Path Parameters:**
- `editionID` (integer) - ID of the competition edition

**Request Body:**
```json
{
    "name": "Group Stage",
    "edition_id": 1
}
```

**Response (`201 Created`):**
```json
{
    "id": 1
}
```

**Possible errors:**
- `400 Bad Request`: invalid or missing request body
- `500 Internal Server Error`: failure during database access

---

## Competitors

### `GET /api/sports/{sportName}/competitors`
**Description:** Returns a list of all competitor names for a given sport

**Path Parameters:**
- `sportName` (string) - Name of the sport

**Response (`200 OK`):**
```json
["Manchester United", "Chelsea", "Arsenal"]
```

**Possible errors:**
- `400 Bad Request`: missing sport name
- `500 Internal Server Error`: failure during database access

---

### `POST /api/competitors`
**Description:** Adds a new competitor to the system

**Request Body:**
```json
{
    "name": "Manchester United",
    "country_code": "ENG",
    "sport_id": 1
}
```

> [!NOTE]
> `country_code` is optional. If provided, it must be a valid code present in the countries table and no longer than 3 characters.

**Response (`201 Created`):**
```json
{
    "competitor_id": 1
}
```

**Possible errors:**
- `400 Bad Request`: invalid JSON, invalid/non-existent `country_code`, or non-existent `sport_id`
- `500 Internal Server Error`: failure during database access

---

## Venues

### `GET /api/sports/{sportName}/venues`
**Description:** Returns a list of all venues for a given sport

**Path Parameters:**
- `sportName` (string) - Name of the sport

**Response (`200 OK`):**
```json
[
    {
        "venue_id": 1,
        "name": "Old Trafford"
    },
    ...
]
```

**Possible errors:**
- `400 Bad Request`: missing sport name
- `500 Internal Server Error`: failure during database access

---

### `POST /api/venues`
**Description:** Adds a new venue to the system

**Request Body:**
```json
{
    "name": "Old Trafford",
    "city": "Manchester",
    "country_code": "ENG",
    "sport_id": 1
}
```

> [!NOTE]
> `country_code` is optional. If provided, it must be a valid code present in the countries table.

**Response (`201 Created`):**
```json
{
    "competitor_id": 1
}
```

**Possible errors:**
- `400 Bad Request`: invalid JSON, invalid/non-existent `country_code`, or non-existent `sport_id`
- `500 Internal Server Error`: failure during database access

---

## Events

### `GET /api/events`
**Description:** Returns a list of all events, optionally filtered by date and/or sport

**Query Parameters:**
- `date` (string, not required) - Filter by date in `YYYY-MM-DD` format
- `sport` (string, not required) - Filter by sport name

**Response (`200 OK`):**
```json
[
    {
        "event_id": 1,
        "date": "2025/03/15",
        "time": "20:45",
        "status": "finished",
        "competition": "Premier League",
        "season": "2024/25",
        "stage": "Matchday 29",
        "home": "Manchester United",
        "away": "Chelsea",
        "result": {
            "home_points": 2,
            "away_points": 1,
            "winner": "home"
        }
    },
    ...
]
```

**Possible errors:**
- `400 Bad Request`: invalid date format
- `500 Internal Server Error`: failure during database access

---

### `GET /api/events/{eventID}`
**Description:** Returns detailed information about a single event by its ID

**Path Parameters:**
- `eventID` (integer) - ID of the event to retrieve

**Response (`200 OK`):**
```json
{
    "date": "2025/03/15",
    "time": "20:45",
    "venue": {
        "country_code": "ENG",
        "city": "Manchester",
        "place": "Old Trafford"
    },
    "status": "finished",
    "competition": "Premier League",
    "season": "2024/25",
    "stage": "Matchday 29",
    "home_country": "ENG",
    "home": "Manchester United",
    "away_country": "ENG",
    "away": "Chelsea",
    "result": {
        "home_points": 2,
        "away_points": 1,
        "outcome": null,
        "forfeit_by": null,
        "winner": "home",
        "result_details": "Goals: Rashford 12', Fernandes 67' | Cole Palmer 45'"
    },
    "event_details": "Derby match"
}
```

**Possible errors:**
- `400 Bad Request`: missing or non-integer event ID
- `404 Not Found`: event with the given ID does not exist
- `500 Internal Server Error`: failure during database access

---

### `POST /api/events`
**Description:** Creates a new event in the system, optionally with a result

**Request Body:**
```json
{
    "status": "finished",
    "venue_date": "2025-03-15T00:00:00Z",
    "venue_time": "2025-03-15T20:45:00Z",
    "venue_id": 1,
    "home_name": "Manchester United",
    "away_name": "Chelsea",
    "stage_id": 5,
    "details": "Derby match",
    "result": {
        "home_points": 2,
        "away_points": 1,
        "outcome": null,
        "forfeit_by": null,
        "result_details": "Goals: Rashford 12', Fernandes 67' | Cole Palmer 45'"
    }
}
```

> [!NOTE]
> - `venue_time`, `venue_id`, `home_name`, `away_name`, `details` and `result` are optional
> - `result` should be included when `status` is `finished` or `live`
> - `forfeit_by` is required when `outcome` is `"forfeit"`, and must be either `"home"` or `"away"`
> - Both competitors must belong to the same sport as the venue and competition stage
> - The combination of `home`, `away`, `venue_date` and `stage_id` must be unique

**Response (`201 Created`):**
```json
{
    "event_id": 1
}
```

**Possible errors:**
- `400 Bad Request`: invalid JSON, invalid status, missing date, same competitors, mismatched sports, invalid venue/competitor/stage
- `409 Conflict`: an event with the same home, away, date and stage already exists
- `500 Internal Server Error`: failure during database access or transaction

---

## Notes
- Error response format:
```json
{
    "code": "ERROR_CODE",
    "message": "Clear error message"
}
```
- Available event statuses: `scheduled`, `finished`, `live`, `cancelled`, `postponed`
- Available result outcomes: `no contest`, `forfeit`
- Available forfeit sides: `home`, `away`
- Winner field in responses is computed server-side and can be: `"home"`, `"away"`, or `""` (empty string, for draws and live events)
- Country codes must be at most **3 characters** long
- Request body size is limited to **1 KB** for most endpoints and **2 KB** for event creation
