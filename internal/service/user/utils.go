package userservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func (s Service) apiRequest(ctx context.Context, method, url string, token string, obj any) error {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		s.logger.Error("failed to prepare request", slog.String("err", err.Error()))
		return err
	}

	if token != "" {
		req.Header.Add("X-Token", token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("failed to request to get info", slog.String("err", err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		s.logger.Error("non-200 response", slog.Int("status_code", resp.StatusCode))
		return err
	}

	return json.NewDecoder(resp.Body).Decode(obj)
}
