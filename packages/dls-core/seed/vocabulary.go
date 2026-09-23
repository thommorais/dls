// Package seed holds the development dataset: a fake channel with enough
// history that the charts have something to say. It is deliberately not
// about any real show, and it is deterministic, so two people running the
// seed see the same numbers and can talk about the same episode.
package seed

import "dls/dls-core/domain"

// The channel. Everything below belongs to it.
const (
	ChannelName   = "CAFÉ COM CAOS"
	ChannelHandle = "@cafecomcaos"
)

var momentTypes = []domain.MomentType{
	{Slug: "pergunta-pra-neide", Label: "Pergunta pra Dona Neide", Color: "#6366F1", Position: 1,
		Description: "Toda vez que alguém no estúdio grita uma pergunta para a Dona Neide, que produz o programa e sabe de tudo."},
	{Slug: "saiu-cantando", Label: "Saiu cantando", Color: "#10B981", Position: 2,
		Description: "Uma palavra qualquer lembra uma música e o episódio para por dois minutos."},
	{Slug: "teoria-maluca", Label: "Teoria maluca", Color: "#F59E0B", Position: 3,
		Description: "Uma explicação sem nenhuma fonte, defendida com absoluta convicção."},
	{Slug: "briga-de-bar", Label: "Briga de bar", Color: "#F43F5E", Position: 4,
		Description: "Discussão que começa sobre o assunto do dia e termina sobre outra coisa."},
	{Slug: "chorou-ao-vivo", Label: "Chorou ao vivo", Color: "#8B5CF6", Position: 5,
		Description: "Raro, sempre inesperado, sempre por causa de história de infância."},
	{Slug: "propaganda-nao-paga", Label: "Propaganda não paga", Color: "#06B6D4", Position: 6,
		Description: "Elogio espontâneo e longo a uma marca que nunca patrocinou nada."},
}

var people = []domain.Person{
	{Slug: "tico", Name: "Tico Nunes", Kind: domain.PersonHost, Gender: domain.GenderMan},
	{Slug: "belinha", Name: "Belinha Rocha", Kind: domain.PersonHost, Gender: domain.GenderWoman},
	{Slug: "dona-neide", Name: "Dona Neide", Kind: domain.PersonStaff, Gender: domain.GenderWoman},
	{Slug: "leo-camera", Name: "Léo Câmera", Kind: domain.PersonStaff, Gender: domain.GenderMan},

	{Slug: "rita-sampaio", Name: "Rita Sampaio", Kind: domain.PersonGuest, Gender: domain.GenderWoman},
	{Slug: "juca-bala", Name: "Juca Bala", Kind: domain.PersonGuest, Gender: domain.GenderMan},
	{Slug: "mari-tavares", Name: "Mari Tavares", Kind: domain.PersonGuest, Gender: domain.GenderWoman},
	{Slug: "bruno-kiko", Name: "Bruno Kiko", Kind: domain.PersonGuest, Gender: domain.GenderMan},
	{Slug: "nanda-prado", Name: "Nanda Prado", Kind: domain.PersonGuest, Gender: domain.GenderWoman},
	{Slug: "seu-valdo", Name: "Seu Valdo", Kind: domain.PersonGuest, Gender: domain.GenderMan},
	{Slug: "cacau-lins", Name: "Cacau Lins", Kind: domain.PersonGuest, Gender: domain.GenderNonBinary},
	{Slug: "dani-moraes", Name: "Dani Moraes", Kind: domain.PersonGuest, Gender: domain.GenderWoman},
	{Slug: "toninho-jazz", Name: "Toninho Jazz", Kind: domain.PersonGuest, Gender: domain.GenderMan},
	{Slug: "paty-figueiredo", Name: "Paty Figueiredo", Kind: domain.PersonGuest, Gender: domain.GenderWoman},
	{Slug: "gordo-silveira", Name: "Gordo Silveira", Kind: domain.PersonGuest, Gender: domain.GenderMan},
	{Slug: "iara-melo", Name: "Iara Melo", Kind: domain.PersonGuest, Gender: domain.GenderWoman},
	{Slug: "zeca-do-pneu", Name: "Zeca do Pneu", Kind: domain.PersonGuest, Gender: domain.GenderMan},
	{Slug: "lu-andrade", Name: "Lu Andrade", Kind: domain.PersonGuest, Gender: domain.GenderWoman},
	// Recorded before anyone thought to ask, and never counted as a woman.
	{Slug: "convidado-misterioso", Name: "Convidado Misterioso", Kind: domain.PersonGuest, Gender: domain.GenderUnknown},
}

