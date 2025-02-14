package config

// Config is a basic struct to hold config data.
type Config struct {
    URL string `json:"url,omitempty"`
    // Add more fields as needed for your use case
}

// Load reads config from filePath and returns a Config struct.
// Currently just returns a placeholder.
func Load(filePath string) (*Config, error) {
    // In future, parse JSON from filePath
    return &Config{
        URL: "http://example.com", // Stub data
    }, nil
}
