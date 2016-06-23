package blog

func rssHandler(handler *Handlder) {
	handler.renderXML(handler.HTMLFile, map[string]interface{}{
		"items": getRSSItems(handler)})
}
