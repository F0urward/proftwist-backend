package dto

func ModerationResultToDTO(allowed bool, categories []string) ModerationResult {
	return ModerationResult{
		Allowed:    allowed,
		Categories: categories,
	}
}

func EmptyModerationResult() ModerationResult {
	return ModerationResult{
		Allowed:    true,
		Categories: []string{},
	}
}

func (mr *ModerationResult) IsModeratedSuccessfully() bool {
	return mr.Allowed && len(mr.Categories) == 0
}
