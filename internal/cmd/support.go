package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cameltech/sendafrica-cli/internal/api"
	"github.com/cameltech/sendafrica-cli/internal/client"
	"github.com/cameltech/sendafrica-cli/internal/output"
)

var supportChatCmd = &cobra.Command{
	Use:   "support-chat",
	Short: "Chat with the support assistant",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}

		message, _ := cmd.Flags().GetString("message")
		sessionID, _ := cmd.Flags().GetString("session-id")
		filePath, _ := cmd.Flags().GetString("chat-file")

		messages := []api.SupportChatMessage{}

		if message != "" {
			messages = append(messages, api.SupportChatMessage{
				Role:    "user",
				Content: message,
			})
		}

		if filePath != "" {
			msgs, err := loadChatFile(filePath)
			if err != nil {
				return fmt.Errorf("reading chat file: %w", err)
			}
			messages = append(messages, msgs...)
		}

		if len(messages) == 0 {
			return fmt.Errorf("provide --message or --chat-file with conversation messages")
		}

		confirm, _ := cmd.Flags().GetBool("confirm")
		req := api.SupportChatRequest{
			SessionID:        sessionID,
			Messages:         messages,
			UserConfirmation: confirm,
		}

		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/support/chat",
			Body:   req,
		})
		if err != nil {
			return err
		}

		var result api.SupportChatResponse
		if err := decodeData(data, &result); err != nil {
			return err
		}

		if printer().Format() == output.FormatTable {
			fmt.Printf("Session: %s\n", result.SessionID)
			fmt.Printf("Status:  %s\n", result.Status)
			if result.Reply != "" {
				fmt.Printf("\n%s\n", result.Reply)
			}
			if result.ConfirmationRequired != nil {
				fmt.Fprintf(os.Stderr, "\nConfirmation required: %s\n", *result.ConfirmationRequired)
			}
			return nil
		}

		return printer().Print(result)
	},
}

// loadChatFile reads a JSON or newline-delimited file of messages.
// JSON format: [{"role":"user","content":"..."}]
// NDJSON format: one JSON object per line.
func loadChatFile(path string) ([]api.SupportChatMessage, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(content))

	// Try JSON array first.
	var arr []api.SupportChatMessage
	if err := json.Unmarshal([]byte(trimmed), &arr); err == nil {
		return arr, nil
	}

	// Fall back to NDJSON (one object per line).
	var messages []api.SupportChatMessage
	scanner := bufio.NewScanner(strings.NewReader(trimmed))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var msg api.SupportChatMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			return nil, fmt.Errorf("parsing line: %w", err)
		}
		messages = append(messages, msg)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func init() {
	supportChatCmd.Flags().String("message", "", "message to send to the support assistant")
	supportChatCmd.Flags().String("session-id", "", "existing chat session ID")
	supportChatCmd.Flags().String("chat-file", "", "file with conversation messages (JSON array or NDJSON)")
	supportChatCmd.Flags().Bool("confirm", false, "confirm a pending side effect (use after a confirmation_required response)")
}
