// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package handlers

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	"github.com/ONSdigital/dp-renderer/model"
	"io"
	"net/url"
	"sync"
)

var (
	lockRenderClientMockBuildPage        sync.RWMutex
	lockRenderClientMockNewBasePageModel sync.RWMutex
)

// Ensure, that RenderClientMock does implement RenderClient.
// If this is not the case, regenerate this file with moq.
var _ RenderClient = &RenderClientMock{}

// RenderClientMock is a mock implementation of RenderClient.
//
//     func TestSomethingThatUsesRenderClient(t *testing.T) {
//
//         // make and configure a mocked RenderClient
//         mockedRenderClient := &RenderClientMock{
//             BuildPageFunc: func(w io.Writer, pageModel interface{}, templateName string)  {
// 	               panic("mock out the BuildPage method")
//             },
//             NewBasePageModelFunc: func() model.Page {
// 	               panic("mock out the NewBasePageModel method")
//             },
//         }
//
//         // use mockedRenderClient in code that requires RenderClient
//         // and then make assertions.
//
//     }
type RenderClientMock struct {
	// BuildPageFunc mocks the BuildPage method.
	BuildPageFunc func(w io.Writer, pageModel interface{}, templateName string)

	// NewBasePageModelFunc mocks the NewBasePageModel method.
	NewBasePageModelFunc func() model.Page

	// calls tracks calls to the methods.
	calls struct {
		// BuildPage holds details about calls to the BuildPage method.
		BuildPage []struct {
			// W is the w argument value.
			W io.Writer
			// PageModel is the pageModel argument value.
			PageModel interface{}
			// TemplateName is the templateName argument value.
			TemplateName string
		}
		// NewBasePageModel holds details about calls to the NewBasePageModel method.
		NewBasePageModel []struct {
		}
	}
}

// BuildPage calls BuildPageFunc.
func (mock *RenderClientMock) BuildPage(w io.Writer, pageModel interface{}, templateName string) {
	if mock.BuildPageFunc == nil {
		panic("RenderClientMock.BuildPageFunc: method is nil but RenderClient.BuildPage was just called")
	}
	callInfo := struct {
		W            io.Writer
		PageModel    interface{}
		TemplateName string
	}{
		W:            w,
		PageModel:    pageModel,
		TemplateName: templateName,
	}
	lockRenderClientMockBuildPage.Lock()
	mock.calls.BuildPage = append(mock.calls.BuildPage, callInfo)
	lockRenderClientMockBuildPage.Unlock()
	mock.BuildPageFunc(w, pageModel, templateName)
}

// BuildPageCalls gets all the calls that were made to BuildPage.
// Check the length with:
//     len(mockedRenderClient.BuildPageCalls())
func (mock *RenderClientMock) BuildPageCalls() []struct {
	W            io.Writer
	PageModel    interface{}
	TemplateName string
} {
	var calls []struct {
		W            io.Writer
		PageModel    interface{}
		TemplateName string
	}
	lockRenderClientMockBuildPage.RLock()
	calls = mock.calls.BuildPage
	lockRenderClientMockBuildPage.RUnlock()
	return calls
}

// NewBasePageModel calls NewBasePageModelFunc.
func (mock *RenderClientMock) NewBasePageModel() model.Page {
	if mock.NewBasePageModelFunc == nil {
		panic("RenderClientMock.NewBasePageModelFunc: method is nil but RenderClient.NewBasePageModel was just called")
	}
	callInfo := struct {
	}{}
	lockRenderClientMockNewBasePageModel.Lock()
	mock.calls.NewBasePageModel = append(mock.calls.NewBasePageModel, callInfo)
	lockRenderClientMockNewBasePageModel.Unlock()
	return mock.NewBasePageModelFunc()
}

// NewBasePageModelCalls gets all the calls that were made to NewBasePageModel.
// Check the length with:
//     len(mockedRenderClient.NewBasePageModelCalls())
func (mock *RenderClientMock) NewBasePageModelCalls() []struct {
} {
	var calls []struct {
	}
	lockRenderClientMockNewBasePageModel.RLock()
	calls = mock.calls.NewBasePageModel
	lockRenderClientMockNewBasePageModel.RUnlock()
	return calls
}

var (
	lockSearchClientMockGetDepartments sync.RWMutex
	lockSearchClientMockGetSearch      sync.RWMutex
)

// Ensure, that SearchClientMock does implement SearchClient.
// If this is not the case, regenerate this file with moq.
var _ SearchClient = &SearchClientMock{}

