package blog

//	"fmt"

type TMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type TArticleMsg struct {
	TMessage
	InsertId int64 `json:"id"`
}

func indexHandler(handler *Handlder) {
	handler.renderTemplate(handler.HTMLFile, BASEHTML, map[string]interface{}{
		"aboutMe":    queryAboutMe(handler),
		"titles":     queryTitlesTop10(handler),
		"articles":   queryArticlesTop5(handler),
		"categroys":  queryCategroys(handler),
		"showsearch": true,
		"showPages":  true})
}

func commonHandler(handler *Handlder) {
	handler.renderTemplate(handler.HTMLFile, ADMINHTML, map[string]interface{}{})
}

func editorHandler(handler *Handlder) {
	handler.renderTemplate(handler.HTMLFile, BASEHTML, map[string]interface{}{
		"categroys": queryCategroys(handler)})
}

func aboutMeHandler(handler *Handlder) {
	handler.renderTemplate(handler.HTMLFile, BASEHTML, map[string]interface{}{
		"aboutMe":    queryAboutMe(handler),
		"showsearch": false})
}

func articleHandler(handler *Handlder) {
	id := stoi(handler.param("id"))
	if id == -1 {
		//			handler.redirect()
		return
	}
	article := queryArticleById(handler, id)
	if article == nil {
		handler.redirectError()
		return
	}
	handler.renderTemplate(handler.HTMLFile, BASEHTML, map[string]interface{}{
		"article":    article,
		"comments":   queryArticleCommentsById(handler /*id*/, 13),
		"showsearch": false})
}

func errorHandler(handler *Handlder) {
	handler.render(handler.HTMLFile, map[string]interface{}{})
}

func robotsHandler(handler *Handlder) {
	handler.render(handler.HTMLFile)
}
