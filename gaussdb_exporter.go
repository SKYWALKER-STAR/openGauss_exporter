package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	_ "github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"

	"collector"
	"config"

)

var (
	metricsPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()

	timeoutOffset = kingpin.Flag(
		"timeout-offset",
		"Offset to subtract from timeout in seconds.",
	).Default("0.25").Float64()

	configGaussDB = kingpin.Flag(
		"config.gaussDB.yaml",
		"Path to .gaussDB.yaml file to read GaussDB credentials from.",
	).Default(".gaussDB.yaml").String()

	gaussDBAddress = kingpin.Flag(
		"gaussDB.address",
		"Address to use for connecting to gaussDB",
	).Default("localhost:5432").String()
	gaussDBUser = kingpin.Flag(
		"gaussDB.username",
		"Hostname to use for connecting to gaussDB",
	).String()

	tlsInsecureSkipVerify = kingpin.Flag(
		"tls.insecure-skip-verify",
		"Ignore certificate and server verification when using a tls connection.",
	).Bool()

	toolkitFlags = webflag.AddFlags(kingpin.CommandLine, ":9194")
	c            = config.Handler{
		Config: &config.Config{},
	}
)

var scrapers = map[collector.Scraper]bool {
	collector.PGRolesCollector{}:			true,
	collector.PGDatabaseCollector{}:		true,
	collector.PGLocksCollector{}:			true,
	collector.PGLongRunningTransactionsCollector{}:	true,
	collector.PGPostmasterCollector{}:		true,
	collector.PGReplicationCollector{}:		true,
	collector.PGReplicationSlotCollector{}:		true,
	collector.PGStatActivityAutovacuumCollector{}:	true,
	collector.PGStatBGWriterCollector{}:		true,
	collector.PGStatDatabaseCollector{}:		true,
	collector.PGStatioUserIndexesCollector{}:	true,
	collector.PGStatIOUserTablesCollector{}:	true,
	collector.PGStatStatementsCollector{}:		true,
	collector.PGStatUserTablesCollector{}:		true,
	collector.PGStatWalReceiverCollector{}:		true,
	collector.PGXlogLocationCollector{}:		true,

}

func filterScrapers(scrapers []collector.Scraper, collectParams []string) []collector.Scraper {
	var filteredScrapers []collector.Scraper

	// Check if we have some "collect[]" query parameters.
	if len(collectParams) > 0 {
		filters := make(map[string]bool)
		for _, param := range collectParams {
			filters[param] = true
		}

		for _, scraper := range scrapers {
			if filters[scraper.Name()] {
				filteredScrapers = append(filteredScrapers, scraper)
			}
		}
	}
	if len(filteredScrapers) == 0 {
		return scrapers
	}
	return filteredScrapers
}

func init() {
	prometheus.MustRegister(version.NewCollector("gaussDB_exporter"))
}

func newHandler(scrapers []collector.Scraper, logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dsn config.DSN
		var err error
		q := r.URL.Query()
		var target string
		target = ""
		if q.Has("target") {
			target = q.Get("target")
		}
		cfg := c.GetConfig()
		level.Info(logger).Log("Ming",cfg.AuthModules["client"].UserPass.Username)
		cfgauthmodule, ok := cfg.AuthModules["client"]
		if !ok {
			level.Error(logger).Log("msg", "Failed to parse section [client] from config file", "err", err)
		}
		if dsn, err = cfgauthmodule.ConfigureTarget(target); err != nil {
			level.Error(logger).Log("msg", "Failed to form dsn from section [client]", "err", err)
		}

		level.Info(logger).Log("Ming",dsn.String())
		collect := q["collect[]"]
		// Use request context for cancellation when connection gets closed.
		ctx := r.Context()
		// If a timeout is configured via the Prometheus header, add it to the context.
		if v := r.Header.Get("X-Prometheus-Scrape-Timeout-Seconds"); v != "" {
			timeoutSeconds, err := strconv.ParseFloat(v, 64)
			if err != nil {
				level.Error(logger).Log("msg", "Failed to parse timeout from Prometheus header", "err", err)
			} else {
				if *timeoutOffset >= timeoutSeconds {
					// Ignore timeout offset if it doesn't leave time to scrape.
					level.Error(logger).Log("msg", "Timeout offset should be lower than prometheus scrape timeout", "offset", *timeoutOffset, "prometheus_scrape_timeout", timeoutSeconds)
				} else {
					// Subtract timeout offset from timeout.
					timeoutSeconds -= *timeoutOffset
				}
				// Create new timeout context with request context as parent.
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutSeconds*float64(time.Second)))
				defer cancel()
				// Overwrite request with timeout context.
				r = r.WithContext(ctx)
			}
		}

		filteredScrapers := filterScrapers(scrapers, collect)

		registry := prometheus.NewRegistry()

		registry.MustRegister(collector.New(ctx, dsn.String(), filteredScrapers, logger))

		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			registry,
		}
		// Delegate http serving to Prometheus client library, which will call collector.Collect.
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}

func main() {
	scraperFlags := map[collector.Scraper] *bool{}

	for scraper,enableByDefault := range scrapers {
		defaultOn := "false"
		if enableByDefault {
			defaultOn = "true"
		}
		f := kingpin.Flag (
			"collect." + scraper.Name(),
			scraper.Help(),
		).Default(defaultOn).Bool()

		scraperFlags[scraper] = f

	}

	promlogConfig := &promlog.Config{}
	kingpin.Version(version.Print("gaussdb_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg","Starting gaussDB_exporter","version",version.Info())
	level.Info(logger).Log("msg","Build context","build_context",version.BuildContext())

	var err error
	if err = c.ReloadConfig(*configGaussDB,logger); err != nil {
		level.Info(logger).Log("msg", "Error parsing host config", "file", *configGaussDB, "err", err)
		os.Exit(1)
	}

	// Register only scrapers enabled by flag.
	enabledScrapers := []collector.Scraper{}
	for scraper, enabled := range scraperFlags {
		if *enabled {
			level.Info(logger).Log("msg", "Scraper enabled", "scraper", scraper.Name())
			enabledScrapers = append(enabledScrapers, scraper)
		}
	}
	handlerFunc := newHandler(enabledScrapers, logger)
	http.Handle(*metricsPath, promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handlerFunc))
	if *metricsPath != "/" && *metricsPath != ""  && *metricsPath!="query"{
		landingConfig := web.LandingConfig{
			Name:        "GaussDB Exporter",
			Description: "Prometheus Exporter for GaussDB servers",
			Version:     version.Info(),
			Links: []web.LandingLinks{
				{
					Address: *metricsPath,
					Text:    "Metrics",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			level.Error(logger).Log("err", err)
			os.Exit(1)
		}
		http.Handle("/", landingPage)
	}

	http.HandleFunc("/query",collector.QueryTest)
	srv := &http.Server{}
	if err := web.ListenAndServe(srv, toolkitFlags, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
