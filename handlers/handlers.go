package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	searchCli "github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	zebedeeCli "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	errs "github.com/ONSdigital/dp-frontend-search-controller/apperrors"
	"github.com/ONSdigital/dp-frontend-search-controller/cache"
	"github.com/ONSdigital/dp-frontend-search-controller/config"
	"github.com/ONSdigital/dp-frontend-search-controller/data"
	"github.com/ONSdigital/dp-frontend-search-controller/mapper"
	dphandlers "github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
)

// Constants...
const (
	homepagePath = "/"
)

// Read Handler
func Read(cfg *config.Config, zc ZebedeeClient, rend RenderClient, searchC SearchClient, cacheList cache.CacheList) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		read(w, req, cfg, zc, rend, searchC, accessToken, collectionID, lang, cacheList)
	})
}

func read(w http.ResponseWriter, req *http.Request, cfg *config.Config, zc ZebedeeClient, rend RenderClient, searchC SearchClient,
	accessToken, collectionID, lang string, cacheList cache.CacheList) {

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	urlQuery := req.URL.Query()

	// get cached census topic and its subtopics
	censusTopicCache, err := cacheList.CensusTopic.GetCensusData(ctx)
	if err != nil {
		log.Error(ctx, "failed to get census topic cache", err)
		setStatusCode(w, req, err)
		return
	}

	validatedQueryParams, err := data.ReviewQuery(ctx, cfg, urlQuery, censusTopicCache)
	if err != nil && !errs.ErrMapForRenderBeforeAPICalls[err] {
		log.Error(ctx, "unable to review query", err)
		setStatusCode(w, req, err)
		return
	}

	apiQuery := data.GetSearchAPIQuery(validatedQueryParams, censusTopicCache)

	var homepageResponse zebedeeCli.HomepageContent
	var searchResp searchCli.Response
	var respErr error
	var departmentResp searchCli.Department

	if errs.ErrMapForRenderBeforeAPICalls[err] {
		// avoid making any API calls
		basePage := rend.NewBasePageModel()
		m := mapper.CreateSearchPage(cfg, req, basePage, validatedQueryParams, []data.Category{}, []data.Topic{}, searchResp, departmentResp, lang, homepageResponse, err.Error())
		rend.BuildPage(w, m, "search")
		return
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		homepageResponse, err = zc.GetHomepageContent(ctx, accessToken, collectionID, lang, homepagePath)
		if err != nil {
			logData := log.Data{"homepage_content": err}
			log.Error(ctx, "unable to get homepage content", err, logData)
			cancel()
			return
		}
	}()
	go func() {
		defer wg.Done()
		searchResp, respErr = searchC.GetSearch(ctx, accessToken, "", collectionID, apiQuery)
		if respErr != nil {
			logData := log.Data{"api query passed to search-api": apiQuery}
			log.Error(ctx, "getting search response from client failed", respErr, logData)
			cancel()
			return
		}
	}()
	go func() {
		defer wg.Done()
		var deptErr error
		departmentResp, deptErr = searchC.GetDepartments(ctx, accessToken, "", collectionID, apiQuery)
		if deptErr != nil {
			logData := log.Data{"api query passed to search-api": apiQuery}
			log.Error(ctx, "getting deartment response from client failed", deptErr, logData)
			return
		}
	}()

	wg.Wait()
	if respErr != nil {
		setStatusCode(w, req, respErr)
		return
	}

	// TO-DO: Until API handles aggregration on datatypes (e.g. bulletins, article), we need to make a second request

	err = validateCurrentPage(ctx, cfg, validatedQueryParams, searchResp.Count)
	if err != nil {
		log.Error(ctx, "unable to validate current page", err)
		setStatusCode(w, req, err)
		return
	}

	categories, topicCategories, err := getCategoriesTypesCount(ctx, accessToken, collectionID, apiQuery, searchC, censusTopicCache)
	if err != nil {
		log.Error(ctx, "getting categories, types and its counts failed", err)
		setStatusCode(w, req, err)
		return
	}

	basePage := rend.NewBasePageModel()
	m := mapper.CreateSearchPage(cfg, req, basePage, validatedQueryParams, categories, topicCategories, searchResp, departmentResp, lang, homepageResponse, "")
	rend.BuildPage(w, m, "search")
}

