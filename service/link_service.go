package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/heronhoga/shortener-be/model"
	"github.com/heronhoga/shortener-be/repository"
	"github.com/heronhoga/shortener-be/util"
)

type LinkService struct {
	repo repository.LinkRepository
}

func NewLinkService(repo repository.LinkRepository) *LinkService {
	return &LinkService{repo: repo}
}

func (s *LinkService) CreateShortLink(ctx context.Context, requests *model.CreateLink, userID string) error {
	// check existing url shortener link if requests name is not null
	if requests.Name != "" {

		// check valid link name
		isValidName := util.CheckValidLinkName(requests.Name)
		if !isValidName {
			return errors.New("Invalid short link's name")
		}

		// check existing link with the name
		_, err := s.repo.CheckExistingLink(ctx, requests.Name)
		if err != nil {
			return errors.New("There's already a short link with that name")
		}
	}

	// generate random url shortname
	linkName, err := util.GenerateRandomName(10)
	if err != nil {
		return errors.New("Error generating random link name")
	}

	// generate link id
	linkID, err := uuid.NewV7()
	if err != nil {
		return errors.New("Error generating link id")
	}

	// parse user id from fiber context
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("Error parsing user id")
	}

	// create new short link data
	newLink := &model.Link{
		ID:        linkID,
		UserID:    parsedUserID,
		Name:      func() string {
			if requests.Name != "" {
				return requests.Name
			}
			return linkName
		}(),
		Url:       requests.Url,
		CreatedAt: time.Now().UTC(),
	}

	// create new link
	err = s.repo.CreateNewLink(ctx, newLink)
	if err != nil {
		return errors.New("Error creating new short link")
	}

	return nil
}