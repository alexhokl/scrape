package scraper

import (
	"strings"
	"testing"
)

func TestCreateLinkScraper_Guardian(t *testing.T) {
	scraper, err := CreateLinkScraper("guardian")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if scraper == nil {
		t.Fatal("expected scraper to be non-nil")
	}

	if _, ok := scraper.(*GuardianScraper); !ok {
		t.Errorf("expected *GuardianScraper, got %T", scraper)
	}
}

func TestCreateLinkScraper_UnsupportedSource(t *testing.T) {
	scraper, err := CreateLinkScraper("unsupported")

	if err == nil {
		t.Fatal("expected error for unsupported source, got nil")
	}

	if scraper != nil {
		t.Errorf("expected nil scraper, got %v", scraper)
	}

	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected error message to contain 'unsupported', got: %v", err)
	}

	if !strings.Contains(err.Error(), "is not supported") {
		t.Errorf("expected error message to contain 'is not supported', got: %v", err)
	}
}

func TestCreateLinkScraper_EmptySource(t *testing.T) {
	scraper, err := CreateLinkScraper("")

	if err == nil {
		t.Fatal("expected error for empty source, got nil")
	}

	if scraper != nil {
		t.Errorf("expected nil scraper, got %v", scraper)
	}
}

func TestCreateArticleScraper_Guardian(t *testing.T) {
	scraper, err := CreateArticleScraper("guardian")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if scraper == nil {
		t.Fatal("expected scraper to be non-nil")
	}

	if _, ok := scraper.(*GuardianScraper); !ok {
		t.Errorf("expected *GuardianScraper, got %T", scraper)
	}
}

func TestCreateArticleScraper_Microsoft(t *testing.T) {
	scraper, err := CreateArticleScraper("microsoft")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if scraper == nil {
		t.Fatal("expected scraper to be non-nil")
	}

	if _, ok := scraper.(*MicrosoftLearnScraper); !ok {
		t.Errorf("expected *MicrosoftLearnScraper, got %T", scraper)
	}
}

func TestCreateArticleScraper_Go(t *testing.T) {
	scraper, err := CreateArticleScraper("go")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if scraper == nil {
		t.Fatal("expected scraper to be non-nil")
	}

	if _, ok := scraper.(*GoDocScraper); !ok {
		t.Errorf("expected *GoDocScraper, got %T", scraper)
	}
}

func TestCreateArticleScraper_Tofugu(t *testing.T) {
	scraper, err := CreateArticleScraper("tofugu")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if scraper == nil {
		t.Fatal("expected scraper to be non-nil")
	}

	if _, ok := scraper.(*TofuguScraper); !ok {
		t.Errorf("expected *TofuguScraper, got %T", scraper)
	}
}

func TestCreateArticleScraper_UnsupportedSource(t *testing.T) {
	scraper, err := CreateArticleScraper("unsupported")

	if err == nil {
		t.Fatal("expected error for unsupported source, got nil")
	}

	if scraper != nil {
		t.Errorf("expected nil scraper, got %v", scraper)
	}

	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected error message to contain 'unsupported', got: %v", err)
	}

	if !strings.Contains(err.Error(), "is not supported") {
		t.Errorf("expected error message to contain 'is not supported', got: %v", err)
	}
}

func TestCreateArticleScraper_EmptySource(t *testing.T) {
	scraper, err := CreateArticleScraper("")

	if err == nil {
		t.Fatal("expected error for empty source, got nil")
	}

	if scraper != nil {
		t.Errorf("expected nil scraper, got %v", scraper)
	}
}

func TestCreateArticleScraper_CaseSensitive(t *testing.T) {
	testCases := []string{"Guardian", "GUARDIAN", "Microsoft", "MICROSOFT", "Go", "GO", "Tofugu", "TOFUGU"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			scraper, err := CreateArticleScraper(tc)

			if err == nil {
				t.Errorf("expected error for case-sensitive source %q, got nil", tc)
			}

			if scraper != nil {
				t.Errorf("expected nil scraper for %q, got %v", tc, scraper)
			}
		})
	}
}

func TestCreateLinkScraper_CaseSensitive(t *testing.T) {
	testCases := []string{"Guardian", "GUARDIAN"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			scraper, err := CreateLinkScraper(tc)

			if err == nil {
				t.Errorf("expected error for case-sensitive source %q, got nil", tc)
			}

			if scraper != nil {
				t.Errorf("expected nil scraper for %q, got %v", tc, scraper)
			}
		})
	}
}
