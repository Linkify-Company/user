package service

import (
	"context"
	auth "github.com/Linkify-Company/auth-client"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
	"time"
	"user/internal/domain"
	"user/internal/repository"
)

type Person interface {
	AddPerson(ctx context.Context, person *domain.Person) (int, errify.IError)
	GetPersonByUserID(ctx context.Context, userID int) (*domain.Person, errify.IError)
}

type Reboot interface {
	RebootSystem(ctx context.Context) errify.IError
}

type Service struct {
	Person
	Reboot

	log         logger.Logger
	AuthService *auth.Service
}

func NewService(
	log logger.Logger,
	pool *pgxpool.Pool,
	repos *repository.Repository,
	authService *auth.Service,
) *Service {
	transaction := repository.NewTransactionsRepos(pool)

	s := Service{
		Person: NewPersonService(log, transaction, repos),
		Reboot: NewRebootService(log, transaction, repos),

		log:         log,
		AuthService: authService,
	}

	go s.authPing()
	go s.pingPSQL(pool)

	return &s
}

func (s *Service) authPing() {
	c := cron.New()

	_, err := c.AddFunc("* * * * *", func() {
		pong, err := s.AuthService.Ping(s.log)
		if err != nil {
			s.log.Errorf("Auth server not found")
			s.log.Warn(err)
			return
		}
		s.log.Infof("Auth server found: %s", pong)
	})
	if err != nil {
		s.log.Warn(errify.NewInternalServerError("error launching the cron to check the connection to the authorization service", "authPing/AddFunc"))
		return
	}

	c.Run()
}

func (s *Service) pingPSQL(pool *pgxpool.Pool) {
	c := cron.New()

	_, err := c.AddFunc("* * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := pool.Ping(ctx)
		if err != nil {
			s.log.Error(errify.NewInternalServerError("postgres connection error", "loss of connection to postgresql"))
			return
		}
		s.log.Debugf("the connection to postgresql is stable")
	})
	if err != nil {
		s.log.Warn(errify.NewInternalServerError("error launching the cron to check the connection to the authorization service", "authPing/AddFunc"))
		return
	}

	c.Run()
}
