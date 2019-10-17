package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"net/http"

	"strings"

	"strconv"

	gce "cloud.google.com/go/compute/metadata"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	monitoring "google.golang.org/api/monitoring/v3"
)



var count int = 1

var podId string
var metricName string
var metricValue int64

var labels2 map[string]string
var stackdriverService *monitoring.Service


func main() {


	http.HandleFunc("/", handler)
    	log.Fatal(http.ListenAndServe(":8083", nil))


}

func handler(w http.ResponseWriter, r *http.Request) {
	count++
	fmt.Fprintf(w, "Hi there, my number is %s!", r.URL.Path[1:])

	// converting the r.URL.Path[1:] variable into an int using Atoi method
	custom_metric_value, err := strconv.Atoi(r.URL.Path[1:])
	if err != nil {
		log.Fatalf("Error getting Stackdriver service: %v", err)
	}
 	custom_metric_value_int64 := int64(custom_metric_value)
	

	// Gather pod information
	podId := flag.String("pod-id", "", "pod id")
	
	metricName := flag.String("metric-name", "response_time1", "custom metric name")
	
	metricValue := flag.Int64("metric-value", custom_metric_value_int64, "custom metric value")
	flag.Parse()

	if *podId == "" {
		log.Fatalf("No pod id specified.")
	}

	//var labels2 map[string]string //:= getResourceLabels(*podId)
	//
	//var err2
	labels2 := getResourceLabels(*podId)

	if count < 2 {
		

		stackdriverService, err := getStackDriverService()
		if err != nil {
			log.Fatalf("Error getting Stackdriver service: %v", err)
		}

		


		//-----------
		err2 := exportMetric(stackdriverService, *metricName, *metricValue, labels2)
		if err2 != nil {
			log.Printf("Failed to write time series data: %v\n", err2)
		} else {
			log.Printf("Finished writing time series with value: %v\n", metricValue)
		}
		time.Sleep(5000 * time.Millisecond)
	//-----------
	} else {
		fmt.Fprintf(w, "Hi there, my number2 is %s!", r.URL.Path[1:])

		//labels2 := getResourceLabels(*podId)
		//-----------
		/*err2 := exportMetric(stackdriverService, *metricName, *metricValue, labels2)
		if err2 != nil {
			log.Printf("Failed to write time series data: %v\n", err2)
		} else {
			log.Printf("Finished writing time series with value: %v\n", metricValue)
		}
		time.Sleep(5000 * time.Millisecond)*/
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


		//var metricValue2 int64 = 1411
		dataPoint := &monitoring.Point{
			Interval: &monitoring.TimeInterval{
				EndTime: time.Now().Format(time.RFC3339),
			},
			Value: &monitoring.TypedValue{
				Int64Value: &metricValue,
			},
			// Value: &monitoring.TypedValue{
			// 	Int64Value: &metricValue2,
			// },
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
