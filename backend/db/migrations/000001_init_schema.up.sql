CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username      VARCHAR(50) UNIQUE NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password      TEXT NOT NULL,
    bio           TEXT,
    avatar_url    TEXT,
    is_banned     BOOLEAN DEFAULT FALSE,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE teams (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(50) UNIQUE NOT NULL,
    title       VARCHAR(100) NOT NULL,
    description TEXT,
    banner_url  TEXT,
    icon_url    TEXT,
    is_private  BOOLEAN DEFAULT FALSE,
    created_by  UUID REFERENCES users(id),
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE team_members (
    user_id    UUID REFERENCES users(id) ON DELETE CASCADE,
    team_id    UUID REFERENCES teams(id) ON DELETE CASCADE,
    role       VARCHAR(20) DEFAULT 'member',
    joined_at  TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, team_id)
);

CREATE TABLE agents (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(100) NOT NULL,
    description   TEXT,
    system_prompt TEXT NOT NULL,
    model         VARCHAR(50) DEFAULT 'gpt-4o',
    tools         JSONB DEFAULT '[]',
    team_id       UUID REFERENCES teams(id) ON DELETE CASCADE,
    created_by    UUID REFERENCES users(id),
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE runs (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    agent_id     UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    team_id      UUID NOT NULL REFERENCES teams(id),
    started_by   UUID NOT NULL REFERENCES users(id),
    status       VARCHAR(20) DEFAULT 'pending',
    input        TEXT NOT NULL,
    output       TEXT,
    started_at   TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE TABLE events (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    run_id     UUID NOT NULL REFERENCES runs(id) ON DELETE CASCADE,
    sequence   INT NOT NULL,
    type       VARCHAR(50) NOT NULL,
    payload    JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_team_members_team_id ON team_members(team_id);
CREATE INDEX idx_agents_team_id ON agents(team_id);
CREATE INDEX idx_runs_agent_id ON runs(agent_id, started_at DESC);
CREATE INDEX idx_runs_team_id ON runs(team_id, started_at DESC);
CREATE INDEX idx_runs_status ON runs(status);
CREATE INDEX idx_events_run_id ON events(run_id, sequence);
