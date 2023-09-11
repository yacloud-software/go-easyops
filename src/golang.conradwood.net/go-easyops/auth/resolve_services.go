package auth

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"gopkg.in/yaml.v2"
	"sync"
)

type serviceToUserIDMap struct {
	Mapping map[string]string // servicename->userid
}

const (
	service_mapping_filename = "/etc/yacloud/config/service_map.yaml"
)

var (
	service_mapping *serviceToUserIDMap
	svcmaplock      sync.Mutex
	// this is the default service mapping and is valid ONLY for the yacloud
	default_service_mapping = map[string]string{
		"alerting.AlertingService":                         "882",
		"am43controller.AM43Controller":                    "4418",
		"antidos.AntiDOS":                                  "6890",
		"apitest.ApiTestService":                           "67",
		"artefact.ArtefactService":                         "998",
		"auth.AuthenticationService":                       "",
		"auth.AuthManagerService":                          "25",
		"autodeployer.AutoDeployer":                        "18",
		"banking.Banking":                                  "4529",
		"binaryversions.BinaryVersions":                    "32160",
		"bitfolk.Bitfolk":                                  "22736",
		"buildrepoarchive.BuildRepoArchive":                "19108",
		"buildrepo.BuildRepoManager":                       "2313",
		"calendarwrapper.CalendarWrapper":                  "8961",
		"callgraph.CallGraphService":                       "642",
		"certmanager.CertManager":                          "1341",
		"cnwemails.CNWEmails":                              "32845",
		"cnwnotification.CNWNotificationService":           "35",
		"codeanalyser.CodeAnalyserService":                 "39",
		"deploymonkey.DeployMonkey":                        "20",
		"dirsizemonitor.DirSizeMonitor":                    "28560",
		"documents.DocumentProcessor":                      "4343",
		"documents.Documents":                              "4231",
		"email.EmailService":                               "242",
		"emailserver.EmailServer":                          "52930",
		"errorlogger.ErrorLogger":                          "",
		"espota.ESPOtaService":                             "833",
		"firewallmgr.FirewallMgr":                          "6900",
		"firmwaretracker.FirmwareTracker":                  "60757",
		"flightlookup.FlightLookup":                        "78593",
		"gdrive.GDrive":                                    "9073",
		"geoip.GeoIPService":                               "33",
		"gitbuilder.GitBuilder":                            "11083",
		"gitdiffsync.GitDiffSync":                          "10346",
		"github.GitHub":                                    "11121",
		"gitserver.GIT2":                                   "158",
		"gitserver.GITCredentials":                         "158",
		"goasterisk.GoAsteriskService":                     "",
		"gomodule.GoModuleService":                         "893",
		"googleactions.GoogleActions":                      "4907",
		"googlecast.GoogleCast":                            "4122",
		"goproxy.GoProxy":                                  "11070",
		"goproxy.GoProxyTestRunner":                        "11070",
		"gotools.GoTools":                                  "42195",
		"groupemail.GroupEmail":                            "10969",
		"h2gproxy.H2GProxyService":                         "37",
		"heating.HeatingService":                           "",
		"heatingschedule.HeatingScheduleService":           "1018",
		"helloworld.HelloWorld":                            "1244",
		"homeconfig.HomeConfig":                            "17481",
		"htmlserver.HTMLServerService":                     "143",
		"htmluserapp.HTMLUserApp":                          "14867",
		"httpdebug.HTTPDebug":                              "100768",
		"ifttt.IFTTTService":                               "22",
		"imagerecorder.ImageRecorder":                      "29473",
		"images.Images":                                    "15795",
		"ipexporter.IPExporterService":                     "1074",
		"ipmanager.IPManagerService":                       "4849",
		"javarepo.JavaRepo":                                "1141",
		"jokercomserver.JokerCom":                          "52930",
		"jsonapimultiplexer.JSONApiMultiplexer":            "59",
		"lockmanager.LockManager":                          "17319",
		"logservice.LogService":                            "",
		"mailshot.MailShot":                                "35242",
		"marantz.Marantz":                                  "3746",
		"mkdb.MKDB":                                        "29",
		"modmetrics.ModMetrics":                            "21645",
		"moduleprober.ModuleProber":                        "5754",
		"moduletime.ModuleTime":                            "28943",
		"objectauth.ObjectAuthService":                     "2923",
		"objectstorearchive.ObjectStoreArchive":            "23268",
		"objectstore.ObjectStore":                          "1112",
		"openweather.OpenWeatherService":                   "1076",
		"pairing.PairingService":                           "",
		"panasonic.PanasonicService":                       "",
		"payments.Payments":                                "19893",
		"pcbtype.PCBType":                                  "60350",
		"personalisedwebsite.PersonalisedWebsite":          "35470",
		"pinger.Pinger":                                    "6637",
		"pinger.PingerList":                                "6637",
		"postgresmgr.PostgresMgr":                          "1070",
		"prober.ProberService":                             "61",
		"promconfig.PromConfigService":                     "65",
		"protorenderer.ProtoRendererService":               "1114",
		"quota.QuotaService":                               "127",
		"registrymultiplexer.RegistryMultiplexerService":   "987",
		"registrymultiplexer.RegistryMultiplexerServiceRP": "",
		"registry.Registry":                                "",
		"repobuilder.RepoBuilder":                          "3539",
		"scacl.SCAclService":                               "258",
		"scapi.SCApiService":                               "753",
		"scapply.Apply":                                    "6139",
		"scautoupdate.SCAutoUpdate":                        "8588",
		"scbluetooth.SCBluetooth":                          "22417",
		"scfunctions.SCFunctionsServer":                    "151",
		"scmodcomms.SCModCommsService":                     "264",
		"scrouter.SCRouter":                                "5303",
		"scserver.SCServer":                                "149",
		"scupdate.SCUpdateService":                         "284",
		"scutils.SCUtilsServer":                            "155",
		"scvuehtml.SCVueHTML":                              "34174",
		"scweb.SCWebService":                               "145",
		"secureargs.SecureArgsService":                     "",
		"sensorapi.SensorAPIService":                       "408",
		"sensors.SensorServer":                             "147",
		"sessionmanager.SessionManager":                    "25721",
		"shellypoller.ShellyPoller":                        "33423",
		"shop.Shop":                                        "20130",
		"slackgateway.SlackGateway":                        "57",
		"sms.SMSService":                                   "176",
		"soundservice.Sound":                               "28357",
		"spamtracker.SpamTracker":                          "10645",
		"speaktome.SpeakToMeService":                       "45",
		"starling.StarlingService":                         "4568",
		"themes.Themes":                                    "5773",
		"threedprintermanager.ThreeDPrinter":               "31754",
		"urlcacher.URLCacher":                              "37004",
		"urlmapper.URLMapper":                              "4805",
		"userappcontroller.UserAppController":              "10292",
		"usercommand.UserCommandService":                   "421",
		"vuehelper.VueHelper":                              "35561",
		"webcammixer.WebCamMixer":                          "",
		"weblogin.Weblogin":                                "43",
		"weekett.Weekett":                                  "35994",
		"wiki.Wiki":                                        "16773",
		"yatools.YATools":                                  "89043",
		// YACLOUD-DEVS only -  extend list here...
	}
)

