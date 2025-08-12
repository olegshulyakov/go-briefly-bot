-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

PRAGMA foreign_keys = ON; -- Ensure foreign key constraints are enforced

-- Client Applications
CREATE TABLE IF NOT EXISTS DictClientApps (
    ID TINYINT NOT NULL,
    App VARCHAR(64) NOT NULL,
    PRIMARY KEY (ID, App)
);

-- Processing Statuses
CREATE TABLE IF NOT EXISTS ProcessingStatus (
    ID TINYINT NOT NULL PRIMARY KEY,
    DisplayName VARCHAR(32) NOT NULL
);

-- Message History (Sharded per ClientAppID)
-- Note: Sharding logic will be handled in Go code when connecting/creating DB files.
CREATE TABLE IF NOT EXISTS MessageHistory (
    ClientAppID TINYINT NOT NULL,
    MessageID BIGINT NOT NULL,
    UserID BIGINT NOT NULL,
    UserName VARCHAR(256),
    UserLanguage VARCHAR(2) NOT NULL,
    MessageContent VARCHAR(4096) NOT NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (ClientAppID, MessageID, UserID)
    -- Foreign key constraint handled at application level due to sharding
    -- FOREIGN KEY (ClientAppID) REFERENCES DictClientApps(ID)
);

-- Processing Queue
CREATE TABLE IF NOT EXISTS ProcessingQueue (
    ClientAppID TINYINT NOT NULL,
    MessageID BIGINT NOT NULL,
    UserID BIGINT NOT NULL,
    Url VARCHAR(256) NOT NULL,
    Language VARCHAR(2) NOT NULL,
    StatusID TINYINT NOT NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ProcessedAt TIMESTAMP NULL,
    RetryCount TINYINT NULL DEFAULT 0,
    ErrorMessage VARCHAR(1024) NULL,

    PRIMARY KEY (ClientAppID, MessageID, UserID, Url),
    -- Foreign key constraints handled at application level due to sharding
    -- FOREIGN KEY (ClientAppID) REFERENCES DictClientApps(ID),
    FOREIGN KEY (StatusID) REFERENCES ProcessingStatus(ID)
);

-- Sources Cache
CREATE TABLE IF NOT EXISTS Sources (
    Url VARCHAR(256) NOT NULL,
    Language VARCHAR(2) NOT NULL,
    Title VARCHAR(256) NOT NULL,
    Text TEXT NOT NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (Url, Language)
);

-- Summaries Cache
CREATE TABLE IF NOT EXISTS Summaries (
    Url VARCHAR(256) NOT NULL,
    Language VARCHAR(2) NOT NULL,
    Summary TEXT NOT NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (Url, Language)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS Summaries;
DROP TABLE IF EXISTS Sources;
DROP TABLE IF EXISTS ProcessingQueue;
DROP TABLE IF EXISTS MessageHistory;
DROP TABLE IF EXISTS ProcessingStatus;
DROP TABLE IF EXISTS DictClientApps;
