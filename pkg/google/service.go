package google

import "google.golang.org/api/gmail/v1"

type Service interface {
	GetLabelsService() Labeler
	GetMessagesService() Messenger
}

type GmailService struct {
	Users    *gmail.UsersService
	Labels   *LabelsService
	Drafts   *DraftsService
	Messages *MessagesService
	History  *HistoryService
	Settings *SettingsService
	Threads  *ThreadsService
}

func (g *GmailService) GetLabelsService() Labeler {
	return g.Labels
}

func (g *GmailService) GetMessagesService() Messenger {
	return g.Messages
}

type DraftsService struct {
	*gmail.UsersDraftsService
}

type HistoryService struct {
	*gmail.UsersHistoryService
}

type SettingsService struct {
	*gmail.UsersSettingsService
}

type ThreadsService struct {
	*gmail.UsersThreadsService
}
