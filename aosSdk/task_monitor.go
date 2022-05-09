package aosSdk

const (
	apiUrlTaskPrefix = apiUrlBlueprints + "/"
	apiUrlTaskSuffix = "/tasks/"
)

type allTasksResponse struct {
	Items []struct {
		Status      string `json:"status"`
		BeginAt     string `json:"begin_at"`
		RequestData struct {
			Url    string `json:"url"`
			Method string `json:"method"`
		} `json:"request_data"`
		UserId              string `json:"user_id"`
		LastUpdatedAt       string `json:"last_updated_at"`
		UserName            string `json:"user_name"`
		CreatedAt           string `json:"created_at"`
		ConfigLastUpdatedAt string `json:"config_last_updated_at"`
		UserIp              string `json:"user_ip"`
		Type                string `json:"type"`
		Id                  string `json:"id"`
	} `json:"items"`
}

type GetTaskResponse struct {
	Status      string `json:"status"`
	BeginAt     string `json:"begin_at"`
	RequestData struct {
		Url     string `json:"url"`
		Headers struct {
			Origin          string `json:"Origin"`
			ContentLength   string `json:"Content-Length"`
			Host            string `json:"Host"`
			AcceptLanguage  string `json:"Accept-Language"`
			AcceptEncoding  string `json:"Accept-Encoding"`
			XForwardedProto string `json:"X-Forwarded-Proto"`
			SecFetchSite    string `json:"Sec-Fetch-Site"`
			XForwardedFor   string `json:"X-Forwarded-For"`
			XUser           string `json:"X-User"`
			SecFetchMode    string `json:"Sec-Fetch-Mode"`
			UserAgent       string `json:"User-Agent"`
			Connection      string `json:"Connection"`
			XUserId         string `json:"X-User-Id"`
			Referer         string `json:"Referer"`
			Accept          string `json:"Accept"`
			SecChUaPlatform string `json:"Sec-Ch-Ua-Platform"`
			SecChUaMobile   string `json:"Sec-Ch-Ua-Mobile"`
			XRealIp         string `json:"X-Real-Ip"`
			ContentType     string `json:"Content-Type"`
			SecChUa         string `json:"Sec-Ch-Ua"`
			SecFetchDest    string `json:"Sec-Fetch-Dest"`
		} `json:"headers"`
		Args struct {
		} `json:"args"`
		Data struct {
			PoolIds []string `json:"pool_ids"`
		} `json:"data"`
		Method string `json:"method"`
	} `json:"request_data"`
	UserId         string `json:"user_id"`
	LastUpdatedAt  string `json:"last_updated_at"`
	UserName       string `json:"user_name"`
	CreatedAt      string `json:"created_at"`
	DetailedStatus struct {
		ApiResponse            string `json:"api_response"`
		ConfigBlueprintVersion int    `json:"config_blueprint_version"`
	} `json:"detailed_status"`
	ConfigLastUpdatedAt string `json:"config_last_updated_at"`
	UserIp              string `json:"user_ip"`
	Type                string `json:"type"`
	Id                  string `json:"id"`
}

// a task uniquely identifies an Apstra task which can be tracked using the
// /api/blueprint/<id>/tasks and /api/blueprint/<id>/tasks/<id> API endpoints.
type task struct {
	blueprint_id string
	task_id      string
	done         chan struct{}
}

// a taskMonitor runs as an independent goroutine, accepts task{}s to monitor
// via taskInChan, closes the task's `done` channel when it detects apstra
// has completed the task.
type taskMonitor struct {
	taskInChan         chan task
	mapTaskToChan      map[string]chan struct{}
	mapTaskToBlueprint map[string]string
	quitChan           chan struct{}
}

func newTaskMonitor() *taskMonitor {
	monitor := taskMonitor{
		taskInChan:         make(chan task),
		mapTaskToChan:      make(map[string]chan struct{}),
		mapTaskToBlueprint: make(map[string]string),
	}
	_ = monitor
	return &monitor
}

func (o *taskMonitor) run() {
	for {
		select {
		case newTask := <-o.taskInChan:
			_ = newTask
		case <-o.quitChan:
			return
		}
	}
}

//func (o Client) GetTaskByBlueprintIdAndTaskId(blueprintId string, taskId string) (*GetTaskResponse, error) {
//	response := GetTaskResponse{}
//	_, err := o.talkToAos(&talkToAosIn{
//		method:        httpMethodGet,
//		url:           apiUrlTaskPrefix + blueprintId + apiUrlTaskSuffix + taskId,
//		fromServerPtr: &response,
//	})
//	return &response, err
//}