/*
sometimes a service needs to verify if it is being called by a specific service. This often implies permissions to access certain privileged bits of information. The assumption is, that service A authenticates a user and calls service B, either immediately after or some time later. In this case service B "trusts" service A. The security implication of this model is, that service B must be able to ensure service A really is who they say they are (the auth server signature should be used for this purpose) and service A has not been replaced with a different service of the same name.
For this purpose, in this case, the service to userid mappings are hardcoded into the file so to match the "yacloud" default. If someone wishes to run their own yacloud the mapping can be overriden. This then is not a programmatic option, but a configuration (administrator) option.
A file in /opt/yacloud/config/service_map.yaml, if exists, will be parsed on startup and used to provide this information.
Any lookup for a servicename that does not exist will lead to a panic() (because it is a fatal error!).
The intention of this function is to provide a means to create a common method of looking up this information, so that, in future, perhaps a good and secure way can be found to automatically map this through a combination of registry/auth-server lookups or similar.
*/
func GetServiceIDByName(servicename string) string {
	svc_to_user_load_mapping()
	uid, found := service_mapping.Mapping[servicename]
	if found {
		return uid
	}
	panic(fmt.Sprintf("[go-easyops] Application requested service \"%s\", which is not mapped to a userid", servicename))
}
func svc_to_user_load_mapping() {
	if service_mapping != nil {
		return
	}
	svcmaplock.Lock()
	defer svcmaplock.Unlock()
	if service_mapping != nil {
		return
	}
	res := &serviceToUserIDMap{
		Mapping: default_service_mapping,
	}
	b, err := utils.ReadFile(service_mapping_filename)
	if err == nil {
		fmt.Printf("[go-easyops] mapping from file %s applied\n", service_mapping_filename)
		res = &serviceToUserIDMap{}
		err = yaml.Unmarshal(b, res)
		if err != nil {
			panic(fmt.Sprintf("File %s cannot be parsed: %s\n", service_mapping_filename, err))
		}
	}
	service_mapping = res
}

func ServiceMapToYaml(m map[string]string) []byte {
	xmap := make(map[string]string)
	for k, v := range default_service_mapping {
		xmap[k] = v
	}
	for k, v := range m {
		if v != "" {
			xmap[k] = v
		} else {
			delete(xmap, k)
		}
	}
	sum := &serviceToUserIDMap{Mapping: xmap}
	b, err := yaml.Marshal(sum)
	if err != nil {
		return []byte(fmt.Sprintf("failed to yaml: %s", err))
	}
	return b
}
