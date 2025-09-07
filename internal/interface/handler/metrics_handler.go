package handler

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// MCPリクエストメトリクス
	mcpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rule_mcp_requests_total",
			Help: "Total number of MCP requests",
		},
		[]string{"method", "status"},
	)

	mcpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "rule_mcp_request_duration_seconds",
			Help:    "Duration of MCP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	// システム統計メトリクス
	totalUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "rule_mcp_total_users",
			Help: "Total number of users",
		},
	)

	totalProjects = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "rule_mcp_total_projects",
			Help: "Total number of projects",
		},
	)

	totalRules = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "rule_mcp_total_rules",
			Help: "Total number of rules",
		},
	)

	activeSessions = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "rule_mcp_active_sessions",
			Help: "Number of active sessions",
		},
	)

	activeApiKeys = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "rule_mcp_active_api_keys",
			Help: "Number of active API keys",
		},
	)

	systemLoadPercent = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "rule_mcp_system_load_percent",
			Help: "System load percentage",
		},
	)
)

type MetricsHandler struct {
	userRepo    domain.UserRepository
	projectRepo domain.ProjectRepository
	ruleRepo    domain.RuleRepository
	metricsRepo domain.MetricsRepository
}

func NewMetricsHandler(userRepo domain.UserRepository, projectRepo domain.ProjectRepository, ruleRepo domain.RuleRepository, metricsRepo domain.MetricsRepository) *MetricsHandler {
	return &MetricsHandler{
		userRepo:    userRepo,
		projectRepo: projectRepo,
		ruleRepo:    ruleRepo,
		metricsRepo: metricsRepo,
	}
}

// RecordMCPRequest MCPリクエストを記録
func (h *MetricsHandler) RecordMCPRequest(method, status string, duration time.Duration) {
	mcpRequestsTotal.WithLabelValues(method, status).Inc()
	mcpRequestDuration.WithLabelValues(method, status).Observe(duration.Seconds())
}

// UpdateSystemMetrics システムメトリクスを更新
func (h *MetricsHandler) UpdateSystemMetrics() {
	// ユーザー数
	if h.userRepo != nil {
		if users, err := h.userRepo.GetAll(); err == nil {
			totalUsers.Set(float64(len(users)))
		}
	}

	// プロジェクト数
	if h.projectRepo != nil {
		if projects, err := h.projectRepo.GetAll(); err == nil {
			totalProjects.Set(float64(len(projects)))
		}
	}

	// ルール数
	if h.ruleRepo != nil && h.projectRepo != nil {
		if projects, err := h.projectRepo.GetAll(); err == nil {
			totalRulesCount := 0
			for _, p := range projects {
				if rules, err := h.ruleRepo.GetByProjectID(p.ProjectID); err == nil {
					totalRulesCount += len(rules)
				}
			}
			totalRules.Set(float64(totalRulesCount))
		}
	}

	// システム負荷
	if load, err := getSystemLoad(); err == nil {
		systemLoadPercent.Set(load)
	}
}

// SetActiveSessions アクティブセッション数を設定
func (h *MetricsHandler) SetActiveSessions(count int) {
	activeSessions.Set(float64(count))
}

// SetActiveApiKeys アクティブAPIキー数を設定
func (h *MetricsHandler) SetActiveApiKeys(count int) {
	activeApiKeys.Set(float64(count))
}

// Metrics Prometheusメトリクスエンドポイント
func (h *MetricsHandler) Metrics(c *gin.Context) {
	h.UpdateSystemMetrics()
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

// getSystemLoad システム負荷を取得
func getSystemLoad() (float64, error) {
	// Linuxの/proc/loadavgから負荷を取得
	if b, err := os.ReadFile("/proc/loadavg"); err == nil {
		parts := strings.Fields(string(b))
		if len(parts) > 0 {
			if load1, err := strconv.ParseFloat(parts[0], 64); err == nil {
				cores := float64(runtime.NumCPU())
				return (load1 / cores) * 100, nil
			}
		}
	}
	return 0, nil
}
