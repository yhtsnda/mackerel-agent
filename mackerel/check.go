package mackerel

import (
	"encoding/json"

	"github.com/mackerelio/mackerel-agent/checks"
)

type monitoringChecksPayload struct {
	Checks []*checkReport `json:"checks"`
}

type checkReport struct {
	Target     monitorTargetHost `json:"target"`
	Name       string            `json:"name"`
	Status     checks.Status     `json:"status"`
	Message    string            `json:"message"`
	OccurredAt Time              `json:"occurredAt"`
}

type monitorTargetHost struct {
	HostID string
}

// MarshalJSON implements json.Marshaler.
func (h monitorTargetHost) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":   "host",
		"hostId": h.HostID,
	})
}

// ReportCheckMonitors sends reports of checks.Checker() to Mackrel API server.
func (api *API) ReportCheckMonitors(hostID string, reports []*checks.Report) error {
	payload := &monitoringChecksPayload{
		Checks: make([]*checkReport, len(reports)),
	}
	for i, report := range reports {
		payload.Checks[i] = &checkReport{
			Target:     monitorTargetHost{HostID: hostID},
			Name:       report.Name,
			Status:     report.Status,
			Message:    report.Message,
			OccurredAt: Time(report.OccurredAt),
		}
	}
	return api.postJSON("/api/v0/monitoring/checks", payload, nil)
}