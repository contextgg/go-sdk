package es

// EventMatcher is a func that can match event to a criteria.
type EventMatcher func(*Event) bool

// MatchAny matches any event.
func MatchAny() EventMatcher {
	return func(*Event) bool {
		return true
	}
}

// MatchEvent matches a specific event type, nil events never match.
func MatchEvent(t string) EventMatcher {
	return func(e *Event) bool {
		return e != nil && e.Type == t
	}
}

// MatchAnyEventOf matches if any of several matchers matches.
func MatchAnyEventOf(types ...interface{}) EventMatcher {
	all := make(map[string]interface{})
	for _, t := range types {
		_, name := GetTypeName(t)
		all[name] = t
	}

	return func(e *Event) bool {
		for _, t := range types {
			_, name := GetTypeName(t)
			if MatchEvent(name)(e) {
				return true
			}
		}
		return false
	}
}