// Read Handler
func MockRead(cfg *config.Config, zc ZebedeeClient, rend RenderClient, searchC SearchClient, cacheList cache.CacheList) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		mockRead(w, req, cfg, zc, rend, searchC, accessToken, collectionID, lang, cacheList)
	})
}

func mockRead(w http.ResponseWriter, req *http.Request, cfg *config.Config, zc ZebedeeClient, rend RenderClient, searchC SearchClient,
	accessToken, collectionID, lang string, cacheList cache.CacheList) {

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	urlQuery := req.URL.Query()

	// MOCK CENSUS TOPIC CACHE - needs to talk to topic-api and get connected to mongodb
	censusTopicCache := &cache.Topic{
		ID:              "4445",
		LocaliseKeyName: "Census",
		Query:           "4445,1234,5645",
		List:            cache.NewSubTopicsMap(),
	}
	censusTopicCache.List.AppendSubtopicID("4445")
	censusTopicCache.List.AppendSubtopicID("1234")

	validatedQueryParams, err := data.ReviewQuery(ctx, cfg, urlQuery, censusTopicCache)
	if err != nil && err != errs.ErrInvalidQueryString && err != errs.ErrFilterNotFound {
		log.Error(ctx, "unable to review query", err)
		setStatusCode(w, req, err)
		return
	}

	// MOCK SEARCH RESPONSE - needs to talk to search api
	var searchResp searchCli.Response
	// get census results only
	if validatedQueryParams.TopicFilter == "4445" && len(validatedQueryParams.Filter.Query) <= 0 {
		searchResp, err = getMockSearchResponse("handlers/mock_search_census_response.json")
		if err != nil {
			log.Error(ctx, "failed to get mock search response", err)
			setStatusCode(w, req, err)
			return
		}
		// get articles only
	} else if len(validatedQueryParams.Filter.Query) > 0 && validatedQueryParams.TopicFilter == "" {
		searchResp, err = getMockSearchResponse("handlers/mock_search_pub_response.json")
		if err != nil {
			log.Error(ctx, "failed to get mock search response", err)
			setStatusCode(w, req, err)
			return
		}
		// get all
	} else {
		searchResp, err = getMockSearchResponse("handlers/mock_search_all_response.json")
		if err != nil {
			log.Error(ctx, "failed to get mock search response", err)
			setStatusCode(w, req, err)
			return
		}
	}

	departmentResp := searchCli.Department{}
	homepageResponse := zebedeeCli.HomepageContent{}

	if err == errs.ErrInvalidQueryString || err == errs.ErrFilterNotFound {
		// avoid making any API calls
		basePage := rend.NewBasePageModel()
		m := mapper.CreateSearchPage(cfg, req, basePage, validatedQueryParams, []data.Category{}, []data.Topic{}, searchResp, departmentResp, lang, homepageResponse, err.Error())
		rend.BuildPage(w, m, "search")
		return
	}

	err = validateCurrentPage(ctx, cfg, validatedQueryParams, searchResp.Count)
	if err != nil {
		log.Error(ctx, "unable to validate current page", err)
		setStatusCode(w, req, err)
		return
	}

	// MOCK CATEGORIES - needs to talk to search api to get its count
	categories := []data.Category{
		{
			LocaliseKeyName: "Publication",
			ContentTypes: []data.ContentType{
				{
					LocaliseKeyName: "Article",
					Group:           "article",
					Types:           []string{"article", "article_download"},
					ShowInWebUI:     true,
					Count:           3,
				},
			},
			Count: 3,
		},
		{
			LocaliseKeyName: "Data",
			ContentTypes: []data.ContentType{
				{
					LocaliseKeyName: "TimeSeries",
					Group:           "time_series",
					Types:           []string{"timeseries"},
					ShowInWebUI:     true,
					Count:           2,
				},
			},
			Count: 2,
		},
		{
			LocaliseKeyName: "Other",
			ContentTypes:    []data.ContentType{},
			Count:           0,
		},
	}

	topicCategories := []data.Topic{
		{
			LocaliseKeyName: "Census",
			Count:           2,
			Query:           "4445",
			ShowInWebUI:     true,
		},
	}

	basePage := rend.NewBasePageModel()
	m := mapper.CreateSearchPage(cfg, req, basePage, validatedQueryParams, categories, topicCategories, searchResp, departmentResp, lang, homepageResponse, "")
	rend.BuildPage(w, m, "search")
}

