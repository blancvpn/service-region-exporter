package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	Name() string
	CheckRegion() (string, error)
	IsEnabled() bool
}

type BaseService struct {
	name    string
	client  *http.Client
	enabled bool
}

func (b *BaseService) Name() string {
	return b.name
}

func (b *BaseService) IsEnabled() bool {
	return b.enabled
}

func (b *BaseService) makeRequest(req *http.Request) ([]byte, error) {
	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limited or blocked: status %d", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}

func getCountryCode(countryName string) (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("https://restcountries.com/v3.1/name/%s?fields=cca2", countryName)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get country code: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var countries []struct {
		CCA2 string `json:"cca2"`
	}

	if err := json.Unmarshal(body, &countries); err != nil {
		return "", err
	}

	if len(countries) == 0 {
		return "", fmt.Errorf("country not found")
	}

	return countries[0].CCA2, nil
}

var countryNameToCode = map[string]string{
	"afghanistan":                           "AF",
	"albania":                               "AL",
	"algeria":                               "DZ",
	"andorra":                               "AD",
	"angola":                                "AO",
	"antigua and barbuda":                   "AG",
	"argentina":                             "AR",
	"armenia":                               "AM",
	"australia":                             "AU",
	"austria":                               "AT",
	"azerbaijan":                            "AZ",
	"bahamas":                               "BS",
	"bahrain":                               "BH",
	"bangladesh":                            "BD",
	"barbados":                              "BB",
	"belarus":                               "BY",
	"belgium":                               "BE",
	"belize":                                "BZ",
	"benin":                                 "BJ",
	"bhutan":                                "BT",
	"bolivia":                               "BO",
	"bosnia and herzegovina":                "BA",
	"botswana":                              "BW",
	"brazil":                                "BR",
	"brunei":                                "BN",
	"bulgaria":                              "BG",
	"burkina faso":                          "BF",
	"burundi":                               "BI",
	"cabo verde":                            "CV",
	"cape verde":                            "CV",
	"cambodia":                              "KH",
	"cameroon":                              "CM",
	"canada":                                "CA",
	"central african republic":              "CF",
	"chad":                                  "TD",
	"chile":                                 "CL",
	"china":                                 "CN",
	"colombia":                              "CO",
	"comoros":                               "KM",
	"congo":                                 "CG",
	"congo (brazzaville)":                   "CG",
	"congo (kinshasa)":                      "CD",
	"democratic republic of the congo":      "CD",
	"cook islands":                          "CK",
	"costa rica":                            "CR",
	"cote d'ivoire":                         "CI",
	"ivory coast":                           "CI",
	"croatia":                               "HR",
	"cuba":                                  "CU",
	"cyprus":                                "CY",
	"czech republic":                        "CZ",
	"czechia":                               "CZ",
	"denmark":                               "DK",
	"djibouti":                              "DJ",
	"dominica":                              "DM",
	"dominican republic":                    "DO",
	"east timor":                            "TL",
	"timor-leste":                           "TL",
	"ecuador":                               "EC",
	"egypt":                                 "EG",
	"el salvador":                           "SV",
	"equatorial guinea":                     "GQ",
	"eritrea":                               "ER",
	"estonia":                               "EE",
	"eswatini":                              "SZ",
	"swaziland":                             "SZ",
	"ethiopia":                              "ET",
	"fiji":                                  "FJ",
	"finland":                               "FI",
	"france":                                "FR",
	"gabon":                                 "GA",
	"gambia":                                "GM",
	"georgia":                               "GE",
	"germany":                               "DE",
	"ghana":                                 "GH",
	"greece":                                "GR",
	"grenada":                               "GD",
	"guatemala":                             "GT",
	"guinea":                                "GN",
	"guinea-bissau":                         "GW",
	"guyana":                                "GY",
	"haiti":                                 "HT",
	"honduras":                              "HN",
	"hungary":                               "HU",
	"iceland":                               "IS",
	"india":                                 "IN",
	"indonesia":                             "ID",
	"iran":                                  "IR",
	"iraq":                                  "IQ",
	"ireland":                               "IE",
	"israel":                                "IL",
	"italy":                                 "IT",
	"jamaica":                               "JM",
	"japan":                                 "JP",
	"jordan":                                "JO",
	"kazakhstan":                            "KZ",
	"kenya":                                 "KE",
	"kiribati":                              "KI",
	"korea":                                 "KR",
	"south korea":                           "KR",
	"republic of korea":                     "KR",
	"north korea":                           "KP",
	"democratic people's republic of korea": "KP",
	"kosovo":                                "XK",
	"kuwait":                                "KW",
	"kyrgyzstan":                            "KG",
	"laos":                                  "LA",
	"latvia":                                "LV",
	"lebanon":                               "LB",
	"lesotho":                               "LS",
	"liberia":                               "LR",
	"libya":                                 "LY",
	"liechtenstein":                         "LI",
	"lithuania":                             "LT",
	"luxembourg":                            "LU",
	"madagascar":                            "MG",
	"malawi":                                "MW",
	"malaysia":                              "MY",
	"maldives":                              "MV",
	"mali":                                  "ML",
	"malta":                                 "MT",
	"marshall islands":                      "MH",
	"mauritania":                            "MR",
	"mauritius":                             "MU",
	"mexico":                                "MX",
	"micronesia":                            "FM",
	"moldova":                               "MD",
	"monaco":                                "MC",
	"mongolia":                              "MN",
	"montenegro":                            "ME",
	"morocco":                               "MA",
	"mozambique":                            "MZ",
	"myanmar":                               "MM",
	"burma":                                 "MM",
	"namibia":                               "NA",
	"nauru":                                 "NR",
	"nepal":                                 "NP",
	"netherlands":                           "NL",
	"new zealand":                           "NZ",
	"nicaragua":                             "NI",
	"niger":                                 "NE",
	"nigeria":                               "NG",
	"north macedonia":                       "MK",
	"macedonia":                             "MK",
	"norway":                                "NO",
	"oman":                                  "OM",
	"pakistan":                              "PK",
	"palau":                                 "PW",
	"palestine":                             "PS",
	"state of palestine":                    "PS",
	"panama":                                "PA",
	"papua new guinea":                      "PG",
	"paraguay":                              "PY",
	"peru":                                  "PE",
	"philippines":                           "PH",
	"poland":                                "PL",
	"portugal":                              "PT",
	"qatar":                                 "QA",
	"romania":                               "RO",
	"russia":                                "RU",
	"russian federation":                    "RU",
	"rwanda":                                "RW",
	"saint kitts and nevis":                 "KN",
	"saint lucia":                           "LC",
	"saint vincent and the grenadines":      "VC",
	"samoa":                                 "WS",
	"san marino":                            "SM",
	"sao tome and principe":                 "ST",
	"saudi arabia":                          "SA",
	"senegal":                               "SN",
	"serbia":                                "RS",
	"seychelles":                            "SC",
	"sierra leone":                          "SL",
	"singapore":                             "SG",
	"slovakia":                              "SK",
	"slovenia":                              "SI",
	"solomon islands":                       "SB",
	"somalia":                               "SO",
	"south africa":                          "ZA",
	"south sudan":                           "SS",
	"spain":                                 "ES",
	"sri lanka":                             "LK",
	"sudan":                                 "SD",
	"suriname":                              "SR",
	"sweden":                                "SE",
	"switzerland":                           "CH",
	"syria":                                 "SY",
	"syrian arab republic":                  "SY",
	"taiwan":                                "TW",
	"tajikistan":                            "TJ",
	"tanzania":                              "TZ",
	"thailand":                              "TH",
	"togo":                                  "TG",
	"tonga":                                 "TO",
	"trinidad and tobago":                   "TT",
	"tunisia":                               "TN",
	"turkey":                                "TR",
	"turkmenistan":                          "TM",
	"tuvalu":                                "TV",
	"uganda":                                "UG",
	"ukraine":                               "UA",
	"united arab emirates":                  "AE",
	"uae":                                   "AE",
	"united kingdom":                        "GB",
	"great britain":                         "GB",
	"uk":                                    "GB",
	"united states":                         "US",
	"usa":                                   "US",
	"united states of america":              "US",
	"uruguay":                               "UY",
	"uzbekistan":                            "UZ",
	"vanuatu":                               "VU",
	"vatican city":                          "VA",
	"venezuela":                             "VE",
	"vietnam":                               "VN",
	"yemen":                                 "YE",
	"zambia":                                "ZM",
	"zimbabwe":                              "ZW",
}

func NormalizeCountry(country string) string {
	normalized := strings.ToLower(strings.TrimSpace(country))

	if len(normalized) == 2 {
		return strings.ToUpper(normalized)
	}

	if code, ok := countryNameToCode[normalized]; ok {
		return code
	}

	return country
}

func GetAllServices(client *http.Client) []Service {
	return []Service{
		NewGoogleService(client),
		NewYouTubeService(client),
		NewChatGPTService(client),
		NewNetflixService(client),
		NewTwitchService(client),
		NewSpotifyService(client),
		NewDeezerService(client),
		NewRedditService(client),
		NewAmazonPrimeService(client),
		NewAppleService(client),
		NewSteamService(client),
		NewPlayStationService(client),
		NewTikTokService(client),
		NewJetBrainsService(client),
		NewBingService(client),
	}
}
