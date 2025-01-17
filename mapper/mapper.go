package mapper

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"

	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-search-controller/config"
	"github.com/ONSdigital/dp-frontend-search-controller/data"
	model "github.com/ONSdigital/dp-frontend-search-controller/model"
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
	searchModels "github.com/ONSdigital/dp-search-api/models"
	topicModel "github.com/ONSdigital/dp-topic-api/models"
)

// CreateSearchPage maps type searchC.Response to model.Page
func CreateSearchPage(cfg *config.Config, req *http.Request, basePage coreModel.Page,
	validatedQueryParams data.SearchURLParams, categories []data.Category, topicCategories []data.Topic,
	respC *searchModels.SearchResponse, lang string, homepageResponse zebedee.HomepageContent, errorMessage string,
	navigationContent *topicModel.Navigation,
) model.SearchPage {
	page := model.SearchPage{
		Page: basePage,
	}

	MapCookiePreferences(req, &page.Page.CookiesPreferencesSet, &page.Page.CookiesPolicy)

	page.Metadata.Title = "Search" //nolint:goconst //The strings aren't actually the same.
	page.Type = "search"           //nolint:goconst //The strings aren't actually the same.
	page.Title.LocaliseKeyName = "SearchResults"
	page.Data.TermLocalKey = "Results"
	page.Count = respC.Count
	page.Language = lang
	page.BetaBannerEnabled = true
	page.SearchDisabled = false
	page.URI = req.URL.RequestURI()
	page.PatternLibraryAssetsPath = cfg.PatternLibraryAssetsPath
	page.Pagination.CurrentPage = validatedQueryParams.CurrentPage
	page.ServiceMessage = homepageResponse.ServiceMessage
	page.EmergencyBanner = mapEmergencyBanner(homepageResponse)
	page.SearchNoIndexEnabled = true
	page.FeatureFlags.IsPublishing = cfg.IsPublishing
	if navigationContent != nil {
		page.NavigationContent = mapNavigationContent(*navigationContent)
	}

	if validatedQueryParams.NLPWeightingEnabled {
		page.ABTest.GTMKey = "nlpSearch"
	} else {
		page.ABTest.GTMKey = "search"
	}

	mapQuery(cfg, &page, validatedQueryParams, respC, *req, errorMessage)

	mapResponse(&page, respC, categories)

	mapFilters(&page, categories, validatedQueryParams)

	mapTopicFilters(cfg, &page, topicCategories, validatedQueryParams)

	return page
}

// CreateDataAggregationPage maps type searchC.Response to model.Page
func CreateDataAggregationPage(cfg *config.Config, req *http.Request, basePage coreModel.Page,
	validatedQueryParams data.SearchURLParams, categories []data.Category, topicCategories []data.Topic, _ []data.PopulationTypes, _ []data.Dimensions,
	respC *searchModels.SearchResponse, lang string, homepageResponse zebedee.HomepageContent, errorMessage string,
	navigationContent *topicModel.Navigation,
	template string,
) model.SearchPage {
	page := model.SearchPage{
		Page: basePage,
	}

	MapCookiePreferences(req, &page.Page.CookiesPreferencesSet, &page.Page.CookiesPolicy)

	mapDataPage(&page, respC, lang, req, cfg, validatedQueryParams, homepageResponse, navigationContent, template)

	mapQuery(cfg, &page, validatedQueryParams, respC, *req, errorMessage)

	mapResponse(&page, respC, categories)

	mapFilters(&page, categories, validatedQueryParams)

	mapTopicFilters(cfg, &page, topicCategories, validatedQueryParams)

	return page
}

