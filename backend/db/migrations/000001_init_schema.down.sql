DROP INDEX IF EXISTS idx_events_run_id;
DROP INDEX IF EXISTS idx_runs_status;
DROP INDEX IF EXISTS idx_runs_team_id;
DROP INDEX IF EXISTS idx_runs_agent_id;
DROP INDEX IF EXISTS idx_agents_team_id;
DROP INDEX IF EXISTS idx_team_members_team_id;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;

DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS runs;
DROP TABLE IF EXISTS agents;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS vector;
DROP EXTENSION IF EXISTS "uuid-ossp";
