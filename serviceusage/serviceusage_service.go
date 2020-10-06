package serviceusage

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/api/serviceusage/v1"
)

const (
	StateEnabled  = "ENABLED"
	StateDisabled = "DISABLED"
)

type ServiceUsageService struct {
	client *serviceusage.Service
}

func NewService(ctx context.Context, client *serviceusage.Service) (*ServiceUsageService, error) {
	return &ServiceUsageService{
		client: client,
	}, nil
}

type ServiceUsage struct {
	Name  string // ex.bigquery.googleapis.com
	State string // ENABLED or DISABLED
}

func (s *ServiceUsageService) ListAll(ctx context.Context, projectNumber int64) ([]*ServiceUsage, error) {
	return s.list(ctx, projectNumber, "")
}

func (s *ServiceUsageService) ListByStateEnabled(ctx context.Context, projectNumber int64) ([]*ServiceUsage, error) {
	return s.list(ctx, projectNumber, StateEnabled)
}

func (s *ServiceUsageService) ListByStateDisabled(ctx context.Context, projectNumber int64) ([]*ServiceUsage, error) {
	return s.list(ctx, projectNumber, StateDisabled)
}

// ListByDiff is base と target のServiceUsageを比較して、差がある場合、target の ServiceUsage を返す
func (s *ServiceUsageService) ListByDiff(ctx context.Context, baseProjectNumber int64, targetProjectNumber int64) ([]*ServiceUsage, error) {
	baseList, err := s.ListAll(ctx, baseProjectNumber)
	if err != nil {
		return nil, fmt.Errorf("failed List ServiceUsage. projectNumber:%d : %w", baseProjectNumber, err)
	}

	target, err := s.ListAll(ctx, targetProjectNumber)
	if err != nil {
		return nil, fmt.Errorf("failed List ServiceUsage. projectNumber:%d : %w", targetProjectNumber, err)
	}
	targetMap := map[string]*ServiceUsage{}
	for _, v := range target {
		targetMap[v.Name] = v
	}

	var result []*ServiceUsage
	for _, base := range baseList {
		target, ok := targetMap[base.Name]
		if !ok {
			result = append(result, target)
			continue
		}
		if base.State != target.State {
			result = append(result, target)
		}
	}
	return result, nil
}

func (s *ServiceUsageService) list(ctx context.Context, projectNumber int64, state string) ([]*ServiceUsage, error) {
	var results []*ServiceUsage
	var nextPageToken string
	for {
		call := s.client.Services.List(fmt.Sprintf("projects/%d", projectNumber)).Context(ctx)
		if state != "" {
			call.Filter(fmt.Sprintf("state:%s", state))
		}
		if nextPageToken != "" {
			call.PageToken(nextPageToken)
		}
		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for _, v := range resp.Services {
			nl := strings.Split(v.Name, "/")
			results = append(results, &ServiceUsage{
				Name:  nl[len(nl)-1],
				State: v.State,
			})
		}
		nextPageToken = resp.NextPageToken
		if resp.NextPageToken == "" {
			break
		}
	}

	return results, nil
}
