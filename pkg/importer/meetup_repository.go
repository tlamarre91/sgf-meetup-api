package importer

import (
	"context"
	"encoding/json"
	"github.com/google/wire"
	"log/slog"
	"sgf-meetup-api/pkg/shared/models"
	"time"
)

type MeetupRepository interface {
	GetEventsUntilDateForGroup(ctx context.Context, group string, beforeDate time.Time) ([]models.MeetupEvent, error)
}

type GraphQLHandler interface {
	ExecuteQuery(ctx context.Context, query string, variables map[string]any) ([]byte, error)
}

type GraphQLMeetupRepository struct {
	handler GraphQLHandler
	logger  *slog.Logger
}

func NewGraphQLMeetupRepository(handler GraphQLHandler, logger *slog.Logger) *GraphQLMeetupRepository {
	return &GraphQLMeetupRepository{
		handler: handler,
		logger:  logger,
	}
}

const getFutureEventsQuery = `
  query ($urlname: String!, $itemsNum: Int!, $cursor: String) {
	events: groupByUrlname(urlname: $urlname) {
	  unifiedEvents(input: { first: $itemsNum, after: $cursor }) {
		count
		pageInfo {
		  endCursor
		  hasNextPage
		}
		edges {
		  node {
			id
			title
			eventUrl
			description
			dateTime
			duration
			venue {
			  name
			  address
			  city
			  state
			  postalCode
			}
			group {
			  name
			  urlname
			}
			host {
			  name
			}
			images {
			  baseUrl
			  preview
			}
		  }
		}
	  }
	}
  }
`

type MeetupFutureEventsResponse struct {
	Data struct {
		Events struct {
			UnifiedEvents struct {
				Count    int `json:"count"`
				PageInfo struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
				Edges []MeetupEdge `json:"edges"`
			} `json:"unifiedEvents"`
		} `json:"events"`
	} `json:"data"`
}

type MeetupEdge struct {
	Node models.MeetupEvent `json:"node"`
}

func (r *GraphQLMeetupRepository) GetEventsUntilDateForGroup(ctx context.Context, group string, beforeDate time.Time) ([]models.MeetupEvent, error) {
	events := make([]models.MeetupEvent, 0)
	cursor := ""
	var maxFutureDate time.Time

	for {
		variables := map[string]any{
			"urlname":  group,
			"itemsNum": 50,
		}

		if cursor != "" {
			variables["cursor"] = cursor
		}

		response, err := executeGraphQLQuery[MeetupFutureEventsResponse](r, ctx, getFutureEventsQuery, variables)

		if err != nil {
			return nil, err
		}

		for _, edge := range response.Data.Events.UnifiedEvents.Edges {
			event := edge.Node
			events = append(events, event)

			if event.DateTime.After(maxFutureDate) {
				maxFutureDate = event.DateTime.Time
			}
		}

		if maxFutureDate.After(beforeDate) {
			break
		}

		pageInfo := response.Data.Events.UnifiedEvents.PageInfo

		if !pageInfo.HasNextPage {
			break
		}

		cursor = pageInfo.EndCursor
	}

	return events, nil
}

func executeGraphQLQuery[T any](r *GraphQLMeetupRepository, ctx context.Context, query string, variables map[string]any) (*T, error) {
	responseBytes, err := r.handler.ExecuteQuery(ctx, query, variables)

	if err != nil {
		return nil, err
	}

	var response T
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

var MeetupRepositoryProviders = wire.NewSet(wire.Bind(new(MeetupRepository), new(*GraphQLMeetupRepository)), NewGraphQLMeetupRepository)
