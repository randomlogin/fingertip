package auto

import (
	"bufio"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

type ProxyStatus int

const (
	ProxyStatusNone      ProxyStatus = 0
	ProxyStatusInstalled             = 1
	ProxyStatusConflict              = 2
)

// firefox
const (
	userProfileHeader = "// AUTOCONFIG:FINGERTIP"
	userPref          = userProfileHeader +
		` (autogenerated - remove this line to disable auto config for this profile and edit the file) 
user_pref("security.enterprise_roots.enabled", true); 
user_pref("network.proxy.type", 5);
`
)

func equalURL(a, b string) bool {
	a = strings.TrimSuffix(strings.TrimSpace(a), "/")
	b = strings.TrimSuffix(strings.TrimSpace(b), "/")
	return strings.EqualFold(a, b)
}

// ConfigureFirefox instructs firefox to read system certs and proxy
// settings
func ConfigureFirefox() error {
	profiles, err := getProfilePaths()
	if err != nil {
		return err
	}

	var lastErr error
	for _, profilePath := range profiles {
		if err := writeUserConfig(profilePath); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// UndoFirefoxConfiguration undoes all actions made by ConfigureFirefox
func UndoFirefoxConfiguration() {
	profiles, err := getProfilePaths()
	if err != nil {
		return
	}

	for _, profilePath := range profiles {
		userjs := path.Join(profilePath, "user.js")
		if ok, _, _ := fileLineContains(userjs, userProfileHeader); ok {
			_ = os.Remove(userjs)
		}
	}
}

func writeUserConfig(profilePath string) (err error) {
	prefs := path.Join(profilePath, "prefs.js")
	// if user has existing proxy configuration
	// ignore this profile
	if ok, line, _ := fileLineContains(prefs, `"network.proxy.type"`); ok {
		// loosely check if its type 0 (no proxy) or 5 (system proxy)
		if !strings.ContainsRune(line, '0') && !strings.ContainsRune(line, '5') {
			return errors.New("profile with existing proxy configuration")
		}
	}

	userjs := path.Join(profilePath, "user.js")
	if _, err = os.Stat(userjs); err == nil {
		if ok, _, _ := fileLineContains(userjs, userProfileHeader); !ok {
			return nil
		}
	}

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("failed checking whether file exists: %v", err)
	}

	return os.WriteFile(userjs, []byte(userPref), 0644)
}

func fileLineContains(file, substr string) (bool, string, error) {
	f, err := os.Open(file)
	if err != nil {
		return false, "", err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.Contains(line, substr) {
			return true, line, nil
		}
	}
	return false, "", nil
}

func getProfilePaths() ([]string, error) {
	c, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed reading user config dir: %v", err)
	}

	pp := path.Join(c, relativeProfilesPath)
	dirs, err := os.ReadDir(pp)
	if err != nil {
		return nil, fmt.Errorf("failed listing profiles: %v", err)
	}

	var paths []string
	for _, dir := range dirs {
		if strings.Contains(dir.Name(), "default") {
			paths = append(paths, path.Join(pp, dir.Name()))
		}
	}

	return paths, nil
}

func readPEM(certPath string) ([]byte, error) {
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed reading root certificate: %v", err)
	}

	// Decode PEM
	certBlock, _ := pem.Decode(cert)
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed decoding cert invalid PEM data")
	}

	return certBlock.Bytes, nil
}

func readX509Cert(certPath string) (*x509.Certificate, error) {
	cert, err := readPEM(certPath)
	if err != nil {
		return nil, err
	}
	c, err := x509.ParseCertificate(cert)
	if err != nil {
		return nil, err
	}

	return c, nil
}
