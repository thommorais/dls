// Package seed holds the development dataset. Content comes from the real
// show's own records (see apps/site/public for the source spreadsheet) as it
// is transcribed; until then, this stays empty rather than inventing numbers.
package seed

import "dls/dls-core/domain"

// The channel. Everything below belongs to it.
const (
	ChannelName   = "Desce a Letra Show"
	ChannelHandle = "@descealetrashow"
)

// Each type here is backed by a real sheet in the source spreadsheet
// (apps/site/public), transcribed from the show's own episodes.
var momentTypes = []domain.MomentType{
	{Slug: "saiu-cantando", Label: "Saiu cantando", Color: "#10B981", Position: 1,
		Description: "Uma palavra ou assunto qualquer lembra uma música e o episódio para para cantar."},
	{Slug: "gags-e-manias", Label: "Gags & manias", Color: "#F59E0B", Position: 2,
		Description: "Os cacoetes recorrentes do estúdio: intervenções do Bulbasauro, memes do Load e afins."},
	{Slug: "comidas-e-bebidas", Label: "Comidas & bebidas", Color: "#06B6D4", Position: 3,
		Description: "O que a bancada come, bebe ou defende comer em pleno programa."},
	{Slug: "presentes-e-mimos", Label: "Presentes & mimos", Color: "#8B5CF6", Position: 4,
		Description: "O que o público manda pro estúdio e como a bancada reage."},
	{Slug: "mencoes-pop", Label: "Menções de cultura pop", Color: "#F43F5E", Position: 5,
		Description: "Craque Neto, celebridades e outras figuras que viram pauta no programa."},
}

var people = []domain.Person{}

type Song struct {
	Key    string
	Title  string
	Artist string
	Genre  string
	Year   int
}

var songs = []Song{}
