CREATE TABLE IF NOT EXISTS sources (
    uuid VARCHAR(36) NOT NULL,
    created_at INT(11) NOT NULL, -- UNIX time
    created_by_uuid VARCHAR(36) NOT NULL,
    updated_at INT(11), -- UNIX time
    updated_by_uuid VARCHAR(36),
    name VARCHAR(255) NOT NULL,
    organization VARCHAR(1023) NOT NULL,
    PRIMARY KEY(uuid)
);

-- CREATE UNIQUE INDEX idx_source_name_org ON srcabl_sources.sources(name, organization);

CREATE TABLE IF NOT EXISTS source_credibility (
    source_uuid VARCHAR(36) NOT NULL,
    positive_count INT NOT NULL,
    negative_count INT NOT NULL,
    FOREIGN KEY(source_uuid) REFERENCES srcabl_sources.sources(uuid)
);
