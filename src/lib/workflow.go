package lib

import (
	"Yearning-go/src/model"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/cookieY/yee/logger"
	"io/ioutil"
	"net/http"
	"strings"
)

//{
//    request_id: "xxxxxx",
//    result_code: "xxx",
//    data: {
//    }
//}

const (
	WORKFLOW_API string = "https://rocket.nioint.com"
)

type WorkflowResponse struct {
	RequestID  string      `json:"request_id"`
	ResultCode string      `json:"result_code"`
	Data       interface{} `json:"data"`
}

type WorkflowRequest struct {
	Creator  string            `json:"creator"`
	FlowCode string            `json:"flow_code"`
	Context  map[string]string `json:"context"`
}

func NewWorkflowRequest() *WorkflowRequest {
	return &WorkflowRequest{
		Creator:  "zhendong.pan",
		FlowCode: "22792e0b-5aa7-4c7c-b7de-c5fd6ca6fede",
		Context:  make(map[string]string),
	}
}

func WorkflowSign(url, appID string, appSecret string) string {
	return url + "?operator_worker_user_id="
}

func SendWorkflowRequest(msg *WorkflowRequest, path string, username string) (*WorkflowResponse, error) {
	//请求地址模板

	//创建一个请求
	url := WorkflowSign(WORKFLOW_API+path, "100678", "appSecret") + username
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		logger.DefaultLogger.Errorf("request:", err)
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-Domain-Id", msg.Creator)
	//发送请求
	resp, err := client.Do(req)

	if err != nil {
		logger.DefaultLogger.Errorf("resp:", err)
		return nil, err
	}

	//关闭请求
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("http code != 200 error")
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response WorkflowResponse
	err = json.Unmarshal(b, &response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func CreateWorkflowInstance(order *model.CoreSqlOrder, user *Token) (workflowID string, err error) {
	req := NewWorkflowRequest()

	req.Creator = user.Username
	var orderType string
	switch order.Type {
	case 0:
		orderType = "DDL"
	case 1:
		orderType = "DML"
	case 2:
		orderType = "QUERY"
	default:
		orderType = "COMMON"
	}
	req.Context["flowType"] = orderType
	req.Context["flowID"] = order.WorkId
	req.Context["applier"] = user.Username
	var area string
	if strings.Contains(area, "欧洲") {
		area = "欧洲"
	} else {
		area = "国内"
	}
	req.Context["area"] = area

	var env string
	if strings.Contains(order.Source, "prod") {
		env = "prod"
	}
	req.Context["env"] = env

	req.Context["dataSource"] = order.Source

	resp, err := SendWorkflowRequest(req, "/api/v1/instance/create", user.Username)

	if err != nil {
		logger.DefaultLogger.Errorf("send workflow request failed with:", err)
		return "", err
	}

	if resp.ResultCode != "success" {
		return "", errors.New("create workflow instance failed")
	}

	type CreateWorkflowInstanceResp struct {
		WorkflowInstanceID string `json:"flow_instance_id"`
	}

	if x, ok := resp.Data.(CreateWorkflowInstanceResp); ok {
		return x.WorkflowInstanceID, nil
	} else {
		return "", errors.New("typo: not type of CreateWorkflowInstanceResp")
	}
}

//{
//        "flow_code": "7a5ebe44-0ee2-40b1-add4-cd8e315d302f",
//        "flow_name": "小欣的测试5-13",
//        "flow_description": "小欣的测试5-13",
//        "creator": "shingu.gu",
//        "flow_instance_id": "XXDCSH-2022052800798-7900",
//        "status": "processing",
//        "created_time": 1653667283000,
//        "status_name": "processing",
//        "context": {
//            "day": "122",
//            "age": "",
//            "workflow_created_time": "1653667283624",
//            "workflow_creator": "shingu.gu"
//        },
//        "meta_data": [
//            {
//                "name": "day",
//                "type": "string",
//                "field": "day",
//                "type_code": "Input"
//            },
//            {
//                "name": "age",
//                "type": "string",
//                "field": "age",
//                "type_code": "Input"
//            }
//        ],
//        "nodes": [
//            {
//                "name": "start",
//                "next": [
//                    "bnzR3EspmPpwvFWdQPPR"
//                ],
//                "flowInstanceNodeType": "start",
//                "flowTemplateNodeId": "start",
//                "flowInstanceId": "XXDCSH-2022052800798-7900",
//                "id": 568,
//                "flow_instace_id": "XXDCSH-2022052800798-7900",
//                "flow_instance_node_id": "d148cf05-361c-45a3-9d90-a6930fdea470",
//                "status": "accept",
//                "flow_template_node_id": "start",
//                "flow_instance_node_type": "start",
//                "start_time": 1653667283000,
//                "end_time": 0,
//                "tasks": [],
//                "extra_info": {},
//                "flow_instance_node_name": "start",
//                "operator": []
//            },
//            {
//                "name": "bnzR3EspmPpwvFWdQPPR",
//                "next": [
//                    "end"
//                ],
//                "type": "normal",
//                "operator": [
//                    {
//                        "user_name": "Shingu GU （顾小欣）",
//                        "worker_user_id": "shingu.gu"
//                    }
//                ],
//                "flowInstanceNodeRule": "or",
//                "flowInstanceNodeTitle": "审批人bnzR3",
//
//                "flowInstanceNodeType": "normal",
//                "flowTemplateNodeId": "bnzR3EspmPpwvFWdQPPR",
//                "flowInstanceId": "XXDCSH-2022052800798-7900",
//                "id": 569,
//                "flow_instace_id": "XXDCSH-2022052800798-7900",
//                "flow_instance_node_id": "2c4a8ae5-6bf4-4732-9a90-46d1fdfb8e0f",
//                "status": "processing",
//                "flow_template_node_id": "bnzR3EspmPpwvFWdQPPR",
//                "flow_instance_node_type": "normal",
//                "start_time": 1653667283000,
//                "end_time": 0,
//                "flow_instance_node_rule": "or",
//                "flow_instance_node_title": "审批人bnzR3",
//                "tasks": [
//                    {
//                        "operator": {
//                            "user_name": "Shingu GU （顾小欣）",
//                            "worker_user_id": "shingu.gu"
//                        },
//                        "status": "processing",
//                        "origin_operator": "",
//                        "task_id": "2647d55d-dc33-4b8f-a7f0-7080b089fe72",
//                        "extra_info": null,
//                        "start_time": 1653667284000,
//                        "end_time": 0
//                    }
//                ],
//                "extra_info": {},
//                "flow_instance_node_name": "bnzR3EspmPpwvFWdQPPR"
//            },
//            {
//                "name": "end",
//                "next": [],
//                "flowInstanceNodeType": "end",
//                "flowTemplateNodeId": "end",
//                "flowInstanceId": "XXDCSH-2022052800798-7900",
//                "flow_instace_id": "XXDCSH-2022052800798-7900",
//                "flow_instance_node_id": "cc14a584-dcd7-4f00-8b16-ad73286bf71e",
//                "flow_template_node_id": "end",
//                "flow_instance_node_type": "end",
//                "is_forcast": true,
//                "tasks": [],
//                "extra_info": {},
//                "flow_instance_node_name": "end",
//                "operator": []
//            }
//        ],
//        "attachments": [
//            {
//                "type": "",
//                "url": "http://www.nio.com/1.jpg",
//                "filename": "icon",
//                "id": "jpg"
//            }
//        ]
//    }
type FlowDetail struct {
	FlowName       string            `json:"flow_name"`
	FlowInstanceID string            `json:"flow_instance_id"`
	Status         string            `json:"status"`
	Context        map[string]string `json:"context"`

	Nodes []FlowNode `json:"nodes"`
}

type FlowNode struct {
	Name string `json:"name"`
	//"审批"
	FlowInstanceNodeTitle string `json:"flowInstanceNodeTitle"`
	//eg. "or"
	FlowInstanceNodeRule string `json:"flow_instance_node_rule"`

	Tasks []FlowTask `json:"tasks"`
}

type FlowTask struct {
	Operator FlowOperator `json:"operator"`
	Status   string       `json:"status"`
}

type FlowOperator struct {
	UserName string `json:"user_name"`
	UserID   string `json:"worker_user_id"`
}

func CallBackWorkflowInstance(workflowInstanceID string, username string) (*FlowDetail, error) {

	url := WORKFLOW_API + "/api/v1/instance/" + workflowInstanceID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.DefaultLogger.Errorf("request:", err)
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	//设置请求头
	//req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-Domain-Id", username)
	//发送请求
	resp, err := client.Do(req)

	if err != nil {
		logger.DefaultLogger.Errorf("resp:", err)
		return nil, err
	}

	//关闭请求
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var flowDetail FlowDetail
	err = json.Unmarshal(body, &flowDetail)
	if err != nil {
		logger.DefaultLogger.Errorf("unmarshal err:", err)
		return nil, err
	}

	return &flowDetail, nil
}
