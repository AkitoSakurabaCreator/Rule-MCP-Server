package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AkitoSakurabaCreator/RuleMCPServer/internal/domain"
)

// ProjectDetector プロジェクト自動検出ユースケース
type ProjectDetector struct {
	projectRepo domain.ProjectRepository
	ruleRepo    domain.RuleRepository
}

// NewProjectDetector プロジェクト検出器を作成
func NewProjectDetector(projectRepo domain.ProjectRepository, ruleRepo domain.RuleRepository) *ProjectDetector {
	return &ProjectDetector{
		projectRepo: projectRepo,
		ruleRepo:    ruleRepo,
	}
}

// DetectionResult 検出結果
type DetectionResult struct {
	Project         *domain.Project `json:"project"`
	Rules           []*domain.Rule  `json:"rules"`
	DetectionMethod string          `json:"detection_method"`
	Confidence      float64         `json:"confidence"`
	Message         string          `json:"message"`
}

// AutoDetectProject プロジェクトを自動検出
func (pd *ProjectDetector) AutoDetectProject(path string) (*DetectionResult, error) {
	// 1. ディレクトリ名から検索（最優先）
	if project, err := pd.detectFromDirectoryName(path); err == nil {
		rules, _ := pd.ruleRepo.GetByProjectID(project.ProjectID)
		return &DetectionResult{
			Project:         project,
			Rules:           rules,
			DetectionMethod: "directory_name",
			Confidence:      0.95,
			Message:         fmt.Sprintf("ディレクトリ名 '%s' からプロジェクトを検出しました", filepath.Base(path)),
		}, nil
	}

	// 2. Gitリポジトリ名から検索
	if project, err := pd.detectFromGit(path); err == nil {
		rules, _ := pd.ruleRepo.GetByProjectID(project.ProjectID)
		return &DetectionResult{
			Project:         project,
			Rules:           rules,
			DetectionMethod: "git_repository",
			Confidence:      0.90,
			Message:         fmt.Sprintf("Gitリポジトリ名からプロジェクトを検出しました"),
		}, nil
	}

	// 3. 言語固有ファイルから検索
	if project, err := pd.detectFromLanguageFiles(path); err == nil {
		rules, _ := pd.ruleRepo.GetByProjectID(project.ProjectID)
		return &DetectionResult{
			Project:         project,
			Rules:           rules,
			DetectionMethod: "language_files",
			Confidence:      0.85,
			Message:         fmt.Sprintf("言語固有ファイルからプロジェクトを検出しました"),
		}, nil
	}

	// 4. デフォルトプロジェクトを返す
	if project, err := pd.getDefaultProject(); err == nil {
		rules, _ := pd.ruleRepo.GetByProjectID(project.ProjectID)
		return &DetectionResult{
			Project:         project,
			Rules:           rules,
			DetectionMethod: "default_project",
			Confidence:      0.70,
			Message:         "デフォルトプロジェクトを使用します",
		}, nil
	}

	return nil, fmt.Errorf("プロジェクトを検出できませんでした: %s", path)
}

// detectFromDirectoryName ディレクトリ名からプロジェクトを検出
func (pd *ProjectDetector) detectFromDirectoryName(path string) (*domain.Project, error) {
	dirName := filepath.Base(path)

	// 一般的な除外ディレクトリ
	excludeDirs := []string{"node_modules", "vendor", "dist", "build", "target", ".git", ".vscode"}
	for _, exclude := range excludeDirs {
		if dirName == exclude {
			return nil, fmt.Errorf("除外ディレクトリ: %s", dirName)
		}
	}

	// プロジェクトIDとしてディレクトリ名を検索
	project, err := pd.projectRepo.GetByID(dirName)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリ名 '%s' のプロジェクトが見つかりません", dirName)
	}

	return project, nil
}

