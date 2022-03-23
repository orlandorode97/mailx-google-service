package google

import (
	"context"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
)

type LabelsService struct {
	s   *gmail.UsersLabelsService
	ctx context.Context
}

func NewLabelsService(labelSvc *gmail.UsersLabelsService) *LabelsService {
	return &LabelsService{
		ctx: context.Background(),
		s:   labelSvc,
	}
}

func (l *LabelsService) Create(userID string, label *gmail.Label) LabelerClient {
	createCall := l.s.Create(userID, label)
	createCall.Context(l.ctx)
	return createCall
}

func (l *LabelsService) Delete(userID string, labelID string) LabelerClientDelete {
	deleteCall := l.s.Delete(userID, labelID)
	deleteCall.Context(l.ctx)
	return deleteCall
}

func (l *LabelsService) Get(userID string, labelID string) LabelerClient {
	getCall := l.s.Get(userID, labelID)
	getCall.Context(l.ctx)
	return getCall
}

func (l *LabelsService) List(userID string) LabelerClientList {
	listCall := l.s.List(userID)
	listCall.Context(l.ctx)
	return listCall
}

func (l *LabelsService) Patch(userID string, labelID string, label *gmail.Label) LabelerClient {
	patchCall := l.s.Patch(userID, labelID, label)
	patchCall.Context(l.ctx)
	return patchCall
}

func (l *LabelsService) Update(userID string, labelID string, label *gmail.Label) LabelerClient {
	updateCall := l.s.Update(userID, labelID, label)
	updateCall.Context(l.ctx)
	return updateCall
}

/*
 The listed interfaces represents an abstraction of the LabelerClient
	Create -> Do() (*gmail.Label, error)
	Delete -> Do() error
	Get -> Do() (*gmail.Label, error)
	List -> Do() (*gmail.ListLabelsResponse, error)
	Patch -> Do() (*gmail.Label, error)
	Update -> Do() (*gmail.Label, error)
*/

type LabelerClientDelete interface {
	Do(opts ...googleapi.CallOption) error
}

type LabelerClient interface {
	Do(opts ...googleapi.CallOption) (*gmail.Label, error)
}

type LabelerClientList interface {
	Do(opts ...googleapi.CallOption) (*gmail.ListLabelsResponse, error)
}

type LabelCreatorCall interface {
	Create(string, *gmail.Label) LabelerClient
}

type LabelDeletorCall interface {
	Delete(string, string) LabelerClientDelete
}

type LabelGetterCall interface {
	Get(string, string) LabelerClient
}

type LabelListerCall interface {
	List(string) LabelerClientList
}

type LabelPatcherCall interface {
	Patch(string, string, *gmail.Label) LabelerClient
}

type LabelUpdaterCall interface {
	Update(string, string, *gmail.Label) LabelerClient
}

type Labeler interface {
	LabelCreatorCall
	LabelDeletorCall
	LabelGetterCall
	LabelListerCall
	LabelPatcherCall
	LabelUpdaterCall
}
