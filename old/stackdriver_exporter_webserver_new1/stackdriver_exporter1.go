package main

import (
	"flag"
	"fmt"
	"log"
	//"time"

	"net/http"

	//"strings"

	//"strconv"

	//gce "cloud.google.com/go/compute/metadata"
	/*"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	monitoring "google.golang.org/api/monitoring/v3"*/
)


var count int = 0

//var podId *string
//var metricName string
//var metricValue int64

//var labels map[string]string
//var stackdriverService *monitoring.Service

func main() {

	//count = count + 1

	// Gather pod information
	//podId_pointer := flag.String("pod-id", "", "pod id")
	/*if *podId_pointer == "" {
		log.Fatalf("No pod id specified.")
		//podId = "No pod id specified111."
	} else {
	   	//podId = *podId_pointer
	}

	podId = "value1" */
	/*stackdriverService1, err := getStackDriverService()
	if err != nil {
		log.Fatalf("Error getting Stackdriver service: %v", err)
	}
	stackdriverService = stackdriverService1 */
	
	/*metricName_pointer := flag.String("metric-name", "response_time1", "custom metric name")
	metricName = *metricName_pointer*/

	http.HandleFunc("/", handler)
    	log.Fatal(http.ListenAndServe(":8083", nil))


}

func handler(w http.ResponseWriter, r *http.Request) {

	
	//fmt.Fprintf(w, "Hi there, my number is %s count is %s! my podId is %s, my metricName is %s", r.URL.Path[1:], count, podId, metricName)
	fmt.Fprintf(w, "Hi there, my number is %s count is %s!", r.URL.Path[1:], count)
	if count < 1 {

		// Gather pod information
		/*podId1 := flag.String("pod-id", "", "pod id")
		if *podId1 == "" {
			log.Fatalf("No pod id specified.")
			//podId = "No pod id specified111."
		} 

		podId = podId1*/
		//metricName := flag.String("metric-name", "response_time1", "custom metric name")
		//flag.Parse()
		/*stackdriverService1, err := getStackDriverService()
		if err != nil {
			log.Fatalf("Error getting Stackdriver service: %v", err)
		}
		stackdriverService = stackdriverService1*/

		//fmt.Fprintf(w, "count2222 is %s!, my metricName is %s", count)//, metricName)
		fmt.Fprintf(w, "count2222 is %s!", count)
			
	} else {
		fmt.Fprintf(w, "count3333 is %s!", count)//, metricName)
	}
	count = count + 1

}

/*func getStackDriverService() (*monitoring.Service, error) {
	oauthClient := oauth2.NewClient(context.Background(), google.ComputeTokenSource(""))
	return monitoring.New(oauthClient)
}*/


