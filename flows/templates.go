package flows

func EnumerateTemplateArray(templates []string, callback func(string)) {
	for _, template := range templates {
		callback(template)
	}
}

func RewriteTemplateArray(templates []string, rewrite func(string) string) {
	for t := range templates {
		templates[t] = rewrite(templates[t])
	}
}

func EnumerateTemplateTranslations(localization Localization, localizable Localizable, key string, callback func(string)) {
	for _, lang := range localization.Languages() {
		translations := localization.GetTranslations(lang)
		for _, tpl := range translations.GetTextArray(localizable.LocalizationUUID(), key) {
			callback(tpl)
		}
	}
}

func RewriteTemplateTranslations(localization Localization, localizable Localizable, key string, rewrite func(string) string) {
	for _, lang := range localization.Languages() {
		translations := localization.GetTranslations(lang)

		templates := translations.GetTextArray(localizable.LocalizationUUID(), key)
		rewritten := make([]string, len(templates))
		for t := range templates {
			rewritten[t] = rewrite(templates[t])
		}
		translations.SetTextArray(localizable.LocalizationUUID(), key, rewritten)
	}
}