type Song struct {
	Key    string
	Title  string
	Artist string
	Genre  string
	Year   int
}

var songs = []Song{
	{Key: "evidencias", Title: "Evidências", Artist: "Chitãozinho & Xororó", Genre: "sertanejo", Year: 1990},
	{Key: "anunciacao", Title: "Anunciação", Artist: "Alceu Valença", Genre: "mpb", Year: 1983},
	{Key: "deixa-a-vida", Title: "Deixa a Vida Me Levar", Artist: "Zeca Pagodinho", Genre: "samba", Year: 2002},
	{Key: "malandragem", Title: "Malandragem", Artist: "Cássia Eller", Genre: "rock", Year: 1994},
	{Key: "sozinho", Title: "Sozinho", Artist: "Caetano Veloso", Genre: "mpb", Year: 1998},
	{Key: "eduardo-e-monica", Title: "Eduardo e Mônica", Artist: "Legião Urbana", Genre: "rock", Year: 1986},
	{Key: "tempo-perdido", Title: "Tempo Perdido", Artist: "Legião Urbana", Genre: "rock", Year: 1986},
	{Key: "garota-de-ipanema", Title: "Garota de Ipanema", Artist: "Tom Jobim", Genre: "bossa nova", Year: 1962},
	{Key: "e-o-amor", Title: "É o Amor", Artist: "Zezé Di Camargo & Luciano", Genre: "sertanejo", Year: 1991},
	{Key: "vaca-profana", Title: "Vaca Profana", Artist: "Gal Costa", Genre: "mpb", Year: 1984},
	{Key: "trem-bala", Title: "Trem-Bala", Artist: "Ana Vilela", Genre: "pop", Year: 2017},
	{Key: "country-roads", Title: "Country Roads", Artist: "John Denver", Genre: "folk", Year: 1971},
	{Key: "bohemian", Title: "Bohemian Rhapsody", Artist: "Queen", Genre: "rock", Year: 1975},
	{Key: "o-portao", Title: "O Portão", Artist: "Roberto Carlos", Genre: "romântico", Year: 1974},
}

// The words that set them off. Half the fun of the stat is that they are
// completely ordinary.
var triggerWords = []string{
	"pizza", "saudade", "segunda-feira", "boleto", "cachorro", "chuva",
	"aniversário", "farofa", "dentista", "karaokê", "futebol", "café",
	"elevador", "praia", "aluguel", "avião", "sogra", "churrasco",
	"feriado", "vizinho", "academia", "pix",
}

var episodeTitles = []string{
	"O primeiro episódio, gravado no corredor",
	"Ninguém sabia que o microfone estava ligado",
	"A entrevista que virou terapia de casal",
	"Duas horas discutindo se pastel é salgado",
	"O dia em que a internet caiu no meio da live",
	"Especial: as piores ideias de negócio do chat",
	"A convidada trouxe bolo e mudou o programa",
	"Debate: existe fila justa em padaria?",
	"O episódio do apagão",
	"Alguém trouxe um cachorro e ninguém trabalhou",
	"A teoria do elevador que nunca sobe",
	"Especial de aniversário, com bolo errado",
	"O dia em que o Tico perdeu uma aposta",
	"Entrevista com quem já dormiu no aeroporto",
	"A briga sobre o jeito certo de fazer café",
	"Episódio gravado depois do churrasco",
	"O recorde: dezoito músicas em duas horas",
	"Quem paga a conta quando ninguém pediu entrada",
	"A convidada que sabia todas as respostas",
	"Especial: histórias de segunda-feira",
	"O dia em que a Belinha chorou com propaganda",
	"Duas horas sobre o preço do aluguel",
	"O episódio em que ninguém concordou com nada",
	"Entrevista interrompida por obra no prédio",
	"A lista definitiva de músicas de karaokê",
	"O dia em que o convidado não apareceu",
	"Especial de férias, gravado na praia",
	"A discussão sobre sogras que durou o programa inteiro",
	"O episódio do bolo de rolo",
	"Quem inventou a fila do banco",
	"A entrevista mais curta da história do canal",
	"O dia em que todo mundo cantou junto",
	"Especial: as melhores desculpas para faltar",
	"A teoria de que segunda-feira começa no domingo",
	"O episódio sem nenhum assunto",
	"Entrevista com quem já foi expulso de academia",
	"A noite em que o estúdio ficou sem luz",
	"O dia em que o Pix salvou o programa",
	"Especial de fim de ano, com todo mundo junto",
	"Cem episódios de conversa e nenhuma conclusão",
}

