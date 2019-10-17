package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"strings"

	gce "cloud.google.com/go/compute/metadata"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	monitoring "google.golang.org/api/monitoring/v3"
)

func main() {
	// Gather pod information
	podId := flag.String("pod-id", "", "pod id")
	metricName := flag.String("metric-name", "response_time1", "custom metric name")
	//metricValue := flag.Int64("metric-value", 20, "custom metric value")
	custom_metric_value := 20
	custom_metric_value_int64 := int64(custom_metric_value)
	metricValue := flag.Int64("metric-value", custom_metric_value_int64, "custom metric value")
	
	//metricValue := 20
	flag.Parse()

	if *podId == "" {
		log.Fatalf("No pod id specified.")
	}

	stackdriverService, err := getStackDriverService()
	if err != nil {
		log.Fatalf("Error getting Stackdriver service: %v", err)
	}

	labels := getResourceLabels(*podId)
	for {
		//err := exportMetric(stackdriverService, *metricName, *metricValue, labels)
		custom_metric_value++
		custom_metric_value_int64 := int64(custom_metric_value)
		err := exportMetric(stackdriverService, *metricName, custom_metric_value_int64, labels)
		if err != nil {
			log.Printf("Failed to write time series data: %v\n", err)
		} else {
			log.Printf("Finished writing time series with value: %v\n", metricValue)
		}
		time.Sleep(15000 * time.Millisecond)
		
		
	}
}

func getStackDriverService() (*monitoring.Service, error) {
	oauthClient := oauth2.NewClient(context.Background(), google.ComputeTokenSource(""))
	return monitoring.New(oauthClient)
}

// getResourceLabels returns resource labels needed to correctly label metric data
// exported to StackDriver. Labels contain details on the cluster (name, zone, project id)
// and pod for which the metric is exported (id)
func getResourceLabels(podId string) map[string]string {
	projectId, _ := gce.ProjectID()
	zone, _ := gce.Zone()
	clusterName, _ := gce.InstanceAttributeValue("cluster-name")
	clusterName = strings.TrimSpace(clusterName)
	return map[string]string{
		"project_id":   projectId,
		"zone":         zone,
		"cluster_name": clusterName,
		// container name doesn't matter here, because the metric is exported for
		// the pod, not the container
		"container_name": "",
		"pod_id":         podId,
		// namespace_id and instance_id don't matter
		"namespace_id": "default",
		"instance_id":  "",
	}
}

func exportMetric(stackdriverService *monitoring.Service, metricName string,
	metricValue int64, resourceLabels map[string]string) error {
	dataPoint := &monitoring.Point{
		Interval: &monitoring.TimeInterval{
			EndTime: time.Now().Format(time.RFC3339),
		},
		Value: &monitoring.TypedValue{
			Int64Value: &metricValue,
		},
	}
	// Write time series data.
	request := &monitoring.CreateTimeSeriesRequest{
		TimeSeries: []*monitoring.TimeSeries{
			{
				Metric: &monitoring.Metric{
					Type: "custom.googleapis.com/" + metricName,
				},
				Resource: &monitoring.MonitoredResource{
					Type:   "gke_container",
					Labels: resourceLabels,
				},
				Points: []*monitoring.Point{
					dataPoint,
				},
			},
		},
	}
	projectName := fmt.Sprintf("projects/%s", resourceLabels["project_id"])
	_, err := stackdriverService.Projects.TimeSeries.Create(projectName, request).Do()
	return err
}
