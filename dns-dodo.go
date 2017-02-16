/*
dns-dodo's main purpose is to update a single dns 'A' record to the public IP address of the system dns-dojo is run on.
*/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"regexp"
	"strings"
	"time"
)

const (
	defaultUserConfigFilePath string = "~/.dns-dodo.conf"
	defaultRootConfigFile     string = "/etc/dns-dodo.conf"
	defaultExternalIPService  string = "http://myexternalip.com/raw"
)

type dnsDodo struct {
	DOPersonalAccessToken string
	ExternalIPServiceUrl  string
	OAuthClient           *http.Client
	client                *godo.Client
}

// Stores settings for the update-dns command
type updateSettings struct {
	ExternalIPServiceURL string   `json:"externalIpServiceUrl"`
	PersonalAccessToken  string   `json:"personalAccessToken"`
	Domain               string   `json:"domain"`
	Subdomain            string   `json:"subdomain"`
	PollFreq             Duration `json:"pollFreq"`
}

// displays the supplied message and exits the application with an error level of 1
func exitWithError(message string) {
	fmt.Println(message)
	os.Exit(1)
}

// Create a DNS-dodo instance.
func NewDnsDoDO() *dnsDodo {
	d := dnsDodo{}
	return &d
}

// Displays the About dns-dodo text
func (d *dnsDodo) About() {
	fmt.Println("")
	fmt.Println("                   .`.....")
	fmt.Println("                :::;++:;;+`")
	fmt.Println("              .:+@@@@@@:@+`.``")
	fmt.Println("             ::+@@;:.'@@;';'#,;.")
	fmt.Println("            .:+@;;;;;,.#':@")
	fmt.Println("            :+@,;'#',;@++#      ..")
	fmt.Println("            ;+,`#`      @'     `..    ``````")
	fmt.Println("           .'#;#        ...'++:`,: `.......```")
	fmt.Println("           ,+@'         ..'::;+`.,`.,:.,,,,.```")
	fmt.Println("           .+##   :#+#  ..,;@#+;.,.,;;;:,;;,,.``")
	fmt.Println("           ,+#+. ;`,;;#....+@.`.:,:;;;;..,;;:,```")
	fmt.Println("           ,+@,.`+;,@:+'.,;,;::::;;;;;,``.;;;,.```")
	fmt.Println("           ,#+...+:@@;+@.:::::;;;;;;;;,..,;;;;,.``")
	fmt.Println("           :'+#..@+;#+#;,;:';;;+';;;;;::,:;;';,```")
	fmt.Println("          `.++;;::@++;#::::';;;;;';''';;;;;'';:.`.")
	fmt.Println("            :+++;:;@@+;:::;;;;;;;;'''''''''''';,..")
	fmt.Println("            .'+;;'';;;:::#;;;;;;'''''''''''''';,`.")
	fmt.Println("             ;+':+.:;;:;++++++++'''''''''''''';,..")
	fmt.Println("              :'+'::'+''++++++++++'''''''''''';,``")
	fmt.Println("               .+@+++++#++++++++#++''''''''''';,.`")
	fmt.Println("                .+@@@++@@+` `:+++;++'''''''''';,.")
	fmt.Println("                  ;@@@@@@+    .++;;'#''''''''';.")
	fmt.Println("                   ;@@@###     ,#+++##''''''+':`")
	fmt.Println("                    '@@++++`    :###@;+''''+';.")
	fmt.Println("                    :#@;;+++:   `.,,` `+++#++.")
	fmt.Println("                    :++;;::;,;         ###++.")
	fmt.Println("                    :++;,,,,...        ,++:`")
	fmt.Println("                   .+++;,,....,.        ,`")
	fmt.Println("   .. .     ...,..;'+++:,.....``.")
	fmt.Println("  ..:.` `` ....,;++++++;,......`")
	fmt.Println(" .,   .`. ..;:':;++++++;,....... `")
	fmt.Println(" `...`....:'++.`.;++++'+........")
	fmt.Println("  ;,..;;.:'++:`,:::++++++,...... .")
	fmt.Println(" .';,.:.;.++++`,,,:@#++++':.....")
	fmt.Println(" `;++,,:.  +  .,,,;@++++++++,,..")
	fmt.Println("  ,;;:''.. : .,,,,:++++++##@@,..")
	fmt.Println("    .+++,. ``,,,,,;:+++++@@@@,..")
	fmt.Println("     @@; ,...,,,::;;+++#@@@@@,.")
	fmt.Println("     #@+ ,::,,,,,;;,###@@@@@@,`")
	fmt.Println("     +@@.:;,,,,,:;'+@@@@@@@@+.")
	fmt.Println("      +@@';:::;;'':@@@@@@@@@,")
	fmt.Println("       .@@;+'++'';@@@@@@@@@@")
	fmt.Println("          @@@++@@@+@@##@''+'")
	fmt.Println("           ;@@@###@+':   ';'     ,`")
	fmt.Println("              @##       ,'':    +';")
	fmt.Println("              #';`+:+::#;+` `:;;+:+;.")
	fmt.Println("           ..:;:':;'###'#@:'';::':;:+")
	fmt.Println(" ..::::',;::;;;#:;::::::.:';'#++++':,")
	fmt.Println(":::::;:;;;':;+:::,..;:;+:;''+;'+,'.::::;::,.")
	fmt.Println("::::''';:;#'##:;'++'','':`::,''''#::::::::::::..")
	fmt.Println(":::;:;:;:;::;;:+#'':;;;::::;:::;:::::;:::::::;;..")
	fmt.Println("..:::::;::;;,;'':;,,,.;;:;;:::;::::::::::::::::.")
	fmt.Println("    ;:::;;:;:::;;:::,'.;:;;;;;:::::::::::::..")
	fmt.Println("        ..;;:::;;;'#+';:::.")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("      _                     _           _")
	fmt.Println("     | |                   | |         | |")
	fmt.Println("   __| |_ __  ___ ______ __| | ___   __| | ___")
	fmt.Println("  / _` | '_ \\/ __|______/ _` |/ _ \\ / _` |/ _ \\")
	fmt.Println(" | (_| | | | \\__ \\     | (_| | (_) | (_| | (_) |")
	fmt.Println("  \\__,_|_| |_|___/      \\__,_|\\___/ \\__,_|\\___/")
	fmt.Println("")
	fmt.Println("               dns Digital Ocean do  ")
	fmt.Println("")
	fmt.Println(" Dynamic DNS sub-domain updater for Digital Ocean  ")
	fmt.Println(" by lummie  - http://lummie.github.io ")
	fmt.Println("")
}

