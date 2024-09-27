package main

import "snippet-box/pkg/models"

// To set the holding structure for any dynamic data to be passed to HTML templates
type templateData struct {
	Snippet *models.Snippet
}
