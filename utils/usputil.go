package utils

import (
	"context"
	"tr369-wss-client/common"
	logger "tr369-wss-client/log"
	"tr369-wss-client/pkg/api"

	"google.golang.org/protobuf/proto"
)

func DecodeUSPMessage(binary []byte) (uspMsg *api.Msg, err error) {
	uspMsg = new(api.Msg)

	opts := proto.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	err = opts.Unmarshal(binary, uspMsg)
	if err != nil {
		logger.Warnf("DecodeUSPMessage failed, err: %s", err)
		return nil, err
	}
	return uspMsg, nil
}

func DecodeUSPRecord(binary []byte) (uspRecord *api.Record, err error) {
	uspRecord = new(api.Record)
	err = proto.Unmarshal(binary, uspRecord)
	if err != nil {
		logger.Warnf("Unmarshal failed, err: %s", err)
		return nil, err
	}
	return uspRecord, nil
}

func EncodeUspMessage(msg *api.Msg) []byte {
	result, _ := proto.Marshal(msg)
	return result
}

func EncodeUspRecord(rec *api.Record) ([]byte, error) {
	result, err := proto.Marshal(rec)
	return result, err
}

func CreateGetResponseMessage(msgId string, resp api.Response_GetResp) (result *api.Msg) {
	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_GET_RESP,
			MsgId:   msgId,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Response{
				Response: &api.Response{
					RespType: &resp,
				},
			},
		},
	}

	return
}

func CreateSetResponseMessage(msgId string, requestPath []string, affectedPath []string, updatedParams []map[string]string) (result *api.Msg) {
	var updatedObjResults []*api.SetResp_UpdatedObjectResult
	for k, path := range requestPath {
		updatedObjResult := &api.SetResp_UpdatedObjectResult{
			RequestedPath: path,
			OperStatus: &api.SetResp_UpdatedObjectResult_OperationStatus{
				OperStatus: &api.SetResp_UpdatedObjectResult_OperationStatus_OperSuccess{
					OperSuccess: &api.SetResp_UpdatedObjectResult_OperationStatus_OperationSuccess{
						UpdatedInstResults: []*api.SetResp_UpdatedInstanceResult{
							{
								AffectedPath:  affectedPath[k],
								UpdatedParams: updatedParams[k],
							},
						},
					},
				},
			},
		}
		updatedObjResults = append(updatedObjResults, updatedObjResult)
	}
	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_SET_RESP,
			MsgId:   msgId,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Response{
				Response: &api.Response{
					RespType: &api.Response_SetResp{
						SetResp: &api.SetResp{
							UpdatedObjResults: updatedObjResults,
						},
					},
				},
			},
		},
	}
	return
}

func CreateOperateResponseMessage(msgId string, operationResults []*api.OperateResp_OperationResult) (result *api.Msg) {

	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_OPERATE_RESP,
			MsgId:   msgId,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Response{
				Response: &api.Response{
					RespType: &api.Response_OperateResp{
						OperateResp: &api.OperateResp{
							OperationResults: operationResults,
						},
					},
				},
			},
		},
	}
	return
}

func CreateAddResponseMessage(msgId string, requestPath []string, affectedPath []string, updatedParams []map[string]string) (result *api.Msg) {
	var createdObjResults []*api.AddResp_CreatedObjectResult
	for k, path := range requestPath {
		createdObjResult := &api.AddResp_CreatedObjectResult{
			RequestedPath: path,
			OperStatus: &api.AddResp_CreatedObjectResult_OperationStatus{
				OperStatus: &api.AddResp_CreatedObjectResult_OperationStatus_OperSuccess{
					OperSuccess: &api.AddResp_CreatedObjectResult_OperationStatus_OperationSuccess{
						InstantiatedPath: affectedPath[k],
						UniqueKeys:       updatedParams[k],
					},
				},
			},
		}
		createdObjResults = append(createdObjResults, createdObjResult)
	}

	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_ADD_RESP,
			MsgId:   msgId,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Response{
				Response: &api.Response{
					RespType: &api.Response_AddResp{
						AddResp: &api.AddResp{
							CreatedObjResults: createdObjResults,
						},
					},
				},
			},
		},
	}
	return
}

