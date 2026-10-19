package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/DPLakmal/guidenow-platform/go/internal/database"
	"github.com/DPLakmal/guidenow-platform/go/pkg/models"
)

type Parser struct {
	apiKey string
	client *http.Client
}

func NewParser(apiKey string) *Parser {
	return &Parser{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Parse converts a natural-language tourist request into structured filters
func (p *Parser) Parse(ctx context.Context, query string) (*models.AISearchFilter, error) {
	cleanQuery := strings.TrimSpace(query)
	if cleanQuery == "" {
		return &models.AISearchFilter{
			OriginalQuery: query,
			Confidence:    0.0,
			Explanation:   "Empty search query provided.",
		}, nil
	}

	// Try Gemini API if API key is provided
	if p.apiKey != "" {
		if result, err := p.callGemini(ctx, cleanQuery); err == nil && result != nil {
			return result, nil
		}
	}

	// Fallback to high-accuracy heuristic NLP parser
	return p.heuristicParse(cleanQuery), nil
}

func (p *Parser) heuristicParse(query string) *models.AISearchFilter {
	lower := strings.ToLower(query)
	filter := &models.AISearchFilter{
		OriginalQuery: query,
		Confidence:    0.92,
		Languages:     []string{},
		Specialties:   []string{},
		RadiusKm:      20.0,
	}

	// 1. Detect Sri Lankan Hubs & Locations
	for hub, coords := range database.HubCoordinates {
		if strings.Contains(lower, hub) {
			filter.Location = strings.Title(hub)
			lat := coords[0]
			lng := coords[1]
			filter.Latitude = &lat
			filter.Longitude = &lng
			break
		}
	}

	// Special location synonyms
	if filter.Location == "" {
		if strings.Contains(lower, "galle fort") {
			coords := database.HubCoordinates["galle"]
			filter.Location = "Galle"
			filter.Latitude = &coords[0]
			filter.Longitude = &coords[1]
		} else if strings.Contains(lower, "lion rock") || strings.Contains(lower, "pidurangala") {
			coords := database.HubCoordinates["sigiriya"]
			filter.Location = "Sigiriya"
			filter.Latitude = &coords[0]
			filter.Longitude = &coords[1]
		} else if strings.Contains(lower, "temple of the tooth") {
			coords := database.HubCoordinates["kandy"]
			filter.Location = "Kandy"
			filter.Latitude = &coords[0]
			filter.Longitude = &coords[1]
		} else if strings.Contains(lower, "nine arch") || strings.Contains(lower, "adams peak") {
			coords := database.HubCoordinates["ella"]
			filter.Location = "Ella"
			filter.Latitude = &coords[0]
			filter.Longitude = &coords[1]
		}
	}

	// 2. Detect Languages
	langMap := map[string]string{
		"english":  "English",
		"german":   "German",
		"french":   "French",
		"spanish":  "Spanish",
		"italian":  "Italian",
		"russian":  "Russian",
		"japanese": "Japanese",
		"chinese":  "Chinese",
		"sinhala":  "Sinhala",
		"tamil":    "Tamil",
	}
	for word, canonical := range langMap {
		if strings.Contains(lower, word) {
			filter.Languages = append(filter.Languages, canonical)
		}
	}

	// 3. Detect Specialties / Activities
	specialtyRules := map[string][]string{
		"Surfing":        {"surf", "surfing", "waves", "beach"},
		"Colonial History": {"colonial", "fort", "architecture", "dutch", "ramparts"},
		"Street Food":    {"food", "culinary", "street food", "pettah food", "cooking", "cuisine"},
		"Ancient Heritage": {"heritage", "archaeology", "ruins", "frescoes", "temple", "monastery"},
		"Hiking & Treks": {"hiking", "trek", "mountain", "climb", "bridge", "trail"},
		"Wildlife & Marine": {"wildlife", "whale", "whales", "birding", "safari", "nature"},
		"Tea Plantation": {"tea", "plantation", "factory", "ceylon tea"},
	}

	for canonical, keywords := range specialtyRules {
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				filter.Specialties = append(filter.Specialties, canonical)
				break
			}
		}
	}

	// 4. Detect Price Constraints ("under $25", "below 30", "max 40")
	priceRegex := regexp.MustCompile(`(?i)(?:under|below|less than|max(?:imum)?)\s*\$?(\d+(?:\.\d+)?)`)
	if match := priceRegex.FindStringSubmatch(query); len(match) > 1 {
		if val, err := strconv.ParseFloat(match[1], 64); err == nil {
			filter.MaxRate = val
		}
	}

	// 5. Detect Rating Constraints ("top rated", "4.5+", "best")
	if strings.Contains(lower, "top rated") || strings.Contains(lower, "best") || strings.Contains(lower, "highly rated") {
		filter.MinRating = 4.8
	} else {
		ratingRegex := regexp.MustCompile(`(?i)(\d(?:\.\d)?)\s*(?:\+|stars?|rating)`)
		if match := ratingRegex.FindStringSubmatch(query); len(match) > 1 {
			if val, err := strconv.ParseFloat(match[1], 64); err == nil && val <= 5.0 {
				filter.MinRating = val
			}
		}
	}

	// 6. Detect Immediate Availability
	if strings.Contains(lower, "now") || strings.Contains(lower, "immediate") || strings.Contains(lower, "right now") || strings.Contains(lower, "today") {
		filter.ImmediateOnly = true
	}

	// Explanation
	var explParts []string
	if filter.Location != "" {
		explParts = append(explParts, fmt.Sprintf("in %s", filter.Location))
	}
	if len(filter.Languages) > 0 {
		explParts = append(explParts, fmt.Sprintf("speaking %s", strings.Join(filter.Languages, "/")))
	}
	if len(filter.Specialties) > 0 {
		explParts = append(explParts, fmt.Sprintf("specializing in %s", strings.Join(filter.Specialties, ", ")))
	}
	if filter.MaxRate > 0 {
		explParts = append(explParts, fmt.Sprintf("under $%.0f/hr", filter.MaxRate))
	}
	if filter.ImmediateOnly {
		explParts = append(explParts, "available immediately")
	}

	if len(explParts) > 0 {
		filter.Explanation = fmt.Sprintf("Searching for guides %s.", strings.Join(explParts, ", "))
	} else {
		filter.Explanation = "Showing available verified guides nearby."
	}

	return filter
}

