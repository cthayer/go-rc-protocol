package rc_protocol

import (
    "crypto"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "github.com/pkg/errors"
    "hash"
    "io/ioutil"
    "path"
    "regexp"
    "strings"
    "time"
)

const AUTH_HEADER_NAME = "authorization"
const ISO_8601_FORMAT = "2006-01-02T15:04:05-0700"

type Auth interface {
    CreateSig(name string, keyDir string) (string, error)
    CheckSig(header string, certDir string) (bool, error)
    GetHeaderName() string
    ParseHeader(header string) []string
}

type auth struct {
    headerPattern *regexp.Regexp
    HeaderName string
    sigAlgorithm hash.Hash
    sigAlgorithmCrypto crypto.Hash
    sigEncoding *base64.Encoding
}

func newAuth() Auth {
    auth := auth{
        headerPattern: regexp.MustCompile(`^RC [^\s;]+;[^\s;]+;.+$`),
        sigAlgorithm: sha256.New(),
        sigAlgorithmCrypto: crypto.SHA256,
        sigEncoding: base64.StdEncoding,
    }

    return auth
}

func (a auth) GetHeaderName() string {
    return AUTH_HEADER_NAME
}

func (a auth) CreateSig(name string, keyDir string) (string, error) {
    key, err := a.loadKey(keyDir, name + ".key")

    if err != nil {
        return "", err
    }

    iso8601 := time.Now().Format(ISO_8601_FORMAT)

    a.sigAlgorithm.Reset()
    bytesWritten, err := a.sigAlgorithm.Write([]byte(iso8601))

    if err != nil || bytesWritten < 1 {
        return "", err
    }

    hashSum := a.sigAlgorithm.Sum(nil)

    sig, err := rsa.SignPKCS1v15(nil, key, a.sigAlgorithmCrypto, hashSum)

    if err != nil {
        return "", err
    }

    sigStr := a.sigEncoding.EncodeToString(sig)

    return "RC " + strings.Join([]string{name, iso8601, sigStr}, ";"), nil
}

func (a auth) CheckSig(header string, certDir string) (bool, error) {
    if !a.headerPattern.MatchString(header) {
        // Log "Invalid format for authorization header: <header_value>
        return false, errors.New("Invalid format for authorization header")
    }

    sigArray := a.ParseHeader(header)

    cert, err := a.loadCert(certDir, sigArray[0])

    if err != nil {
        return false, err
    }

    signature, err := a.sigEncoding.DecodeString(sigArray[2])

    if err != nil {
        return false, err
    }

    a.sigAlgorithm.Reset()

    bytesWritten, err := a.sigAlgorithm.Write([]byte(sigArray[1]))

    if err != nil || bytesWritten < 1 {
        return false, err
    }

    hashSum := a.sigAlgorithm.Sum(nil)

    err = rsa.VerifyPKCS1v15(cert, a.sigAlgorithmCrypto, hashSum, signature)

    if err != nil {
        return false, err
    }

    return true, nil
}

func (a *auth) loadPEM(dir string, name string)  (*pem.Block, error) {
    pemStr, err := ioutil.ReadFile(path.Join(dir, name))

    if err != nil {
        // failed to read the pem file
        return nil, err
    }

    block, _ := pem.Decode(pemStr)

    if block == nil {
        // failed to parse pem string
        return nil, errors.New("Failed to parse certificate pem")
    }

    return block, nil
}

func (a *auth) loadKey(dir string, name string) (*rsa.PrivateKey, error) {
    block, err := a.loadPEM(dir, name)

    if err != nil {
        // Log error
        return nil, err
    }

    key, err := x509.ParsePKCS1PrivateKey(block.Bytes)

    if err != nil {
        // failed to load the key (must be RSA private key)
        return nil, err
    }

    return key, nil
}

func (a *auth) loadCert(dir string, name string) (*rsa.PublicKey, error) {
    block, err := a.loadPEM(dir, name)

    if err != nil {
        // Log error
        return nil, err
    }

    cert, err := x509.ParsePKIXPublicKey(block.Bytes)

    if err != nil {
        return nil, err
    }

    switch cert.(type) {
    case *rsa.PublicKey:
        return cert.(*rsa.PublicKey), nil
    default:
        return nil, errors.New("Invalid pem file.  Must be an RSA Public Key.")
    }
}

func (a auth) ParseHeader(header string) []string {
     return strings.SplitN(strings.Replace(header, "RC ", "", 1), ";", 3)
}
