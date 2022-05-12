package aosSdk

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	apiUrlTasksPrefix = apiUrlBlueprintsPrefix
	apiUrlTasksSuffix = "/tasks/"

	taskMonFirstCheckDelay = 100 * time.Millisecond
	taskMonPollInterval    = 500 * time.Millisecond

	taskStatusOngoing = "in_progress"
	taskStatusInit    = "init"
	taskStatusFail    = "failed"
	taskStatusSuccess = "succeeded"
	taskStatusTimeout = "timeout"
)

// getAllTasksResponse is sent by Apstra in response to GET at
// 'apiUrlTaskPrefix + blueprintId + apiUrlTaskSuffix'
type getAllTasksResponse struct {
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
		Id                  TaskId `json:"id"`
	} `json:"items"`
}

// getTaskResponse is sent by Apstra in response to GET at
// 'apiUrlTaskPrefix + blueprintId + apiUrlTaskSuffix + taskId'
type getTaskResponse struct {
	Status      string `json:"status"`
	BeginAt     string `json:"begin_at"`
	RequestData struct {
		Url     string `json:"url"`
		Headers struct {
			ContentLength   string `json:"Content-Length"`
			AcceptEncoding  string `json:"Accept-Encoding"`
			XForwardedProto string `json:"X-Forwarded-Proto"`
			XForwardedFor   string `json:"X-Forwarded-For"`
			Connection      string `json:"Connection"`
			XUser           string `json:"X-User"`
			Accept          string `json:"Accept"`
			UserAgent       string `json:"User-Agent"`
			Host            string `json:"Host"`
			XUserId         string `json:"X-User-Id"`
			XRealIp         string `json:"X-Real-Ip"`
			ContentType     string `json:"Content-Type"`
		} `json:"headers"`
		Args struct {
			Async string `json:"async"`
		} `json:"args"`
		Data struct {
			SzType   string `json:"sz_type"`
			VrfName  string `json:"vrf_name"`
			RtPolicy struct {
			} `json:"rt_policy"`
			Label string `json:"label"`
		} `json:"data"`
		Method string `json:"method"`
	} `json:"request_data"`
	UserId         string `json:"user_id"`
	LastUpdatedAt  string `json:"last_updated_at"`
	UserName       string `json:"user_name"`
	CreatedAt      string `json:"created_at"`
	DetailedStatus struct {
		ApiResponse            json.RawMessage `json:"api_response"`
		ConfigBlueprintVersion int             `json:"config_blueprint_version"`
	} `json:"detailed_status"`
	ConfigLastUpdatedAt string `json:"config_last_updated_at"`
	UserIp              string `json:"user_ip"`
	Type                string `json:"type"`
	Id                  TaskId `json:"id"`
}

// taskMonitorMonReq uniquely identifies an Apstra task which can be tracked
// using the // /api/blueprint/<id>/tasks and /api/blueprint/<id>/tasks/<id> API
// endpoints.
// This structure is submitted by a caller using taskMonitor.taskInChan. When
// the task is no longer outstanding (success, timeout, failed), taskMonitor
// responds via responseChan with the complete getTaskResponse structure received
// from Apstra (or any errors encountered along the way)
type taskMonitorMonReq struct {
	bluePrintId  ObjectId                 // task API calls must reference a blueprint
	taskId       TaskId                   // tracks the task
	responseChan chan<- *taskCompleteInfo // talk here when the task is complete
}

// taskCompleteInfo is generated by taskMonitor when a TaskId exits pending modes
// ("init" or "in_progress"). The id field represents the object created by
// the task, and result can be any status representing a completed task
// ("succeeded", "failed", "timeout")
type taskCompleteInfo struct {
	status *getTaskResponse
	err    error
}

// a taskMonitor runs as an independent goroutine, accepts task{}s to monitor
// via taskInChan, closes the task's `done` channel when it detects apstra
// has completed the task.
type taskMonitor struct {
	client               *Client                             // for making Apstra API calls
	mapBpIdToSliceTaskId map[ObjectId][]TaskId               // track tasks by blueprint
	mapTaskIdToChan      map[TaskId]chan<- *taskCompleteInfo // track reply channels by task
	taskInChan           <-chan *taskMonitorMonReq           // for learning about new tasks
	timer                *time.Timer                         // triggers check()
	errChan              chan<- error                        // error feedback to main loop
	lock                 sync.Mutex                          // control access to mapBpIdToTask
	tmQuit               <-chan struct{}
}

// newTaskMonitor creates a new taskMonitor, but does not start it.
func newTaskMonitor(c *Client) *taskMonitor {
	monitor := taskMonitor{
		timer:                time.NewTimer(0),
		client:               c,
		taskInChan:           c.taskMonChan,
		errChan:              c.cfg.errChan,
		mapBpIdToSliceTaskId: make(map[ObjectId][]TaskId),
		mapTaskIdToChan:      make(map[TaskId]chan<- *taskCompleteInfo),
	}
	<-monitor.timer.C
	return &monitor
}

