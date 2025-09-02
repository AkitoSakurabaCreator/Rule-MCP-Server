export interface Project {
  project_id: string;
  name: string;
  description: string;
  language: string;
  apply_global_rules: boolean;
  created_at: string;
  updated_at: string;
  rule_count: number;
}

export interface Rule {
  id: number;
  project_id: string;
  rule_id: string;
  name: string;
  description: string;
  type: string;
  severity: string;
  pattern: string;
  message: string;
  is_active: boolean;
}

export interface GlobalRule {
  id: number;
  language: string;
  rule_id: string;
  name: string;
  description: string;
  type: string;
  severity: string;
  pattern: string;
  message: string;
  is_active: boolean;
}

export interface ValidationResult {
  valid: boolean;
  errors: string[];
  warnings: string[];
}

export interface ProjectRules {
  project_id: string;
  rules: Rule[];
}
