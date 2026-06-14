package service

import (
	"example/todo/internal/domain"
	"example/todo/internal/repository"
	"log/slog"
)

type listService struct {
	repo repository.ListsRepository
	log  *slog.Logger
}

func NewListService(repo repository.ListsRepository, log *slog.Logger) *listService {
	return &listService{repo: repo, log: log}
}

func (s *listService) Create(list domain.CreateList) (int, error) {
	list_id, err := s.repo.Create(list)
	if err != nil {
		s.log.Error(err.Error())
		return 0, err
	}
	return list_id, nil
}
func (s *listService) Get(user_id, list_id int) (domain.List, error) {
	list, err := s.repo.Get(user_id, list_id)
	if err != nil {
		s.log.Error(err.Error())
		return domain.List{}, err
	}
	return list, nil
}
func (s *listService) GetAll(user_id int) ([]domain.List, error) {
	lists, err := s.repo.GetAll(user_id)
	if err != nil {
		s.log.Error(err.Error())
		return []domain.List{}, err
	}
	return lists, nil
}
func (s *listService) Update(user_id, list_id int, title, description *string) error {
	if err := s.repo.Update(user_id, list_id, title, description); err != nil {
		s.log.Error(err.Error())
		return err
	}
	return nil

}
func (s *listService) Delete(user_id, list_id int) error {
	if err := s.repo.Delete(user_id, list_id); err != nil {
		s.log.Error(err.Error())
		return err
	}
	return nil
}