func mapDataPage(page *model.SearchPage, respC *searchModels.SearchResponse, lang string, req *http.Request, cfg *config.Config, validatedQueryParams data.SearchURLParams, homepageResponse zebedee.HomepageContent, navigationContent *topicModel.Navigation, template string) {
	switch template {
	case "all-adhocs":
		page.Metadata.Title = "User requested data"
		page.Title.LocaliseKeyName = "UserRequestedData"
		page.Data.KeywordFilterEnabled = true
		page.Data.DateFilterEnabled = true
	case "home-datalist":
		page.Metadata.Title = "All data related to home"
		page.Title.LocaliseKeyName = "DataList"
		page.Data.KeywordFilterEnabled = true
		page.Data.ContentTypeFilterEnabled = true
		page.Data.DateFilterEnabled = true
		page.Data.EnableHomeSwitch = true
	case "home-publications":
		page.Metadata.Title = "All publications related to home"
		page.Title.LocaliseKeyName = "HomePublications"
		page.Data.EnableHomeSwitch = true
		page.Data.KeywordFilterEnabled = true
		page.Data.ContentTypeFilterEnabled = true
	case "all-methodologies":
		page.Metadata.Title = "All methodology"
		page.Title.LocaliseKeyName = "AllMethodology"
		page.Data.KeywordFilterEnabled = true
		page.Data.TopicFilterEnabled = true
	case "published-requests":
		page.Metadata.Title = "Freedom of Information (FOI) requests"
		page.Title.LocaliseKeyName = "FOIRequests"
		page.Data.KeywordFilterEnabled = true
		page.Data.DateFilterEnabled = true
	case "home-list":
		page.Metadata.Title = "List of all home"
		page.Title.LocaliseKeyName = "HomeList"
		page.Data.KeywordFilterEnabled = true
	case "home-methodology":
		page.Metadata.Title = "Methodology related to home"
		page.Title.LocaliseKeyName = "HomeMethodology"
		page.Data.KeywordFilterEnabled = true
	case "time-series-tool":
		page.Metadata.Title = "Time series explorer"
		page.Title.LocaliseKeyName = "TimeSeriesExplorer"
		page.Data.KeywordFilterEnabled = true
		page.Data.UpdatedFilterEnabled = true
		page.Data.DateFilterEnabled = true
		page.Data.TopicFilterEnabled = true
		page.Data.EnableTimeSeriesExport = true
	default:
		page.Metadata.Title = template
	}
	page.Type = "Data Aggregation Page"
	page.Data.TermLocalKey = "Results"
	page.Count = respC.Count
	page.Language = lang
	page.BetaBannerEnabled = true
	page.SearchDisabled = false
	page.URI = req.URL.RequestURI()
	page.PatternLibraryAssetsPath = cfg.PatternLibraryAssetsPath
	page.Pagination.CurrentPage = validatedQueryParams.CurrentPage
	page.ServiceMessage = homepageResponse.ServiceMessage
	page.EmergencyBanner = mapEmergencyBanner(homepageResponse)
	page.SearchNoIndexEnabled = true
	page.FeatureFlags.IsPublishing = cfg.IsPublishing
	if navigationContent != nil {
		page.NavigationContent = mapNavigationContent(*navigationContent)
	}

	page.AfterDate = coreModel.InputDate{
		Language:        page.Language,
		Id:              "after-date",
		InputNameDay:    "after-day",
		InputNameMonth:  "after-month",
		InputNameYear:   "after-year",
		InputValueDay:   validatedQueryParams.AfterDate.DayString(),
		InputValueMonth: validatedQueryParams.AfterDate.MonthString(),
		InputValueYear:  validatedQueryParams.AfterDate.YearString(),
		Title: coreModel.Localisation{
			LocaleKey: "ReleasedAfter",
			Plural:    1,
		},
		Description: coreModel.Localisation{
			LocaleKey: "DateFilterDescription",
			Plural:    1,
		},
	}
}

// CreateSearchPage maps type searchC.Response to model.Page
func CreateDataFinderPage(cfg *config.Config, req *http.Request, basePage coreModel.Page,
	validatedQueryParams data.SearchURLParams, categories []data.Category, topicCategories []data.Topic, populationTypes []data.PopulationTypes, dimensions []data.Dimensions,
	respC *searchModels.SearchResponse, lang string, homepageResponse zebedee.HomepageContent, errorMessage string,
	navigationContent *topicModel.Navigation,
) model.SearchPage {
	page := model.SearchPage{
		Page: basePage,
	}

	MapCookiePreferences(req, &page.Page.CookiesPreferencesSet, &page.Page.CookiesPolicy)

	page.Metadata.Title = "Search"
	page.Type = "search"
	page.Title.LocaliseKeyName = "FindCensusData"
	page.Data.TermLocalKey = "DatasetsLower"
	page.Count = respC.Count
	page.Language = lang
	page.BetaBannerEnabled = true
	page.Page.Breadcrumb = []coreModel.TaxonomyNode{{Title: "Home", URI: "/"}, {Title: "Census", URI: "/census"}, {Title: "Find census data"}}
	page.SearchDisabled = false
	page.URI = req.URL.RequestURI()
	page.PatternLibraryAssetsPath = cfg.PatternLibraryAssetsPath
	page.Pagination.CurrentPage = validatedQueryParams.CurrentPage
	page.ServiceMessage = homepageResponse.ServiceMessage
	page.EmergencyBanner = mapEmergencyBanner(homepageResponse)
	page.SearchNoIndexEnabled = true
	if navigationContent != nil {
		page.NavigationContent = mapNavigationContent(*navigationContent)
	}
	mapDatasetQuery(cfg, &page, validatedQueryParams, respC, *req, errorMessage)

	mapResponse(&page, respC, categories)

	mapCensusTopicFilters(cfg, &page, topicCategories, validatedQueryParams)

	mapPopulationTypesFilters(cfg, &page, populationTypes, validatedQueryParams)

	mapDimensionsFilters(cfg, &page, dimensions, validatedQueryParams)

	return page
}

