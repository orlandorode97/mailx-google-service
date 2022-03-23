package google

import (
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
)

type MessagesService struct {
	s *gmail.UsersMessagesService
}

func NewMessagesService(messagesSvc *gmail.UsersMessagesService) *MessagesService {
	return &MessagesService{
		s: messagesSvc,
	}
}

func (m *MessagesService) BatchDelete(userID string, req *gmail.BatchDeleteMessagesRequest) MessengerClient {
	return m.s.BatchDelete(userID, req)
}
func (m *MessagesService) BatchModify(userID string, req *gmail.BatchModifyMessagesRequest) MessengerClient {
	return m.s.BatchModify(userID, req)
}
func (m *MessagesService) Delete(userID string, messageID string) MessengerClient {
	return m.s.Delete(userID, messageID)
}
func (m *MessagesService) Get(userID string, messageID string) MessengerClientResp {
	return m.s.Get(userID, messageID)
}
func (m *MessagesService) Import(userID string, message *gmail.Message) MessengerClientResp {
	return m.s.Import(userID, message)
}
func (m *MessagesService) Insert(userID string, message *gmail.Message) MessengerClientResp {
	return m.s.Insert(userID, message)
}
func (m *MessagesService) List(userID string) MessengerClientList {
	return m.s.List(userID)
}
func (m *MessagesService) Modify(userID string, messageID string, req *gmail.ModifyMessageRequest) MessengerClientResp {
	return m.s.Modify(userID, messageID, req)
}
func (m *MessagesService) Send(userID string, message *gmail.Message) MessengerClientResp {
	return m.s.Send(userID, message)
}
func (m *MessagesService) Trash(userID string, messageID string) MessengerClientResp {
	return m.s.Trash(userID, messageID)
}
func (m *MessagesService) Untrash(userID string, messageID string) MessengerClientResp {
	return m.s.Untrash(userID, messageID)
}

/*
 The listed interfaces represents an abstraction of the *gmail.UsersMessagesService and its methods and actioners:
	BatchDelete -> Do()
	BatchModify -> Do()
	Delete -> Do()
	Get -> Do()
	Import -> Do()
	Insert -> Do()
	List -> Do()
	Modify -> Do()
	Send -> Do()
	Trash -> Do()
	Untrash -> Do()
*/

type MessengerClient interface {
	Do(opts ...googleapi.CallOption) error
}

type MessengerClientResp interface {
	Do(opts ...googleapi.CallOption) (*gmail.Message, error)
}

type MessengerClientList interface {
	Do(opts ...googleapi.CallOption) (*gmail.ListMessagesResponse, error)
}

type MessageBatchDeletorCall interface {
	BatchDelete(string, *gmail.BatchDeleteMessagesRequest) MessengerClient
}

type MessageBatchModifierCall interface {
	BatchModify(string, *gmail.BatchModifyMessagesRequest) MessengerClient
}

type MessageDeletorCall interface {
	Delete(string, string) MessengerClient
}

type MessageGetterCall interface {
	Get(string, string) MessengerClientResp
}

type MessageImporterCall interface {
	Import(string, *gmail.Message) MessengerClientResp
}

type MessageInserterCall interface {
	Insert(string, *gmail.Message) MessengerClientResp
}

type MessageListerCall interface {
	List(string) MessengerClientList
}

type MessageModifierCall interface {
	Modify(string, string, *gmail.ModifyMessageRequest) MessengerClientResp
}

type MessageSenderCall interface {
	Send(string, *gmail.Message) MessengerClientResp
}

type MessageTrasherCall interface {
	Trash(string, string) MessengerClientResp
}

type MessageUntrasherCall interface {
	Untrash(string, string) MessengerClientResp
}

type Messenger interface {
	MessageBatchDeletorCall
	MessageBatchModifierCall
	MessageDeletorCall
	MessageGetterCall
	MessageImporterCall
	MessageInserterCall
	MessageListerCall
	MessageModifierCall
	MessageSenderCall
	MessageTrasherCall
	MessageUntrasherCall
}
