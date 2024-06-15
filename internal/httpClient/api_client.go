package httpClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const BasePath = "/game"

// NewClient creates a new API client
func NewClient(baseURL string, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// getRequest creates a new GET request
func (c *Client) getRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", c.Token)
	return req, nil
}

// postRequest creates a new POST request
func (c *Client) postRequest(url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, c.BaseURL+url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", c.Token)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// handleResponse handles the server response
func handleResponse(resp *http.Response, successCode int, result interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == successCode {
		return json.Unmarshal(body, result)
	}

	var errResp error
	switch resp.StatusCode {
	case 401:
		errResp = UnauthorizedError{}
	case 403:
		errResp = ForbiddenError{}
	case 429:
		errResp = RateLimitExceededError{}
	case 400:
		errResp = BadRequestError{}
	case 404:
		errResp = NotFoundError{}
	default:
		errResp = errors.New("unexpected API error")
	}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("failed to process error response: %w", err)
	}
	return errResp
}

// GetGameStatus retrieves the game status
func (c *Client) GetGameStatus() (GameStatus, error) {
	req, err := c.getRequest(BasePath)
	if err != nil {
		return GameStatus{}, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return GameStatus{}, err
	}

	var gameState GameStatus
	if err := handleResponse(resp, http.StatusOK, &gameState); err != nil {
		return GameStatus{}, err
	}
	return gameState, nil
}

// StartGame starts a new game
func (c *Client) StartGame(nick, desc, targetNick string, coords []string, botGame bool) (string, error) {
	bodyData := map[string]interface{}{
		"coords":      coords,
		"desc":        desc,
		"nick":        nick,
		"target_nick": targetNick,
		"wpbot":       botGame,
	}

	body, err := json.Marshal(bodyData)
	if err != nil {
		return "", err
	}

	req, err := c.postRequest(BasePath, body)
	if err != nil {
		return "", err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK {
		token := resp.Header.Get("X-Auth-Token")
		c.Token = token
		return token, nil
	}

	return "", handleResponse(resp, http.StatusOK, nil)
}

// GetGameBoard retrieves the game board
func (c *Client) GetGameBoard() (*GameBoard, error) {
	req, err := c.getRequest(BasePath + "/board")
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var gameBoard GameBoard
	if err := handleResponse(resp, http.StatusOK, &gameBoard); err != nil {
		return nil, err
	}
	return &gameBoard, nil
}

// Fire executes a shot
func (c *Client) Fire(data FireData) (FireResult, error) {
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return FireResult{}, err
	}

	req, err := c.postRequest(BasePath+"/fire", bodyBytes)
	if err != nil {
		return FireResult{}, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return FireResult{}, err
	}

	var fireResult FireResult
	if err := handleResponse(resp, http.StatusOK, &fireResult); err != nil {
		return FireResult{}, err
	}
	return fireResult, nil
}

// AbandonGame abandons the game
func (c *Client) AbandonGame() error {
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL+BasePath+"/abandon", nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusOK, nil)
}

// GetGameDescription retrieves the game description
func (c *Client) GetGameDescription() (GameDescription, error) {
	req, err := c.getRequest(BasePath + "/desc")
	if err != nil {
		return GameDescription{}, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return GameDescription{}, err
	}

	var gameDescription GameDescription
	if err := handleResponse(resp, http.StatusOK, &gameDescription); err != nil {
		return GameDescription{}, err
	}
	return gameDescription, nil
}

// RefreshGameSession refreshes the game session
func (c *Client) RefreshGameSession() error {
	req, err := c.getRequest(BasePath + "/refresh")
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusOK, nil)
}

// GetAllGames retrieves all games with a given status
func (c *Client) GetAllGames(status string) (GameList, error) {
	var gameList GameList

	req, err := c.getRequest("/list")
	if err != nil {
		return gameList, err
	}

	q := req.URL.Query()
	q.Add("status", status)
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return gameList, err
	}

	if err := handleResponse(resp, http.StatusOK, &gameList); err != nil {
		return gameList, err
	}
	return gameList, nil
}

// GetLobbyPlayers retrieves players in the lobby
func (c *Client) GetLobbyPlayers() ([]LobbyPlayer, error) {
	var players []LobbyPlayer

	req, err := c.getRequest("/lobby")
	if err != nil {
		return players, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return players, err
	}

	if err := handleResponse(resp, http.StatusOK, &players); err != nil {
		return players, err
	}
	return players, nil
}

// GetTopPlayerStats retrieves statistics of the top players
func (c *Client) GetTopPlayerStats() (TopPlayerStats, error) {
	var topStats TopPlayerStats

	req, err := c.getRequest("/stats")
	if err != nil {
		return topStats, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return topStats, err
	}

	if err := handleResponse(resp, http.StatusOK, &topStats); err != nil {
		return topStats, err
	}
	return topStats, nil
}

// GetPlayerStats retrieves a player's statistics
func (c *Client) GetPlayerStats(nick string) (GameStats, error) {
	url := "https://go-pjatk-server.fly.dev/api/stats/" + strings.TrimSpace(nick)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return GameStats{}, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return GameStats{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GameStats{}, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	var response struct {
		Stats GameStat `json:"stats"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return GameStats{}, err
	}

	return GameStats{response.Stats}, nil
}

// AbortGame aborts the game
func (c *Client) AbortGame() error {
	url := "https://go-pjatk-server.fly.dev/api/game/abandon"

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-Token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}

	return nil
}
