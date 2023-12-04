package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/shubham-yadavv/log-ingestor/database"
	"gorm.io/gorm"
)

func SearchLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params := r.URL.Query()

	startDate := params.Get("start_date")
	endDate := params.Get("end_date")

	var startDateTime, endDateTime time.Time
	var err error

	if startDate != "" {
		startDateTime, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	}

	if endDate != "" {
		endDateTime, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
	}

	query := database.DB
	query = applyFilter(query, "level", params.Get("level"))
	query = applyFilter(query, "message", params.Get("message"))
	query = applyFilter(query, "resource_id", params.Get("resource_id"))
	query = applyFilter(query, "trace_id", params.Get("trace_id"))
	query = applyFilter(query, "span_id", params.Get("span_id"))
	query = applyFilter(query, "commit", params.Get("commit"))
	query = applyFilter(query, "parent_resource_id", params.Get("parent_resource_id"))
	query = applyDateRange(query, "timestamp", startDateTime, endDateTime)

	messageRegex := params.Get("message_regex")

	var logs []database.Log
	if err := query.Find(&logs).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving logs: %v", err), http.StatusInternalServerError)
		return
	}

	var filteredLogs []database.Log
	for _, log := range logs {
		if messageRegex != "" && !regexp.MustCompile(messageRegex).MatchString(log.Message) {
			continue
		}
		filteredLogs = append(filteredLogs, log)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredLogs)
}

func applyFilter(query *gorm.DB, field, value string) *gorm.DB {
	if value != "" {
		return query.Where(fmt.Sprintf("%s = ?", field), value)
	}
	return query
}


func applyDateRange(query *gorm.DB, field string, startDate, endDate time.Time) *gorm.DB {
	if !startDate.IsZero() {
		query = query.Where(fmt.Sprintf("%s >= ?", field), startDate)
	}
	if !endDate.IsZero() {
		query = query.Where(fmt.Sprintf("%s <= ?", field), endDate)
	}
	return query
}
