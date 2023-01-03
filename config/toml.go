package config

type ConfigTOML struct {
	General            GeneralTOML       `toml:"general"`
	Style              StyleTOML         `toml:"style"`
	Media              MediaTOML         `toml:"media"`
	OpenPattern        OpenPatternTOML   `toml:"open-pattern"`
	OpenCustom         OpenCustomTOML    `toml:"open-custom"`
	NotificationConfig NotificationsTOML `toml:"desktop-notification"`
	Input              InputTOML         `toml:"input"`
}

type GeneralTOML struct {
	Editor              *string             `toml:"editor"`
	Confirmation        *bool               `toml:"confirmation"`
	MouseSupport        *bool               `toml:"mouse-support"`
	DateFormat          *string             `toml:"date-format"`
	DateTodayFormat     *string             `toml:"date-today-format"`
	DateRelative        *int                `toml:"date-relative"`
	MaxWidth            *int                `toml:"max-width"`
	QuoteReply          *bool               `toml:"quote-reply"`
	ShortHints          *bool               `toml:"short-hints"`
	ShowFilterPhrase    *bool               `toml:"show-filter-phrase"`
	ListPlacement       *string             `toml:"list-placement"`
	ListSplit           *string             `toml:"list-split"`
	ListProportion      *int                `toml:"list-proportion"`
	ContentProportion   *int                `toml:"content-proportion"`
	TerminalTitle       *int                `toml:"terminal-title"`
	ShowIcons           *bool               `toml:"show-icons"`
	ShowHelp            *bool               `toml:"show-help"`
	RedrawUI            *bool               `toml:"redraw-ui"`
	LeaderKey           *string             `toml:"leader-key"`
	LeaderTimeout       *int64              `toml:"leader-timeout"`
	Timelines           *[]TimelineTOML     `toml:"timelines"`
	LeaderActions       *[]LeaderActionTOML `toml:"leader-actions"`
	StickToTop          *bool               `toml:"stick-to-top"`
	NotificationsToHide *[]string           `toml:"notifications-to-hide"`
	ShowBoostedUser     *bool               `toml:"show-boosted-user"`
}

type TimelineTOML struct {
	Name        *string   `toml:"name"`
	Type        *string   `toml:"type"`
	Data        *string   `toml:"data"`
	Keys        *[]string `toml:"keys"`
	SpecialKeys *[]string `toml:"special-keys"`
	Shortcut    *string   `toml:"shortcut"`
	HideBoosts  *bool     `toml:"hide-boosts"`
	HideReplies *bool     `toml:"hide-replies"`
}

type LeaderActionTOML struct {
	Type     *string `toml:"type"`
	Data     *string `toml:"data"`
	Shortcut *string `toml:"shortcut"`
}

type StyleTOML struct {
	Theme *string `toml:"theme"`

	XrdbPrefix *string `toml:"xrdb-prefix"`

	Background *string `toml:"background"`
	Text       *string `toml:"text"`

	Subtle      *string `toml:"subtle"`
	WarningText *string `toml:"warning-text"`

	TextSpecial1 *string `toml:"text-special-one"`
	TextSpecial2 *string `toml:"text-special-two"`

	TopBarBackground *string `toml:"top-bar-background"`
	TopBarText       *string `toml:"top-bar-text"`

	StatusBarBackground *string `toml:"status-bar-background"`
	StatusBarText       *string `toml:"status-bar-text"`

	StatusBarViewBackground *string `toml:"status-bar-view-background"`
	StatusBarViewText       *string `toml:"status-bar-view-text"`

	ListSelectedBackground *string `toml:"list-selected-background"`
	ListSelectedText       *string `toml:"list-selected-text"`

	ListSelectedInactiveBackground *string `toml:"list-selected-inactive-background"`
	ListSelectedInactiveText       *string `toml:"list-selected-inactive-text"`

	ControlsText      *string `toml:"controls-text"`
	ControlsHighlight *string `toml:"controls-highlight"`

	AutocompleteBackground *string `toml:"autocomplete-background"`
	AutocompleteText       *string `toml:"autocomplete-text"`

	AutocompleteSelectedBackground *string `toml:"autocomplete-selected-background"`
	AutocompleteSelectedText       *string `toml:"autocomplete-selected-text"`

	ButtonColorOne *string `toml:"button-color-one"`
	ButtonColorTwo *string `toml:"button-color-two"`

	TimelineNameBackground *string `toml:"timeline-name-background"`
	TimelineNameText       *string `toml:"timeline-name-text"`

	IconColor *string `toml:"icon-color"`

	CommandText *string `toml:"command-text"`
}

