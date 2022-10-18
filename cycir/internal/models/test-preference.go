package models

// AllPreferences returns a slice of preferences
func (repo *TestRepository) AllPreferences() ([]Preference, error) {
	var preferences []Preference
	return preferences, nil
}

// SetSystemPref updates a system preference setting
func (repo *TestRepository) SetSystemPref(name, value string) error {
	return nil
}

// UpdateSystemPref updates a system preference setting
func (repo *TestRepository) UpdateSystemPref(name, value string) error {
	return nil
}

// InsertOrUpdateSitePreferences inserts or updates all site prefs from map
func (repo *TestRepository) InsertOrUpdateSitePreferences(pm map[string]string) error {
	return nil
}
