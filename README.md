# osrs-db
osrs-db is a Go application that pulls Old School RuneScape game data from the OSRS Wiki
and loads it into a local SQLite database. Each game object is fetched, normalized, and
stored as a single record — so tools like RuneLite plugins can query everything they need
about an item or entity in one place instead of piecing it together from raw wiki responses.

---

## How It Works

Fetches data from the OSRS Wiki API, maps it into consistent Go structs, and writes it to
SQLite. 

## Whats used
- **Go**
- **OSRS Wiki API**
- **SQLite**

---