type ViewerTOML struct {
	Program  *string `toml:"program"`
	Args     *string `toml:"args"`
	Terminal *bool   `toml:"terminal"`
	Single   *bool   `toml:"single"`
	Reverse  *bool   `toml:"reverse"`
}

type MediaTOML struct {
	Image *ViewerTOML `toml:"image"`
	Video *ViewerTOML `toml:"video"`
	Audio *ViewerTOML `toml:"audio"`
	Link  *ViewerTOML `toml:"link"`
}

type PatternTOML struct {
	Matching *string `toml:"matching"`
	Program  *string `toml:"program"`
	Args     *string `toml:"args"`
	Terminal *bool   `toml:"terminal"`
}

type OpenPatternTOML struct {
	Patterns *[]PatternTOML `toml:"patterns"`
}

type CustomTOML struct {
	Program     *string   `toml:"program"`
	Args        *string   `toml:"args"`
	Terminal    *bool     `toml:"terminal"`
	Hint        *string   `toml:"hint"`
	Keys        *[]string `toml:"keys"`
	SpecialKeys *[]string `toml:"special-keys"`
}

type OpenCustomTOML struct {
	Programs *[]CustomTOML `toml:"programs"`
}

type NotificationsTOML struct {
	Followers *bool `toml:"followers"`
	Favorite  *bool `toml:"favorite"`
	Mention   *bool `toml:"mention"`
	Update    *bool `toml:"update"`
	Boost     *bool `toml:"boost"`
	Poll      *bool `toml:"poll"`
	Posts     *bool `toml:"posts"`
}

type KeyHintTOML struct {
	Hint        *string   `toml:"hint"`
	HintAlt     *string   `toml:"hint-alt"`
	Keys        *[]string `toml:"keys"`
	SpecialKeys *[]string `toml:"special-keys"`
}

