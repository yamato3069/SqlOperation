package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// -----------------------------
// LLMへプロンプトを投げて SQL を生成
// -----------------------------
func generateSQLWithPrompt(prompt string) (string, error) {

	payload := map[string]interface{}{
		"model":  "codellama",
		"prompt": prompt,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var b strings.Builder

	for {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			break
		}

		if s, ok := chunk["response"].(string); ok {
			b.WriteString(s)
		}

		if done, ok := chunk["done"].(bool); ok && done {
			break
		}
	}

	sql := b.String()

	// 余計な ``` を除去
	sql = strings.ReplaceAll(sql, "```", "")
	sql = strings.TrimSpace(sql)

	// 先頭に「sql」だけの行が付いてくることがあるので落とす
	lines := strings.Split(sql, "\n")
	if len(lines) > 0 {
		first := strings.TrimSpace(lines[0])
		if strings.EqualFold(first, "sql") || strings.EqualFold(first, "```sql") {
			sql = strings.Join(lines[1:], "\n")
			sql = strings.TrimSpace(sql)
		}
	}

	// 安全チェック（SELECT から始まっているか）
	if !strings.HasPrefix(strings.ToUpper(sql), "SELECT") {
		return "", fmt.Errorf("安全でないSQLです: %s", sql)
	}

	return sql, nil
}
