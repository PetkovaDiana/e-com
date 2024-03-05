package dto

import (
	"github.com/google/uuid"
)

type CategoryData struct {
	Categories   []*Category `json:"data"`
	CategoryPath []*Category `json:"category_path"`
}

type Category struct {
	UUID        uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Image       string        `json:"image"`
	SubCategory []SubCategory `json:"sub_categories"`
}

func (p *Category) ImageMediaRoot(mediaRoot string) {
	if p.Image != "" {
		p.Image = mediaRoot + p.Image
	}
}

type SubCategory struct {
	UUID  uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Image string    `json:"image"`
}

func (p *SubCategory) ImageMediaRoot(mediaRoot string) {
	if p.Image != "" {
		p.Image = mediaRoot + p.Image
	}
}

type CategoryParams struct {
	UUID  uuid.UUID `json:"id"`
	Layer int       `json:"layer"`
}

type CategoryCharacteristic struct {
	UUID  string `json:"id"`
	Title string `json:"title"`
}