type InputTOML struct {
	GlobalDown  *KeyHintTOML `toml:"global-down"`
	GlobalUp    *KeyHintTOML `toml:"global-up"`
	GlobalEnter *KeyHintTOML `toml:"global-enter"`
	GlobalBack  *KeyHintTOML `toml:"global-back"`
	GlobalExit  *KeyHintTOML `toml:"global-exit"`

	MainHome     *KeyHintTOML `toml:"main-home"`
	MainEnd      *KeyHintTOML `toml:"main-end"`
	MainPrevFeed *KeyHintTOML `toml:"main-prev-feed"`
	MainNextFeed *KeyHintTOML `toml:"main-next-feed"`
	MainPrevPane *KeyHintTOML `toml:"main-prev-pane"`
	MainNextPane *KeyHintTOML `toml:"main-next-pane"`
	MainCompose  *KeyHintTOML `toml:"main-compose"`

	StatusAvatar       *KeyHintTOML `toml:"status-avatar"`
	StatusBoost        *KeyHintTOML `toml:"status-boost"`
	StatusDelete       *KeyHintTOML `toml:"status-delete"`
	StatusEdit         *KeyHintTOML `toml:"status-edit"`
	StatusFavorite     *KeyHintTOML `toml:"status-favorite"`
	StatusMedia        *KeyHintTOML `toml:"status-media"`
	StatusLinks        *KeyHintTOML `toml:"status-links"`
	StatusPoll         *KeyHintTOML `toml:"status-poll"`
	StatusReply        *KeyHintTOML `toml:"status-reply"`
	StatusBookmark     *KeyHintTOML `toml:"status-bookmark"`
	StatusThread       *KeyHintTOML `toml:"status-thread"`
	StatusUser         *KeyHintTOML `toml:"status-user"`
	StatusViewFocus    *KeyHintTOML `toml:"status-view-focus"`
	StatusYank         *KeyHintTOML `toml:"status-yank"`
	StatusToggleCW     *KeyHintTOML `toml:"status-toggle-cw"`
	StatusShowFiltered *KeyHintTOML `toml:"status-show-filtered"`

	UserAvatar              *KeyHintTOML `toml:"user-avatar"`
	UserBlock               *KeyHintTOML `toml:"user-block"`
	UserFollow              *KeyHintTOML `toml:"user-follow"`
	UserFollowRequestDecide *KeyHintTOML `toml:"user-follow-request-decide"`
	UserMute                *KeyHintTOML `toml:"user-mute"`
	UserLinks               *KeyHintTOML `toml:"user-links"`
	UserUser                *KeyHintTOML `toml:"user-user"`
	UserViewFocus           *KeyHintTOML `toml:"user-view-focus"`
	UserYank                *KeyHintTOML `toml:"user-yank"`

	ListOpenFeed   *KeyHintTOML `toml:"list-open-feed"`
	ListUserList   *KeyHintTOML `toml:"list-user-list"`
	ListUserAdd    *KeyHintTOML `toml:"list-user-add"`
	ListUserDelete *KeyHintTOML `toml:"list-user-delete"`

	TagOpenFeed *KeyHintTOML `toml:"tag-open-feed"`
	TagFollow   *KeyHintTOML `toml:"tag-follow"`

	LinkOpen *KeyHintTOML `toml:"link-open"`
	LinkYank *KeyHintTOML `toml:"link-yank"`

	ComposeEditCW               *KeyHintTOML `toml:"compose-edit-cw"`
	ComposeEditText             *KeyHintTOML `toml:"compose-edit-text"`
	ComposeIncludeQuote         *KeyHintTOML `toml:"compose-include-quote"`
	ComposeMediaFocus           *KeyHintTOML `toml:"compose-media-focus"`
	ComposePost                 *KeyHintTOML `toml:"compose-post"`
	ComposeToggleContentWarning *KeyHintTOML `toml:"compose-toggle-content-warning"`
	ComposeVisibility           *KeyHintTOML `toml:"compose-visibility"`
	ComposeLanguage             *KeyHintTOML `toml:"compose-language"`
	ComposePoll                 *KeyHintTOML `toml:"compose-poll"`

	MediaDelete   *KeyHintTOML `toml:"media-delete"`
	MediaEditDesc *KeyHintTOML `toml:"media-edit-desc"`
	MediaAdd      *KeyHintTOML `toml:"media-add"`

	VoteVote   *KeyHintTOML `toml:"vote-vote"`
	VoteSelect *KeyHintTOML `toml:"vote-select"`

	PollAdd         *KeyHintTOML `toml:"poll-add"`
	PollEdit        *KeyHintTOML `toml:"poll-edit"`
	PollDelete      *KeyHintTOML `toml:"poll-delete"`
	PollMultiToggle *KeyHintTOML `toml:"poll-multi-toggle"`
	PollExpiration  *KeyHintTOML `toml:"poll-expiration"`

	PreferenceName         *KeyHintTOML `toml:"preference-name"`
	PreferenceVisibility   *KeyHintTOML `toml:"preference-visibility"`
	PreferenceBio          *KeyHintTOML `toml:"preference-bio"`
	PreferenceSave         *KeyHintTOML `toml:"preference-save"`
	PreferenceFields       *KeyHintTOML `toml:"preference-fields"`
	PreferenceFieldsAdd    *KeyHintTOML `toml:"preference-fields-add"`
	PreferenceFieldsEdit   *KeyHintTOML `toml:"preference-fields-edit"`
	PreferenceFieldsDelete *KeyHintTOML `toml:"preference-fields-delete"`
}