func getMockSearchResponse(path string) (searchCli.Response, error) {
	var respC searchCli.Response

	sampleResponse, err := ioutil.ReadFile(path)
	if err != nil {
		return searchCli.Response{}, err
	}

	err = json.Unmarshal(sampleResponse, &respC)
	if err != nil {
		return searchCli.Response{}, err
	}

	return respC, nil
}

// validateCurrentPage checks if the current page exceeds the total pages which is a bad request
func validateCurrentPage(ctx context.Context, cfg *config.Config, validatedQueryParams data.SearchURLParams, resultsCount int) error {

	if resultsCount > 0 {
		totalPages := data.GetTotalPages(cfg, validatedQueryParams.Limit, resultsCount)

		if validatedQueryParams.CurrentPage > totalPages {
			err := errs.ErrPageExceedsTotalPages
			log.Error(ctx, "current page exceeds total pages", err)

			return err
		}
	}

	return nil
}

// getCategoriesTypesCount removes the filters and communicates with the search api again to retrieve the number of search results for each filter categories and subtypes
func getCategoriesTypesCount(ctx context.Context, accessToken, collectionID string, apiQuery url.Values, searchC SearchClient, censusTopicCache *cache.Topic) ([]data.Category, []data.Topic, error) {
	//Remove filter to get count of all types for the query from the client
	apiQuery.Del("content_type")
	apiQuery.Del("topics")

	countResp, err := searchC.GetSearch(ctx, accessToken, "", collectionID, apiQuery)
	if err != nil {
		logData := log.Data{"api query passed to search-api": apiQuery}
		log.Error(ctx, "getting search query count from client failed", err, logData)
		return nil, nil, err
	}

	categories := data.GetCategories()
	topicCategories := data.GetTopicCategories(censusTopicCache, countResp)

	setCountToCategories(ctx, countResp, categories)

	return categories, topicCategories, nil
}

func setCountToCategories(ctx context.Context, countResp searchCli.Response, categories []data.Category) {
	for _, responseType := range countResp.ContentTypes {
		foundFilter := false

	categoryLoop:
		for i, category := range categories {
			for j, contentType := range category.ContentTypes {
				for _, subType := range contentType.Types {
					if responseType.Type == subType {
						categories[i].Count += responseType.Count
						categories[i].ContentTypes[j].Count += responseType.Count

						foundFilter = true

						break categoryLoop
					}
				}
			}
		}

		if !foundFilter {
			log.Warn(ctx, "unrecognised filter type returned from api")
		}
	}
}

func setStatusCode(w http.ResponseWriter, req *http.Request, err error) {
	status := http.StatusInternalServerError

	if err, ok := err.(ClientError); ok {
		if err.Code() == http.StatusNotFound {
			status = err.Code()
		}
	}

	if errs.BadRequestMap[err] {
		status = http.StatusBadRequest
	}

	log.Error(req.Context(), "setting-response-status", err)

	w.WriteHeader(status)
}