// start causes the taskMonitor to run
func (o *taskMonitor) start(quit <-chan struct{}) {
	o.tmQuit = quit
	go o.run()
}

// run is the main taskMonitor loop
func (o *taskMonitor) run() {
	// main loop
	for {
		select {
		// timer event
		case <-o.timer.C:
			go o.check()
		case in := <-o.taskInChan: // new task event
			o.timer.Stop() // timer may be about to fire, but we're already running
			o.lock.Lock()
			// the task referenced by 'in' better not be one we're already tracking
			// big assumption here: task IDs are globally unique
			if _, found := o.mapTaskIdToChan[in.taskId]; found {
				in.responseChan <- &taskCompleteInfo{
					status: nil,
					err:    fmt.Errorf("task ID %s already being tracked", in.taskId),
				}
				close(in.responseChan)
				o.lock.Unlock()
				o.timer.Reset(taskMonFirstCheckDelay)
				continue
			} else {
				o.mapTaskIdToChan[in.taskId] = in.responseChan
			}

			// the blueprint referenced by 'in' may already exist
			if _, found := o.mapBpIdToSliceTaskId[in.bluePrintId]; found {
				// previously known blueprint, append new task to the existing slice
				o.mapBpIdToSliceTaskId[in.bluePrintId] = append(o.mapBpIdToSliceTaskId[in.bluePrintId], in.taskId)
			} else {
				// new blueprint, create the task slice with this ID as the only element
				o.mapBpIdToSliceTaskId[in.bluePrintId] = []TaskId{in.taskId}
			}
			o.lock.Unlock()
			o.timer.Reset(taskMonFirstCheckDelay)
		case <-o.tmQuit: // program exit
			return
		}
	}
}

func (o *taskMonitor) checkBlueprints() {
	// loop over blueprints known to have outstanding tasks
BlueprintLoop:
	for bpId, taskIdList := range o.mapBpIdToSliceTaskId {
		// get task result info from Apstra
		taskIdToStatus, err := o.client.getBlueprintTasksStatus(bpId, taskIdList)
		if err != nil {
			err = fmt.Errorf("error getting tasks for blueprint %s - %w", string(bpId), err)
			// todo: not happy with this error handling
			if o.errChan != nil {
				o.errChan <- err
			} else {
				log.Println(err)
			}
			continue BlueprintLoop
		}
		o.checkTasksInBlueprint(bpId, taskIdToStatus)
	} // BlueprintLoop

	// after processing all monitored tasks, one or more
	// per-blueprint lists might be empty. Delete them.
	for bpId, taskList := range o.mapBpIdToSliceTaskId {
		if len(taskList) == 0 {
			delete(o.mapBpIdToSliceTaskId, bpId)
		}
	}
}

// checkTasksInBlueprint takes a blueprint ID and a [TaskId]status (string)
// fetched using the task summary API endpoint. We'll ask Apstra for the
// detailed API output for any tasks no longer in progress, and return that via
// the caller specified channel.
func (o *taskMonitor) checkTasksInBlueprint(bpId ObjectId, mapTaskIdToStatus map[TaskId]string) {
TaskLoop:
	// loop over *all* outstanding tasks associated with this blueprint
	for i, taskId := range o.mapBpIdToSliceTaskId[bpId] {
		// make sure Apstra response (input to this function) includes our taskId
		if _, found := mapTaskIdToStatus[taskId]; !found {
			// Apstra response doesn't have our task Id:
			//   error, delete it from the slice, next task
			o.mapTaskIdToChan[taskId] <- &taskCompleteInfo{
				err: fmt.Errorf("blueprint '%s' task '%s' unknown to Apstra server", bpId, taskId),
			}
			o.mapBpIdToSliceTaskId[bpId] = append(o.mapBpIdToSliceTaskId[bpId][:i], o.mapBpIdToSliceTaskId[bpId][i+1:]...)
			continue TaskLoop
		}

		// What did Apstra say about our taskId?
		switch mapTaskIdToStatus[taskId] {
		case taskStatusInit:
			continue TaskLoop // still working; skip to next task
		case taskStatusOngoing:
			continue TaskLoop // still working; skip to next task
		case taskStatusFail: // done; fallthrough
		case taskStatusSuccess: // done; fallthrough
		case taskStatusTimeout: // done; fallthrough
		default: // something unexpected
			// Unexpected task status response from Apstra:
			//   error, delete it from the slice, next task
			o.mapTaskIdToChan[taskId] <- &taskCompleteInfo{
				err: fmt.Errorf("blueprint '%s' task '%s' status unexpected: %s", bpId, taskId, mapTaskIdToStatus[taskId]),
			}
			o.mapBpIdToSliceTaskId[bpId] = append(o.mapBpIdToSliceTaskId[bpId][:i], o.mapBpIdToSliceTaskId[bpId][i+1:]...)
			continue TaskLoop
		}

		// if we got here, we're able to return a recognized and conclusive
		// task status result to the caller. Fetch the full details from Apstra.
		taskInfo, err := o.client.getBlueprintTaskStatusById(bpId, taskId)
		o.mapTaskIdToChan[taskId] <- &taskCompleteInfo{
			status: taskInfo,
			err:    err,
		}

		// remove this task from the list of monitored tasks
		o.mapBpIdToSliceTaskId[bpId] = append(o.mapBpIdToSliceTaskId[bpId][:i], o.mapBpIdToSliceTaskId[bpId][i+1:]...)
	} // TaskLoop
}

