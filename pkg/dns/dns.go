package dns

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

// Server is a struct holding all of the state information about your LDAP server
type Server struct {
	port       int            // The port number of the DNS server.
	configFile string         // The configFile containing the configuration for the server
	logger     zerolog.Logger // The logger being used for console printing.
}

// NewServer creates a new server instance which performs DNS resolution according to a config file.
func NewServer(configFile string, port int) *Server {
	server := &Server{configFile: configFile, port: port, logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr})}
	return server
}

// Listen starts the server up on the specified port and begins listening for connections.
func (s *Server) Listen() error {
	// start the server
	listen := fmt.Sprintf("0.0.0.0:%d", s.port)
	s.logger.Info().Msgf("starting DNS server on %s", listen)

	// Start waiting for configuration changes
	go s.WatchForConfigChanges()

	// Load the initial configuration the first time
	s.ReloadConfiguration(s.configFile)

	// Listen
	return nil
}

// WatchForConfigChanges starts watching the configuration file for writes and applies changes automatically.
func (s *Server) WatchForConfigChanges() {
	// Set up the file watcher
	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	w.SetMaxEvents(1)
	go func() {
		for {
			select {
			case event := <-w.Event:
				s.ReloadConfiguration(event.Path)
			case err := <-w.Error:
				s.logger.Err(err)
			case <-w.Closed:
				return
			}
		}
	}()
	// Watch this file for changes.
	if err := w.Add(s.configFile); err != nil {
		s.logger.Error().Err(err)
		return
	}

	if err := w.Start(10 * time.Second); err != nil {
		s.logger.Error().Err(err)
		return
	}
}

// ReloadConfiguration reads the configuration file and applies the changes.
func (s *Server) ReloadConfiguration(filename string) {
	s.logger.Info().Msgf("reloading configuration file '%s'", filename)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to read file")
		return
	}
	var c map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to unmarshall file")
		return
	}
}
