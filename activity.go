package aeroqualaqy1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/matt-doug-davidson/connector"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

// Activity is used to create a custom activity. Add values here to retain them.
// Objects used by the time are defined here.
// Common structure
type Activity struct {
	settings *Settings // Defind in metadata.go in this package
	Mappings map[string]map[string]interface{}
	status   string
}

// Metadata returns the activity's metadata
// Common function
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// The init function is executed after the package is imported. This function
// runs before any other in the package.
func init() {
	//_ = activity.Register(&Activity{})
	_ = activity.Register(&Activity{}, New)
}

// Used when the init function is called. The settings, Input and Output
// structures are optional depends application. These structures are
// defined in the metadata.go file in this package.
var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// New Looks to be used when the Activity structure contains fields that need to be
// configured using the InitContext information.
// New does this
func New(ctx activity.InitContext) (activity.Activity, error) {
	logger := ctx.Logger()
	logger.Info("aeroqualaqy1:New enter")
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		logger.Error("Failed to convert settings")
		return nil, err
	}

	// Declared an empty map interface
	var result map[string]interface{}
	json.Unmarshal([]byte(s.Mappings), &result)

	mm := map[string]map[string]interface{}{}
	for _, mapper := range result {
		mapper1 := mapper.(map[string]interface{})
		for sensor, sensorInfo := range mapper1 {
			si := sensorInfo.(map[string]interface{})
			mm[sensor] = make(map[string]interface{})
			se := map[string]interface{}{}
			f, foundF := si["field"]
			if !foundF {
				continue
			}
			se["field"] = f
			mm[sensor] = se
		}
	}

	// Create the activity with settings as defaut. Set any other field in
	//the activity here as well
	act := &Activity{settings: s, Mappings: mm, status: "RUNNING"}

	logger.Info("aeroqualaqy1:New exit")
	return act, nil
}

func (a *Activity) createErrorStatus(description string) map[string]interface{} {
	payload := map[string]interface{}{}
	payload["status"] = "ERROR"
	payload["datetime"] = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	payload["messageId"] = uuid.New().String()
	payload["description"] = description

	return payload
}

func (a *Activity) sendErrorMessage(ctx activity.Context, entity string, description string) bool {
	if a.status == "ERROR" {
		return false
	}
	a.status = "ERROR"
	output := map[string]interface{}{}
	output["data"] = a.createErrorStatus(description)
	output["entity"] = entity
	err := ctx.SetOutput("connectorMsg", output)
	if err != nil {
		ctx.Logger().Error("Failed to set output oject ", err.Error())
		return false
	}
	return true
}

// Eval evaluates the activity
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	logger := ctx.Logger()
	logger.Info("aeroqualaqy1:Eval enter")

	output := map[string]interface{}{}
	host := a.settings.Host
	port := a.settings.Port
	username := a.settings.Username
	password := a.settings.Password
	instrument := a.settings.Instrument
	mappings := a.Mappings
	entity := a.settings.Entity

	url := "http://" + host + ":" + port + "/api/account/login"
	pld := "UserName=" + username + "&" + "password=" + password
	payload := strings.NewReader(pld)
	// Convert payload to a Reader

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logger.Error("Error:", err.Error())
		return false, err
	}

	client := http.Client{Timeout: 5 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		description := fmt.Sprintf("Connection to %s refused", url)
		rc := a.sendErrorMessage(ctx, entity, description)
		logger.Info("aeroqualaqy1:Eval exit")
		return rc, nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		description := fmt.Sprintf("Connection to %s returned status code %d", url, res.StatusCode)
		rc := a.sendErrorMessage(ctx, entity, description)
		logger.Info("aeroqualaqy1:Eval exit")
		return rc, nil
	}

	url = "http://" + host + ":" + port + "/api/data/" + instrument

	reqGet, _ := http.NewRequest("GET", url, nil)

	// Build the query
	q := reqGet.URL.Query()
	q.Add("duration", "20")
	q.Add("averagingperiod", "10")
	q.Add("includejournal", "false")
	reqGet.URL.RawQuery = q.Encode()

	// Transfer the cookie from the response to the request.
	for _, cookie := range res.Cookies() {
		httpCookie := http.Cookie{Name: cookie.Name, Value: cookie.Value}
		reqGet.AddCookie(&httpCookie)
	}

	resGet, err := client.Do(reqGet)
	if err != nil {
		description := fmt.Sprintf("Connection to %s refused", url)
		rc := a.sendErrorMessage(ctx, entity, description)
		logger.Info("aeroqualaqy1:Eval exit")
		return rc, nil
	}

	body, err := ioutil.ReadAll(resGet.Body)
	if err != nil {
		logger.Error("Error:", err.Error())
		description := fmt.Sprintf("ioutil.ReadAll failed")
		rc := a.sendErrorMessage(ctx, entity, description)
		logger.Info("aeroqualaqy1:Eval exit")
		return rc, nil
	}
	fmt.Println(string(body))
	message := parse(body, mappings)
	if message == nil {
		logger.Error("Message error. Parse failed.")
		return false, nil
	}
	if a.status == "ERROR" {
		message["status"] = "RUNNING"
	}
	a.status = "RUNNING"

	output["data"] = message
	output["entity"] = entity

	err = ctx.SetOutput("connectorMsg", output)
	if err != nil {
		logger.Error("Failed to set output oject ", err.Error())
		return false, err
	}
	logger.Info("aeroqualaqy1:Eval exit")
	return true, nil
}

func parse(body []byte, mappings map[string]map[string]interface{}) map[string]interface{} {

	var e interface{}
	err := json.Unmarshal([]byte(body), &e)
	if err != nil {
		panic(err)
	}
	m1 := e.(map[string]interface{})
	d1 := m1["data"].([]interface{})
	// The first element of the array contains the latest data
	if d1 == nil {
		return nil
	}
	if d1[0] == nil {
		return nil
	}
	d2 := d1[0].(map[string]interface{})
	datetime := connector.FormatESPTime(d2["Time"].(string))

	values := make([]map[string]interface{}, 0, 10)
	message := map[string]interface{}{}

	// Loop over the data, convert the field names, add the amounts while
	// adding to the connector message.
	for k, v := range d2 {
		value := map[string]interface{}{}
		field, found := mappings[k]["field"]
		if !found {
			continue
		}
		if field == "" {
			continue
		}
		value["field"] = field
		value["amount"] = v.(float64)
		values = append(values, value)
	}
	message["values"] = values
	message["datetime"] = datetime

	return message
}
