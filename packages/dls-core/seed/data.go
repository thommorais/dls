// Code generated from apps/site/public/desce-a-letra-show-dados-v7.xlsx. DO NOT EDIT BY HAND.
package seed

type rawEpisode struct {
	Number    int
	YouTubeID string
}

var rawEpisodes = []rawEpisode{
	{529, "aq5-YtgytqY"},
	{531, "Vme7qk9NECM"},
	{562, "4VjfkFmn0IE"},
	{600, "27f7VEe7QVQ"},
	{602, "q02_uyu2QKk"},
	{603, "0iimGpySeJk"},
	{607, "GcrrPgDv1Dc"},
	{608, "rOtFhFxfhTE"},
	{609, "DhTo34z3noI"},
	{611, "5MvEmsrFKsY"},
	{617, "2Z_N3p-Z2aM"},
	{619, "8YBemV5RDRQ"},
	{620, "5creXtJTZpc"},
	{621, "IBjoOiw0iY"},
	{623, "0wsdOI1RkWo"},
	{624, "CY7Ns0gEWNA"},
	{625, "mekYK_2RnA8"},
	{627, "M_Reb5WaIps"},
	{628, "sIOnDXwttW4"},
	{631, "lIhI99LbVvk"},
	{632, "zuvCdkUWr0w"},
	{635, "wvf_rHDNHdk"},
	{636, "-XTVEDJ7k1g"},
	{637, "jr4HNj9k8Nw"},
	{639, "ATx3j20Yxs0"},
	{640, "3nrug68QLdw"},
	{645, "oIJJvsXIOHs"},
	{647, "aSzszw95o8M"},
}

type rawOpening struct {
	EpisodeNumber int
	AtSeconds     int
	Title         string
	Author        string
	Handle        string
	Genre         string
}

var rawOpenings = []rawOpening{
	{531, 720, "Repente / Paródia (Estilo Caju & Castanha)", "Felipe Hildebrand Braguin (Braonha)", "", "Repente / Paródia (Estilo Caju & Castanha)"},
	{647, 930, "Beat Dark / Industrial Metal", "Felipe Pestana", "@pestanafelip", "Beat Dark / Industrial Metal"},
	{645, 1210, "Reggae Roots Original", "Pedro Zanotti", "", "Reggae Roots Original"},
	{640, 720, "System of a Down Paródia", "Gustavo Pinheiros", "", "System of a Down Paródia"},
	{639, 675, "Stoner Rock Hipnótico", "Begano", "@began01", "Stoner Rock Hipnótico"},
	{637, 1400, "Vinheta Temática 5ª Temp", "Ouvinte Colaborador", "", "Vinheta Temática 5ª Temp"},
	{636, 1133, "Vinheta Amostradinho", "Ouvinte Colaborador", "", "Vinheta Amostradinho"},
	{635, 1080, "Reggaeton Etéreo", "Will Tottaro", "", "Reggaeton Etéreo"},
	{609, 940, "Blues Tradicional com Gaita", "Vittor Tibo", "", "Blues Tradicional com Gaita"},
	{632, 850, "Eletro / Synth Experimental", "Douglas Assunção", "", "Eletro / Synth Experimental"},
	{631, 920, "Boom Bap / Rap Hip-Hop", "Ouvinte Beatmaker", "", "Boom Bap / Rap Hip-Hop"},
	{628, 1080, "Western / Country Rock", "Caio Simões", "", "Western / Country Rock"},
	{628, 900, "Acapella / Beatbox Fêmea", "Ouvinte Cantora", "", "Acapella / Beatbox Fêmea"},
	{625, 2040, "Vinheta Shrek / Hip-Hop", "Ouvinte Colaborador", "", "Vinheta Shrek / Hip-Hop"},
	{623, 1080, "Ukulele Acústico Folk", "Ângelo Branco", "", "Ukulele Acústico Folk"},
	{621, 960, "Rock Autoral do Paraná", "Cigarro Mata (Banda)", "", "Rock Autoral do Paraná"},
	{620, 900, "Grindcore / Black Metal", "Dedo Dorfo, André, Rodolfo & Bilo", "", "Grindcore / Black Metal"},
	{619, 2071, "Eletrônica / DLS Máfia", "Igor Lobo, Guilherme, Cadu, Beck & Brian", "", "Eletrônica / DLS Máfia"},
	{529, 960, "Samba / Pagode 'Lula Silva'", "Ouvinte Colaborador", "", "Samba / Pagode 'Lula Silva'"},
	{617, 1080, "Forró / Xote 'Açúcar Doce'", "Vittor Recife", "", "Forró / Xote 'Açúcar Doce'"},
	{611, 960, "Arrastapé / Baião Nordestino", "Diogo do Monte", "", "Arrastapé / Baião Nordestino"},
	{608, 2009, "Vinheta É o Tchan / Axé", "Ouvinte Colaborador", "", "Vinheta É o Tchan / Axé"},
	{607, 900, "Capella / Poesia Falada", "Vin Marques", "", "Capella / Poesia Falada"},
	{531, 2074, "Death Metal / Grindcore", "Ouvinte Metalhead", "", "Death Metal / Grindcore"},
	{562, 3316, "Reggaeton Urbano", "Ouvinte Latino", "", "Reggaeton Urbano"},
}

