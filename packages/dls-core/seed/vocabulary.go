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

// saiu-cantando is the only moment type backed by real data so far, sourced
// from the "Cantorias do DLShow" sheet.
var momentTypes = []domain.MomentType{
	{Slug: "saiu-cantando", Label: "Saiu cantando", Color: "#10B981", Position: 1,
		Description: "Uma palavra ou assunto qualquer lembra uma música e o episódio para para cantar."},
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

type openingSeed struct {
	Title  string
	Author string
	Handle string
	Genre  string
}

// The openings the audience sends in, from the "Aberturas dos Ouvintes"
// sheet.
var openings = []openingSeed{}
