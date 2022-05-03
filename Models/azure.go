package models

import "time"

type Text struct {
	Text string `json:"text"`
}

type ID struct {
	ID string `json:"id"`
}

type Drop struct {
	Location    string `json:"location"`
	Type        string `json:"type"`
	Url         string `json:"url"`
	DownloadUrl string `json:"downloadUrl"`
}

type Log struct {
	Type        string `json:"type"`
	Url         string `json:"url"`
	DownloadUrl string `json:"downloadUrl"`
}

type LastChangedBy struct {
	DisplayName string `json:"displayName"`
	Url         string `json:"url"`
	ID          string `json:"id"`
	UniqueName  string `json:"uniqueName"`
	ImageUrl    string `json:"imageUrl"`
}

type Definition struct {
	BatchSize      int8   `json:"batchSize"`
	TriggerType    string `json:"triggerType"`
	DefinitionType string `json:"definitionType"`
	ID             int8   `json:"id"`
	Name           string `json:"name"`
	Url            string `json:"url"`
}

type Queue struct {
	QueueType string `json:"queueType"`
	ID        int8   `json:"id"`
	Name      string `json:"name"`
	Url       string `json:"url"`
}

type RequestedFor struct {
	DisplayName string `json:"displayName"`
	Url         string `json:"url"`
	ID          string `json:"id"`
	UniqueName  string `json:"uniqueName"`
	ImageUrl    string `json:"imageUrl"`
}

type Requests struct {
	ID           int8         `json:"id"`
	Url          string       `json:"url"`
	RequestedFor RequestedFor `json:"requestedFor"`
}

type Resource struct {
	Uri                string        `json:"uri"`
	ID                 int8          `json:"id"`
	BuildNumber        string        `json:"buildNumber"`
	Url                string        `json:"url"`
	StartTime          time.Time     `json:"startTime"`
	FinishTime         time.Time     `json:"finishTime"`
	Reason             string        `json:"reason"`
	Status             string        `json:"status"`
	DropLocation       string        `json:"dropLocation"`
	Drop               Drop          `json:"drop"`
	Log                Log           `json:"log"`
	SourceGetVersion   string        `json:"sourceGetVersion"`
	LastChangedBy      LastChangedBy `json:"lastChangedBy"`
	RetainIndefinitely bool          `json:"retainIndefinitely"`
	HasDiagnostics     bool          `json:"hasDiagnostics"`
	Definition         Definition    `json:"definition"`
	Queue              Queue         `json:"queue"`
	Requests           []Requests    `json:"requests"`
}

type ResourceContainers struct {
	Collection ID `json:"collection"`
	Account    ID `json:"account"`
	Project    ID `json:"project"`
}

type AzureRequest struct {
	SubscriptionId     string             `json:"subscriptionId"`
	NotificationId     int8               `json:"notificationId"`
	ID                 string             `json:"id"`
	EventType          string             `json:"eventType"`
	PublisherId        string             `json:"publisherId"`
	Message            Text               `json:"message"`
	DetailedMessage    Text               `json:"detailedMessage"`
	Resource           Resource           `json:"resource"`
	ResourceVersion    string             `json:"resourceVersion"`
	ResourceContainers ResourceContainers `json:"resourceContainers"`
	CreatedDate        time.Time          `json:"createdDate"`
}
