package server

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	hpprof "net/http/pprof"
	"os"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
	"strings"

	"golang.conradwood.net/go-easyops/appinfo"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	pp "golang.conradwood.net/go-easyops/profiling"
	"google.golang.org/grpc"
)

func debugHandler(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	z := strings.Split(p, "/")
	if len(z) == 0 {
		fmt.Printf("Invalid debug request: %s\n", p)
		return
	}
	lp := z[len(z)-1]
	fmt.Printf("Last part: %s\n", lp)
	if lp == "cpu" {
		debugCpuHandler(w, req)
		return

	}

	if lp == "heapdump" {
		writeHeap(w, req)
		return
	}
	if lp == "info" {
		hpprof.Index(w, req)
		return
	}
	if lp == "goroutine" { // tested, works
		profile := pprof.Lookup(lp)
		if profile != nil {
			serve_debug_profile(profile, w, req)
			return
		}
	}
	h := hpprof.Handler(lp)
	if h == nil {
		fmt.Printf("[go-easyops] no such handler:%s\n", lp)
		return
	}
	h.ServeHTTP(w, req)
	//todo
}
func serve_debug_profile(p *pprof.Profile, w http.ResponseWriter, req *http.Request) {
	buf := &bytes.Buffer{}
	p.WriteTo(buf, 1)
	b := buf.Bytes()
	b = bytes.ReplaceAll(b, []byte("\n"), []byte("<br/>"))
	bold := []string{"golang.conradwood.net", "golang.singingcat.net", "golang.yacloud.eu"}
	for _, bol := range bold {
		b = bytes.ReplaceAll(b, []byte(bol), []byte("<b>"+bol+"</b>"))
	}
	w.Header()["Content-Type"] = []string{"text/html"}
	w.Write([]byte("<html><body>"))
	w.Write(b)
	w.Write([]byte("</body></html>"))
}
func writeHeap(w http.ResponseWriter, req *http.Request) {
	filename := "dump"
	f, err := os.Create(filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failure opening file %s", err), 404)
		return
	}
	debug.WriteHeapDump(f.Fd())
	f.Close()
	//Check if file exists and open
	fd, err := os.Open(filename)
	if err != nil {
		//File not found, send 404
		http.Error(w, fmt.Sprintf("File not open:%s", err), 500)
		return
	}
	defer fd.Close() //Close after function return

	//Get the file size
	FileStat, _ := fd.Stat()                           //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", FileSize)

	io.Copy(w, fd) //'Copy' the file to the client
	return
}

func debugCpuHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to parse form %s", err)
		return
	}
	if len(req.Form["Download"]) != 0 {
		if pp.IsActive() {
			w.WriteHeader(409)
			fmt.Fprintf(w, "Download unavailable whilst profiling is active")
			return
		}
		b := pp.GetBuf()
		if b.Len() == 0 {
			w.WriteHeader(404)
			fmt.Fprintf(w, "No profiling data available. Enable profiling for a longer period of time perhaps?")
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename=cpuprofile")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", b.Len()))
		w.Write(b.Bytes())
		return
	}
	fmt.Fprintf(w, ("<html><body>"))
	if len(req.Form["toggle"]) != 0 {
		fmt.Fprintf(w, "toggled<br/>")
		pp.Toggle()
	}
	s := "Inactive"
	if pp.IsActive() {
		s = "Active"
	}
	fmt.Fprintf(w, "CPU Profiling: %s\n", s)
	fmt.Fprintf(w, "&nbsp<a href=\"?toggle\">Toggle</a></br>")
	if pp.IsActive() {
		fmt.Fprintf(w, "Download unavailable whilst profiling is active<br/>")
	} else {
		b := pp.GetBuf()
		if b.Len() == 0 {
			fmt.Fprintf(w, "Download unavailable - enable profiling first<br/>")
		} else {
			fmt.Fprintf(w, "<a href=\"?Download\">Download</a> most recent profile</br>")
		}
	}
	fmt.Fprintf(w, "</body></html>")
	return
}
func helpHandler(w http.ResponseWriter, req *http.Request, sd *serverDef) {
	s := "<html><body>"
	s = s + "<a href=\"/internal/pleaseshutdown\">shutdown</a><br/>"
	s = s + "<a href=\"/internal/health\">server health</a><br/>"
	s = s + "<a href=\"/internal/service-info/version\">VersionInfo</a><br/>"
	s = s + "<a href=\"/internal/service-info/metrics\">metrics</a><br/>"
	s = s + "<a href=\"/internal/clearcache\">clearcache</a> (append /name to clear a specific cache)<br/>"
	s = s + "<a href=\"/internal/parameters\">parameters</a><br/>"
	s = s + "<a href=\"/internal/service-info/grpc-connections\">GRPC Connections</a><br/>"
	s = s + "<a href=\"/internal/service-info/grpc-callers\">GRPC Server Caller list (who called this service)</a><br/>"
	s = s + "<a href=\"/internal/service-info/dependencies\">Registered GRPC Dependencies</a><br/>"
	s = s + "<a href=\"/internal/debug/info\">Go-Profiler</a><br/>"
	s = s + "<a href=\"/internal/debug/cpu\">CPU Profiler</a><br/>"
	s = s + "<a href=\"/internal/debug/heapdump\">Download Heap Dump</a><br/>"
	s = s + "</body></html>"
	fmt.Fprintf(w, "%s", s)
}

func healthzHandler(w http.ResponseWriter, req *http.Request, sd *serverDef) {
	fmt.Fprintf(w, getHealthString())
}