// detectFromGit Gitリポジトリ名からプロジェクトを検出
func (pd *ProjectDetector) detectFromGit(path string) (*domain.Project, error) {
	gitConfigPath := filepath.Join(path, ".git", "config")

	if _, err := os.Stat(gitConfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Gitリポジトリではありません")
	}

	// Git設定ファイルを読み込み
	data, err := os.ReadFile(gitConfigPath)
	if err != nil {
		return nil, fmt.Errorf("Git設定ファイルの読み込みに失敗: %v", err)
	}

	// origin URLからリポジトリ名を抽出
	content := string(data)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "url = ") {
			url := strings.TrimPrefix(line, "url = ")
			repoName := pd.extractRepoNameFromURL(url)
			if repoName != "" {
				// リポジトリ名でプロジェクトを検索
				if project, err := pd.projectRepo.GetByID(repoName); err == nil {
					return project, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("Gitリポジトリ名からプロジェクトを特定できませんでした")
}

// extractRepoNameFromURL URLからリポジトリ名を抽出
func (pd *ProjectDetector) extractRepoNameFromURL(url string) string {
	// SSH形式: git@github.com:username/repo-name.git
	if strings.Contains(url, ":") && strings.Contains(url, "@") {
		parts := strings.Split(url, ":")
		if len(parts) == 2 {
			repoPart := parts[1]
			if strings.HasSuffix(repoPart, ".git") {
				repoPart = strings.TrimSuffix(repoPart, ".git")
			}
			return repoPart
		}
	}

	// HTTPS形式: https://github.com/username/repo-name.git
	if strings.HasPrefix(url, "http") {
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			repoPart := parts[len(parts)-1]
			if strings.HasSuffix(repoPart, ".git") {
				repoPart = strings.TrimSuffix(repoPart, ".git")
			}
			return repoPart
		}
	}

	return ""
}

// detectFromLanguageFiles 言語固有ファイルからプロジェクトを検出
func (pd *ProjectDetector) detectFromLanguageFiles(path string) (*domain.Project, error) {
	// 言語別の設定ファイル
	languageFiles := map[string]string{
		"go":     "go.mod",
		"node":   "package.json",
		"python": "requirements.txt",
		"java":   "pom.xml",
		"rust":   "Cargo.toml",
		"php":    "composer.json",
		"ruby":   "Gemfile",
	}

	for language, filename := range languageFiles {
		filePath := filepath.Join(path, filename)
		if _, err := os.Stat(filePath); err == nil {
			// 言語に基づいてプロジェクトを検索
			projects, err := pd.projectRepo.GetByLanguage(language)
			if err == nil && len(projects) > 0 {
				// 最初に見つかったプロジェクトを返す
				return projects[0], nil
			}
		}
	}

	return nil, fmt.Errorf("言語固有ファイルからプロジェクトを特定できませんでした")
}

// getDefaultProject デフォルトプロジェクトを取得
func (pd *ProjectDetector) getDefaultProject() (*domain.Project, error) {
	// "default"プロジェクトを検索
	project, err := pd.projectRepo.GetByID("default")
	if err != nil {
		return nil, fmt.Errorf("デフォルトプロジェクトが見つかりません")
	}
	return project, nil
}

// ScanLocalProjects ローカルディレクトリをスキャンしてプロジェクトを検出
func (pd *ProjectDetector) ScanLocalProjects(basePath string) ([]DetectionResult, error) {
	var results []DetectionResult

	// ベースパス内のディレクトリを再帰的にスキャン
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ディレクトリのみ処理
		if !info.IsDir() {
			return nil
		}

		// 除外ディレクトリをスキップ
		excludeDirs := []string{".git", "node_modules", "vendor", "dist", "build", "target", ".vscode"}
		dirName := filepath.Base(path)
		for _, exclude := range excludeDirs {
			if dirName == exclude {
				return filepath.SkipDir
			}
		}

		// プロジェクトを検出
		if result, err := pd.AutoDetectProject(path); err == nil {
			results = append(results, *result)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ローカルディレクトリのスキャンに失敗: %v", err)
	}

	return results, nil
}