func mapQuery(cfg *config.Config, page *model.SearchPage, validatedQueryParams data.SearchURLParams, respC *searchModels.SearchResponse, req http.Request, errorMessage string) {
	page.Data.Query = validatedQueryParams.Query

	page.Data.Filter = validatedQueryParams.Filter.Query

	page.Data.ErrorMessage = errorMessage

	mapSort(page, validatedQueryParams)

	mapPagination(cfg, req, page, validatedQueryParams, respC)
}

func mapDatasetQuery(cfg *config.Config, page *model.SearchPage, validatedQueryParams data.SearchURLParams, respC *searchModels.SearchResponse, req http.Request, errorMessage string) {
	page.Data.Query = validatedQueryParams.Query

	page.Data.Filter = validatedQueryParams.Filter.Query

	page.Data.ErrorMessage = errorMessage

	mapDatasetSort(page, validatedQueryParams)

	mapPagination(cfg, req, page, validatedQueryParams, respC)
}

func mapSort(page *model.SearchPage, validatedQueryParams data.SearchURLParams) {
	page.Data.Sort.Query = validatedQueryParams.Sort.Query

	page.Data.Sort.LocaliseFilterKeys = validatedQueryParams.Filter.LocaliseKeyName

	page.Data.Sort.LocaliseSortKey = validatedQueryParams.Sort.LocaliseKeyName

	pageSortOptions := make([]model.SortOptions, len(data.SortOptions))
	for i := range data.SortOptions {
		pageSortOptions[i] = model.SortOptions{
			Query:           data.SortOptions[i].Query,
			LocaliseKeyName: data.SortOptions[i].LocaliseKeyName,
		}
	}

	page.Data.Sort.Options = pageSortOptions
}

func mapDatasetSort(page *model.SearchPage, validatedQueryParams data.SearchURLParams) {
	page.Data.Sort.Query = validatedQueryParams.Sort.Query

	page.Data.Sort.LocaliseFilterKeys = validatedQueryParams.Filter.LocaliseKeyName

	page.Data.Sort.LocaliseSortKey = validatedQueryParams.Sort.LocaliseKeyName

	pageSortOptions := make([]model.SortOptions, len(data.DatasetSortOptions))
	for i := range data.DatasetSortOptions {
		pageSortOptions[i] = model.SortOptions{
			Query:           data.DatasetSortOptions[i].Query,
			LocaliseKeyName: data.DatasetSortOptions[i].LocaliseKeyName,
		}
	}

	page.Data.Sort.Options = pageSortOptions
}

func mapPagination(cfg *config.Config, req http.Request, page *model.SearchPage, validatedQueryParams data.SearchURLParams, respC *searchModels.SearchResponse) {
	page.Data.Pagination.Limit = validatedQueryParams.Limit
	page.Data.Pagination.LimitOptions = data.LimitOptions

	page.Data.Pagination.CurrentPage = validatedQueryParams.CurrentPage
	page.Data.Pagination.TotalPages = data.GetTotalPages(cfg, validatedQueryParams.Limit, respC.Count)
	page.Data.Pagination.PagesToDisplay = data.GetPagesToDisplay(cfg, req, validatedQueryParams, page.Data.Pagination.TotalPages)
	page.Data.Pagination.FirstAndLastPages = data.GetFirstAndLastPages(req, validatedQueryParams, page.Data.Pagination.TotalPages)
}

func mapResponse(page *model.SearchPage, respC *searchModels.SearchResponse, categories []data.Category) {
	page.Data.Response.Count = respC.Count

	mapResponseCategories(page, categories)

	mapResponseItems(page, respC)

	page.Data.Response.Suggestions = respC.Suggestions
	page.Data.Response.AdditionalSuggestions = respC.AdditionSuggestions
}