type rawMoment struct {
	EpisodeNumber int
	TypeSlug      string
	AtSeconds     int
	Summary       string
}

var rawMoments = []rawMoment{
	{531, "saiu-cantando", 4980, "Leãozinho (Caetano Veloso / Armandinho) — cantada por Cauê Moura & Load Comics. Pauta sobre a prisão do influencer Wilker Leão. Cauê canta 'Gosto muito de te ver, leãozinho...' logo após comentar a prisão."},
	{635, "saiu-cantando", 8100, "Odisseia (Lucas Inutilismo) — cantada por Cauê Moura & Load. Comentário sobre o show do Lucas Inutilismo no Rock in Rio. Cauê imita os vocais emotivos e melódicos do metalcore de Lucas."},
	{600, "saiu-cantando", 300, "Rap Nacional / Sabotage / Facção Central — cantada por Cauê Moura & Load. Reação ao aviso de direitos autorais do YouTube. Cauê e Load cantam trechos de clássicos do rap nacional no estúdio."},
	{608, "saiu-cantando", 448, "Músicas do É o Tchan — cantada por Cauê & Load. Abertura e nostalgia dos anos 90. Sessão de dança da cordinha e nostalgia no início do programa."},
	{531, "saiu-cantando", 2074, "Gritos Death Metal / Gutural — cantada por Cauê Moura. Reação à vinheta enviada por ouvinte. Cauê imita vocal gutural rasgado acompanhando a vinheta pesada."},
	{531, "saiu-cantando", 1413, "Riffs e Vocais Heavy Metal — cantada por Cauê Moura & Load. Vinheta com guitarras de metal. Interpretação cômica de rock pesado ao vivo na bancada."},
	{531, "gags-e-manias", 15, "Bulbasauro Intervention (Bulbasauro & Load): Problema no cabo XLR do Load. Buba ajeita o microfone do Load ao vivo e avisa que comprou 30m de cabo novo."},
	{531, "gags-e-manias", 90, "Cauê Late & Sono (Cauê Moura): Relato de insônia e cansaço ao vivo. Cauê conta que dormiu às 22:30, acordou 01:30 e depois dormiu até 07:45, chegando estafado."},
	{531, "gags-e-manias", 910, "Bulbasauro Intervention (Bulbasauro & Cauê): Ajuste de volume de retorno ao vivo. Cauê pede pro Buba abaixar o áudio de retorno porque o microfone do Load estava estourando."},
	{531, "gags-e-manias", 1260, "Bulbasauro Intervention (Bip) (Bulbasauro & Cauê): Bip de censura durante rants. Cauê pede abertamente: 'Solta o bip aí pra mim Buba' ao falar de figuras políticas/famosos."},
	{531, "gags-e-manias", 1980, "Perguntas pra Ana / Ciência (Cauê Moura & Load): Fungos, Microbioma e Frutas Picadas. Debate científico sobre bolor em alimentos, microbiota intestinal e o cofre mundial de sementes da Noruega."},
	{531, "gags-e-manias", 4680, "Gesto do Load ('Mamou Fantasma') (Load Comics & Cauê): Brincadeira de mamar o microfone. Load faz o famoso gesto de 'pagar uma gulosa' no microfone e Cauê brinca com a cena na câmera."},
	{637, "gags-e-manias", 0, "Bulbasauro Intervention (Bulbasauro & Bancada): Impressora 3D ao vivo no estúdio. Buba coloca a impressora 3D para rodar ao vivo imprimindo mimos vermelhos para membros."},
	{635, "gags-e-manias", 71, "Cauê Rant / Gultural (Cauê Moura): Grito gutural ao vivo. Cauê manda um vocal estilo death metal assustando a audiência com fone isolado."},
	{628, "gags-e-manias", 10, "Gags visuais do Bala (Balla & Bancada): Animação da vinheta da spraycan. Balla altera a arte da spraycan do Cauê e Load lançando fumaça e efeitos especiais na tela."},
	{627, "gags-e-manias", 0, "Gags da Poltrona / Cadeira (Cauê, Load & Gaiofato): Gag recorrente da escada / poltrona. Load senta numa escada/cadeira de baralho enquanto Cauê usa a poltrona oficial no teatro."},
	{600, "gags-e-manias", 10, "Load Senta na Escada (Load Comics): Gag da poltrona sem braço/escada. Load fica numa posição desconfortável no estúdio provocando piadas sobre sua coluna."},
	{625, "gags-e-manias", 58, "Load Senta na Cadeira Alta (Load Comics & Cauê): Load em estado de 'barril'. Load escolhe sentar numa cadeira alta estilo baralho e Cauê zoa que agora olha Load de cima."},
	{623, "gags-e-manias", 0, "Buba & Troca de Posição (Cauê, Load & Buba): Inversão de câmeras e enquadramento. Troca temporária da poltrona e enquadramento de frente para o Cauê."},
	{617, "gags-e-manias", 970, "Cauê Late / Re-opening (Cauê Moura & Load): Abertura de Novo (Restart de intro). Episódio teve duas aberturas gravadas no mesmo dia: [07:29] Abertura e [16:10] Abertura (de novo)."},
	{607, "gags-e-manias", 60, "Bulbasauro nas Picape (Bulbasauro): Buba como 'Mito das Picaps'. Buba dança e regula o mixer aumentando os efeitos sonoros de corneta e palmas."},
	{531, "gags-e-manias", 3066, "Bulbasauro Intervention (Bulbasauro): Bulba nega bolo pra velhinha. Quirk histórico do Buba sendo pautado e relembrado pela bancada."},
	{562, "gags-e-manias", 8125, "Meme do Load (Load Comics): O Load é ou não é Hipster?!. Debate de mais de 15 minutos analisando o estilo, roupas e hábitos alternativos do Load."},
	{531, "gags-e-manias", 1413, "Perguntas pra Ana / Ciência (Cauê Moura & Load): Física Quântica e Ciência. Discussão acalorada sobre mecânica quântica, átomos e referência ao canal Nunca Vi Um Cientista."},
	{529, "gags-e-manias", 1625, "Meme do Load (Load Comics): O Chinelo do Load. Discussão cômica sobre o icônico chinelo de dedo que o Load usa durante as transmissões no estúdio."},
	{531, "comidas-e-bebidas", 1980, "Manga / Mamão / Frutas Picadas Mofadas — Cauê Moura (Pacto das Frutas Picadas no Mercado). Cauê defende comprar fruta picada mesmo se tiver fungo/larva, alegando que 'é nutriente extra'."},
	{531, "comidas-e-bebidas", 2280, "Triple Stacker do BK do Lixo (2017) — Cauê Moura (Relato) (Era solteiro morando no Bairro do Limão). Cauê conta que jogou um lanche BK extra no lixo às 2h e recuperou para comer 12h depois."},
	{531, "comidas-e-bebidas", 3810, "Pizza Gelada / Arroz Frio / Comida Fria — Cauê Moura & Load (Enquete do dia no programa). 44% da audiência concordou que pizza gelada é boa. Cauê ama arroz frio e comida da geladeira."},
	{531, "comidas-e-bebidas", 3900, "Chá Mate Torrado com Gengibre Quente — Cauê Moura (Preparo ao vivo no microondas do estúdio). Cauê tomou o chá tinindo de quente. Buba e Load recusaram por estar pelando."},
	{531, "comidas-e-bebidas", 4080, "Perfume Kaiak Aventura — Load Comics (Nécessaire do Load). Load tirou a colônia da bolsa para mostrar como se mantém cheiroso nos compromissos."},
	{628, "comidas-e-bebidas", 300, "Doce de Leite de Viçosa vs. Lavanda — Cauê Moura & Buba (Degustação de presentes dos ouvintes). Cauê elogiou o Doce de Leite de Viçosa tradicional e detonou a versão com lavanda ('gosto de sabão')."},
	{623, "comidas-e-bebidas", 1200, "Moranguetes / Chocolates do Teatro — Cauê Moura & Load (Camarim das sessões do teatro). Cauê conta que devorou 19 moranguetes no camarim enquanto Buba fiscalizava doces."},
	{607, "comidas-e-bebidas", 2700, "Merenda de Escola / Danoninho na Caneca — Cauê Moura & Load (Nostalgia de merendeiras anos 2000). Cauê relembrou merenda de escola pública, macarrão de panela de pressão com salsicha e Danoninho na caneca."},
	{611, "comidas-e-bebidas", 2000, "Pão Doce com Margarina — Load Comics & Cauê (Apresentação da receita polêmica do Load). Load revelou que passa margarina no pão doce tradicional de creme. Cauê ficou profundamente horrorizado."},
	{609, "comidas-e-bebidas", 120, "Pote Inteiro de Azeitona / Parmesão — Cauê Moura (Larica) (Hábitos noturnos de belisco). Cauê confessou que quando está na larica come um pote inteiro de azeitonas ou quatro tubos de parmesão."},
	{531, "presentes-e-mimos", 120, "Perfume / Aromatizador de Ambiente — de Ouvinte Camarada para Bancada. Cauê borrifa o perfume no estúdio durante o início do programa."},
	{531, "presentes-e-mimos", 300, "Funko 3D do Load ('Moad') — de Jonathan (@printei_3D) para Load Comics. Estátua 3D customizada do Load com um moedor de cana saindo da traseira."},
	{531, "presentes-e-mimos", 480, "Kit de Chás Variados (Mate, Gengibre, Pêssego) — de Camarada / Membro para Bancada. Cauê preparou o chá quente ao vivo e tomou no microondas."},
	{531, "presentes-e-mimos", 600, "Fantasminhas Vermelhos 3D do Comunismo — de Bulbasauro (Maker) para Membros do Canal. Buba imprimiu em 3D mimos vermelhos com martelo/foice para sortear pros membros."},
	{628, "presentes-e-mimos", 120, "Porta-Temperos 3D de Casinha — de Ouvinte Vaca para Bulbasauro. Buba recebeu e guardou com extremo cuidado um conjunto de casinhas de porcelana para temperos."},
	{625, "presentes-e-mimos", 120, "Meias Customizadas & Cocada de Abacaxi — de Raquel Pinto & Ouvinte para Cauê & Load. Ganharam kit de meias divertidas e cocadas artesanais de abacaxi."},
	{624, "presentes-e-mimos", 120, "Estátua de Cerâmica 'Mão de Rock' — de Três Árvores Cerâmica para Cauê Moura. Cauê recebeu uma peça de cerâmica artesanal em formato de mão de rock e testou no dente."},
	{603, "presentes-e-mimos", 60, "Quadro 'Ky Crusher' (Moad Transformer) — de Azagal (Jovem Nerd) para Load Comics. Quadro emoldurado em alta qualidade desenhado por Marcelo Mattiere apresentando o Load como um Transformer moedor de cana."},
	{531, "mencoes-pop", 6000, "Craque Neto — mencionado por Cauê Moura & Otávio Reis: Superchat chamando Cauê de 'Meu Craque Neto'. Cauê faz defesa apaixonada do título de 90 e gol olímpico no Flamengo: 'Se Neto não é corintiano, não sei quem é!'"},
	{531, "mencoes-pop", 900, "David Jones & Zema — mencionado por Cauê Moura & Load: Zema usando voz de IA do David Jones em campanha. React indignado e cômico ao vídeo onde Zema usou a voz de IA do streamer David Jones sem permissão."},
	{531, "mencoes-pop", 1200, "Luciano Huck — mencionado por Cauê Moura: Plataforma de recomendação eleitoral do Huck. Cauê critica e ridiculariza o 'sapatênis humano' e a indicação de candidatos de direita na ferramenta."},
	{608, "mencoes-pop", 2400, "Craque Neto — mencionado por Load Comics & Cauê: Frases icônicas do Craque Neto. Load cita a clássica frase do Neto: 'Traumatizado fui eu quando vi meu pai comendo minha mãe com 5 anos!'"},
	{607, "mencoes-pop", 600, "Neymar Jr — mencionado por Cauê Moura & Load: Eliminação da Seleção / NeyFracasso. Análise ácida sobre a performance de Neymar na Seleção Brasileira e festas."},
	{602, "mencoes-pop", 1800, "Monark, Igor 3K & Sacani — mencionado por Cauê Moura & Load: Dramas e bastidores do Flow Podcast. Análise dos desdobramentos das polêmicas de Monark e afastamento do Flow."},
	{600, "mencoes-pop", 4753, "Edir Macedo — mencionado por Cauê Moura & Load: Investigação do Banco do Edir Macedo. Comentário afiado da bancada: 'Jesus é o caminho, Edir Macedo é o pedágio'."},
}