type openingSeed struct {
	Title  string
	Author string
	Handle string
	Genre  string
}

// The openings the audience sends in. The genre is what makes the library
// worth filtering.
var openings = []openingSeed{
	{"Abertura em forró", "Marta Queiroz", "@martaq", "forró"},
	{"Café com Caos Metal", "Igor Prates", "@igorprates", "metal"},
	{"Vinheta de novela", "Cida Barros", "@cidabarros", "trilha"},
	{"Versão bossa nova", "Tuca Lemos", "@tucalemos", "bossa nova"},
	{"Abertura de desenho animado", "Pedro Hax", "@pedrohax", "chiptune"},
	{"Rap da Dona Neide", "MC Bilú", "@mcbilu", "rap"},
	{"Abertura com bateria de escola de samba", "Bloco do Zé", "@blocodoze", "samba"},
	{"Trilha de faroeste", "Ana Clara Dias", "@anaclaradias", "country"},
	{"Versão gospel", "Coral São Jorge", "@coralsaojorge", "gospel"},
	{"Tema de videogame", "Juninho 8bits", "@juninho8bits", "chiptune"},
	{"Abertura em pagode", "Roda do Beco", "@rodadobeco", "samba"},
	{"Versão jazz de elevador", "Toninho Jazz", "@toninhojazz", "jazz"},
	{"Abertura brega-funk", "DJ Pipoca", "@djpipoca", "funk"},
	{"Tema épico de filme", "Estúdio Vento", "@estudiovento", "trilha"},
	{"Abertura cantada pelo chat", "Coletivo Ao Vivo", "@coletivoaovivo", "coral"},
	{"Vinheta de rádio AM", "Seu Nilton", "@seunilton", "retrô"},
	{"Abertura em axé", "Banda Pé na Areia", "@penaareia", "axé"},
	{"Versão lo-fi para estudar", "quietkid", "@quietkid", "lo-fi"},
	{"Abertura de telejornal", "Redação Falsa", "@redacaofalsa", "trilha"},
	{"Tema de luta livre", "Careca do Ringue", "@carecaringue", "rock"},
	{"Abertura em choro", "Regional da Esquina", "@regionalesquina", "choro"},
	{"Versão punk de garagem", "As Tampas", "@astampas", "punk"},
	{"Abertura com viola caipira", "Zé da Viola", "@zedaviola", "sertanejo raiz"},
	{"Tema de suspense", "Noite Longa", "@noitelonga", "trilha"},
	{"Abertura em reggae", "Praia Norte", "@praianorte", "reggae"},
	{"Versão ópera", "Tenor do Bairro", "@tenordobairro", "clássico"},
	{"Abertura feita no celular às 3h", "Insônia FM", "@insoniafm", "lo-fi"},
	{"Tema de programa de auditório", "Dona Glória", "@donagloria", "retrô"},
	{"Abertura em maracatu", "Baque do Recife", "@baquerecife", "maracatu"},
	{"Versão eletrônica", "Kaos Beats", "@kaosbeats", "eletrônica"},
	{"Abertura cantada pelo filho do Tico", "Tico Jr.", "@ticojr", "infantil"},
	{"Tema de abertura de podcast sério", "Estúdio Cinza", "@estudiocinza", "corporativo"},
	{"Abertura em bolero", "Trio Madrugada", "@triomadrugada", "bolero"},
	{"Versão a cappella", "Grupo Eco", "@grupoeco", "coral"},
	{"Abertura com buzina de caminhão", "Zeca do Pneu", "@zecadopneu", "experimental"},
	{"Tema de fim de ano", "Coral São Jorge", "@coralsaojorge", "coral"},
}
