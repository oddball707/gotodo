package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"

	"github.com/sirupsen/logrus"

	d "gotodo/dao"
	h "gotodo/handler"
	s "gotodo/service"

	proto "gotodo/proto"

	"google.golang.org/grpc"
)

var logger *logrus.Entry

func main() {
	// read env vars from the docker image
	env := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env[pair[0]] = pair[1]
	}
	// read the cmd line args
	for i := 1; i+1 < len(os.Args); i += 2 {
		env[os.Args[i]] = os.Args[i+1]
	}

	// initialize the logger
	logrus.SetFormatter(&runtime.Formatter{ChildFormatter: &logrus.JSONFormatter{}})
	logrus.SetOutput(os.Stdout)

	logger.Logger.SetLevel(logrus.InfoLevel)
	if l, ok := env["logLevel"]; ok {
		setLogLevel(l)
	}

	// start listening tcp:host:port
	grpcPort := env["grpcPort"]
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		logger.WithError(err).WithField("grpcPort", grpcPort).Panic("failed to create gRPC listener")
	}

	// initialize dao layer with database interface/session
	dao := d.NewDao(logger)
	srv := s.NewService(logger, dao)
	hnd := h.NewHandler(logger, srv)

	// create grpc server and apply middleware
	grpcServer := grpc.NewServer()

	// register PB with grpcServer
	proto.RegisterTodoServiceServer(grpcServer, hnd)

	// register http handlers for health and readiness checks
	httpPort := env["httpPort"]
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/readiness", readinessHandler)
	http.HandleFunc("/log", changeLogLevelHandler)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil); err != nil {
			// cannot panic, because this probably is an intentional close
			logger.WithError(err).WithField("httpPort", httpPort).Error("http: ListenAndServe()")
		}
	}()
	logger.WithField("port", httpPort).Info("http started")

	// start the gRPC server
	err = grpcServer.Serve(listen)
	if err != nil {
		logger.WithError(err).Panic("gRpc Server failed to start")
	}
	logger.WithField("port", grpcPort).Info("gRPC started")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Trace("healthHandler")
	fmt.Fprintln(w, "Healthy")
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	logger.Trace("readinessHandler")
	fmt.Fprintln(w, "Ready")
}

// changeLogLevelHandler http handler for setting the logger's log level
// Expected url format: POST "http://.../log?level=info"
func changeLogLevelHandler(w http.ResponseWriter, r *http.Request) {
	logger.Trace("changeLogLevelHandler")

	// Only handles POST http requests
	if r.Method == "POST" {
		// Get the level from the url
		values := r.URL.Query()["level"]
		if err := setLogLevel(values[0]); err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprintf(w, "log level set to %s", values[0])
		}
	} else {
		fmt.Fprint(w, "only accepts POST calls")
	}
}

// setLogLevel is a helper method for setting the logger's logging level
// Possible values for level are "panic", "fatal", "error", "warn", "warning", "info", "debug", and "trace"
func setLogLevel(level string) error {
	// Get the log level from the given string
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		// Unable to parse the given level
		logger.WithError(err).WithField("level", level).Warnf("Unknown log level. Log level will remain [%s].", logger.Logger.Level.String())
		return err
	}
	// Able to parse the level, set the logger's log level
	logger.Logger.SetLevel(logLevel)
	logger.WithField("level", level).Info("log level set")
	return nil
}

