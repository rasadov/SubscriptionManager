package tests

import (
	"os"
	"testing"

	"github.com/rasadov/subscription-manager/internal/repository"
	"github.com/rasadov/subscription-manager/internal/service"
	"github.com/rasadov/subscription-manager/tests/mocks"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB
	testRepository repository.SubscriptionRepository
	testService    service.SubscriptionService
)

func TestMain(m *testing.M) {
	db = mocks.NewTestDB()
	testRepository = repository.NewSubscriptionRepositiry(db)
	testService = service.NewSubscriptionService(mocks.NewMockSubscriptionRepository())

	os.Exit(m.Run())
}

func ResetDB() {
	mocks.ResetDB(db)
}

func SetupServiceUnitTests(t *testing.T) {
	testService = service.NewSubscriptionService(mocks.NewMockSubscriptionRepository())
}

func SetupServiceIntegrationTests(t *testing.T) {
	testService = service.NewSubscriptionService(testRepository)
}

func SetupRepo(t *testing.T) {
	db = mocks.NewTestDB()
	testRepository = repository.NewSubscriptionRepositiry(db)
	ResetDB()
}
