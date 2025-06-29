package history

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "path/filepath"
    "time"

    _ "modernc.org/sqlite"
)

type Record struct {
    ID        string `json:"id"`
    Prompt    string `json:"prompt"`
    Response  string `json:"response"`
    Model     string `json:"model"`
    Timestamp string `json:"timestamp"`
    Success   bool   `json:"success"`
}

type Store struct {
    db *sql.DB
}

func New(baseDir string) (*Store, error) {
    path := filepath.Join(baseDir, "state", "history.db")
    db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
    if err != nil {
        return nil, fmt.Errorf("open db: %w", err)
    }
    if err := initSchema(db); err != nil {
        db.Close()
        return nil, err
    }
    return &Store{db: db}, nil
}

func initSchema(db *sql.DB) error {
    schema := `
CREATE TABLE IF NOT EXISTS history(
    id TEXT PRIMARY KEY,
    prompt TEXT NOT NULL,
    response TEXT NOT NULL,
    model TEXT NOT NULL,
    timestamp TEXT NOT NULL,
    success INTEGER NOT NULL
);
CREATE VIRTUAL TABLE IF NOT EXISTS history_fts USING fts5(prompt, response, content='history', content_rowid='rowid');
CREATE TRIGGER IF NOT EXISTS history_ai AFTER INSERT ON history BEGIN
    INSERT INTO history_fts(rowid,prompt,response) VALUES(new.rowid,new.prompt,new.response);
END;
CREATE TRIGGER IF NOT EXISTS history_ad AFTER DELETE ON history BEGIN
    INSERT INTO history_fts(history_fts,rowid,prompt,response) VALUES('delete',old.rowid,old.prompt,old.response);
END;
CREATE TRIGGER IF NOT EXISTS history_au AFTER UPDATE ON history BEGIN
    INSERT INTO history_fts(history_fts,rowid,prompt,response) VALUES('delete',old.rowid,old.prompt,old.response);
    INSERT INTO history_fts(rowid,prompt,response) VALUES(new.rowid,new.prompt,new.response);
END;`
    if _, err := db.Exec(schema); err != nil {
        return fmt.Errorf("init schema: %w", err)
    }
    return nil
}

func (s *Store) Close() error {
    return s.db.Close()
}

func boolToInt(b bool) int {
    if b {
        return 1
    }
    return 0
}

func intToBool(i int) bool {
    return i != 0
}

func (s *Store) Add(r Record) error {
    if r.ID == "" {
        return fmt.Errorf("id required")
    }
    if r.Timestamp == "" {
        r.Timestamp = time.Now().UTC().Format(time.RFC3339)
    }
    _, err := s.db.Exec(`INSERT INTO history(id,prompt,response,model,timestamp,success) VALUES(?,?,?,?,?,?)`,
        r.ID, r.Prompt, r.Response, r.Model, r.Timestamp, boolToInt(r.Success))
    if err != nil {
        return fmt.Errorf("insert: %w", err)
    }
    return nil
}

func scanRows(rows *sql.Rows) ([]Record, error) {
    var recs []Record
    for rows.Next() {
        var r Record
        var success int
        if err := rows.Scan(&r.ID, &r.Prompt, &r.Response, &r.Model, &r.Timestamp, &success); err != nil {
            return nil, fmt.Errorf("scan: %w", err)
        }
        r.Success = intToBool(success)
        recs = append(recs, r)
    }
    return recs, nil
}

func (s *Store) All() ([]Record, error) {
    rows, err := s.db.Query(`SELECT id,prompt,response,model,timestamp,success FROM history ORDER BY timestamp DESC`)
    if err != nil {
        return nil, fmt.Errorf("query: %w", err)
    }
    defer rows.Close()
    return scanRows(rows)
}

func (s *Store) Search(q string, limit int) ([]Record, error) {
    if limit <= 0 {
        limit = 20
    }
    rows, err := s.db.Query(`SELECT id,prompt,response,model,timestamp,success FROM history WHERE rowid IN (SELECT rowid FROM history_fts WHERE history_fts MATCH ? LIMIT ?) ORDER BY timestamp DESC`, q, limit)
    if err != nil {
        return nil, fmt.Errorf("search: %w", err)
    }
    defer rows.Close()
    return scanRows(rows)
}

func (s *Store) Export() ([]byte, error) {
    recs, err := s.All()
    if err != nil {
        return nil, err
    }
    data, err := json.MarshalIndent(recs, "", "  ")
    if err != nil {
        return nil, fmt.Errorf("marshal: %w", err)
    }
    return data, nil
}

func (s *Store) Import(data []byte) error {
    var recs []Record
    if err := json.Unmarshal(data, &recs); err != nil {
        return fmt.Errorf("unmarshal: %w", err)
    }
    tx, err := s.db.Begin()
    if err != nil {
        return fmt.Errorf("begin: %w", err)
    }
    stmt, err := tx.Prepare(`INSERT OR REPLACE INTO history(id,prompt,response,model,timestamp,success) VALUES(?,?,?,?,?,?)`)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("prepare: %w", err)
    }
    defer stmt.Close()
    for _, r := range recs {
        if r.Timestamp == "" {
            r.Timestamp = time.Now().UTC().Format(time.RFC3339)
        }
        if _, err := stmt.Exec(r.ID, r.Prompt, r.Response, r.Model, r.Timestamp, boolToInt(r.Success)); err != nil {
            tx.Rollback()
            return fmt.Errorf("exec: %w", err)
        }
    }
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit: %w", err)
    }
    return nil
}

