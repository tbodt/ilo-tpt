package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	TOKEN    = "OTAzODU2OTg2MzI2Njk1OTQ3.YXzEag.omRBfI9wa-Ghtko9swXvbRUb4_M"
	JAN_TEPO = "155417194530996225"
)

// ilo li kute e wile jan sama nasin ni:
// @ilo o weka e ike mi
// @ilo o lukin ala e ike mi sama
// @ilo o mu la ilo li mu

var (
	nimiIlo string
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	siko, pkl := discordgo.New("Bot " + TOKEN)
	if pkl != nil {
		log.Fatal(pkl)
	}
	siko.Identify.Intents = discordgo.IntentsGuildMessages
	siko.AddHandler(func(s *discordgo.Session, t *discordgo.MessageCreate) {
		pkl := tokiLiKama(s, t)
		if pkl != nil {
			log.Printf("mi pali tan toki %#v la ilo li pakala %v", t.Content, pkl)
			s.ChannelMessageSendComplex(t.ChannelID, &discordgo.MessageSend{
				Content: fmt.Sprintf(":warning: ilo li pakala! <@%s> o pona e ni\n```%v```", JAN_TEPO, pkl),
				AllowedMentions: &discordgo.MessageAllowedMentions{
					Users: []string{JAN_TEPO},
				},
			})
		}
	})

	pkl = siko.Open()
	if pkl != nil {
		log.Fatal(pkl)
	}
	ijoNimiIlo, pkl := siko.User("@me")
	if pkl != nil {
		log.Fatal(pkl)
	}
	nimiIlo = "<@!" + ijoNimiIlo.ID + ">"
	log.Println("ilo li open")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	siko.Close()

	// ilo o pali e seme?
}

type wileJan struct {
	WileMa map[string]wileJanMa
}
type wileJanMa struct {
	TokiOPona     bool
	IjoPiLukinAla []kenToki
}

// ijo ni wan la toki wan li ken sama li ken sama ala
type kenToki struct {
	Open string
	Pini string
}

func kenTokiIjo(ijo string) (kenToki, error) {
	insaIjo := strings.SplitN(ijo, "ijo", 2)
	if len(insaIjo) == 0 {
		return kenToki{}, fmt.Errorf("%#v li jo ala e %#v lon insa", ijo, "ijo")
	}
	if len(insaIjo) == 1 {
		return kenToki{Open: insaIjo[0]}, nil
	}
	return kenToki{Open: insaIjo[0], Pini: insaIjo[1]}, nil
}

func (k *kenToki) samaAlaSama(toki string) bool {
	return strings.HasPrefix(toki, k.Open) && strings.HasSuffix(toki, k.Pini)
}

func lipuSonaLaNimi(nimi string) string {
	return path.Join("sona", path.Join("/", nimi)) + ".json"
}
func oJoESona(nimi string, ijo interface{}) error {
	nanpa, pkl := os.ReadFile(lipuSonaLaNimi(nimi))
	if os.IsNotExist(pkl) {
		return nil
	}
	if pkl != nil {
		return pkl
	}
	return json.Unmarshal(nanpa, ijo)
}
func oAwenESona(nimi string, ijo interface{}) error {
	nanpa, pkl := json.Marshal(ijo)
	if pkl != nil {
		return pkl
	}
	return os.WriteFile(lipuSonaLaNimi(nimi), nanpa, 0666)
}

func oAlasaEKenTokiLonKulupu(kulupu []kenToki, kenAlasa kenToki) int {
	for n, ken := range kulupu {
		if ken.Open == kenAlasa.Open && ken.Pini == kenAlasa.Pini {
			return n
		}
	}
	return -1
}

