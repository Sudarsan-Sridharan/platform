// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/mattermost/platform/model"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

type JoinProvider struct {
}

const (
	CMD_JOIN = "join"
)

func init() {
	RegisterCommandProvider(&JoinProvider{})
}

func (me *JoinProvider) GetTrigger() string {
	return CMD_JOIN
}

func (me *JoinProvider) GetCommand(T goi18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CMD_JOIN,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_join.desc"),
		AutoCompleteHint: T("api.command_join.hint"),
		DisplayName:      T("api.command_join.name"),
	}
}

func (me *JoinProvider) DoCommand(args *model.CommandArgs, message string) *model.CommandResponse {
	if result := <-Srv.Store.Channel().GetByName(args.TeamId, message, true); result.Err != nil {
		return &model.CommandResponse{Text: args.T("api.command_join.list.app_error"), ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL}
	} else {
		channel := result.Data.(*model.Channel)

		if channel.Name == message {

			if channel.Type != model.CHANNEL_OPEN {
				return &model.CommandResponse{Text: args.T("api.command_join.fail.app_error"), ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL}
			}

			if err := JoinChannel(channel, args.UserId); err != nil {
				return &model.CommandResponse{Text: args.T("api.command_join.fail.app_error"), ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL}
			}

			team, err := GetTeam(channel.TeamId)
			if err != nil {
				return &model.CommandResponse{Text: args.T("api.command_join.fail.app_error"), ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL}
			}

			return &model.CommandResponse{GotoLocation: args.SiteURL + "/" + team.Name + "/channels/" + channel.Name, Text: args.T("api.command_join.success"), ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL}
		}
	}

	return &model.CommandResponse{ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: args.T("api.command_join.missing.app_error")}
}