// Returns the public ip address as returned by the specified ExternalIPServiceURL
func (d *dnsDodo) getExternalIP() string {
	// check IP service URL exists
	if d.ExternalIPServiceUrl == "" {
		exitWithError("Please supply an External IP Service URL, e.g. " + defaultExternalIPService)
	}

	// request external ip address
	resp, err := http.Get(d.ExternalIPServiceUrl)
	if err != nil {
		exitWithError(fmt.Sprintf("Failed trying External IP Service : %v", err.Error()))
	}

	if resp.StatusCode != 200 {
		exitWithError(fmt.Sprintf("Failed trying External IP Service Response : %v", resp.StatusCode))
	}

	// close the response later
	defer resp.Body.Close()

	// get the whole contents or the response body
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exitWithError(fmt.Sprintf("Failed to read response: %v", err.Error()))
	}

	// convert byte array to string and return
	return strings.TrimSpace(string(contents))
}

// checks the ip address matches a valid IP regex
func (d *dnsDodo) CheckIPV4(ip string) {
	// check that the ip address string conforms to an IPV4 address
	matches, err := regexp.MatchString(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`, ip)
	if err != nil {
		exitWithError(fmt.Sprintf("Error matching IP to an IPV4 address [%s]", ip))
	}
	if matches == false {
		exitWithError(fmt.Sprintf("The ip address does not appear to be valid: %s", ip))
	}
}

// Required Token interface for the oauth2 support
// Sets up a token using the supplied personal access token
func (d *dnsDodo) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: d.DOPersonalAccessToken,
	}
	return token, nil
}

// establishes a client connection to Digital Ocean API and stores it in the client property
func (d *dnsDodo) EstablishGoDoClient() {
	oauthClient := oauth2.NewClient(oauth2.NoContext, d)
	d.client = godo.NewClient(oauthClient)
}

// retrieves all the Domain Records registered with the Digital Ocean account
// As the api is paged, this method collates all page results into a single list
func (d *dnsDodo) GetDNSEntries(domainName string) []godo.DomainRecord {
	list := []godo.DomainRecord{} // to hold the domain records
	opts := &godo.ListOptions{}
	for {
		// get a page of domain records
		if pageOfRecords, resp, err := d.client.Domains.Records(domainName, opts); err != nil {
			exitWithError(err.Error())
		} else {
			// Add the retrieved records to the list
			list = append(list, pageOfRecords...)

			// if there is not a links or is the last page in the response then we are done
			if resp.Links == nil || resp.Links.IsLastPage() {
				break
			}

			// if there are no more api requests per hour then fail
			if resp.Rate.Remaining == 0 {
				exitWithError(fmt.Sprintf("No API requests remaining for this hour, %v used. Next request can be made at %v", resp.Rate.Limit, resp.Rate.Reset))
			}

			// get current Page and increment the request options to the next page
			if page, err := resp.Links.CurrentPage(); err != nil {
				exitWithError(err.Error())
			} else {
				opts.Page = page + 1
			}
		}

	}
	return list
}

func (d *dnsDodo) FilteredRecords(records []godo.DomainRecord, typeFilter, nameFilter string) []godo.DomainRecord {
	result := []godo.DomainRecord{}
	for _, record := range records {
		if (typeFilter == "" || record.Type == typeFilter) && (nameFilter == "" || record.Name == nameFilter) {
			result = append(result, record)
		}
	}
	return result
}

func (d *dnsDodo) OutputDomainRecords(records []godo.DomainRecord) {
	fmt.Println("ID\tName\tType\tData")
	for _, record := range records {
		fmt.Printf("%d\t%s\t%s\t%s\n", record.ID, record.Name, record.Type, record.Data)
	}
}

func (d *dnsDodo) UpdateDNSEntry(domainName, subdomain, ipAddress string, printTimestamp bool) {
	d.CheckIPV4(ipAddress)

	updateMsgFmt := "About to update %s.%s to %s\n"
	if printTimestamp {
		updateMsgFmt = fmt.Sprintf("[%s] %s", time.Now(), updateMsgFmt)
	}
	fmt.Printf(updateMsgFmt, subdomain, domainName, ipAddress)
	records := d.GetDNSEntries(domainName)
	if len(records) == 0 { // Check we have some entires before filtering
		exitWithError(fmt.Sprintf("The dodo failed to retrieve any DNS entries for the domain '%s'", domainName))
	}
	records = d.FilteredRecords(records, "A", subdomain)
	if len(records) == 0 { // Check we have found a record
		exitWithError(fmt.Sprintf("There does not seem to be a DNS entry of type 'A' that has a name of '%s' for the domain '%s'", subdomain, domainName))
	}
	if len(records) > 1 { // Check we have not found more than one matching record
		exitWithError(fmt.Sprintf("There does not seem to be a DNS entry of type 'A' that has a name of '%s' for the domain '%s'", subdomain, domainName))
	}

	// get the single record to update
	record := records[0]
	if record.Name != subdomain || record.Type != "A" {
		exitWithError(fmt.Sprintf("Expected the record %v to be the record to update", record))
	}

	// check the ip has changed
	if record.Data != ipAddress {
		// create an update record with the updated details
		domainER := godo.DomainRecordEditRequest{
			Type:     record.Type,
			Name:     record.Name,
			Priority: record.Priority,
			Port:     record.Port,
			Weight:   record.Weight,
			Data:     ipAddress,
		}

		// update the specified record id
		recordId := record.ID
		_, _, err := d.client.Domains.EditRecord(domainName, recordId, &domainER)
		if err != nil {
			exitWithError(err.Error())
		}
		fmt.Println("Updated.")
	} else {
		fmt.Println("IP has not changed.")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "dns-dodo"
	app.Usage = "Dynamic DNS sub-domain updater for Digital Ocean."
	app.Version = "1.2"

	// setup the default flags that are optional
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "extip, x",
			Usage: "[Optional] External IP service URL",
			Value: defaultExternalIPService,
		},
	}

	// add the commands to the application
	app.Commands = []cli.Command{
		// show the external IP address from the external IP service
		{
			Name:  "show-ip",
			Usage: fmt.Sprintf("Display the external IP Address using the default provider [%v], otherwise specify via --extip", defaultExternalIPService),
			Action: func(c *cli.Context) {
				dnsDodo := NewDnsDoDO()
				dnsDodo.ExternalIPServiceUrl = c.GlobalString("extip")
				ip := dnsDodo.getExternalIP()
				dnsDodo.CheckIPV4(ip)
				fmt.Printf("%s returned from %s\n", ip, dnsDodo.ExternalIPServiceUrl)
			},
		},

		// show the DNS entries on Digital Ocean
		{
			Name:  "show-dns",
			Usage: "Display the DNS entries on the account associated with the specified Personal Access Token and the specified domain name",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "pat, t",
					Usage: "[Required] Personal Access Token from the Digital Ocean API, https://cloud.digitalocean.com/settings/tokens/new",
				},
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "[Required] Domain Name to update or list dns entries for",
				},
				cli.StringFlag{
					Name:  "type",
					Usage: "[Optional] Type of DNS record to display e.g A",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "[Optional] value of the DNS Name property to filter to e.g @ | sub-domain",
				},
			},
			Action: func(c *cli.Context) {
				if !c.IsSet("pat") {
					exitWithError("Error: A Personal Access Token must be supplied using the --pat flag")
				}

				if !c.IsSet("domain") {
					exitWithError("Error: A domain must be specified to retrieve the DNS records for using the --domain flag")
				}

				dnsDodo := NewDnsDoDO()
				dnsDodo.DOPersonalAccessToken = c.String("pat")
				dnsDodo.EstablishGoDoClient()
				records := dnsDodo.GetDNSEntries(c.String("domain"))
				records = dnsDodo.FilteredRecords(records, c.String("type"), c.String("name"))
				dnsDodo.OutputDomainRecords(records)
			},
		},

		// update a DNS entry
		{
			Name:  "update-dns",
			Usage: "Updates the DNS entry type 'A' for the specified [domain] that has a name of [subdomain] to the external IP address",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "pat, t",
					Usage: "[Required] Personal Access Token from the Digital Ocean API, https://cloud.digitalocean.com/settings/tokens/new",
				},
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "[Required] Domain Name to update the DNS 'A' record for",
				},
				cli.StringFlag{
					Name:  "sub-domain, s",
					Usage: "Sub-domain of the Domain Name that will be updated with the IP address",
				},
				cli.BoolFlag{
					Name:  "poll, p",
					Usage: "[Optional] Poll for changes to your external IP to send to Digital Ocean",
				},
				cli.DurationFlag{
					Name:  "pollfreq, f",
					Usage: "[Optional] Polling frequency in standard duration format e.g. 5m (5 minutes). Only applicable with --poll.",
				},
				cli.StringFlag{
					Name:  "config, c",
					Usage: "[Optional] Config file from which to source update settings. Any settings from this file are overriden by flags provided directly. (Default ~/.dns-dodo.conf)",
				},
			},
			Action: func(c *cli.Context) {
				settings := getUpdateSettings(c)

				dnsDodo := NewDnsDoDO()
				dnsDodo.ExternalIPServiceUrl = settings.ExternalIPServiceURL
				dnsDodo.DOPersonalAccessToken = settings.PersonalAccessToken
				dnsDodo.EstablishGoDoClient()

				var lastIP string
				updateFunc := func(polling bool) {
					ip := dnsDodo.getExternalIP()
					dnsDodo.CheckIPV4(ip) // check the ip is valid before we attempt to connect to Digital Ocean

					// If IP hasn't changed since we last polled, don't update
					if ip != lastIP {
						dnsDodo.UpdateDNSEntry(settings.Domain, settings.Subdomain, ip, polling)
					} else {
						fmt.Printf("[%s] IP (%s) hasn't changed since last poll\n", time.Now(), ip)
					}
					lastIP = ip
				}

				polling := c.Bool("poll")
				updateFunc(polling)

				if polling {
					pollFreq := settings.PollFreq.Duration
					if pollFreq == 0 {
						pollFreq = time.Minute
					}

					sigChan := make(chan os.Signal)
					signal.Notify(sigChan)
					updateTicker := time.NewTicker(pollFreq)

					for {
						select {
						case <-updateTicker.C:
							updateFunc(true)
						case <-sigChan:
							fmt.Println("")
							fmt.Println("Going extinct...")
							updateTicker.Stop()
							os.Exit(0)
						}
					}
				}
			},
		},

		// go get the dodo from the dojo
		{
			Name:  "dodo",
			Usage: "Probably pining for the fjords...",
			Action: func(c *cli.Context) {
				NewDnsDoDO().About()
			},
		},

		// display the application version
		{
			Name:  "version",
			Usage: "Display application version",
			Action: func(c *cli.Context) {
				fmt.Println(app.Version)
			},
		},
	}

	// if no commands are supplied just show the help
	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	// go make those dodos extinct
	app.Run(os.Args)
}

// Derives the settings to be used in an update via config file/flags
func getUpdateSettings(c *cli.Context) *updateSettings {
	var configFilePath string
	usingDefaultFile := true

	if usr, _ := user.Current(); usr.Uid == "0" {
		configFilePath = defaultRootConfigFile
	} else {
		configFilePath = defaultUserConfigFilePath
	}

	if c.IsSet("config") {
		configFilePath = c.String("config")
		usingDefaultFile = false
	}

	// First find any settings defined in a config file.
	settings := readConfigFile(configFilePath, usingDefaultFile)

	// Override any settings from the config file with flags on the command
	applyFlags(settings, c)

	// Ensure that all mandatory settings for an update have been provided.
	checkUpdateSettings(settings)

	return settings
}

func readConfigFile(configFilePath string, isDefaultFile bool) *updateSettings {
	if strings.Index(configFilePath, "~") == 0 {
		configFilePath = strings.Replace(configFilePath, "~", os.Getenv("HOME"), 1)
	}

	b, err := ioutil.ReadFile(configFilePath)

	if os.IsNotExist(err) && !isDefaultFile {
		fmt.Println("Config file", configFilePath, "does not exist")
		os.Exit(1)
	}

	var settingsFromFile updateSettings
	json.Unmarshal(b, &settingsFromFile)
	return &settingsFromFile
}

// Puts settings from the CLI context into an updateSettings struct.
func applyFlags(settings *updateSettings, c *cli.Context) {
	if c.IsSet("pat") {
		settings.PersonalAccessToken = c.String("pat")
	}

	if c.IsSet("domain") {
		settings.Domain = c.String("domain")
	}

	if c.IsSet("subdomain") {
		settings.Subdomain = c.String("subdomain")
	}

	if c.IsSet("pollFreq") {
		settings.PollFreq = Duration{c.Duration("pollFreq")}
	}

	if c.GlobalIsSet("extip") {
		settings.ExternalIPServiceURL = c.GlobalString("extip")
	}
}

// Validates that mandatory properties have been provided to do an update.
func checkUpdateSettings(settings *updateSettings) {
	if settings.PersonalAccessToken == "" {
		exitWithError("No Personal Access Token specified")
	}
	if settings.Domain == "" {
		exitWithError("No domain specified")
	}

	if settings.Subdomain == "" {
		exitWithError("No subdomain specified")
	}
}