func mapResponseItems(page *model.SearchPage, respC *searchModels.SearchResponse) {
	itemPage := []model.ContentItem{}
	for i := range respC.Items {
		item := model.ContentItem{}

		mapItemDescription(&item, &respC.Items[i])

		mapItemHighlight(&item, &respC.Items[i])

		item.Type.Type = respC.Items[i].DataType
		item.Type.LocaliseKeyName = data.GetGroupLocaliseKey(respC.Items[i].DataType)

		item.URI = respC.Items[i].URI
		item.Dataset.PopulationType = respC.Items[i].PopulationType

		itemPage = append(itemPage, item)
	}

	page.Data.Response.Items = itemPage
}

func mapItemDescription(item *model.ContentItem, itemC *searchModels.Item) {
	item.Description = model.Description{
		CDID:            itemC.CDID,
		DatasetID:       itemC.DatasetID,
		Language:        itemC.Language,
		MetaDescription: itemC.MetaDescription,
		ReleaseDate:     itemC.ReleaseDate,
		Summary:         itemC.Summary,
		Title:           itemC.Title,
	}

	if len(itemC.Keywords) != 0 {
		item.Description.Keywords = itemC.Keywords
	} else {
		item.Description.Keywords = nil
	}
}

func mapItemHighlight(item *model.ContentItem, itemC *searchModels.Item) {
	itemHighlight := itemC.Highlight
	if !reflect.ValueOf(itemHighlight).IsNil() {
		item.Description.Highlight = model.Highlight{
			DatasetID:       itemHighlight.DatasetID,
			Edition:         itemC.Edition,
			Keywords:        itemHighlight.Keywords,
			MetaDescription: itemHighlight.MetaDescription,
			Summary:         itemHighlight.Summary,
			Title:           itemHighlight.Title,
		}
	} else {
		item.Description.Highlight = model.Highlight{}
	}
}

func mapResponseCategories(page *model.SearchPage, categories []data.Category) {
	pageCategories := []model.Category{}

	for _, category := range categories {
		pageContentType := []model.ContentType{}

		for _, contentType := range category.ContentTypes {
			pageContentType = append(pageContentType, model.ContentType{
				Group:           contentType.Group,
				Count:           contentType.Count,
				LocaliseKeyName: contentType.LocaliseKeyName,
				Types:           contentType.Types,
			})
		}

		pageCategories = append(pageCategories, model.Category{
			Count:           category.Count,
			LocaliseKeyName: category.LocaliseKeyName,
			ContentTypes:    pageContentType,
		})
	}

	page.Data.Response.Categories = pageCategories
}

func mapFilters(page *model.SearchPage, categories []data.Category, queryParams data.SearchURLParams) {
	filters := make([]model.Filter, len(categories))

	for i := range categories {
		var filter model.Filter
		filter.LocaliseKeyName = categories[i].LocaliseKeyName
		filter.NumberOfResults = categories[i].Count

		var keys []string
		var subTypes []model.Filter
		if len(categories[i].ContentTypes) > 0 {
			for _, contentType := range categories[i].ContentTypes {
				if !contentType.ShowInWebUI && contentType.Count > 0 {
					filter.NumberOfResults -= contentType.Count
					continue
				}
				var subType model.Filter
				subType.LocaliseKeyName = contentType.LocaliseKeyName
				subType.NumberOfResults = contentType.Count
				subType.FilterKey = []string{contentType.Group}

				isChecked := mapIsChecked(subType.FilterKey, queryParams)
				subType.IsChecked = isChecked
				subTypes = append(subTypes, subType)

				keys = append(keys, contentType.Group)
			}
		}

		filter.Types = subTypes
		filter.FilterKey = keys
		filter.IsChecked = mapIsChecked(filter.FilterKey, queryParams)
		filters[i] = filter
	}

	page.Data.Filters = filters
}