func CreateDeleteResponseMessage(msgId string, requestPath []string, affectedPath []string) (result *api.Msg) {
	var deletedObjResults []*api.DeleteResp_DeletedObjectResult
	for k, path := range requestPath {
		deletedObjResult := &api.DeleteResp_DeletedObjectResult{
			RequestedPath: path,
			OperStatus: &api.DeleteResp_DeletedObjectResult_OperationStatus{
				OperStatus: &api.DeleteResp_DeletedObjectResult_OperationStatus_OperSuccess{
					OperSuccess: &api.DeleteResp_DeletedObjectResult_OperationStatus_OperationSuccess{
						AffectedPaths: []string{affectedPath[k]},
					},
				},
			},
		}
		deletedObjResults = append(deletedObjResults, deletedObjResult)
	}

	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_DELETE_RESP,
			MsgId:   msgId,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Response{
				Response: &api.Response{
					RespType: &api.Response_DeleteResp{
						DeleteResp: &api.DeleteResp{
							DeletedObjResults: deletedObjResults,
						},
					},
				},
			},
		},
	}
	return
}

func CreateGetSupportedDMMessage(ctx context.Context, paths []string, firstLevelOnly, retCmds, retEvents, retParams bool) (result *api.Msg) {
	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_GET_SUPPORTED_DM,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Request{
				Request: &api.Request{
					ReqType: &api.Request_GetSupportedDm{
						GetSupportedDm: &api.GetSupportedDM{
							ObjPaths:       paths,
							FirstLevelOnly: firstLevelOnly,
							ReturnCommands: retCmds,
							ReturnEvents:   retEvents,
							ReturnParams:   retParams,
						},
					},
				},
			},
		},
	}

	return
}

func CreateGetInstancesMessage(ctx context.Context, paths []string, firstLevelOnly bool) (result *api.Msg) {
	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_GET_INSTANCES,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Request{
				Request: &api.Request{
					ReqType: &api.Request_GetInstances{
						GetInstances: &api.GetInstances{
							ObjPaths:       paths,
							FirstLevelOnly: firstLevelOnly,
						},
					},
				},
			},
		},
	}

	return
}

func CreateOperateCompleteMessage(objPath string, commandName string, commandKey string, outputArgs map[string]string) (result *api.Msg) {

	completeNotify := &api.Notify{
		SubscriptionId: "",
		SendResp:       false,
		Notification: &api.Notify_OperComplete{
			OperComplete: &api.Notify_OperationComplete{
				ObjPath:     objPath,
				CommandName: commandName,
				CommandKey:  commandKey,
				OperationResp: &api.Notify_OperationComplete_ReqOutputArgs{
					ReqOutputArgs: &api.Notify_OperationComplete_OutputArgs{
						OutputArgs: outputArgs,
					},
				},
			},
		},
	}

	result = CreateNotifyMessage(completeNotify)
	return
}

func CreateNotifyMessage(notify *api.Notify) (result *api.Msg) {
	return &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_NOTIFY,
			MsgId:   common.RandStr(10),
		},
		Body: &api.Body{
			MsgBody: &api.Body_Request{
				Request: &api.Request{
					ReqType: &api.Request_Notify{
						Notify: notify,
					},
				},
			},
		},
	}
}

func CreateGetSupportedProtocolMessage(ctx context.Context, controllerSupported string) (result *api.Msg) {
	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_GET_SUPPORTED_PROTO,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Request{
				Request: &api.Request{
					ReqType: &api.Request_GetSupportedProtocol{
						GetSupportedProtocol: &api.GetSupportedProtocol{
							ControllerSupportedProtocolVersions: controllerSupported,
						},
					},
				},
			},
		},
	}

	return
}

func CreateUspRecordNoSession(ver, to, from string, msg *api.Msg) (result *api.Record) {
	result = &api.Record{
		Version: ver,
		ToId:    to,
		FromId:  from,
		RecordType: &api.Record_NoSessionContext{
			NoSessionContext: &api.NoSessionContextRecord{
				Payload: EncodeUspMessage(msg),
			},
		},
	}

	return
}
