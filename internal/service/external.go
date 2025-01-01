package service

import (
	"encoding/json"
	"fmt"

	"io"
	"net/http"

	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/pkg/sys"
)

const baseUrlForAllFilms = "https://kinopoiskapiunofficial.tech/api/v2.2/films/"

type ExternalService struct {
	cfg *config.Config
}

func NewExternalService(cfg *config.Config) *ExternalService {
	return &ExternalService{
		cfg: cfg,
	}
}

func (s *ExternalService) SetExternalRequest(filmId string) (dto.GetExternalFilmDTO, error) {
	apiKey := s.cfg.DB.ApiKey
	urlAllFilms := fmt.Sprintf("%s%s", baseUrlForAllFilms, filmId)
	req, err := http.NewRequest("GET", urlAllFilms, nil)
	if err != nil {
		return dto.GetExternalFilmDTO{},
			sys.NewError(sys.ErrUnknown, err.Error())
	}

	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	resAllFilms, err := client.Do(req)
	if err != nil {
		return dto.GetExternalFilmDTO{},
			sys.NewError(sys.ErrUnknown, err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resAllFilms.Body)

	if resAllFilms.StatusCode == http.StatusNotFound {
		return dto.GetExternalFilmDTO{}, sys.NewError(sys.ErrUnknown, "Фильм не найден в внешнем сервисе")
	}

	if resAllFilms.StatusCode != http.StatusOK {
		return dto.GetExternalFilmDTO{},
			sys.NewError(sys.ErrUnknown, "Фильм не найден в внешнем сервисе")
	}

	bodyAllFilms, err := io.ReadAll(resAllFilms.Body)
	if err != nil {
		return dto.GetExternalFilmDTO{},
			sys.NewError(sys.ErrUnknown, err.Error())
	}

	var externalFilm dto.GetExternalFilmDTO
	err = json.Unmarshal(bodyAllFilms, &externalFilm)
	if err != nil {
		return dto.GetExternalFilmDTO{},
			sys.NewError(sys.ErrUnknown, err.Error())
	}

	return externalFilm, nil
}
