-- Create tables for rule management
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    project_id VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    language VARCHAR(50),
    apply_global_rules BOOLEAN DEFAULT true,
    access_level VARCHAR(20) DEFAULT 'public',
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS global_rules (
    id SERIAL PRIMARY KEY,
    language VARCHAR(50) NOT NULL,
    rule_id VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    pattern TEXT NOT NULL,
    message TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    access_level VARCHAR(20) DEFAULT 'public',
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(language, rule_id)
);

CREATE TABLE IF NOT EXISTS rules (
    id SERIAL PRIMARY KEY,
    project_id VARCHAR(100) NOT NULL,
    rule_id VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    pattern TEXT NOT NULL,
    message TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    access_level VARCHAR(20) DEFAULT 'public',
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
    UNIQUE(project_id, rule_id)
);

CREATE TABLE IF NOT EXISTS rule_violations (
    id SERIAL PRIMARY KEY,
    project_id VARCHAR(100) NOT NULL,
    rule_id INTEGER NOT NULL,
    code_snippet TEXT,
    file_path VARCHAR(500),
    line_number INTEGER,
    severity VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
    FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE CASCADE
);

-- Create API keys table for authentication
CREATE TABLE IF NOT EXISTS api_keys (
    id SERIAL PRIMARY KEY,
    key_hash VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    access_level VARCHAR(20) NOT NULL DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    expires_at TIMESTAMP,
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create users table for team management
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create project_members table for team access control
CREATE TABLE IF NOT EXISTS project_members (
    id SERIAL PRIMARY KEY,
    project_id VARCHAR(100) NOT NULL,
    user_id INTEGER NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    permissions JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(project_id, user_id)
);

-- Insert sample data
INSERT INTO projects (project_id, name, description, language, apply_global_rules, access_level, created_by) VALUES
    ('default', 'Default Project', 'Default project with common rules', 'general', true, 'public', 'system'),
    ('web-app', 'Web Application', 'Web application specific rules', 'javascript', true, 'public', 'system'),
    ('api-service', 'API Service', 'API service specific rules', 'go', true, 'public', 'system'),
    ('team-project', 'Team Project', 'Team collaboration project', 'typescript', true, 'user', 'admin')
ON CONFLICT (project_id) DO NOTHING;

-- Insert initial admin account (default password: admin123)
INSERT INTO users (username, email, full_name, role) VALUES
    ('admin', 'admin@rulemcp.com', 'System Administrator', 'admin')
ON CONFLICT (username) DO NOTHING;

-- Insert sample users
INSERT INTO users (username, email, full_name, role) VALUES
    ('developer1', 'dev1@example.com', 'Developer 1', 'user'),
    ('developer2', 'dev2@example.com', 'Developer 2', 'user')
ON CONFLICT (username) DO NOTHING;

-- Insert sample API keys (hashed values - in production, use proper hashing)
INSERT INTO api_keys (key_hash, name, description, access_level, created_by) VALUES
    ('admin_key_hash_123', 'Admin API Key', 'Full access to all features', 'admin', 'admin'),
    ('user_key_hash_456', 'User API Key', 'Standard user access', 'user', 'admin'),
    ('public_key_hash_789', 'Public API Key', 'Public read-only access', 'public', 'admin')
ON CONFLICT (key_hash) DO NOTHING;

-- Insert project members
INSERT INTO project_members (project_id, user_id, role, permissions) VALUES
    ('team-project', 2, 'member', '{"read": true, "write": true, "delete": false}'),
    ('team-project', 3, 'member', '{"read": true, "write": true, "delete": false}')
ON CONFLICT (project_id, user_id) DO NOTHING;

-- Insert global rules for different languages
INSERT INTO global_rules (language, rule_id, name, description, type, severity, pattern, message, access_level, created_by) VALUES
    ('general', 'no-hardcoded-secrets', 'No Hardcoded Secrets', 'API keys, passwords, and other secrets should not be hardcoded in source code', 'security', 'error', 'api_key', 'Hardcoded API key detected. Use environment variables instead.', 'public', 'system'),
    ('general', 'no-sql-injection', 'No SQL Injection', 'Raw SQL queries should not be constructed by string concatenation', 'security', 'error', 'SELECT * FROM', 'Raw SQL query detected. Use parameterized queries or ORM.', 'public', 'system'),
    ('javascript', 'no-console-log', 'No Console Log', 'Console.log statements should not be in production code', 'style', 'warning', 'console.log', 'Console.log detected. Use proper logging framework in production.', 'public', 'system'),
    ('javascript', 'no-inline-styles', 'No Inline Styles', 'CSS styles should be in separate stylesheets, not inline', 'style', 'warning', 'style=\"', 'Inline styles detected. Move to CSS file.', 'public', 'system'),
    ('go', 'naming-convention', 'Naming Convention', 'Functions and variables should use camelCase', 'style', 'warning', 'function_name', 'Function name should use camelCase (e.g., functionName).', 'public', 'system'),
    ('go', 'error-handling', 'Error Handling', 'All async operations must have proper error handling', 'reliability', 'warning', 'if err != nil', 'Error handling required. Check if err != nil.', 'public', 'system'),
    ('python', 'no-print', 'No Print Statements', 'Print statements should not be in production code', 'style', 'warning', 'print(', 'Print statement detected. Use proper logging framework in production.', 'public', 'system'),
    ('python', 'type-hints', 'Type Hints Required', 'Function parameters should have type hints', 'style', 'warning', 'def ', 'Function definition without type hints detected.', 'public', 'system'),
    ('typescript', 'strict-null-checks', 'Strict Null Checks', 'Enable strict null checks for better type safety', 'quality', 'warning', 'strictNullChecks', 'Enable strict null checks in tsconfig.json', 'user', 'admin')
ON CONFLICT (language, rule_id) DO NOTHING;

INSERT INTO rules (project_id, rule_id, name, description, type, severity, pattern, message, access_level, created_by) VALUES
    ('default', 'no-hardcoded-secrets', 'No Hardcoded Secrets', 'API keys, passwords, and other secrets should not be hardcoded in source code', 'security', 'error', 'api_key', 'Hardcoded API key detected. Use environment variables instead.', 'public', 'system'),
    ('default', 'no-sql-injection', 'No SQL Injection', 'Raw SQL queries should not be constructed by string concatenation', 'security', 'error', 'SELECT * FROM', 'Raw SQL query detected. Use parameterized queries or ORM.', 'public', 'system'),
    ('default', 'naming-convention', 'Naming Convention', 'Functions and variables should use camelCase', 'style', 'warning', 'function_name', 'Function name should use camelCase (e.g., functionName).', 'public', 'system'),
    ('web-app', 'no-console-log', 'No Console Log', 'Console.log statements should not be in production code', 'style', 'warning', 'console.log', 'Console.log detected. Use proper logging framework in production.', 'public', 'system'),
    ('web-app', 'no-inline-styles', 'No Inline Styles', 'CSS styles should be in separate stylesheets, not inline', 'style', 'warning', 'style="', 'Inline styles detected. Move to CSS file.', 'public', 'system'),
    ('api-service', 'input-validation', 'Input Validation', 'All API inputs must be validated', 'security', 'error', 'req.body', 'Direct access to req.body without validation detected.', 'public', 'system'),
    ('api-service', 'error-handling', 'Error Handling', 'All async operations must have proper error handling', 'reliability', 'warning', 'catch (', 'Async operation without proper error handling detected.', 'public', 'system'),
    ('team-project', 'code-review-required', 'Code Review Required', 'All code changes must go through code review', 'process', 'error', 'TODO:', 'Code review required before merging', 'user', 'admin')
ON CONFLICT (project_id, rule_id) DO NOTHING;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_rules_project_id ON rules(project_id);
CREATE INDEX IF NOT EXISTS idx_rule_violations_project_id ON rule_violations(project_id);
CREATE INDEX IF NOT EXISTS idx_rule_violations_rule_id ON rule_violations(rule_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_project_members_project_id ON project_members(project_id);
CREATE INDEX IF NOT EXISTS idx_project_members_user_id ON project_members(user_id);
