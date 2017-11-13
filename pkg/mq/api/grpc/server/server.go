
// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.
 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.
 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
 
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package server

import (
	"fmt"

	"github.com/goodrain/rainbond/pkg/mq/api/grpc/pb"
	"github.com/goodrain/rainbond/pkg/mq/api/mq"

	"github.com/Sirupsen/logrus"

	context "golang.org/x/net/context"

	proto "github.com/golang/protobuf/proto"
	grpc1 "google.golang.org/grpc"
)

type mqServer struct {
	actionMQ mq.ActionMQ
}

func (s *mqServer) Enqueue(ctx context.Context, in *pb.EnqueueRequest) (*pb.TaskReply, error) {
	if in.Topic == "" || !s.actionMQ.TopicIsExist(in.Topic) {
		return nil, fmt.Errorf("topic %s is not support", in.Topic)
	}
	message, err := proto.Marshal(in.Message)
	if err != nil {
		return nil, err
	}
	err = s.actionMQ.Enqueue(ctx, in.Topic, string(message))
	if err != nil {
		return nil, err
	}
	logrus.Debugf("task (%v) enqueue.", in.Message.String())
	return &pb.TaskReply{
		Status: "success",
	}, nil
}
func (s *mqServer) Topics(ctx context.Context, in *pb.TopicRequest) (*pb.TaskReply, error) {
	return &pb.TaskReply{
		Status: "success",
		Topics: s.actionMQ.GetAllTopics(),
	}, nil
}

func (s *mqServer) Dequeue(ctx context.Context, in *pb.DequeueRequest) (*pb.TaskMessage, error) {
	if in.Topic == "" || !s.actionMQ.TopicIsExist(in.Topic) {
		return nil, fmt.Errorf("topic %s is not support", in.Topic)
	}
	message, err := s.actionMQ.Dequeue(ctx, in.Topic)
	if err != nil {
		return nil, err
	}
	var task pb.TaskMessage
	err = proto.Unmarshal([]byte(message), &task)
	if err != nil {
		return nil, err
	}
	logrus.Infof("task (%s) dnqueue by (%s).", task.GetTaskType(), in.ClientHost)
	return &task, nil
}

//RegisterServer 注册服务
func RegisterServer(server *grpc1.Server, actionMQ mq.ActionMQ) {
	pb.RegisterTaskQueueServer(server, &mqServer{actionMQ})
}