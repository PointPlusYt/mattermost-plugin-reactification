package main

import (
    "fmt"
    "strings"

    "github.com/mattermost/mattermost-server/v6/model"
    "github.com/mattermost/mattermost-server/v6/plugin"
)

type Plugin struct {
    plugin.MattermostPlugin
    botUserID string
}

func (p *Plugin) OnActivate() error {
    // Créer le bot via helpers
    botUserID, appErr := p.API.CreateBot(&model.Bot{
        Username:    "reactification",
        DisplayName: "Reactification",
        Description: "Notifications de réactions sur vos messages",
    })
    if appErr != nil {
        // Le bot existe peut-être déjà, on le cherche
        user, err2 := p.API.GetUserByUsername("reactification")
        if err2 != nil {
            return fmt.Errorf("Impossible de créer ou trouver le bot")
        }
        p.botUserID = user.Id
    } else {
        p.botUserID = botUserID.UserId
    }

    return p.API.RegisterCommand(&model.Command{
        Trigger:          "reactification",
        DisplayName:      "Reactification",
        Description:      "Activer/désactiver les notifications de réactions",
        AutoComplete:     true,
        AutoCompleteDesc: "on | off",
        AutoCompleteHint: "[on|off]",
    })
}

func (p *Plugin) ReactionHasBeenAdded(c *plugin.Context, reaction *model.Reaction) {
    // Récupérer le post ciblé
    post, appErr := p.API.GetPost(reaction.PostId)
    if appErr != nil {
        return
    }

    // Pas de notification si on réagit à son propre message
    if post.UserId == reaction.UserId {
        return
    }

    // Vérifier si l'auteur du post a désactivé les notifications
    if p.isDisabled(post.UserId) {
        return
    }

    // Récupérer l'utilisateur qui a réagi
    reactor, appErr := p.API.GetUser(reaction.UserId)
    if appErr != nil {
        return
    }

    // Récupérer le canal via le post
    channel, appErr := p.API.GetChannel(post.ChannelId)
    if appErr != nil {
        return
    }

    // Construire le message de notification
    message := p.buildNotificationMessage(reactor, reaction, post, channel)

    // Ouvrir ou récupérer le canal DM entre le bot et l'auteur du post
    dmChannel, appErr := p.API.GetDirectChannel(p.botUserID, post.UserId)
    if appErr != nil {
        return
    }

    p.API.CreatePost(&model.Post{
        ChannelId: dmChannel.Id,
        UserId:    p.botUserID,
        Message:   message,
    })
}

func (p *Plugin) buildNotificationMessage(
    reactor *model.User,
    reaction *model.Reaction,
    post *model.Post,
    channel *model.Channel,
) string {
    preview := strings.TrimSpace(post.Message)
    var previewLine string
    if preview == "" {
        previewLine = "votre message"
    } else if len([]rune(preview)) <= 80 {
        previewLine = fmt.Sprintf("\"%s\"", preview)
    } else {
        runes := []rune(preview)
        previewLine = fmt.Sprintf("\"%s...\"", string(runes[:80]))
    }

    siteURL := p.API.GetConfig().ServiceSettings.SiteURL
    postLink := ""
    if *siteURL != "" {
        postLink = fmt.Sprintf("[%s](%s/_redirect/pl/%s)", previewLine, *siteURL, post.Id)
    }
    postChannel := ""
    if channel.DisplayName != "" {
        postChannel = fmt.Sprintf("dans **%s**", channel.DisplayName)
    }

    return fmt.Sprintf(
        "**@%s** a réagi :%s: à %s %s",
        reactor.Username,
        reaction.EmojiName,
        postLink,
        postChannel,
    )
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
    parts := strings.Fields(args.Command)
    if len(parts) < 2 {
        return &model.CommandResponse{
            Text: "Usage: `/reactification on` ou `/reactification off`",
        }, nil
    }

    switch strings.ToLower(parts[1]) {
    case "on":
        p.setDisabled(args.UserId, false)
        return &model.CommandResponse{
            Text: "✅ Notifications de réactions **activées**.",
        }, nil
    case "off":
        p.setDisabled(args.UserId, true)
        return &model.CommandResponse{
            Text: "🔕 Notifications de réactions **désactivées**.",
        }, nil
    default:
        return &model.CommandResponse{
            Text: "Usage: `/reactification on` ou `/reactification off`",
        }, nil
    }
}
