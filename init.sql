-- Create tables for rule management
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    project_id VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    language VARCHAR(50),
    apply_global_rules BOOLEAN DEFAULT true,
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

-- Insert sample data
INSERT INTO projects (project_id, name, description, language, apply_global_rules) VALUES
    ('default', 'Default Project', 'Default project with common rules', 'general', true),
    ('web-app', 'Web Application', 'Web application specific rules', 'javascript', true),
    ('api-service', 'API Service', 'API service specific rules', 'go', true)
ON CONFLICT (project_id) DO NOTHING;

-- Insert global rules for different languages
INSERT INTO global_rules (language, rule_id, name, description, type, severity, pattern, message) VALUES
    ('general', 'no-hardcoded-secrets', 'No Hardcoded Secrets', 'API keys, passwords, and other secrets should not be hardcoded in source code', 'security', 'error', 'api_key', 'Hardcoded API key detected. Use environment variables instead.'),
    ('general', 'no-sql-injection', 'No SQL Injection', 'Raw SQL queries should not be constructed by string concatenation', 'security', 'error', 'SELECT * FROM', 'Raw SQL query detected. Use parameterized queries or ORM.'),
    ('javascript', 'no-console-log', 'No Console Log', 'Console.log statements should not be in production code', 'style', 'warning', 'console.log', 'Console.log detected. Use proper logging framework in production.'),
    ('javascript', 'no-inline-styles', 'No Inline Styles', 'CSS styles should be in separate stylesheets, not inline', 'style', 'warning', 'style=\"', 'Inline styles detected. Move to CSS file.'),
    ('go', 'naming-convention', 'Naming Convention', 'Functions and variables should use camelCase', 'style', 'warning', 'function_name', 'Function name should use camelCase (e.g., functionName).'),
    ('go', 'error-handling', 'Error Handling', 'All async operations must have proper error handling', 'reliability', 'warning', 'if err != nil', 'Error handling required. Check if err != nil.'),
    ('python', 'no-print', 'No Print Statements', 'Print statements should not be in production code', 'style', 'warning', 'print(', 'Print statement detected. Use proper logging framework in production.'),
    ('python', 'type-hints', 'Type Hints Required', 'Function parameters should have type hints', 'style', 'warning', 'def ', 'Function definition without type hints detected.')
ON CONFLICT (language, rule_id) DO NOTHING;

INSERT INTO rules (project_id, rule_id, name, description, type, severity, pattern, message) VALUES
    ('default', 'no-hardcoded-secrets', 'No Hardcoded Secrets', 'API keys, passwords, and other secrets should not be hardcoded in source code', 'security', 'error', 'api_key', 'Hardcoded API key detected. Use environment variables instead.'),
    ('default', 'no-sql-injection', 'No SQL Injection', 'Raw SQL queries should not be constructed by string concatenation', 'security', 'error', 'SELECT * FROM', 'Raw SQL query detected. Use parameterized queries or ORM.'),
    ('default', 'naming-convention', 'Naming Convention', 'Functions and variables should use camelCase', 'style', 'warning', 'function_name', 'Function name should use camelCase (e.g., functionName).'),
    ('web-app', 'no-console-log', 'No Console Log', 'Console.log statements should not be in production code', 'style', 'warning', 'console.log', 'Console.log detected. Use proper logging framework in production.'),
    ('web-app', 'no-inline-styles', 'No Inline Styles', 'CSS styles should be in separate stylesheets, not inline', 'style', 'warning', 'style="', 'Inline styles detected. Move to CSS file.'),
    ('api-service', 'input-validation', 'Input Validation', 'All API inputs must be validated', 'security', 'error', 'req.body', 'Direct access to req.body without validation detected.'),
    ('api-service', 'error-handling', 'Error Handling', 'All async operations must have proper error handling', 'reliability', 'warning', 'catch (', 'Async operation without proper error handling detected.')
ON CONFLICT (project_id, rule_id) DO NOTHING;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_rules_project_id ON rules(project_id);
CREATE INDEX IF NOT EXISTS idx_rule_violations_project_id ON rule_violations(project_id);
CREATE INDEX IF NOT EXISTS idx_rule_violations_rule_id ON rule_violations(rule_id);
