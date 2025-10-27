package userservice

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func (s *UserService) apiRequest(method, url string, token string, obj any) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		s.sl.Error("failed to prepare request", slog.String("err", err.Error()))
		return err
	}

	if token != "" {
		req.Header.Add("X-Token", token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.sl.Error("failed to request to get info", slog.String("err", err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		s.sl.Error("non-200 response", slog.Int("status_code", resp.StatusCode))
		return err
	}

	decodeErr := json.NewDecoder(resp.Body).Decode(obj)
	if decodeErr != nil {
		s.sl.Error("failed to unmarshal data to struct", slog.String("err", decodeErr.Error()))
		return err
	}

	return nil
}
