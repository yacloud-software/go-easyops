package server

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	pp "golang.conradwood.net/go-easyops/profiling"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
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
		pprof.Index(w, req)
		return
	}
	h := pprof.Handler(lp)
	if h == nil {
		fmt.Printf("[go-easyops] no such handler:%s\n", lp)
		return
	}
	h.ServeHTTP(w, req)
	//todo
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
	s = s + "<a href=\"/internal/service-info/version\">VersionInfo</a><br/>"
	s = s + "<a href=\"/internal/service-info/metrics\">metrics</a><br/>"
	s = s + "<a href=\"/internal/clearcache\">clearcache</a> (append /name to clear a specific cache)<br/>"
	s = s + "<a href=\"/internal/parameters\">parameters</a><br/>"
	s = s + "<a href=\"/internal/service-info/grpc-connections\">GRPC Connections</a><br/>"
	s = s + "<a href=\"/internal/debug/info\">Go-Profiler</a><br/>"
	s = s + "<a href=\"/internal/debug/cpu\">CPU Profiler</a><br/>"
	s = s + "<a href=\"/internal/debug/heapdump\">Download Heap Dump</a><br/>"
	s = s + "</body></html>"
	fmt.Fprintf(w, "%s", s)
}

func healthzHandler(w http.ResponseWriter, req *http.Request, sd *serverDef) {
	fmt.Fprintf(w, "OK")
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
	} else if strings.HasPrefix(p, "/internal/service-info/metrics") {
		fmt.Printf("Request path: \"%s\"\n", p)
		m := strings.TrimPrefix(p, "/internal/service-info/metrics")
		m = strings.TrimLeft(m, "/")
	} else {
		fmt.Printf("Invalid path: \"%s\"\n", p)
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
	fmt.Fprintf(w, "app_buildid: %d\n", cmdline.APP_BUILD_NUMBER)
	fmt.Fprintf(w, "app_timestamp: %d\n", cmdline.APP_BUILD_TIMESTAMP)
	fmt.Fprintf(w, "app_description: %s\n", cmdline.APP_BUILD_DESCRIPTION)
	fmt.Fprintf(w, "app_repository: %s\n", cmdline.APP_BUILD_REPOSITORY)
	fmt.Fprintf(w, "app_repository_id: %d\n", cmdline.APP_BUILD_REPOSITORY_ID)
	fmt.Fprintf(w, "app_commit: %s\n", cmdline.APP_BUILD_COMMIT)
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
		fmt.Fprintf(w, "No parameter specified (try https://.../internal/parameters?debug=true)\n")
		http.Error(w, "no parameter specified", errno)
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
	}
	pp.ProfilingCheckStart() // make it pick up on changes to flag if any
	fmt.Fprintf(w, "Done")
}

// this services the /pleaseshutdown url
func pleaseShutdown(w http.ResponseWriter, req *http.Request, s *grpc.Server) {
	stopping()
	fmt.Fprintf(w, "OK\n")
	fmt.Printf("Received request to shutdown.\n")
	s.Stop()   // maybe even s.GracefulStop() ?
	os.Exit(0) // i'd prefer not to exit here unless something is relying on it.
}
