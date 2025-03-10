package requests

import (
	"fmt"
	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"github.com/DavidMovas/SpeakUp-Server/internal/shared/model"
	genID "github.com/DavidMovas/SpeakUp-Server/internal/utils/helpers"
	genName "github.com/DavidMovas/SpeakUp-Server/internal/utils/models/helpers"
	"time"
)

var _ model.Requestable[CreatePrivateChatRequest, *v1.CreateChatRequest_PrivateChat] = (*CreatePrivateChatRequest)(nil)

type CreatePrivateChatRequest struct {
	InitiatorID string
	MemberID    string

	ID        string
	Slug      string
	Name      string
	Type      string
	CreatedAt time.Time
}

func (c CreatePrivateChatRequest) Make(req *v1.CreateChatRequest_PrivateChat) (*CreatePrivateChatRequest, error) {
	c.InitiatorID = req.GetInitiatorId()
	c.MemberID = req.GetMemberId()

	pair := fmt.Sprintf("%s_%s", c.InitiatorID, c.MemberID)
	c.ID = genID.GenerateID()
	c.Slug = pair
	c.Name = pair
	c.Type = "private"
	c.CreatedAt = time.Now()

	return &c, nil
}

var _ model.Requestable[CreateGroupChatRequest, *v1.CreateChatRequest_GroupChat] = (*CreateGroupChatRequest)(nil)

type CreateGroupChatRequest struct {
	InitiatorID string
	MemberIDs   []string
	Name        string

	ID        string
	Slug      string
	Type      string
	CreatedAt time.Time
}

func (c CreateGroupChatRequest) Make(req *v1.CreateChatRequest_GroupChat) (*CreateGroupChatRequest, error) {
	c.InitiatorID = req.GetInitiatorId()
	c.MemberIDs = req.GetMembersIds()
	c.Name = req.GetName()

	c.ID = genID.GenerateID()
	c.Slug = genName.GenerateRandomUsername(c.Name)
	c.Type = "group"
	c.CreatedAt = time.Now()

	return &c, nil
}
