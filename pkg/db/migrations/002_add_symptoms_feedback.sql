-- Таблица симптомов (связь человек → симптомы)
CREATE TABLE IF NOT EXISTS person_symptoms (
    person_id TEXT REFERENCES people(id),
    symptom_key TEXT,      -- "constipation", "alopecia"
    custom_label TEXT,     -- "облысение" (если ввёл своё)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (person_id, symptom_key)
);

-- Таблица фидбека (для RL-обучения)
CREATE TABLE IF NOT EXISTS feedback (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    person_id TEXT,
    intention_hash TEXT,   -- хэш сгенерированного текста
    action TEXT,           -- "copied", "printed", "ignored"
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    reward REAL            -- +1, +2, +5, -1
);

-- Таблица аффирмаций (граф знаний)
CREATE TABLE IF NOT EXISTS affirmations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author TEXT,           -- "Louise Hay", "Zhikarentsev"
    old_thought TEXT,
    new_thought TEXT,
    chakra_index INTEGER,  -- 0-6
    symptoms TEXT          -- JSON: ["constipation", "fear"]
);