// SearchClientMock is a mock implementation of SearchClient.
//
//     func TestSomethingThatUsesSearchClient(t *testing.T) {
//
//         // make and configure a mocked SearchClient
//         mockedSearchClient := &SearchClientMock{
//             GetDepartmentsFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, query url.Values) (search.Department, error) {
// 	               panic("mock out the GetDepartments method")
//             },
//             GetSearchFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, query url.Values) (search.Response, error) {
// 	               panic("mock out the GetSearch method")
//             },
//         }
//
//         // use mockedSearchClient in code that requires SearchClient
//         // and then make assertions.
//
//     }
type SearchClientMock struct {
	// GetDepartmentsFunc mocks the GetDepartments method.
	GetDepartmentsFunc func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, query url.Values) (search.Department, error)

	// GetSearchFunc mocks the GetSearch method.
	GetSearchFunc func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, query url.Values) (search.Response, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetDepartments holds details about calls to the GetDepartments method.
		GetDepartments []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserAuthToken is the userAuthToken argument value.
			UserAuthToken string
			// ServiceAuthToken is the serviceAuthToken argument value.
			ServiceAuthToken string
			// CollectionID is the collectionID argument value.
			CollectionID string
			// Query is the query argument value.
			Query url.Values
		}
		// GetSearch holds details about calls to the GetSearch method.
		GetSearch []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserAuthToken is the userAuthToken argument value.
			UserAuthToken string
			// ServiceAuthToken is the serviceAuthToken argument value.
			ServiceAuthToken string
			// CollectionID is the collectionID argument value.
			CollectionID string
			// Query is the query argument value.
			Query url.Values
		}
	}
}

// GetDepartments calls GetDepartmentsFunc.
func (mock *SearchClientMock) GetDepartments(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, query url.Values) (search.Department, error) {
	if mock.GetDepartmentsFunc == nil {
		panic("SearchClientMock.GetDepartmentsFunc: method is nil but SearchClient.GetDepartments was just called")
	}
	callInfo := struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		Query            url.Values
	}{
		Ctx:              ctx,
		UserAuthToken:    userAuthToken,
		ServiceAuthToken: serviceAuthToken,
		CollectionID:     collectionID,
		Query:            query,
	}
	lockSearchClientMockGetDepartments.Lock()
	mock.calls.GetDepartments = append(mock.calls.GetDepartments, callInfo)
	lockSearchClientMockGetDepartments.Unlock()
	return mock.GetDepartmentsFunc(ctx, userAuthToken, serviceAuthToken, collectionID, query)
}

// GetDepartmentsCalls gets all the calls that were made to GetDepartments.
// Check the length with:
//     len(mockedSearchClient.GetDepartmentsCalls())
func (mock *SearchClientMock) GetDepartmentsCalls() []struct {
	Ctx              context.Context
	UserAuthToken    string
	ServiceAuthToken string
	CollectionID     string
	Query            url.Values
} {
	var calls []struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		Query            url.Values
	}
	lockSearchClientMockGetDepartments.RLock()
	calls = mock.calls.GetDepartments
	lockSearchClientMockGetDepartments.RUnlock()
	return calls
}

// GetSearch calls GetSearchFunc.
func (mock *SearchClientMock) GetSearch(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, query url.Values) (search.Response, error) {
	if mock.GetSearchFunc == nil {
		panic("SearchClientMock.GetSearchFunc: method is nil but SearchClient.GetSearch was just called")
	}
	callInfo := struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		Query            url.Values
	}{
		Ctx:              ctx,
		UserAuthToken:    userAuthToken,
		ServiceAuthToken: serviceAuthToken,
		CollectionID:     collectionID,
		Query:            query,
	}
	lockSearchClientMockGetSearch.Lock()
	mock.calls.GetSearch = append(mock.calls.GetSearch, callInfo)
	lockSearchClientMockGetSearch.Unlock()
	return mock.GetSearchFunc(ctx, userAuthToken, serviceAuthToken, collectionID, query)
}

// GetSearchCalls gets all the calls that were made to GetSearch.
// Check the length with:
//     len(mockedSearchClient.GetSearchCalls())
func (mock *SearchClientMock) GetSearchCalls() []struct {
	Ctx              context.Context
	UserAuthToken    string
	ServiceAuthToken string
	CollectionID     string
	Query            url.Values
} {
	var calls []struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		Query            url.Values
	}
	lockSearchClientMockGetSearch.RLock()
	calls = mock.calls.GetSearch
	lockSearchClientMockGetSearch.RUnlock()
	return calls
}