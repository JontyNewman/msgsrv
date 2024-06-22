package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type HttpMessageRepository struct {
	client      http.Client
	baseAddress string
}

func InitHttpMessageRepository(baseAddress string) HttpMessageRepository {
	return HttpMessageRepository{
		client:      http.Client{},
		baseAddress: baseAddress,
	}
}

func (repo *HttpMessageRepository) Add(message string) (uint, error) {

	body := strings.NewReader(message)
	response, err := repo.client.Post(repo.baseAddress, "text/plain", body)

	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	data := struct {
		ID uint `json:"id"`
	}{}

	decoder := json.NewDecoder(response.Body)

	if err := decoder.Decode(&data); err != nil {
		return 0, err
	}

	return data.ID, nil
}

func (repo *HttpMessageRepository) Fetch(id uint) (string, bool, error) {
	response, err := repo.client.Get(repo.baseAddress + strconv.FormatUint(uint64(id), 10))

	if err != nil {
		return "", false, err
	}

	if response.StatusCode == http.StatusNotFound {
		return "", false, nil
	}

	if response.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("unexpected response status code %d", response.StatusCode)
	}

	message, err := io.ReadAll(response.Body)

	if err != nil {
		return "", false, err
	}

	return string(message), true, nil
}
