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

type CreatedBy struct {
	DisplayName string `json:"displayName"`
	Url         string `json:"url"`
	ID          string `json:"id"`
	UniqueName  string `json:"uniqueName"`
	ImageUrl    string `json:"imageUrl"`
}

type Definition struct {
	BatchSize      int32   `json:"batchSize"`
	TriggerType    string  `json:"triggerType"`
	DefinitionType string  `json:"definitionType"`
	ID             int32   `json:"id"`
	Name           string  `json:"name"`
	Url            string  `json:"url"`
	Uri            string  `json:"uri"`
	Type           string  `json:"type"`
	QueueStatus    string  `json:"queueStatus"`
	Revision       int32   `json:"revision"`
	Project        Project `json:"project"`
}

type Queue struct {
	QueueType string `json:"queueType"`
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Url       string `json:"url"`
}

type RequestedFor struct {
	DisplayName string `json:"displayName"`
	Url         string `json:"url"`
	ID          string `json:"id"`
	UniqueName  string `json:"uniqueName"`
	ImageUrl    string `json:"imageUrl"`
	Links       Links  `json:"_links"`
}

type Requests struct {
	ID           int32        `json:"id"`
	Url          string       `json:"url"`
	RequestedFor RequestedFor `json:"requestedFor"`
}

type Project struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Url            string    `json:"url"`
	State          string    `json:"state"`
	Revision       int32     `json:"revision"`
	Visibility     string    `json:"visibility"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
}

type Repository struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Url             string  `json:"url"`
	DefaultBranch   string  `json:"defaultBranch"`
	RemoteUrl       string  `json:"remoteUrl"`
	Project         Project `json:"project"`
	Size            int32   `json:"size"`
	SshUrl          string  `json:"sshUrl"`
	WebUrl          string  `json:"webUrl"`
	IsDisabled      bool    `json:"isDisabled"`
	IsInMaintenance bool    `json:"isInMaintenance"`
}

type Reviewers struct {
	ReviewerUrl string `json:"reviewerUrl"`
	Vote        int32  `json:"vote"` //10 approved | 5 approved with suggestions | 0 no vote | -5 waiting for author | -10 rejected
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	UniqueName  string `json:"uniqueName"`
	Url         string `json:"url"`
	ImageUrl    string `json:"imageUrl"`
	IsContainer bool   `json:"isContainer"`
}

type UserAuthor struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type MergeCommit struct {
	CommitID  string     `json:"commitId"`
	Url       string     `json:"url"`
	Author    UserAuthor `json:"author"`
	Committer UserAuthor `json:"committer"`
	Comment   string     `json:"comment"`
}

type LinkHref struct {
	Href string `json:"href"`
}

type Links struct {
	Self                    LinkHref `json:"self"`
	Web                     LinkHref `json:"web"`
	SourceVersionDisplayUri LinkHref `json:"sourceVersionDisplayUri"`
	Timeline                LinkHref `json:"timeline"`
	Badge                   LinkHref `json:"badge"`
	Avatar                  LinkHref `json:"avatar"`
}

type TriggerInfo struct {
	SourceBranch      string `json:"ci.sourceBranch"`
	SourceSha         string `json:"ci.sourceSha"`
	Message           string `json:"ci.message"`
	TriggerRepository string `json:"ci.triggerRepository"`
}

type Resource struct {
	Uri                   string        `json:"uri"`
	ID                    int32         `json:"id"`
	BuildNumber           string        `json:"buildNumber"`
	Url                   string        `json:"url"`
	Title                 string        `json:"title"`
	StartTime             time.Time     `json:"startTime"`
	FinishTime            time.Time     `json:"finishTime"`
	Reason                string        `json:"reason"`
	Status                string        `json:"status"`
	DropLocation          string        `json:"dropLocation"`
	Drop                  Drop          `json:"drop"`
	Log                   Log           `json:"log"`
	SourceGetVersion      string        `json:"sourceGetVersion"`
	CreatedBy             CreatedBy     `json:"createdBy"`
	LastChangedBy         LastChangedBy `json:"lastChangedBy"`
	RetainIndefinitely    bool          `json:"retainIndefinitely"`
	HasDiagnostics        bool          `json:"hasDiagnostics"`
	Definition            Definition    `json:"definition"`
	Queue                 Queue         `json:"queue"`
	Requests              []Requests    `json:"requests"`
	Reviewers             []Reviewers   `json:"reviewers"`
	Repository            Repository    `json:"repository"`
	CodeReviewId          int32         `json:"codeReviewId"`
	PullRequestId         int32         `json:"pullRequestId"`
	CreationDate          time.Time     `json:"creationDate"`
	Description           string        `json:"description"`
	SourceRefName         string        `json:"sourceRefName"`
	TargetRefName         string        `json:"targetRefName"`
	MergeStatus           string        `json:"mergeStatus"`
	IsDraft               bool          `json:"isDraft"`
	MergeID               string        `json:"mergeId"`
	LastMergeSourceCommit MergeCommit   `json:"lastMergeSourceCommit"`
	LastMergeTargetCommit MergeCommit   `json:"lastMergeTargetCommit"`
	LastMergeCommit       MergeCommit   `json:"lastMergeCommit"`
	SupportsIterations    bool          `json:"supportsIterations"`
	ArtifactId            string        `json:"ArtifactId"`
	Links                 Links         `json:"_links"`
	TriggerInfo           TriggerInfo   `json:"triggerInfo"`
	Result                string        `json:"result"`
	QueueTime             time.Time     `json:"queueTime"`
	Project               Project       `json:"project"`
	RequestedFor          RequestedFor  `json:"requestedFor"`
}

type ResourceContainers struct {
	Collection ID `json:"collection"`
	Account    ID `json:"account"`
	Project    ID `json:"project"`
}

type AzureRequest struct {
	SubscriptionId     string             `json:"subscriptionId"`
	NotificationId     int32              `json:"notificationId"`
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

type Message struct {
	Text     string `json:"text"`
	Html     string `json:"html"`
	Markdown string `json:"markdown"`
}

type AzurePipeline struct {
	SubscriptionId     string             `json:"subscriptionId"`
	NotificationId     int32              `json:"notificationId"`
	ID                 string             `json:"id"`
	EventType          string             `json:"eventType"`
	PublisherId        string             `json:"publisherId"`
	Message            Message            `json:"message"`
	DetailedMessage    Message            `json:"detailedMessage"`
	Resource           Resource           `json:"resource"`
	ResourceVersion    string             `json:"resourceVersion"`
	ResourceContainers ResourceContainers `json:"resourceContainers"`
	CreatedDate        time.Time          `json:"createdDate"`
}
