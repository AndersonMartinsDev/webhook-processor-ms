package model

// Main struct to hold the entire JSON payload
type WWEBJSPayload struct {
	Data struct {
		ID struct {
			FromMe     bool   `json:"fromMe"`
			Remote     string `json:"remote"`
			ID         string `json:"id"`
			Serialized string `json:"_serialized"`
		} `json:"id"`
		Viewed                      bool     `json:"viewed"`
		Body                        string   `json:"body"`
		Type                        string   `json:"type"`
		T                           int64    `json:"t"`
		ClientReceivedTsMillis      int64    `json:"clientReceivedTsMillis"`
		NotifyName                  string   `json:"notifyName"`
		From                        string   `json:"from"`
		To                          string   `json:"to"`
		Ack                         int      `json:"ack"`
		Invis                       bool     `json:"invis"`
		IsNewMsg                    bool     `json:"isNewMsg"`
		Star                        bool     `json:"star"`
		KicNotified                 bool     `json:"kicNotified"`
		RecvFresh                   bool     `json:"recvFresh"`
		IsFromTemplate              bool     `json:"isFromTemplate"`
		PollInvalidated             bool     `json:"pollInvalidated"`
		IsSentCagPollCreation       bool     `json:"isSentCagPollCreation"`
		LatestEditMsgKey            *string  `json:"latestEditMsgKey"`
		LatestEditSenderTimestampMs *int64   `json:"latestEditSenderTimestampMs"`
		MentionedJidList            []string `json:"mentionedJidList"`
		GroupMentions               []string `json:"groupMentions"`
		IsEventCanceled             bool     `json:"isEventCanceled"`
		EventInvalidated            bool     `json:"eventInvalidated"`
		IsVcardOverMmsDocument      bool     `json:"isVcardOverMmsDocument"`
		IsForwarded                 bool     `json:"isForwarded"`
		IsQuestion                  bool     `json:"isQuestion"`
		QuestionReplyQuotedMessage  *string  `json:"questionReplyQuotedMessage"`
		QuestionResponsesCount      int      `json:"questionResponsesCount"`
		ReadQuestionResponsesCount  int      `json:"readQuestionResponsesCount"`
		HasReaction                 bool     `json:"hasReaction"`
		ViewMode                    string   `json:"viewMode"`
		// MessageSecret                []int    `json:"messageSecret"`
		ProductHeaderImageRejected   bool    `json:"productHeaderImageRejected"`
		LastPlaybackProgress         int     `json:"lastPlaybackProgress"`
		IsDynamicReplyButtonsMsg     bool    `json:"isDynamicReplyButtonsMsg"`
		IsCarouselCard               bool    `json:"isCarouselCard"`
		ParentMsgId                  *string `json:"parentMsgId"`
		CallSilenceReason            *string `json:"callSilenceReason"`
		IsVideoCall                  bool    `json:"isVideoCall"`
		CallDuration                 *string `json:"callDuration"`
		CallCreator                  *string `json:"callCreator"`
		CallParticipants             *string `json:"callParticipants"`
		IsCallLink                   *string `json:"isCallLink"`
		CallLinkToken                *string `json:"callLinkToken"`
		IsMdHistoryMsg               bool    `json:"isMdHistoryMsg"`
		StickerSentTs                int     `json:"stickerSentTs"`
		IsAvatar                     bool    `json:"isAvatar"`
		LastUpdateFromServerTs       int     `json:"lastUpdateFromServerTs"`
		InvokedBotWid                *string `json:"invokedBotWid"`
		BizBotType                   *string `json:"bizBotType"`
		BotResponseTargetId          *string `json:"botResponseTargetId"`
		BotPluginType                *string `json:"botPluginType"`
		BotPluginReferenceIndex      *string `json:"botPluginReferenceIndex"`
		BotPluginSearchProvider      *string `json:"botPluginSearchProvider"`
		BotPluginSearchUrl           *string `json:"botPluginSearchUrl"`
		BotPluginSearchQuery         *string `json:"botPluginSearchQuery"`
		BotPluginMaybeParent         bool    `json:"botPluginMaybeParent"`
		BotReelPluginThumbnailCdnUrl *string `json:"botReelPluginThumbnailCdnUrl"`
		BotMessageDisclaimerText     *string `json:"botMessageDisclaimerText"`
		BotMsgBodyType               *string `json:"botMsgBodyType"`
		// ReportingTokenInfo           struct {
		// 	ReportingToken []int `json:"reportingToken"`
		// 	Version        int   `json:"version"`
		// 	ReportingTag   []int `json:"reportingTag"`
		// } `json:"reportingTokenInfo"`
		RequiresDirectConnection              *string  `json:"requiresDirectConnection"`
		BizContentPlaceholderType             *string  `json:"bizContentPlaceholderType"`
		HostedBizEncStateMismatch             bool     `json:"hostedBizEncStateMismatch"`
		SenderOrRecipientAccountTypeHosted    bool     `json:"senderOrRecipientAccountTypeHosted"`
		PlaceholderCreatedWhenAccountIsHosted bool     `json:"placeholderCreatedWhenAccountIsHosted"`
		GalaxyFlowDisabled                    bool     `json:"galaxyFlowDisabled"`
		Links                                 []string `json:"links"`
	} `json:"_data"`
	ID struct {
		FromMe     bool   `json:"fromMe"`
		Remote     string `json:"remote"`
		ID         string `json:"id"`
		Serialized string `json:"_serialized"`
	} `json:"id"`
	Ack             int      `json:"ack"`
	HasMedia        bool     `json:"hasMedia"`
	Body            string   `json:"body"`
	Type            string   `json:"type"`
	Timestamp       int64    `json:"timestamp"`
	From            string   `json:"from"`
	To              string   `json:"to"`
	DeviceType      string   `json:"deviceType"`
	IsForwarded     bool     `json:"isForwarded"`
	ForwardingScore int      `json:"forwardingScore"`
	IsStatus        bool     `json:"isStatus"`
	IsStarred       bool     `json:"isStarred"`
	FromMe          bool     `json:"fromMe"`
	HasQuotedMsg    bool     `json:"hasQuotedMsg"`
	HasReaction     bool     `json:"hasReaction"`
	VCards          []string `json:"vCards"`
	MentionedIds    []string `json:"mentionedIds"`
	GroupMentions   []string `json:"groupMentions"`
	IsGif           bool     `json:"isGif"`
	Links           []string `json:"links"`
}
