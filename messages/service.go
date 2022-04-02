package messages

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/go-kit/log"
	"github.com/orlandorode97/mailx-google-service"
	"github.com/orlandorode97/mailx-google-service/pkg/google"
	"github.com/orlandorode97/mailx-google-service/pkg/models"
	"github.com/orlandorode97/mailx-google-service/pkg/repos"
	"google.golang.org/api/gmail/v1"
)

const (
	messagesLimit int64 = 10
)

type Service interface {
	GetMessages(context.Context, string) ([]*models.Message, error)
	GetMessageByID(context.Context, string, string) (*models.Message, error)
}

type service struct {
	logger      log.Logger
	repo        repos.Repository
	mailxSvc    mailx.Service
	messagesSvc google.Messenger
}

func New(logger log.Logger, repo repos.Repository, mailx mailx.Service) Service {
	return &service{
		logger:   logger,
		repo:     repo,
		mailxSvc: mailx,
	}
}

func (s *service) addMessagesService(svc google.Messenger) {
	s.messagesSvc = svc
}

func (s *service) recreateMessageService(ctx context.Context, userID string) error {
	svc, err := s.mailxSvc.RecreateGmailService(ctx, userID)
	if err != nil {
		return err
	}
	s.addMessagesService(svc.GetMessagesService())
	return nil
}

func (s *service) GetMessages(ctx context.Context, userID string) ([]*models.Message, error) {

	if s.messagesSvc == nil {
		s.recreateMessageService(context.Background(), userID)
	}

	messagesResp, err := s.messagesSvc.List(userID, messagesLimit).Do()
	if err != nil {
		s.logger.Log(
			"message", fmt.Sprintf("error getting messages for user=%s", userID),
			"error", err.Error(),
			"severity", "ERROR",
		)
		return nil, err
	}
	fmt.Printf("length alv %v \n\n\n\n", len(messagesResp.Messages))
	s.logger.Log(
		"message", fmt.Sprintf("get messages for user=%s", userID),
		"severity", "INFO",
	)

	messageChan := make(chan *models.Message)

	messages := make([]*models.Message, 0)
	var counter int32 = 0
	var wg sync.WaitGroup
	go func() {
		for message := range messageChan {
			messages = append(messages, message)
		}
	}()

	for _, message := range messagesResp.Messages {
		wg.Add(1)
		go func(userID, messageID string) {
			msg, err := s.GetMessageByID(ctx, userID, messageID)
			if err != nil {
				return
			}
			messageChan <- msg
			atomic.AddInt32(&counter, 1)
			wg.Done()
		}(userID, message.Id)
	}

	wg.Wait()
	fmt.Println("total of items", counter)
	fmt.Printf("length alv of messages %v \n\n\n\n", len(messages))
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].InternalDate > messages[j].InternalDate
	})
	return messages, nil
}

func (s *service) GetMessageByID(ctx context.Context, userID string, messageID string) (*models.Message, error) {

	if s.messagesSvc == nil {
		s.recreateMessageService(context.Background(), userID)
	}

	message, err := s.messagesSvc.Get(userID, messageID).Do()
	if err != nil {
		s.logger.Log(
			"message", fmt.Sprintf("error message=%s for user= %s", messageID, userID),
			"error", err.Error(),
			"severity", "ERROR",
		)
		return nil, err
	}

	s.logger.Log(
		"message", fmt.Sprintf("get message=%s for user=%s", messageID, userID),
		"severity", "INFO",
	)

	msg, err := s.hydrateMessage(message)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *service) hydrateMessage(message *gmail.Message) (*models.Message, error) {
	var (
		data string
		html []byte
		err  error
	)

	for _, part := range message.Payload.Parts {
		if data != "" {
			break
		}
		// Some messages include the html to decode in []*gmail.MessagePart from the parent *gmail.MessagePart
		if part.Parts != nil {
			for _, p := range part.Parts {
				if p.MimeType == "text/html" {
					data = p.Body.Data
					break
				}
			}
		}

		if part.MimeType == "text/html" {
			data = part.Body.Data
			break
		}
	}

	if data != "" {
		html, err = base64.URLEncoding.DecodeString(data)
		if err != nil {
			return nil, err
		}
	}

	return &models.Message{
		ID:           message.Id,
		HistoryID:    message.HistoryId,
		InternalDate: message.InternalDate,
		LabelIDS:     message.LabelIds,
		Payload:      message.Payload,
		SizeEstimate: message.SizeEstimate,
		Snippet:      message.Snippet,
		ThreadID:     message.ThreadId,
		HTML:         string(html),
	}, nil
}
