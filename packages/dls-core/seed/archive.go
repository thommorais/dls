package seed

import "dls/dls-core/domain"

// archiveEntries is the knowledge base: what the channel is, what each
// running joke means, and the dates it remembers. A counter on the site
// points at its own glossary entry through moment_type, which is what turns
// a number into something a new viewer can read.
func archiveEntries() []ArchiveEntry {
	published := firstEpisode.AddDate(0, 0, -30)

	entries := []ArchiveEntry{
		{
			Slug:    "sobre-o-canal",
			Title:   "Sobre o " + ChannelName,
			Kind:    domain.ArchiveAbout,
			Summary: "Programa semanal de conversa, gravado às quartas, apresentado por Tico Nunes e Belinha Rocha.",
			Body: "O " + ChannelName + " (" + ChannelHandle + ") vai ao ar toda quarta-feira desde janeiro de 2025. " +
				"O formato nunca foi decidido: começou como uma conversa de vinte minutos sobre a semana e virou duas horas de assunto nenhum, " +
				"com convidados que às vezes nem chegam a ser entrevistados.\n\n" +
				"A produção é da Dona Neide, que responde perguntas gritadas do estúdio há mais de quarenta episódios, " +
				"e a câmera é do Léo, que aparece mais do que gostaria.\n\n" +
				"Este site conta o que acontece no programa: quantas vezes alguém saiu cantando, quantas perguntas a Dona Neide já respondeu, " +
				"quem chora mais e quais aberturas o público mandou.",
			Tags:        []string{"canal", "formato"},
			PublishedAt: published,
		},
		{
			Slug:        "como-as-estatisticas-sao-contadas",
			Title:       "Como as estatísticas são contadas",
			Kind:        domain.ArchiveAbout,
			Summary:     "O que entra na conta, o que fica de fora, e por que dois números parecidos não são a mesma coisa.",
			Body:        "Cada momento é uma linha: o episódio, o tipo, o segundo do vídeo e um resumo. Um tipo novo é um cadastro, não uma mudança de código.\n\nMulheres entrevistadas conta entrevistas, não pessoas: a mesma convidada em dois episódios conta duas vezes. Quem só aparece no estúdio e nunca é entrevistado não entra. Gênero não preenchido nunca é contado como mulher.\n\nAs palavras que puxam música são normalizadas em minúsculas, porque quem digita erra o caixa alta e o contador não deveria se importar.",
			Tags:        []string{"metodologia"},
			PublishedAt: published,
		},
		{
			Slug:           "a-regra-da-dona-neide",
			Title:          "A regra da Dona Neide",
			Kind:           domain.ArchiveSegment,
			Summary:        "Toda pergunta gritada para a produção vira estatística, mesmo quando ninguém responde.",
			Body:           "Nasceu no episódio 4, quando o Tico gritou \"Dona Neide, tem café?\" no meio de uma entrevista séria. A convidada esperou. A resposta veio do corredor: \"acabou\".\n\nDesde então vale a regra: se a pergunta é gritada e a Dona Neide está fora do quadro, conta.",
			Tags:           []string{"quadro", "dona neide"},
			MomentTypeSlug: "pergunta-pra-neide",
			PublishedAt:    published,
		},
		{
			Slug:           "saiu-cantando",
			Title:          "Saiu cantando",
			Kind:           domain.ArchiveGlossary,
			Summary:        "Uma palavra qualquer lembra uma música e o programa para.",
			Body:           "Funciona assim: alguém diz \"pizza\", alguém emenda uma música que tem pizza na letra, e os dois cantam até esquecerem o assunto. A palavra que puxou é registrada junto, e é por isso que o ranking de palavras existe.\n\nO recorde é do episódio 17: dezoito músicas em duas horas, nenhuma terminada.",
			Tags:           []string{"música", "quadro"},
			MomentTypeSlug: "saiu-cantando",
			PublishedAt:    published,
		},
		{
			Slug:           "teoria-maluca",
			Title:          "Teoria maluca",
			Kind:           domain.ArchiveGlossary,
			Summary:        "Explicação sem fonte, defendida com convicção total.",
			Body:           "Uma teoria maluca precisa de três coisas: nenhuma fonte, uma certeza absoluta e pelo menos trinta segundos de defesa. A mais citada até hoje é a do elevador que só demora quando alguém está atrasado.",
			Tags:           []string{"quadro"},
			MomentTypeSlug: "teoria-maluca",
			PublishedAt:    published,
		},
		{
			Slug:           "briga-de-bar",
			Title:          "Briga de bar",
			Kind:           domain.ArchiveGlossary,
			Summary:        "Discussão que começa no assunto do dia e termina em outro.",
			Body:           "O critério é o desvio: se a discussão começou sobre o tema do episódio e terminou sobre comida, transporte ou família, é briga de bar. O episódio 23 tem sete, um recorde que ninguém reivindica.",
			Tags:           []string{"quadro"},
			MomentTypeSlug: "briga-de-bar",
			PublishedAt:    published,
		},
		{
			Slug:           "chorou-ao-vivo",
			Title:          "Chorou ao vivo",
			Kind:           domain.ArchiveGlossary,
			Summary:        "Raro, sempre inesperado, quase sempre por história de infância.",
			Body:           "Conta quando é visível no vídeo. Olho marejado sem lágrima não conta, e essa regra já foi contestada duas vezes no próprio programa.",
			Tags:           []string{"quadro"},
			MomentTypeSlug: "chorou-ao-vivo",
			PublishedAt:    published,
		},
		{
			Slug:           "propaganda-nao-paga",
			Title:          "Propaganda não paga",
			Kind:           domain.ArchiveGlossary,
			Summary:        "Elogio longo e espontâneo a uma marca que nunca patrocinou nada.",
			Body:           "Precisa passar de um minuto e não pode ter contrato. A padaria da esquina lidera com folga e nunca mandou um pão.",
			Tags:           []string{"quadro", "patrocínio"},
			MomentTypeSlug: "propaganda-nao-paga",
			PublishedAt:    published,
		},
		{
			Slug:        "aberturas-do-publico",
			Title:       "Aberturas do público",
			Kind:        domain.ArchiveSegment,
			Summary:     "Desde o episódio 2, a abertura é sempre de alguém que assiste.",
			Body:        "Qualquer pessoa manda uma abertura. Entra no ar sem edição, com o nome de quem fez, e fica no acervo com gênero e a data em que foi ao ar.\n\nJá entraram forró, metal, chiptune, maracatu e uma feita inteiramente com buzina de caminhão.",
			Tags:        []string{"público", "abertura"},
			PublishedAt: published,
		},
		{
			Slug:        "o-quadro-da-cadeira-vazia",
			Title:       "O quadro da cadeira vazia",
			Kind:        domain.ArchiveSegment,
			Summary:     "Quando o convidado não aparece, a cadeira fica e a entrevista acontece do mesmo jeito.",
			Body:        "Estreou no episódio 26, por acidente. O convidado avisou em cima da hora, e os dois decidiram entrevistar a cadeira. Voltou outras vezes por pedido do público.",
			Tags:        []string{"quadro"},
			PublishedAt: published,
		},
		{
			Slug:        "dona-neide",
			Title:       "Dona Neide",
			Kind:        domain.ArchiveBio,
			Summary:     "Produtora do programa e a pessoa mais perguntada do estúdio.",
			Body:        "Produz o " + ChannelName + " desde o primeiro episódio. Não aparece no quadro por escolha própria, responde tudo do corredor e já resolveu ao vivo um problema de energia que quase derrubou a gravação.",
			Tags:        []string{"equipe"},
			PersonSlug:  "dona-neide",
			PublishedAt: published,
		},
		{
			Slug:        "leo-camera",
			Title:       "Léo Câmera",
			Kind:        domain.ArchiveBio,
			Summary:     "Operador de câmera, participante involuntário.",
			Body:        "Entrou para cobrir uma semana e ficou. Aparece no quadro toda vez que ri alto, o que acontece com frequência suficiente para o público reconhecer a risada.",
			Tags:        []string{"equipe"},
			PersonSlug:  "leo-camera",
			PublishedAt: published,
		},
		{
			Slug:        "o-primeiro-episodio",
			Title:       "O primeiro episódio",
			Kind:        domain.ArchiveMilestone,
			Summary:     "8 de janeiro de 2025, gravado no corredor porque o estúdio estava ocupado.",
			Body:        "Duração de quarenta minutos, uma câmera, nenhum convidado. A abertura era uma música de biblioteca gratuita, e é o único episódio sem abertura do público.",
			Tags:        []string{"história"},
			EpisodeSlug: "ep-001",
			PublishedAt: firstEpisode,
		},
		{
			Slug:        "o-recorde-das-dezoito-musicas",
			Title:       "O recorde das dezoito músicas",
			Kind:        domain.ArchiveMilestone,
			Summary:     "Episódio 17: dezoito músicas em duas horas, nenhuma cantada até o fim.",
			Body:        "Começou com a palavra \"saudade\" nos primeiros cinco minutos e nunca se recuperou. É o recorde que o programa mais tenta bater e o único que a produção pediu para não baterem de novo.",
			Tags:        []string{"história", "recorde"},
			EpisodeSlug: "ep-017",
			PublishedAt: firstEpisode.AddDate(0, 0, 7*16),
		},
		{
			Slug:        "o-apagao",
			Title:       "O apagão",
			Kind:        domain.ArchiveTrivia,
			Summary:     "O episódio 9 foi gravado com luz de celular depois que o prédio inteiro ficou sem energia.",
			Body:        "Dezessete minutos no escuro, gravados com dois celulares apoiados em copos. O áudio sobreviveu, o vídeo não. Foi o episódio mais assistido do primeiro semestre.",
			Tags:        []string{"história"},
			EpisodeSlug: "ep-009",
			PublishedAt: firstEpisode.AddDate(0, 0, 7*8),
		},
		{
			Slug:        "o-bolo-errado",
			Title:       "O bolo errado",
			Kind:        domain.ArchiveTrivia,
			Summary:     "No especial de aniversário, a confeitaria escreveu o nome do canal errado.",
			Body:        "Chegou escrito \"CAFÉ COM CAUS\". Ninguém corrigiu, e a grafia virou a piada do chat por três meses.",
			Tags:        []string{"história"},
			PublishedAt: firstEpisode.AddDate(0, 0, 7*11),
		},
		{
			Slug:        "por-que-quarta-feira",
			Title:       "Por que quarta-feira",
			Kind:        domain.ArchiveTrivia,
			Summary:     "A escolha do dia foi feita por eliminação, e nunca mudou.",
			Body:        "Segunda todo mundo está de mau humor, terça ninguém lembra, quinta já é quase sexta e sexta ninguém assiste. Sobrou quarta.",
			Tags:        []string{"formato"},
			PublishedAt: published,
		},
	}

	return entries
}