func tokiLiKama(s *discordgo.Session, t *discordgo.MessageCreate) error {
	//tenpoOpen := time.Now()

	if t.Author.Bot {
		return nil
	}

	var wile wileJan
	pkl := oJoESona(t.Author.ID, &wile)
	if pkl != nil {
		return pkl
	}
	wileMa := wile.WileMa[t.GuildID]

	// toki li tawa ala tawa ilo?
	if strings.HasPrefix(t.Content, nimiIlo+" o ") || strings.HasPrefix(t.Content, nimiIlo+" li ") {
		toki := strings.TrimSpace(strings.TrimPrefix(t.Content, nimiIlo+" "))
		var tokiKama string
		wileLiAnte := false
		if toki == "o mu" {
			tokiKama = "mu"
		} else if toki == "li seme e mi" {
			if wileMa.TokiOPona {
				tokiKama += "ma ni la mi weka e toki ike sina.\n"
			} else {
				tokiKama += "ma ni la mi weka ala e toki ike sina.\n"
			}
			if len(wileMa.IjoPiLukinAla) > 0 {
				tokiKama += "toki sina li sama ni la mi lukin ala:\n"
				for _, ijo := range wileMa.IjoPiLukinAla {
					tokiKama += fmt.Sprintf("• %sijo%s\n", ijo.Open, ijo.Pini)
				}
			}
		} else if toki == "o pona e mi" {
			wileMa.TokiOPona = true
			wileLiAnte = true
			tokiKama = "mi kama ni"
		} else if toki == "o pona ala e mi" {
			wileMa.TokiOPona = false
			wileLiAnte = true
			tokiKama = "mi kama ni"
		} else if strings.HasPrefix(toki, "o lukin ala e toki mi sama ") {
			kenToki, pkl := kenTokiIjo(strings.TrimPrefix(toki, "o lukin ala e toki mi sama "))
			if pkl != nil {
				tokiKama = fmt.Sprintf("pakala! %v", pkl)
			} else if oAlasaEKenTokiLonKulupu(wileMa.IjoPiLukinAla, kenToki) != -1 {
				tokiKama = "mi awen ni"
			} else {
				tokiKama = "mi kama ni"
				wileMa.IjoPiLukinAla = append(wileMa.IjoPiLukinAla, kenToki)
				wileLiAnte = true
			}
		} else if strings.HasPrefix(toki, "o lukin e toki mi sama ") {
			kenToki, pkl := kenTokiIjo(strings.TrimPrefix(toki, "o lukin e toki mi sama "))
			if pkl != nil {
				tokiKama = fmt.Sprintf("pakala! %v", pkl)
			} else {
				nKen := oAlasaEKenTokiLonKulupu(wileMa.IjoPiLukinAla, kenToki)
				if nKen == -1 {
					tokiKama = "mi awen ni"
				} else {
					tokiKama = "mi kama ni"
					wileMa.IjoPiLukinAla = append(wileMa.IjoPiLukinAla[:nKen], wileMa.IjoPiLukinAla[nKen+1:]...)
					wileLiAnte = true
				}
			}
		} else {
			tokiKama = fmt.Sprintf(`
toki! mi o seme? sina ken toki tawa mi sama ni:
**%[1]s o pona e mi** la ni li kama: sina toki kepeken toki pona ala lon ma ni la mi weka e toki sina.
**%[1]s o pona ala e mi** la mi kama lukin ala e toki sina.
**%[1]s o lukin ala e toki mi sama __ijo__** la ni li kama: toki sina li sama __ijo__ la mi lukin ala e ona.
    __ijo__ la o sitelen e `+"`"+`ijo`+"`"+` o pana e sitelen nasa lon open lon pini.  sitelen nasa ni li lon toki la mi lukin ala e toki, sama ni:
        __ijo__ li `+"`"+`./ijo`+"`"+` la sitelen `+"`"+`./`+"`"+` li lon open toki la mi lukin ala.
    ilo Pulaki li kute e sitelen nasa lon nasin sama.
**%[1]s o lukin e toki mi sama __ijo__** la mi lukin ala e toki sina sama __ijo__ lon tenpo pini la mi weka e ni li kama ni ala.
**%[1]s li seme e mi** la mi toki e ni: mi lukin ala lukin e toki seme sina lon ma ni.
**%[1]s o mu** la mi mu.
`, nimiIlo)
		}
		if wileLiAnte {
			if wile.WileMa == nil {
				wile.WileMa = make(map[string]wileJanMa)
			}
			wile.WileMa[t.GuildID] = wileMa
			pkl := oAwenESona(t.Author.ID, &wile)
			if pkl != nil {
				return pkl
			}
		}
		_, pkl = s.ChannelMessageSendReply(t.ChannelID, tokiKama, t.Reference())
		if pkl != nil {
			return pkl
		}
		return nil
	}

	// toki li tawa ilo ala.

	// jan toki li wile ala wile e weka pi toki ike? ilo li sona e wile tan seme? jan li toki e ni:
	//	ona li wile toki pona taso (ma ante la wile li ante)
	//	sitelen nasa ni li lon toki la ilo o lukin ala e toki
	//		jan kepeken li wile e ijo lon ma ale tan ilo pulaki.
	//		jan lawa ma li wile e ijo lon ma wan lon jan ale tan ni: ilo ante li kute kepeken toki ike. ken la ni li suli ala.
	//		ken la jan kepeken li wile e ijo lon ma wan. ken la ni li suli ala.
	//	mi wile pali lili la nasin li ni taso lon tenpo ni: sitelen nasa li tawa jan wan tawa ma wan.

	if !wileMa.TokiOPona {
		return nil
	}
	for _, ijo := range wileMa.IjoPiLukinAla {
		if ijo.samaAlaSama(t.Content) {
			return nil
		}
	}

	if !tokiLiPonaAlaPona(t.Content) {
		pkl := s.ChannelMessageDelete(t.ChannelID, t.ID)
		if pkl != nil {
			return pkl
		}
		tomoPiJanNi, pkl := s.UserChannelCreate(t.Author.ID)
		if pkl != nil {
			return pkl
		}
		_, pkl = s.ChannelMessageSend(tomoPiJanNi.ID, fmt.Sprintf("toki sina ni li pona ala. mi weka e ona.\n>>> %s", t.Content))
		if pkl != nil {
			return pkl
		}
	}

	return nil
}