func (p *Parser) callGemini(ctx context.Context, query string) (*models.AISearchFilter, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", p.apiKey)

	systemPrompt := `You are an AI intent parser for GuideNow, a Sri Lanka tourist guide matching platform.
Parse the user's travel request into a structured JSON filter matching this exact schema:
{
  "location": "City or Region name (e.g., Galle, Colombo, Kandy, Sigiriya, Weligama, Ella)",
  "languages": ["English", "German", etc],
  "specialties": ["Surfing", "Colonial History", "Street Food", etc],
  "max_rate": 25.0,
  "min_rating": 4.5,
  "immediate_only": true,
  "explanation": "Short summary of extracted criteria"
}
Return ONLY valid raw JSON with no markdown formatting.`

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": fmt.Sprintf("%s\n\nUser Query: %s", systemPrompt, query)},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API responded with status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil || len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("invalid gemini response format")
	}

	rawText := geminiResp.Candidates[0].Content.Parts[0].Text
	rawText = strings.TrimPrefix(rawText, "```json")
	rawText = strings.TrimPrefix(rawText, "```")
	rawText = strings.TrimSuffix(rawText, "```")
	rawText = strings.TrimSpace(rawText)

	var filter models.AISearchFilter
	if err := json.Unmarshal([]byte(rawText), &filter); err != nil {
		return nil, err
	}

	filter.OriginalQuery = query
	filter.Confidence = 0.98

	// Resolve location to lat/lng
	locLower := strings.ToLower(filter.Location)
	if coords, ok := database.HubCoordinates[locLower]; ok {
		lat := coords[0]
		lng := coords[1]
		filter.Latitude = &lat
		filter.Longitude = &lng
		filter.RadiusKm = 20.0
	}

	return &filter, nil
}
