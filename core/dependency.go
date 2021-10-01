package core

import (
	"github.com/moritiza/url-shortener/app/handler"
	"github.com/moritiza/url-shortener/app/repository"
	"github.com/moritiza/url-shortener/app/service"
	"github.com/moritiza/url-shortener/config"
)

// Dependencies store all dependencies
type Dependencies struct {
	Repositories Repositories
	Services     Services
	Handlers     Handlers
}

// Repositories store all repositories
type Repositories struct {
	UrlRepository repository.UrlRepository
}

// Services store all services
type Services struct {
	UrlService service.UrlService
}

// Handlers store all handlers
type Handlers struct {
	UrlHandler handler.UrlHandler
}

// PrepareDependensies prepare application necessary dependencies
func PrepareDependensies(cfg config.Config) *Dependencies {
	var (
		repositories Repositories
		services     Services
		handlers     Handlers
	)

	repositories.UrlRepository = repository.NewUrlRepository(cfg.Database)
	services.UrlService = service.NewUrlService(*cfg.Logger, repositories.UrlRepository)
	handlers.UrlHandler = handler.NewUrlHandler(*cfg.Logger, *cfg.Validator, services.UrlService)

	return &Dependencies{
		Repositories: repositories,
		Services:     services,
		Handlers:     handlers,
	}
}
