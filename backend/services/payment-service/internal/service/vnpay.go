package service

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"payment-service/internal/config"
	"payment-service/internal/model"
)

// GenerateVNPayURL creates the payment URL with secure HMAC SHA512 signature
func GenerateVNPayURL(cfg *config.Config, p *model.Payment, ipAddr string) string {
	vnpParams := map[string]string{
		"vnp_Version":    "2.1.0",
		"vnp_Command":    "pay",
		"vnp_TmnCode":    cfg.VnpTmnCode,
		"vnp_Amount":     fmt.Sprintf("%.0f", p.Amount*100), // VNPay requires amount * 100
		"vnp_CreateDate": time.Now().Format("20060102150405"),
		"vnp_CurrCode":   "VND",
		"vnp_IpAddr":     ipAddr,
		"vnp_Locale":     "vn",
		"vnp_OrderInfo":  fmt.Sprintf("Thanh toan don hang %s", p.OrderID.String()),
		"vnp_OrderType":  "other",
		"vnp_ReturnUrl":  cfg.VnpReturnUrl,
		"vnp_TxnRef":     p.ID.String(),
	}

	// Sort keys alphabetically (VNPay requirement)
	var keys []string
	for k := range vnpParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signData strings.Builder
	var queryData strings.Builder

	for i, k := range keys {
		v := vnpParams[k]
		if i > 0 {
			signData.WriteString("&")
			queryData.WriteString("&")
		}
		// Hash data doesn't encode values, but query string does
		signData.WriteString(fmt.Sprintf("%s=%s", k, v))
		queryData.WriteString(fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}

	// Generate HMAC SHA512 signature
	mac := hmac.New(sha512.New, []byte(cfg.VnpHashSecret))
	mac.Write([]byte(signData.String()))
	secureHash := hex.EncodeToString(mac.Sum(nil))

	// Append signature to query
	queryData.WriteString(fmt.Sprintf("&vnp_SecureHash=%s", secureHash))

	return fmt.Sprintf("%s?%s", cfg.VnpUrl, queryData.String())
}

// VerifyVNPaySignature checks the integrity of VNPay callbacks
func VerifyVNPaySignature(secret string, params url.Values) bool {
	secureHash := params.Get("vnp_SecureHash")
	if secureHash == "" {
		return false
	}

	// Extract and sort keys excluding vnp_SecureHash and vnp_SecureHashType
	var keys []string
	for k := range params {
		if k != "vnp_SecureHash" && k != "vnp_SecureHashType" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var signData strings.Builder
	for i, k := range keys {
		if i > 0 {
			signData.WriteString("&")
		}
		// In callback, url.Values contains unescaped values. VNPay requires hashing unescaped.
		signData.WriteString(fmt.Sprintf("%s=%s", k, url.QueryEscape(params.Get(k)))) 
		// Note: VNPay's doc for IPN/Return can be tricky. Usually it hashes URL-encoded values for callback
	}

	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(signData.String()))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	return secureHash == expectedHash
}