func mapTopicFilters(cfg *config.Config, page *model.SearchPage, topicCategories []data.Topic, queryParams data.SearchURLParams) {
	if !cfg.EnableCensusTopicFilterOption {
		return
	}

	var topicsQueryParam []string
	if queryParams.TopicFilter != "" {
		topicsQueryParam = strings.Split(queryParams.TopicFilter, ",")
	}

	mapTopicQueryParams := make(map[string]bool)
	for i := range topicsQueryParam {
		mapTopicQueryParams[topicsQueryParam[i]] = true
	}

	topicFilters := make([]model.TopicFilter, len(topicCategories))

	for i := range topicCategories {
		if !topicCategories[i].ShowInWebUI {
			continue
		}

		var topicFilter model.TopicFilter

		topicFilter.LocaliseKeyName = topicCategories[i].LocaliseKeyName
		topicFilter.NumberOfResults = topicCategories[i].Count
		topicFilter.Query = topicCategories[i].Query
		topicFilter.DistinctItemsCount = topicCategories[i].DistinctItemsCount

		if len(topicsQueryParam) > 0 {
			topicFilter.IsChecked = true
		}

		topicFilters[i] = topicFilter

		for j := range topicCategories[i].Subtopics {
			if !topicCategories[i].Subtopics[j].ShowInWebUI {
				continue
			}
			var subtopicFilter model.TopicFilter

			subtopicFilter.LocaliseKeyName = topicCategories[i].Subtopics[j].LocaliseKeyName
			subtopicFilter.NumberOfResults = topicCategories[i].Subtopics[j].Count
			subtopicFilter.Query = topicCategories[i].Subtopics[j].Query

			if mapTopicQueryParams[topicCategories[i].Subtopics[j].Query] {
				subtopicFilter.IsChecked = true
			}

			topicFilters[i].Types = append(topicFilters[i].Types, subtopicFilter)
		}
	}

	page.Data.TopicFilters = topicFilters
}

func mapCensusTopicFilters(cfg *config.Config, page *model.SearchPage, topicCategories []data.Topic, queryParams data.SearchURLParams) {
	if !cfg.EnableCensusTopicFilterOption {
		return
	}

	var topicsQueryParam []string
	if queryParams.TopicFilter != "" {
		topicsQueryParam = strings.Split(queryParams.TopicFilter, ",")
	}

	mapTopicQueryParams := make(map[string]bool)
	for i := range topicsQueryParam {
		mapTopicQueryParams[topicsQueryParam[i]] = true
	}

	topicFilters := make([]model.TopicFilter, len(topicCategories))

	for i := range topicCategories {
		if !topicCategories[i].ShowInWebUI {
			continue
		}

		var topicFilter model.TopicFilter

		topicFilter.LocaliseKeyName = topicCategories[i].LocaliseKeyName
		topicFilter.NumberOfResults = topicCategories[i].Count
		topicFilter.Query = topicCategories[i].Query
		topicFilter.DistinctItemsCount = topicCategories[i].DistinctItemsCount

		if len(topicsQueryParam) > 0 {
			topicFilter.IsChecked = true
		}

		topicFilters[i] = topicFilter

		for j := range topicCategories[i].Subtopics {
			if !topicCategories[i].Subtopics[j].ShowInWebUI {
				continue
			}
			var subtopicFilter model.TopicFilter

			subtopicFilter.LocaliseKeyName = topicCategories[i].Subtopics[j].LocaliseKeyName
			subtopicFilter.NumberOfResults = topicCategories[i].Subtopics[j].Count
			subtopicFilter.Query = topicCategories[i].Subtopics[j].Query

			if mapTopicQueryParams[topicCategories[i].Subtopics[j].Query] {
				subtopicFilter.IsChecked = true
			}

			topicFilters[i].Types = append(topicFilters[i].Types, subtopicFilter)
		}
	}

	page.Data.CensusFilters = topicFilters[0].Types
}

func mapPopulationTypesFilters(cfg *config.Config, page *model.SearchPage, populationTypes []data.PopulationTypes, queryParams data.SearchURLParams) {
	if !cfg.EnableCensusPopulationTypesFilterOption {
		return
	}

	var popultationTypesQueryParam []string
	if queryParams.PopulationTypeFilter != "" {
		popultationTypesQueryParam = strings.Split(queryParams.PopulationTypeFilter, ",")
	}

	mapPopultationTypesQueryParams := make(map[string]bool)
	for i := range popultationTypesQueryParam {
		mapPopultationTypesQueryParams[popultationTypesQueryParam[i]] = true
	}

	populationTypeFilters := make([]model.PopulationTypeFilter, len(populationTypes))

	for i := range populationTypes {
		if !populationTypes[i].ShowInWebUI {
			continue
		}

		var populationTypesFilter model.PopulationTypeFilter

		populationTypesFilter.LocaliseKeyName = populationTypes[i].LocaliseKeyName
		populationTypesFilter.NumberOfResults = populationTypes[i].Count
		populationTypesFilter.Query = queryParams.Query
		populationTypesFilter.Count = populationTypes[i].Count
		populationTypesFilter.Type = populationTypes[i].Type

		if len(popultationTypesQueryParam) > 0 {
			for _, v := range popultationTypesQueryParam {
				if v == populationTypesFilter.LocaliseKeyName {
					populationTypesFilter.IsChecked = true
				}
			}
		}

		populationTypeFilters[i] = populationTypesFilter
	}
	page.Data.PopulationTypeFilter = populationTypeFilters
}

