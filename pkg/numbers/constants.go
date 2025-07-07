package numbers

// Default values for new documents and tables
const (
	DefaultRowCount    = 12
	DefaultColCount    = 8
	DefaultHeaderRows  = 1
	DefaultHeaderCols  = 1
	DefaultFontSize    = 12.0
	DefaultFontName    = "Helvetica"
	DefaultTextInset   = 1.0
	DefaultBorderWidth = 0.35
)

// Maximum limits
const (
	MaxRowCount    = 1000000
	MaxColCount    = 1000
	MaxSheetCount  = 250
	MaxTableCount  = 100
)

// Supported Numbers versions (placeholder)
var SupportedNumbersVersions = []string{
	"10.3",
	"11.0", 
	"11.1",
	"11.2",
	"12.0",
	"12.1",
	"13.0",
	"13.1",
	"14.0",
	"14.1",
}

// Error messages
const (
	ErrInvalidCellAddress    = "invalid cell address"
	ErrCellOutOfBounds      = "cell position out of bounds" 
	ErrInvalidSheetName     = "invalid sheet name"
	ErrInvalidTableName     = "invalid table name"
	ErrSheetNotFound        = "sheet not found"
	ErrTableNotFound        = "table not found"
	ErrInsufficientArgs     = "insufficient arguments"
	ErrInvalidRange         = "invalid range"
	ErrUnsupportedFormat    = "unsupported file format"
	ErrEncryptedDocument    = "encrypted documents are not supported"
)

// MIME types
const (
	NumbersMimeType = "application/vnd.apple.numbers"
	ZipMimeType     = "application/zip"
)

// File extensions
const (
	NumbersExtension = ".numbers"
	ZipExtension     = ".zip"
)

// Default colors
var (
	ColorBlack = RGB{R: 0, G: 0, B: 0}
	ColorWhite = RGB{R: 255, G: 255, B: 255}
	ColorRed   = RGB{R: 255, G: 0, B: 0}
	ColorGreen = RGB{R: 0, G: 255, B: 0}
	ColorBlue  = RGB{R: 0, G: 0, B: 255}
)

// Common border styles
var (
	DefaultBorder = Border{
		Width: DefaultBorderWidth,
		Color: ColorBlack,
		Style: BorderSolid,
	}
	
	NoBorder = Border{
		Width: 0,
		Color: ColorBlack,
		Style: BorderNone,
	}
)