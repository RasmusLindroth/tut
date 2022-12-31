package api

func (ac *AccountClient) GetCharLimit() int {
	if ac.Instance != nil {
		return ac.Instance.Configuration.Statuses.MaxCharacters
	}
	if ac.InstanceOld == nil || ac.InstanceOld.Configuration == nil || ac.InstanceOld.Configuration.Statuses == nil {
		return 500
	}
	s := ac.InstanceOld.Configuration.Statuses
	if val, ok := (*s)["max_characters"]; ok {
		switch v := val.(type) {
		case int:
			return v
		}
	}
	return 500
}

func (ac *AccountClient) GetLengthURL() int {
	if ac.Instance != nil {
		return ac.Instance.Configuration.Statuses.CharactersReservedPerURL
	}
	if ac.InstanceOld == nil || ac.InstanceOld.Configuration == nil || ac.InstanceOld.Configuration.Statuses == nil {
		return 23
	}
	s := ac.InstanceOld.Configuration.Statuses
	if val, ok := (*s)["characters_reserved_per_url"]; ok {
		switch v := val.(type) {
		case int:
			return v
		}
	}
	return 23
}

func (ac *AccountClient) GetPollOptions() (options, chars int) {
	if ac.Instance != nil {
		return ac.Instance.Configuration.Polls.MaxOptions, ac.Instance.Configuration.Polls.MaxCharactersPerOption
	}
	if ac.InstanceOld == nil || ac.InstanceOld.Configuration == nil || ac.InstanceOld.Configuration.Polls == nil {
		return 4, 50
	}
	s := ac.InstanceOld.Configuration.Polls
	opts, okOne := (*s)["max_options"]
	c, okTwo := (*s)["max_characters_per_option"]
	if okOne && okTwo {
		a, b := 4, 50
		switch v := opts.(type) {
		case int:
			a = v
		}
		switch v := c.(type) {
		case int:
			b = v
		}
		return a, b
	}
	return 4, 50
}