var (
	IJO_PI_IKE_ALA = []*regexp.Regexp{
		regexp.MustCompile(`\|\|[^\|]+\|\|`),
		regexp.MustCompile(`"[^"]+"`),
		regexp.MustCompile(`<a?:\w+:\d+>`),
		regexp.MustCompile(`<t?:\w+:\d+(:\w)?>`),
		regexp.MustCompile(`https?:\/\/\S+`),
	}
	SITELEN_NIMI = regexp.MustCompile(`\pL+`)
	NIMI_PONA    = regexp.MustCompile(`^([jklmnpstw]?[aeiou]n?([jklmnpstw][aeiou]n?)*|\p{Lu}.*|.{1}|n|msa|cw)$`)
)

func tokiLiPonaAlaPona(toki string) bool {
	// toki li ike ala ike?
	// o weka e ijo pi ike ala. ni li ken nasa ilo li ken toki len. ijo li lon poki pi sitelen "" la ni kin li pona.
	for _, ijo := range IJO_PI_IKE_ALA {
		toki = ijo.ReplaceAllString(toki, " ")
	}
	// sitelen sama mute li lon poka la o wan e ona.
	sitelenToki := []rune(toki)
	toki = ""
	for n, s := range sitelenToki {
		if n > 0 && s == sitelenToki[n-1] {
			continue
		}
		toki += string(s)
	}
	sitelenToki = []rune(toki)
	// o lukin e nimi wan ale.
	mutePiNimiIke := 0
	muteNimi := 0
	for _, nimi := range SITELEN_NIMI.FindAllString(toki, -1) {
		muteNimi++
		if !NIMI_PONA.MatchString(nimi) {
			mutePiNimiIke++
		}
	}
	if float64(mutePiNimiIke) > float64(muteNimi)*0.1 {
		return false
	}
	return true
}
