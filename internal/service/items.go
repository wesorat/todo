package service

import (
	"example/todo/internal/domain"
	"example/todo/internal/repository"
	"log/slog"
)

type itemService struct {
	repo repository.ItemsRepository
	log  *slog.Logger
}

func NewItemService(repo repository.ItemsRepository, log *slog.Logger) *itemService {
	return &itemService{repo: repo, log: log}
}

func (s *itemService) Create(user_id int, item domain.CreateItem) (int, error) {
	item_id, err := s.repo.Create(user_id, item)
	if err != nil {
		s.log.Error(err.Error())
		return 0, err
	}
	return item_id, nil
}
func (s *itemService) Get(user_id, item_id int) (domain.Item, error) {
	item, err := s.repo.Get(user_id, item_id)
	if err != nil {
		s.log.Error(err.Error())
		return domain.Item{}, err
	}
	return item, nil

}
func (s *itemService) GetAll(user_id, list_id int) ([]domain.Item, error){
	items, err := s.repo.GetAll(user_id, list_id)
	if err != nil {
		s.log.Error(err.Error())
		return []domain.Item{}, err
	}
	return items, nil
}
func (s *itemService) Update(user_id, item_id int, title, description *string, done *bool) error{
	if err := s.repo.Update(user_id, item_id, title, description, done); err != nil {
		s.log.Error(err.Error())
		return err
	}
	return nil
}
func (s *itemService) Delete(user_id, item_id int) error {
	if err := s.repo.Delete(user_id, item_id); err != nil {
		s.log.Error(err.Error())
		return err
	}
	return nil
}
