package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DailyService struct {
	vault *Vault
}

func NewDaily(vault *Vault) *DailyService {
	return &DailyService{vault: vault}
}

func DailyNotePath(date time.Time) string {
	return "daily/" + date.Format("2006-01-02") + ".md"
}

func (d *DailyService) EnsureDailyNote(date time.Time) (*VaultNote, error) {
	path := DailyNotePath(date)

	if d.vault.NoteExists(path) {
		return d.vault.ReadNote(path)
	}

	template := d.getTemplate()
	content := applyTemplate(template, date)

	// Parse frontmatter from template
	fm, body := parseFrontmatter([]byte(content))
	if fm.ID == "" {
		fm.ID = newULID()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if fm.Created == "" {
		fm.Created = now
	}
	if fm.Modified == "" {
		fm.Modified = now
	}
	if fm.Title == "" {
		fm.Title = date.Format("2006-01-02")
	}
	if fm.Tags == nil {
		fm.Tags = []string{"daily"}
	}

	return d.vault.WriteNote(path, body, fm)
}

func (d *DailyService) getTemplate() string {
	configPath := filepath.Join(d.vault.Path(), ".prometheus", "config.json")
	data, err := os.ReadFile(configPath)
	if err == nil {
		var config struct {
			DailyNoteTemplate string `json:"dailyNoteTemplate"`
		}
		if json.Unmarshal(data, &config) == nil && config.DailyNoteTemplate != "" {
			return config.DailyNoteTemplate
		}
	}

	return `---
title: "{{date}}"
tags: [daily]
---

# {{dateJa}}

## タスク

- [ ]

## メモ

`
}

func applyTemplate(template string, date time.Time) string {
	dateStr := date.Format("2006-01-02")

	// Japanese date format
	weekdays := []string{"日曜日", "月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日"}
	dateJa := date.Format("2006年1月2日") + weekdays[date.Weekday()]

	result := strings.ReplaceAll(template, "{{date}}", dateStr)
	result = strings.ReplaceAll(result, "{{dateJa}}", dateJa)
	return result
}

func RecentDailyDates(count int) []time.Time {
	dates := make([]time.Time, count)
	today := time.Now()
	for i := 0; i < count; i++ {
		dates[i] = today.AddDate(0, 0, -i)
	}
	return dates
}
