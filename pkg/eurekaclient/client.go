package eurekaclient

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"github.com/ArthurHlt/go-eureka-client/eureka"
)

// patchClient safely replaces the internal HTTP client used by the Eureka client
// with a custom, more robust configuration.
//
// This function uses reflection and unsafe pointers to modify the unexported
// `httpClient` field inside the `eureka.Client` struct, since the library does not
// expose a setter. The new HTTP client includes custom timeouts, TLS settings, and
// a Dialer with keep-alive support to improve connection reliability.
//
// Parameters:
//   - c: A pointer to an initialized *eureka.Client instance whose internal
//     HTTP client will be replaced.
//
// Behavior:
//   - Enables TLS connections with `InsecureSkipVerify: true` (use only if you trust the network).
//   - Sets reasonable connection and request timeouts.
//   - Ensures long-lived keep-alive connections to reduce overhead on frequent requests.
//
// Note:
//
//	This function relies on `unsafe` and reflection to modify unexported fields.
//	Use it with caution, as internal library changes may break compatibility in
//	future versions.
//
// Example:
//
//	client := eureka.NewClient([]string{"https://eureka.local:8761/eureka"})
//	patchClient(client)
func patchClient(c *eureka.Client) {
	v := reflect.ValueOf(c).Elem()
	httpClientField := v.FieldByName("httpClient")
	if !httpClientField.IsValid() {
		panic("httpClient field not found in eureka.Client")
	}

	// Build our proper client
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	newClient := &http.Client{
		Transport: tr,
		Timeout:   60 * time.Second,
	}

	// Safely set pointer (handle both nil and non-nil)
	ptrToField := unsafe.Pointer(httpClientField.UnsafeAddr())
	realPtr := (**http.Client)(ptrToField)
	*realPtr = newClient
}

// InitEurekaClient initializes and registers a service instance with a Eureka server.
//
// It handles client creation, registration, periodic heartbeats, and graceful
// de-registration on shutdown signals (SIGINT, SIGTERM). Internally, it patches
// the Eureka client’s HTTP transport to allow TLS connections with self-signed
// certificates and reasonable timeouts.
//
// Parameters:
//   - machines: A slice of Eureka server URLs (e.g. []string{"https://host:port/eureka"}).
//   - hostName: The hostname of the service instance (unique identifier for Eureka).
//   - app: The Eureka application name (service ID to register under).
//   - ip: The service’s internal IP address or host used for registration.
//   - port: The port on which the service is listening.
//   - ttl: Lease renewal interval (in seconds) — determines heartbeat frequency.
//   - isSsl: Whether the service endpoint uses HTTPS (true) or HTTP (false).
//
// Example:
//
//	eurekaclient.InitEurekaClient(
//	    []string{"https://eureka.local:8761/eureka"},
//	    "my-service-host",
//	    "MY-SERVICE",
//	    "10.0.0.5",
//	    8080,
//	    30,
//	    false,
//	)
//
// On startup, the function registers the instance with Eureka, then sends heartbeats
// periodically. When the process receives an interrupt or terminate signal, it
// automatically de-registers the instance before exiting.
func InitEurekaClient(machines []string, hostName string, app string, ip string, port int, ttl uint, isSsl bool) {
	client := eureka.NewClient(
		machines,
	)

	patchClient(client)

	instance := eureka.NewInstanceInfo(
		hostName, // The hostname your service is running on
		app,      // The service ID (APP_NAME) in Eureka
		ip,       // Hostname for internal IP
		port,     // The port your Gin app will run on
		ttl,      // Lease renewal interval (seconds)
		isSsl,    // Is this a secure port? (SSL)
	)

	instance.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	instance.Metadata.Map["management.port"] = strconv.Itoa(port)

	// 2. Register the application with Eureka
	err := client.RegisterInstance(app, instance)
	if err != nil {
		log.Fatalf("Failed to register with Eureka: %s", err)
	}
	log.Printf("Successfully registered with Eureka as %s", app)

	// 3. Start sending heartbeats in a separate goroutine
	go func() {
		for {
			err := client.SendHeartbeat(instance.App, instance.HostName)
			if err != nil {
				log.Printf("Failed to send heartbeat: %s", err)
				// If heartbeat fails, you might want to try re-registering
			} else {
				log.Println("Sent heartbeat to Eureka")
			}
			// Wait for the duration of the lease renewal interval
			// (You can also use time.Ticker for a cleaner approach)
			time.Sleep(time.Duration(ttl) * time.Second)
		}
	}()

	// 4. Set up a channel to listen for OS interrupt signals (like Ctrl+C)
	// This is crucial for graceful shutdown and de-registering from Eureka.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Block until a signal is received
		<-sigChan

		log.Println("Shutting down... de-registering from Eureka.")

		// De-register the instance from Eureka
		err := client.UnregisterInstance(instance.App, instance.HostName)
		if err != nil {
			log.Printf("Failed to de-register from Eureka: %s", err)
		} else {
			log.Println("Successfully de-registered from Eureka.")
		}

		// Exit the application
		os.Exit(0)
	}()
}
