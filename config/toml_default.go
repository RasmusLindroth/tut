package config

var tvar = true
var fvar = false

var bt = &tvar
var bf = &fvar

func sp(s string) *string {
	return &s
}
func ip(i int) *int {
	return &i
}

func ip64(i int64) *int64 {
	return &i
}

var ConfigDefault = ConfigTOML{
	General: GeneralTOML{
		Editor:              sp("USE_TUT_INTERNAL"),
		Confirmation:        bt,
		MouseSupport:        bf,
		DateFormat:          sp("2006-01-02 15:04"),
		DateTodayFormat:     sp("15:04"),
		DateRelative:        ip(-1),
		QuoteReply:          bf,
		MaxWidth:            ip(0),
		ShortHints:          bf,
		ShowFilterPhrase:    bt,
		ShowIcons:           bt,
		ShowHelp:            bt,
		RedrawUI:            bt,
		StickToTop:          bf,
		ShowBoostedUser:     bf,
		ListPlacement:       sp("left"),
		ListSplit:           sp("row"),
		ListProportion:      ip(1),
		ContentProportion:   ip(2),
		TerminalTitle:       ip(0),
		LeaderKey:           sp(""),
		LeaderTimeout:       ip64(1000),
		NotificationsToHide: &[]string{},
		Timelines: &[]TimelineTOML{
			{
				Name:        sp("Home"),
				Type:        sp("home"),
				HideBoosts:  bf,
				HideReplies: bf,
			},
			{
				Name: sp("Notifications"),
				Type: sp("notifications"),
				Keys: &[]string{"n", "N"},
			},
		},
	},
	Style: StyleTOML{
		Theme:                          sp("none"),
		XrdbPrefix:                     sp("guess"),
		Background:                     sp("#272822"),
		Text:                           sp("#f8f8f2"),
		Subtle:                         sp("#808080"),
		WarningText:                    sp("#f92672"),
		TextSpecial1:                   sp("#ae81ff"),
		TextSpecial2:                   sp("#a6e22e"),
		TopBarBackground:               sp("#f92672"),
		TopBarText:                     sp("#f8f8f2"),
		StatusBarBackground:            sp("#f92672"),
		StatusBarText:                  sp("#f8f8f2"),
		StatusBarViewBackground:        sp("#ae81ff"),
		StatusBarViewText:              sp("#f8f8f2"),
		CommandText:                    sp("#f8f8f2"),
		ListSelectedBackground:         sp("#f92672"),
		ListSelectedText:               sp("#f8f8f2"),
		ListSelectedInactiveBackground: sp("#ae81ff"),
		ListSelectedInactiveText:       sp("#f8f8f2"),
		ControlsText:                   sp("#f8f8f2"),
		ControlsHighlight:              sp("#a6e22e"),
		AutocompleteBackground:         sp("#272822"),
		AutocompleteText:               sp("#f8f8f2"),
		AutocompleteSelectedBackground: sp("#ae81ff"),
		AutocompleteSelectedText:       sp("#f8f8f2"),
		ButtonColorOne:                 sp("#f92672"),
		ButtonColorTwo:                 sp("#272822"),
		TimelineNameBackground:         sp("#272822"),
		TimelineNameText:               sp("#808080"),
	},
	Media: MediaTOML{
		Image: &ViewerTOML{
			Program:  sp("xdg-open"),
			Args:     sp(""),
			Terminal: bf,
			Single:   bt,
			Reverse:  bf,
		},
		Video: &ViewerTOML{
			Program:  sp("xdg-open"),
			Args:     sp(""),
			Terminal: bf,
			Single:   bt,
			Reverse:  bf,
		},
		Audio: &ViewerTOML{
			Program:  sp("xdg-open"),
			Args:     sp(""),
			Terminal: bf,
			Single:   bt,
			Reverse:  bf,
		},
		Link: &ViewerTOML{
			Program:  sp("xdg-open"),
			Args:     sp(""),
			Terminal: bf,
			Single:   bt,
			Reverse:  bf,
		},
	},
	NotificationConfig: NotificationsTOML{
		Followers: bf,
		Favorite:  bf,
		Mention:   bf,
		Update:    bf,
		Boost:     bf,
		Poll:      bf,
		Posts:     bf,
	},
	Input: InputTOML{
		GlobalDown: &KeyHintTOML{
			Keys:        &[]string{"j", "J"},
			SpecialKeys: &[]string{"Down"},
		},
		GlobalUp: &KeyHintTOML{
			Keys:        &[]string{"k", "K"},
			SpecialKeys: &[]string{"Up"},
		},
		GlobalEnter: &KeyHintTOML{
			SpecialKeys: &[]string{"Enter"},
		},
		GlobalBack: &KeyHintTOML{
			Hint:        sp("[Esc]"),
			SpecialKeys: &[]string{"Esc"},
		},
		GlobalExit: &KeyHintTOML{
			Hint: sp("[Q]uit"),
			Keys: &[]string{"q", "Q"},
		},
		MainHome: &KeyHintTOML{
			Hint:        sp(""),
			Keys:        &[]string{"g"},
			SpecialKeys: &[]string{"Home"},
		},
		MainEnd: &KeyHintTOML{
			Hint:        sp(""),
			Keys:        &[]string{"G"},
			SpecialKeys: &[]string{"End"},
		},
		MainPrevFeed: &KeyHintTOML{
			Hint:        sp(""),
			Keys:        &[]string{"h", "H"},
			SpecialKeys: &[]string{"Left"},
		},
		MainNextFeed: &KeyHintTOML{
			Hint:        sp(""),
			Keys:        &[]string{"l", "L"},
			SpecialKeys: &[]string{"Right"},
		},
		MainPrevPane: &KeyHintTOML{
			Hint:        sp(""),
			SpecialKeys: &[]string{"Backtab"},
		},
		MainNextPane: &KeyHintTOML{
			Hint:        sp(""),
			SpecialKeys: &[]string{"Tab"},
		},
		MainCompose: &KeyHintTOML{
			Hint: sp(""),
			Keys: &[]string{"c", "C"},
		},
		StatusAvatar: &KeyHintTOML{
			Hint: sp("[A]vatar"),
			Keys: &[]string{"a", "A"},
		},
		StatusBoost: &KeyHintTOML{
			Hint:    sp("[B]oost"),
			HintAlt: sp("Un[B]oost"),
			Keys:    &[]string{"b", "B"},
		},
		StatusEdit: &KeyHintTOML{
			Hint: sp("[E]dit"),
			Keys: &[]string{"e", "E"},
		},
		StatusDelete: &KeyHintTOML{
			Hint: sp("[D]elete"),
			Keys: &[]string{"d", "D"},
		},
		StatusFavorite: &KeyHintTOML{
			Hint:    sp("[F]avorite"),
			HintAlt: sp("Un[F]avorite"),
			Keys:    &[]string{"f", "F"},
		},
		StatusMedia: &KeyHintTOML{
			Hint: sp("[M]edia"),
			Keys: &[]string{"m", "M"},
		},
		StatusLinks: &KeyHintTOML{
			Hint: sp("[O]pen"),
			Keys: &[]string{"o", "O"},
		},
		StatusPoll: &KeyHintTOML{
			Hint: sp("[P]oll"),
			Keys: &[]string{"p", "P"},
		},
		StatusReply: &KeyHintTOML{
			Hint: sp("[R]eply"),
			Keys: &[]string{"r", "R"},
		},
		StatusBookmark: &KeyHintTOML{
			Hint:    sp("[S]ave"),
			HintAlt: sp("Un[S]ave"),
			Keys:    &[]string{"s", "S"},
		},
		StatusThread: &KeyHintTOML{
			Hint: sp("[T]hread"),
			Keys: &[]string{"t", "T"},
		},
		StatusUser: &KeyHintTOML{
			Hint: sp("[U]ser"),
			Keys: &[]string{"u", "U"},
		},
		StatusViewFocus: &KeyHintTOML{
			Hint: sp("[V]iew"),
			Keys: &[]string{"v", "V"},
		},
		StatusYank: &KeyHintTOML{
			Hint: sp("[Y]ank"),
			Keys: &[]string{"y", "Y"},
		},
		StatusToggleCW: &KeyHintTOML{
			Hint: sp("Press [Z] to toggle cw"),
			Keys: &[]string{"z", "Z"},
		},
		StatusShowFiltered: &KeyHintTOML{
			Hint: sp("Press [Z] to view filtered toot"),
			Keys: &[]string{"z", "Z"},
		},
		UserAvatar: &KeyHintTOML{
			Hint: sp("[A]vatar"),
			Keys: &[]string{"a", "A"},
		},
		UserBlock: &KeyHintTOML{
			Hint:    sp("[B]lock"),
			HintAlt: sp("Un[B]lock"),
			Keys:    &[]string{"b", "B"},
		},
		UserFollow: &KeyHintTOML{
			Hint:    sp("[F]ollow"),
			HintAlt: sp("Un[F]ollow"),
			Keys:    &[]string{"f", "F"},
		},
		UserFollowRequestDecide: &KeyHintTOML{
			Hint:    sp("Follow [R]equest"),
			HintAlt: sp("Follow [R]equest"),
			Keys:    &[]string{"r", "R"},
		},
		UserMute: &KeyHintTOML{
			Hint:    sp("[M]ute"),
			HintAlt: sp("Un[M]ute"),
			Keys:    &[]string{"m", "M"},
		},
		UserLinks: &KeyHintTOML{
			Hint: sp("[O]pen"),
			Keys: &[]string{"o", "O"},
		},
		UserUser: &KeyHintTOML{
			Hint: sp("[U]ser"),
			Keys: &[]string{"u", "U"},
		},
		UserViewFocus: &KeyHintTOML{
			Hint: sp("[V]iew"),
			Keys: &[]string{"v", "V"},
		},
		UserYank: &KeyHintTOML{
			Hint: sp("[Y]ank"),
			Keys: &[]string{"y", "Y"},
		},
		ListOpenFeed: &KeyHintTOML{
			Hint: sp("[O]pen"),
			Keys: &[]string{"o", "O"},
		},
		ListUserList: &KeyHintTOML{
			Hint: sp("[U]sers"),
			Keys: &[]string{"u", "U"},
		},
		ListUserAdd: &KeyHintTOML{
			Hint: sp("[A]dd"),
			Keys: &[]string{"a", "A"},
		},
		ListUserDelete: &KeyHintTOML{
			Hint: sp("[D]elete"),
			Keys: &[]string{"d", "D"},
		},
		LinkOpen: &KeyHintTOML{
			Hint: sp("[O]pen"),
			Keys: &[]string{"o", "O"},
		},
		LinkYank: &KeyHintTOML{
			Hint: sp("[Y]ank"),
			Keys: &[]string{"y", "Y"},
		},
		TagOpenFeed: &KeyHintTOML{
			Hint: sp("[O]pen"),
			Keys: &[]string{"o", "O"},
		},
		TagFollow: &KeyHintTOML{
			Hint:    sp("[F]ollow"),
			HintAlt: sp("Un[F]ollow"),
			Keys:    &[]string{"f", "F"},
		},
		ComposeEditCW: &KeyHintTOML{
			Hint: sp("[C]W text"),
			Keys: &[]string{"c", "C"},
		},
		ComposeEditText: &KeyHintTOML{
			Hint: sp("[E]dit text"),
			Keys: &[]string{"e", "E"},
		},
		ComposeIncludeQuote: &KeyHintTOML{
			Hint: sp("[I]nclude quote"),
			Keys: &[]string{"i", "I"},
		},
		ComposeMediaFocus: &KeyHintTOML{
			Hint: sp("[M]edia"),
			Keys: &[]string{"m", "M"},
		},
		ComposePost: &KeyHintTOML{
			Hint: sp("[P]ost"),
			Keys: &[]string{"p", "P"},
		},
		ComposeToggleContentWarning: &KeyHintTOML{
			Hint: sp("[T]oggle CW"),
			Keys: &[]string{"t", "T"},
		},
		ComposeVisibility: &KeyHintTOML{
			Hint: sp("[V]isibility"),
			Keys: &[]string{"v", "V"},
		},
		ComposeLanguage: &KeyHintTOML{
			Hint: sp("[L]ang"),
			Keys: &[]string{"l", "L"},
		},
		ComposePoll: &KeyHintTOML{
			Hint: sp("P[O]ll"),
			Keys: &[]string{"o", "O"},
		},
		MediaDelete: &KeyHintTOML{
			Hint: sp("[D]elete"),
			Keys: &[]string{"d", "D"},
		},
		MediaEditDesc: &KeyHintTOML{
			Hint: sp("[E]dit desc"),
			Keys: &[]string{"e", "E"},
		},
		MediaAdd: &KeyHintTOML{
			Hint: sp("[A]dd"),
			Keys: &[]string{"a", "A"},
		},
		VoteVote: &KeyHintTOML{
			Hint: sp("[V]ote"),
			Keys: &[]string{"v", "V"},
		},
		VoteSelect: &KeyHintTOML{
			Hint:        sp("[Enter] to select"),
			Keys:        &[]string{" "},
			SpecialKeys: &[]string{"Enter"},
		},
		PollAdd: &KeyHintTOML{
			Hint: sp("[A]dd"),
			Keys: &[]string{"a", "A"},
		},
		PollEdit: &KeyHintTOML{
			Hint: sp("[E]dit"),
			Keys: &[]string{"e", "E"},
		},
		PollDelete: &KeyHintTOML{
			Hint: sp("[D]elete"),
			Keys: &[]string{"d", "D"},
		},
		PollMultiToggle: &KeyHintTOML{
			Hint: sp("Toggle [M]ultiple"),
			Keys: &[]string{"m", "M"},
		},
		PollExpiration: &KeyHintTOML{
			Hint: sp("E[X]pires"),
			Keys: &[]string{"x", "X"},
		},
		PreferenceName: &KeyHintTOML{
			Hint: sp("[N]ame"),
			Keys: &[]string{"n", "N"},
		},
		PreferenceVisibility: &KeyHintTOML{
			Hint: sp("[V]isibility"),
			Keys: &[]string{"v", "V"},
		},
		PreferenceBio: &KeyHintTOML{
			Hint: sp("[B]io"),
			Keys: &[]string{"b", "B"},
		},
		PreferenceSave: &KeyHintTOML{
			Hint: sp("[S]ave"),
			Keys: &[]string{"s", "S"},
		},
		PreferenceFields: &KeyHintTOML{
			Hint: sp("[F]ields"),
			Keys: &[]string{"f", "F"},
		},
		PreferenceFieldsAdd: &KeyHintTOML{
			Hint: sp("[A]dd"),
			Keys: &[]string{"a", "A"},
		},
		PreferenceFieldsEdit: &KeyHintTOML{
			Hint: sp("[E]dit"),
			Keys: &[]string{"e", "E"},
		},
		PreferenceFieldsDelete: &KeyHintTOML{
			Hint: sp("[D]elete"),
			Keys: &[]string{"d", "D"},
		},
	},
}
