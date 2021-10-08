package dto

// CreateShortUrl store received data for create short url
type CreateShortUrl struct {
	OriginalUrl string `json:"original_url" validate:"required,url"`
	ShortUrl    string `json:"short_url,omitempty"`
}
