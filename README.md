# osrs-db

osrs-db is a REST API built in Go that aggregates and normalizes Old School RuneScape game 
data from the official OSRS Wiki into clean, structured records. Rather than consuming raw, 
inconsistently formatted wiki data directly, osrs-db runs an ETL pipeline that fetches, 
transforms, and stores game data into unified structs backed by SQLite making it easy for 
tools like RuneLite plugins to query reliable, well-shaped data.

---

## How It Works

### ETL Pipeline
Fetches raw game data from the OSRS Wiki API and transforms it into a consistent, normalized 
struct format. Rather than requiring multiple API calls across different wiki endpoints to 
assemble a complete picture of a single game object, the pipeline consolidates all relevant 
data into a single unified record so RuneLite plugins and other tools can retrieve 
everything they need about a game object in a single query.

### Data Models
Unified Go structs define the shape of each game entity, ensuring that data stored in the 
database follows a predictable schema regardless of how it was represented on the wiki.

### REST API *(planned)*
A REST layer will expose the normalized data through clean JSON endpoints, allowing external 
tools and plugins to query OSRS game data without dealing with raw wiki formatting.

---

## Tech Stack

- **Go** — Core application logic and API server
- **OSRS Wiki API** — Source of raw game data
- **SQLite** — Lightweight persistent storage for normalized game data

---

## Status

Active development. ETL pipeline and data models are in progress. Database integration and 
REST endpoints are planned next.

---