// this services the /service-info/ url
func serveServiceInfo(w http.ResponseWriter, req *http.Request, sd *serverDef) {
	p := req.URL.Path
	if strings.HasPrefix(p, "/internal/service-info/name") {
		fmt.Fprintf(w, (sd.name))
	} else if strings.HasPrefix(p, "/internal/service-info/version") {
		serveVersion(w, req, sd)
	} else if strings.HasPrefix(p, "/internal/service-info/grpc-connections") {
		serveGRPCConnections(w, req, sd)
	} else if strings.HasPrefix(p, "/internal/service-info/grpc-callers") {
		serveGRPCCallers(w, req, sd)
	} else if strings.HasPrefix(p, "/internal/service-info/dependencies") {
		serveDependencies(w, req, sd)
	} else if strings.HasPrefix(p, "/internal/service-info/metrics") {
		fmt.Printf("Request path: \"%s\"\n", p)
		m := strings.TrimPrefix(p, "/internal/service-info/metrics")
		m = strings.TrimLeft(m, "/")
	} else {
		fmt.Printf("Invalid path: \"%s\"\n", p)
	}
}

// serve /internal/service-info/grpc-callers
func serveGRPCCallers(w http.ResponseWriter, req *http.Request, sd *serverDef) {
	usage_info := GetUsageInfo()
	for _, service := range usage_info.Services() {
		for _, method := range service.Methods() {
			for _, callers := range method.Callers() {
				fmt.Fprintf(w, "%s.%s %s\n", service.Name(), method.Name(), callers.String())
			}
		}
	}
}

// serve /internal/service-info/dependencies
func serveDependencies(w http.ResponseWriter, req *http.Request, sd *serverDef) {
	s := client.GetDependencies()
	fmt.Fprintf(w, "# %d registered dependencies\n", len(s))
	for _, r := range s {
		fmt.Fprintf(w, "%s\n", r)
	}
}

// serve /internal/service-info/grpc-connections
func serveGRPCConnections(w http.ResponseWriter, req *http.Request, sd *serverDef) {
	s := common.GetConnectionNames()
	fmt.Fprintf(w, "# %d requested connections\n", len(s))
	for _, r := range s {
		fmt.Fprintf(w, "%s\n", r.Name)
	}
	bs := common.GetBlockedConnectionNames()
	fmt.Fprintf(w, "# %d blocked connections\n", len(bs))
	for _, r := range bs {
		fmt.Fprintf(w, "%s\n", r.Name)
	}
}

// services the version url /internal/version/go-framework
func serveVersion(w http.ResponseWriter, req *http.Request, sd *serverDef) {
	fmt.Fprintf(w, "go_framework_buildid: %d\n", cmdline.BUILD_NUMBER)
	fmt.Fprintf(w, "go_framework_timestamp: %d\n", cmdline.BUILD_TIMESTAMP)
	fmt.Fprintf(w, "go_framework_description: %s\n", cmdline.BUILD_DESCRIPTION)
	fmt.Fprintf(w, "app_buildid: %d\n", appinfo.AppInfo().Number)
	fmt.Fprintf(w, "app_timestamp: %d\n", appinfo.AppInfo().Timestamp)
	fmt.Fprintf(w, "app_description: %s\n", appinfo.AppInfo().Description)
	fmt.Fprintf(w, "app_repository: %s\n", appinfo.AppInfo().RepositoryName)
	fmt.Fprintf(w, "app_repository_id: %d\n", appinfo.AppInfo().RepositoryID)
	fmt.Fprintf(w, "app_artefact_id: %d\n", appinfo.AppInfo().ArtefactID)
	fmt.Fprintf(w, "app_commit: %s\n", appinfo.AppInfo().CommitID)

}

// this servers /internal/parameters url
func paraHandler(w http.ResponseWriter, req *http.Request, sd *serverDef) {
	errno := 402
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Failed to parse request: %s\n", err)
		return
	}
	if len(req.Form) == 0 {
		flag.VisitAll(func(f *flag.Flag) {
			s := "SET"
			if fmt.Sprintf("%v", f.Value) == fmt.Sprintf("%v", f.DefValue) {
				s = "DEFAULT"
			}
			fmt.Fprintf(w, "%s %s %s %s\n", "STRING", s, f.Name, f.Value)
		})
		return
	}
	for name, value := range req.Form {
		if len(value) != 1 {
			http.Error(w, fmt.Sprintf("odd number of values for %s: %d (expected 1)\n", name, len(value)), errno)
			return
		}
		//fmt.Fprintf(w, "Setting %s to %s\n", name, value)
		f := flag.Lookup(name)
		if f == nil {
			http.Error(w, "No such flag\n", errno)
			return
		}
		err = f.Value.Set(value[0])
		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot set value of %s to %s: %s\n", name, value, err), errno)
			return
		}
		err = ipc_send_new_para(sd, name, value[0])
		if err != nil {
			fmt.Printf("[go-easyops] failed to send parameter change via ipc (%s=%s): %s\n", name, value[0], err)
			// no further action, considering this somewhat optional for now
		}

	}
	pp.ProfilingCheckStart() // make it pick up on changes to flag if any
	fmt.Fprintf(w, "Done")
}

// this services the /pleaseshutdown url
func pleaseShutdown(w http.ResponseWriter, req *http.Request, s *grpc.Server) {
	stopping(make(chan bool, 10))
	fmt.Fprintf(w, "OK\n")
	fmt.Printf("Received request to shutdown.\n")
	s.Stop()   // maybe even s.GracefulStop() ?
	os.Exit(0) // i'd prefer not to exit here unless something is relying on it.
}
