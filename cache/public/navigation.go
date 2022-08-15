package public

import (
	"context"

	topicModel "github.com/ONSdigital/dp-topic-api/models"
	topicCli "github.com/ONSdigital/dp-topic-api/sdk"
	"github.com/ONSdigital/log.go/v2/log"
)

func UpdateNavigationData(ctx context.Context, lang string, topicClient topicCli.Clienter) func() *topicModel.Navigation {
	return func() *topicModel.Navigation {
		headers := topicCli.Headers{}
		options := topicCli.Options{}

		switch lang {
		case "cy":
			options.Lang = topicCli.Welsh
		default:
			options.Lang = topicCli.English
		}

		navigationData, err := topicClient.GetNavigationPublic(ctx, headers, options)
		if err != nil {
			logData := log.Data{
				"headers": headers,
				"options": options,
			}
			log.Error(ctx, "failed to get navigation data from client", err, logData)
		}

		return navigationData
	}
}