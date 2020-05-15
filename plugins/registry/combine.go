package registry

import (
	"github.com/laszlocph/woodpecker/model"
)

type combined struct {
	registries []model.RegistryService
}

func Combined(registries ...model.RegistryService) model.RegistryService {
	return &combined{registries}
}

func (c combined) RegistryFind(repo *model.Repo, name string) (*model.Registry, error) {
	for _, registry := range c.registries {
		res, err := registry.RegistryFind(repo, name)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
	}
	return nil, nil
}

func (c combined) RegistryList(repo *model.Repo) ([]*model.Registry, error) {
	var registeries []*model.Registry
	for _, registory := range c.registries {
		list, err := registory.RegistryList(repo)
		if err != nil {
			return nil, err
		}
		registeries = append(registeries, list...)
	}
	return registeries, nil
}

func (c combined) RegistryCreate(repo *model.Repo, registry *model.Registry) error {
	for _, reg := range c.registries {
		err := reg.RegistryCreate(repo, registry)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c combined) RegistryUpdate(repo *model.Repo, registry *model.Registry) error {
	for _, reg := range c.registries {
		err := reg.RegistryUpdate(repo, registry)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c combined) RegistryDelete(repo *model.Repo, name string) error {
	for _, registry := range c.registries {
		err := registry.RegistryDelete(repo, name)
		if err != nil {
			return err
		}
	}
	return nil
}
