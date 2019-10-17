package main


import (
    "fmt"
    "log"
    "net/http"
    "strconv"
    "sync"


	"flag"
	
	"time"
	

	"strings"

	gce "cloud.google.com/go/compute/metadata"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	monitoring "google.golang.org/api/monitoring/v3"

)

var counter int
var mutex = &sync.Mutex{}


var podId string		//30
var custom_metric_value int
var custom_metric_value_int64 int64
var metricName *string
var metricValue int64
var labels map[string]string
var stackdriverService *monitoring.Service


func echoString(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello")
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {

	custom_metric_value, err := strconv.Atoi(r.URL.Path[1:])
	if err != nil {
		log.Fatalf("Error getting Stackdriver service: %v", err)
	}

 	custom_metric_value_int64 := int64(custom_metric_value)
	//metricValue := flag.Int64("metric-value", custom_metric_value_int64_1, "custom metric value")

	if counter < 1 {
		// Gather pod information
		podId := flag.String("pod-id", "", "pod id")
		metricName1 := flag.String("metric-name", "response_time1", "custom metric name")
		//metricValue := flag.Int64("metric-value", 20, "custom metric value")
		//custom_metric_value := 20
		//custom_metric_value = counter
		//custom_metric_value_int64 := int64(custom_metric_value)
		metricValue := flag.Int64("metric-value", custom_metric_value_int64, "custom metric value")
	
		metricName = metricName1
		
		//metricValue := 20
		flag.Parse()

		if *podId == "" {
			log.Fatalf("No pod id specified.")
		}			//61

		stackdriverService1, err := getStackDriverService()
		stackdriverService = stackdriverService1
		if err != nil {
			log.Fatalf("Error getting Stackdriver service: %v", err)
		}

		labels1 := getResourceLabels(*podId)
		labels = labels1
		
		err2 := exportMetric(stackdriverService, *metricName1, custom_metric_value_int64, labels)
		if err2 != nil {
			log.Printf("Failed to write time series data: %v\n", err2)
		} else {
			log.Printf("Finished writing time series with value: %v\n", metricValue)
		}
	
	} else {
		//custom_metric_value = counter
		//custom_metric_value_int64_2 := int64(custom_metric_value)
		
		fmt.Fprintf(w, strconv.Itoa(custom_metric_value))
		fmt.Fprintf(w, "Hi there, my metric name is %s!", *metricName)
		//fmt.Fprintf(w, "Hi there, my number is %s!", r.URL.Path[1:])
		
		err3 := exportMetric(stackdriverService, *metricName, custom_metric_value_int64, labels)
		if err3 != nil {
			log.Printf("Failed to write time series data: %v\n", err3)
		} else {
			log.Printf("Finished writing time series with value: %v\n", metricValue)
		}
	}


    mutex.Lock()
    counter++
    fmt.Fprintf(w, strconv.Itoa(counter))
    mutex.Unlock()
}

func main() {
    //http.HandleFunc("/", echoString)
	http.HandleFunc("/", incrementCounter)

    http.HandleFunc("/increment", incrementCounter)

    http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hi")
    })

    log.Fatal(http.ListenAndServe(":8081", nil))

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
