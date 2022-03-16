package google

import (
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
)

/*
 The listed interfaces represents an abstraction of the *gmail.UserLabelsService and its methods and actioners:
	Create -> Do()
	Delete -> Do()
	Get -> Do()
	Patch -> Do()
	Update -> Do()
*/

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
	Delete(string, string) LabelerClient
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