// check
//   invokes checkBlueprints
//             invokes checkTasksInBlueprint
//   resets timer (maybe)
func (o *taskMonitor) check() {
	//todo blueprint lock should wrap this whole function?
	o.lock.Lock()
	o.checkBlueprints()
	if len(o.mapBpIdToSliceTaskId) > 0 {
		o.timer.Reset(taskMonPollInterval)
	}
	o.lock.Unlock()
}

// taskListToFilterExpr returns the supplied []ObjectId as a string prepped for
// Apstra's API response filter. Something like: "id in ['abc','def']"
func taskListToFilterExpr(in []TaskId) string {
	var quotedList []string
	for i := range in {
		quotedList = append(quotedList, "'"+string(in[i])+"'")
	}
	return "id in [" + strings.Join(quotedList, ",") + "]"
}

// getBlueprintTasksStatus returns a map of TaskId to status (strings like
// "succeeded", "init", etc...)
func (o Client) getBlueprintTasksStatus(bpid ObjectId, taskIdList []TaskId) (map[TaskId]string, error) {
	aosUrl, err := url.Parse(apiUrlTasksPrefix + string(bpid) + apiUrlTasksSuffix)
	aosUrl.RawQuery = url.Values{"filter": []string{taskListToFilterExpr(taskIdList)}}.Encode()
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w",
			apiUrlTasksPrefix+string(bpid)+apiUrlTasksSuffix, err)
	}
	response := &getAllTasksResponse{}
	err = o.talkToAos(&talkToAosIn{
		method:      httpMethodGet,
		url:         aosUrl,
		apiInput:    nil,
		apiResponse: response,
		doNotLogin:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting getAllTasksResponse for blueprint '%s' - %w", bpid, err)
	}

	// todo: compare the submitted taskIdList against the server's response.
	//  what if we're trying to track tasks the server doesn't know about?

	result := make(map[TaskId]string)
	for _, i := range response.Items {
		if i.Status == "" {
			return nil, fmt.Errorf("server resopnse included empty task status")
		}
		if i.Id == "" {
			return nil, fmt.Errorf("server resopnse included empty task id")
		}
		result[i.Id] = i.Status
	}
	return result, nil
}

func (o Client) getBlueprintTaskStatusById(bpid ObjectId, tid TaskId) (*getTaskResponse, error) {
	aosUrl, err := url.Parse(apiUrlTasksPrefix + string(bpid) + apiUrlTasksSuffix + string(tid))
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w",
			apiUrlTasksPrefix+string(bpid)+apiUrlTasksSuffix+string(tid), err)
	}
	result := &getTaskResponse{}
	return result, o.talkToAos(&talkToAosIn{
		method:      httpMethodGet,
		url:         aosUrl,
		apiInput:    nil,
		apiResponse: result,
		doNotLogin:  false,
	})
}

func blueprintIdFromUrl(in *url.URL) (ObjectId, error) {
	split1 := strings.Split(in.String(), apiUrlBlueprintsPrefix)
	if len(split1) != 2 {
		return "", fmt.Errorf("failed to extract blueprint ID from URL '%s' at step 1", in.String())
	}

	split2 := strings.Split(split1[1], apiUrlPathDelim)
	if len(split1) == 0 {
		return "", fmt.Errorf("failed to extract blueprint ID from URL '%s' at step 2", in.String())
	}

	return ObjectId(split2[0]), nil
}

// waitForTaskCompletion interacts with the taskMonitor, returns the Apstra API
// *getTaskResponse
func waitForTaskCompletion(bId ObjectId, tId TaskId, mon chan *taskMonitorMonReq) (*getTaskResponse, error) {
	// task status update channel (how we'll learn the task is complete
	reply := make(chan *taskCompleteInfo, 1) // Task Complete Info Channel

	// submit our task to the task monitor
	mon <- &taskMonitorMonReq{
		bluePrintId:  bId,
		taskId:       tId,
		responseChan: reply,
	}

	tci := <-reply
	return tci.status, tci.err

}