func mapDimensionsFilters(cfg *config.Config, page *model.SearchPage, dimensions []data.Dimensions, queryParams data.SearchURLParams) {
	if !cfg.EnableCensusDimensionsFilterOption {
		return
	}

	var dimensionsQueryParam []string
	if queryParams.PopulationTypeFilter != "" {
		dimensionsQueryParam = strings.Split(queryParams.DimensionsFilter, ",")
	}

	mapPopultationTypesQueryParams := make(map[string]bool)
	for i := range dimensionsQueryParam {
		mapPopultationTypesQueryParams[dimensionsQueryParam[i]] = true
	}

	dimensionsFilters := make([]model.DimensionsFilter, len(dimensions))

	for i := range dimensions {
		if !dimensions[i].ShowInWebUI {
			continue
		}

		var dimensionsFilter model.DimensionsFilter

		dimensionsFilter.LocaliseKeyName = dimensions[i].LocaliseKeyName
		dimensionsFilter.NumberOfResults = dimensions[i].Count
		dimensionsFilter.Query = queryParams.Query
		dimensionsFilter.Count = dimensions[i].Count
		dimensionsFilter.Type = dimensions[i].Type

		if len(dimensionsQueryParam) > 0 {
			for _, v := range dimensionsQueryParam {
				if v == dimensionsFilter.LocaliseKeyName {
					dimensionsFilter.IsChecked = true
				}
			}
		}

		dimensionsFilters[i] = dimensionsFilter
	}

	page.Data.DimensionsFilter = dimensionsFilters
}

func mapIsChecked(contentTypes []string, queryParams data.SearchURLParams) bool {
	for _, query := range queryParams.Filter.Query {
		for _, contentType := range contentTypes {
			if query == contentType {
				return true
			}
		}
	}

	return false
}

// MapCookiePreferences reads cookie policy and preferences cookies and then maps the values to the page model
func MapCookiePreferences(req *http.Request, preferencesIsSet *bool, policy *coreModel.CookiesPolicy) {
	preferencesCookie := cookies.GetCookiePreferences(req)
	*preferencesIsSet = preferencesCookie.IsPreferenceSet

	*policy = coreModel.CookiesPolicy{
		Essential: preferencesCookie.Policy.Essential,
		Usage:     preferencesCookie.Policy.Usage,
	}
}

func mapEmergencyBanner(hpc zebedee.HomepageContent) coreModel.EmergencyBanner {
	var mappedEmergencyBanner coreModel.EmergencyBanner
	emptyBannerObj := zebedee.EmergencyBanner{}
	bannerData := hpc.EmergencyBanner

	if bannerData != emptyBannerObj {
		mappedEmergencyBanner.Title = bannerData.Title
		mappedEmergencyBanner.Type = strings.Replace(bannerData.Type, "_", "-", -1)
		mappedEmergencyBanner.Description = bannerData.Description
		mappedEmergencyBanner.URI = bannerData.URI
		mappedEmergencyBanner.LinkText = bannerData.LinkText
	}

	return mappedEmergencyBanner
}

// mapNavigationContent takes navigationContent as returned from the client and returns information needed for the navigation bar
func mapNavigationContent(navigationContent topicModel.Navigation) []coreModel.NavigationItem {
	var mappedNavigationContent []coreModel.NavigationItem

	if navigationContent.Items != nil {
		for _, rootContent := range *navigationContent.Items {
			var subItems []coreModel.NavigationItem

			if rootContent.SubtopicItems != nil {
				for _, subtopicContent := range *rootContent.SubtopicItems {
					subItems = append(subItems, coreModel.NavigationItem{
						Uri:   subtopicContent.URI,
						Label: subtopicContent.Label,
					})
				}
			}

			mappedNavigationContent = append(mappedNavigationContent, coreModel.NavigationItem{
				Uri:      rootContent.URI,
				Label:    rootContent.Label,
				SubItems: subItems,
			})
		}
	}

	return mappedNavigationContent
}
