-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Seed DictClientApps
INSERT OR IGNORE INTO DictClientApps(ID, App) VALUES (1, 'telegram');
-- Add more clients as needed: (2, 'whatsapp'), (3, 'vk'), etc.

-- Seed ProcessingStatus
INSERT OR IGNORE INTO ProcessingStatus(ID, DisplayName) VALUES (0, 'queued');
INSERT OR IGNORE INTO ProcessingStatus(ID, DisplayName) VALUES (10, 'failed');
INSERT OR IGNORE INTO ProcessingStatus(ID, DisplayName) VALUES (20, 'loading');
INSERT OR IGNORE INTO ProcessingStatus(ID, DisplayName) VALUES (30, 'loaded');
INSERT OR IGNORE INTO ProcessingStatus(ID, DisplayName) VALUES (40, 'summarizing');
INSERT OR IGNORE INTO ProcessingStatus(ID, DisplayName) VALUES (50, 'summarized');
INSERT OR IGNORE INTO ProcessingStatus(ID, DisplayName) VALUES (60, 'completed');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- Note: Down migrations for seeded data are tricky.
-- Deleting seeded data might remove user data.
-- For simplicity, we'll leave the data. A real system might need a more complex approach.
-- DELETE FROM ProcessingStatus WHERE ID IN (0, 10, 20, 30, 40, 50, 60);
-- DELETE FROM DictClientApps WHERE ID = 1 AND App = 'telegram';
