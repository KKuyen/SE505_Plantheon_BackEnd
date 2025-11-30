-- Create activity_keywords table
CREATE TABLE IF NOT EXISTS activity_keywords (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(100) NOT NULL DEFAULT 'GENERAL',
    base_days_offset INTEGER NOT NULL DEFAULT 0,
    is_free_time BOOLEAN NOT NULL DEFAULT false,
    hour_time INTEGER,
    time_duration INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create disease_activity_keywords junction table
CREATE TABLE IF NOT EXISTS disease_activity_keywords (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disease_id UUID NOT NULL,
    activity_keyword_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_disease FOREIGN KEY (disease_id) REFERENCES diseases(id) ON DELETE CASCADE,
    CONSTRAINT fk_activity_keyword FOREIGN KEY (activity_keyword_id) REFERENCES activity_keywords(id) ON DELETE CASCADE,
    UNIQUE(disease_id, activity_keyword_id)
);

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_activity_keywords_name ON activity_keywords(name);
CREATE INDEX IF NOT EXISTS idx_disease_activity_keywords_disease_id ON disease_activity_keywords(disease_id);
CREATE INDEX IF NOT EXISTS idx_disease_activity_keywords_activity_keyword_id ON disease_activity_keywords(activity_keyword_id);

-- Remove activity_keywords_ids column from diseases table if it exists
ALTER TABLE diseases DROP COLUMN IF EXISTS activity_keywords_ids;
