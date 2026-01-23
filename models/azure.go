package models

import "time"

// Text represents a simple text message structure from Azure DevOps webhooks.
type Text struct {
	Text string `json:"text"`
}

// ID represents a simple identifier structure used in Azure DevOps containers.
type ID struct {
	ID string `json:"id"`
}

// Drop represents artifact drop information from Azure DevOps builds.
type Drop struct {
	Location    string `json:"location"`
	Type        string `json:"type"`
	Url         string `json:"url"`
	DownloadUrl string `json:"downloadUrl"`
}

// Log represents log file information from Azure DevOps builds.
type Log struct {
	Type        string `json:"type"`
	Url         string `json:"url"`
	DownloadUrl string `json:"downloadUrl"`
}

// LastChangedBy represents user information for the last person who modified a resource.
type LastChangedBy struct {
	DisplayName string `json:"displayName"`
	Url         string `json:"url"`
	ID          string `json:"id"`
	UniqueName  string `json:"uniqueName"`
	ImageUrl    string `json:"imageUrl"`
}

// CreatedBy represents user information for the person who created a resource.
type CreatedBy struct {
	DisplayName string `json:"displayName"`
	Url         string `json:"url"`
	ID          string `json:"id"`
	UniqueName  string `json:"uniqueName"`
	ImageUrl    string `json:"imageUrl"`
}

// Definition represents a pipeline or build definition in Azure DevOps.
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

// Queue represents an agent queue used for running builds or pipelines.
type Queue struct {
	QueueType string `json:"queueType"`
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Url       string `json:"url"`
}

// RequestedFor represents user information for the person who requested a build or deployment.
type RequestedFor struct {
	DisplayName string `json:"displayName"`
	Url         string `json:"url"`
	ID          string `json:"id"`
	UniqueName  string `json:"uniqueName"`
	ImageUrl    string `json:"imageUrl"`
	Links       Links  `json:"_links"`
}

// Requests represents a build request in Azure DevOps.
type Requests struct {
	ID           int32        `json:"id"`
	Url          string       `json:"url"`
	RequestedFor RequestedFor `json:"requestedFor"`
}

// Project represents an Azure DevOps project.
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

// Repository represents a Git repository in Azure DevOps.
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

// Reviewers represents a pull request reviewer with their vote status.
//
// Vote values:
//   - 10: Approved
//   - 5: Approved with suggestions
//   - 0: No vote
//   - -5: Waiting for author
//   - -10: Rejected
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

// UserAuthor represents Git commit author or committer information.
type UserAuthor struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

// MergeCommit represents a merge commit in a pull request.
type MergeCommit struct {
	CommitID  string     `json:"commitId"`
	Url       string     `json:"url"`
	Author    UserAuthor `json:"author"`
	Committer UserAuthor `json:"committer"`
	Comment   string     `json:"comment"`
}

// LinkHref represents a single link with an href URL.
type LinkHref struct {
	Href string `json:"href"`
}

// Links contains various URLs related to a resource in Azure DevOps.
type Links struct {
	Self                    LinkHref `json:"self"`
	Web                     LinkHref `json:"web"`
	SourceVersionDisplayUri LinkHref `json:"sourceVersionDisplayUri"`
	Timeline                LinkHref `json:"timeline"`
	Badge                   LinkHref `json:"badge"`
	Avatar                  LinkHref `json:"avatar"`
}

// TriggerInfo contains information about what triggered a pipeline build.
type TriggerInfo struct {
	SourceBranch      string `json:"ci.sourceBranch"`
	SourceSha         string `json:"ci.sourceSha"`
	Message           string `json:"ci.message"`
	TriggerRepository string `json:"ci.triggerRepository"`
}

// Release represents a release in Azure DevOps.
type Release struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Url   string `json:"url"`
	Links Links  `json:"_links"`
}

// Deployment represents deployment information for a release.
//
// DeploymentStatus values include: "succeeded", "failed", "stopped", etc.
type Deployment struct {
	ID               int32        `json:"id"`
	Reason           string       `json:"reason"`
	DeploymentStatus string       `json:"deploymentStatus"`
	OperationStatus  string       `json:"operationStatus"`
	RequestedBy      RequestedFor `json:"requestedBy"`
	RequestedFor     RequestedFor `json:"requestedFor"`
	QueuedOn         time.Time    `json:"queuedOn"`
	StartedOn        time.Time    `json:"startedOn"`
	CompletedOn      time.Time    `json:"completedOn"`
}

// Environment represents a deployment environment in a release pipeline.
type Environment struct {
	ID                int32        `json:"id"`
	ReleaseID         int32        `json:"releaseId"`
	Name              string       `json:"name"`
	Status            string       `json:"status"`
	CreatedOn         time.Time    `json:"createdOn"`
	ModifiedOn        time.Time    `json:"modifiedOn"`
	Release           Release      `json:"release"`
	ReleaseDefinition Release      `json:"releaseDefinition"`
	ReleaseCreatedBy  RequestedFor `json:"releaseCreatedBy"`
	TriggerReason     string       `json:"triggerReason"`
	TimeToDeploy      float64      `json:"timeToDeploy"`
}

// Resource represents the main resource data in an Azure DevOps webhook event.
//
// This struct is polymorphic and contains fields for different types of resources:
// pipelines, pull requests, and releases. Not all fields will be populated
// depending on the event type.
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
	Environment           Environment   `json:"environment"`
	Deployment            Deployment    `json:"deployment"`
	StageName             string        `json:"stageName"`
	AttemptId             int16         `json:"attemptId"`
}

// ResourceContainers contains identifiers for the Azure DevOps containers
// (collection, account, project) associated with the webhook event.
type ResourceContainers struct {
	Collection ID `json:"collection"`
	Account    ID `json:"account"`
	Project    ID `json:"project"`
}

// AzureRequest represents the complete webhook payload structure from Azure DevOps.
//
// This is the main structure used to deserialize Azure DevOps webhook events
// for pipelines, pull requests, and releases. The Resource field contains
// the specific event data, which varies based on the EventType.
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

// Message represents a formatted message with text, HTML, and Markdown versions.
type Message struct {
	Text     string `json:"text"`
	Html     string `json:"html"`
	Markdown string `json:"markdown"`
}

// AzurePipeline represents an alternative pipeline webhook payload structure
// that uses Message instead of Text for message fields.
//
// This structure may be used for certain Azure DevOps pipeline events
// that provide richer message formatting options.
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

// AzureRepository represents a Git repository returned from the Azure DevOps REST API.
//
// This structure is used when fetching repository details via the API,
// such as when retrieving repository information for pipeline notifications.
type AzureRepository struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	URL             string  `json:"url"`
	Project         Project `json:"project"`
	DefaultBranch   string  `json:"defaultBranch"`
	Size            int64   `json:"size"`
	RemoteURL       string  `json:"remoteUrl"`
	SSHURL          string  `json:"sshUrl"`
	WebURL          string  `json:"webUrl"`
	Links           Links   `json:"_links"`
	IsDisabled      bool    `json:"isDisabled"`
	IsInMaintenance bool    `json:"isInMaintenance"`
}
