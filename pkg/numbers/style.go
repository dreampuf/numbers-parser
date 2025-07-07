package numbers

// RGB represents a color in RGB format
type RGB struct {
	R, G, B uint8
}

// NewRGB creates a new RGB color
func NewRGB(r, g, b uint8) RGB {
	return RGB{R: r, G: g, B: b}
}

// HorizontalAlignment represents horizontal text alignment
type HorizontalAlignment int

const (
	AlignLeft HorizontalAlignment = iota
	AlignRight
	AlignCenter
	AlignJustified
	AlignAuto
)

// VerticalAlignment represents vertical text alignment
type VerticalAlignment int

const (
	AlignTop VerticalAlignment = iota
	AlignMiddle
	AlignBottom
)

// Alignment represents text alignment
type Alignment struct {
	Horizontal HorizontalAlignment
	Vertical   VerticalAlignment
}

// NewAlignment creates a new alignment
func NewAlignment(horizontal HorizontalAlignment, vertical VerticalAlignment) *Alignment {
	return &Alignment{
		Horizontal: horizontal,
		Vertical:   vertical,
	}
}

// BackgroundImage represents a background image for a cell
type BackgroundImage struct {
	Data     []byte
	Filename string
}

// NewBackgroundImage creates a new background image
func NewBackgroundImage(data []byte, filename string) *BackgroundImage {
	return &BackgroundImage{
		Data:     data,
		Filename: filename,
	}
}

// Style represents cell formatting and styling
type Style struct {
	Name           string
	Alignment      *Alignment
	BackgroundImage *BackgroundImage
	BackgroundColor *RGB
	FontColor      RGB
	FontSize       float64
	FontName       string
	Bold           bool
	Italic         bool
	Strikethrough  bool
	Underline      bool
	FirstIndent    float64
	LeftIndent     float64
	RightIndent    float64
	TextInset      float64
	TextWrap       bool
}

// NewStyle creates a new style with default values
func NewStyle() *Style {
	return &Style{
		Alignment:      NewAlignment(AlignAuto, AlignTop),
		FontColor:      NewRGB(0, 0, 0),
		FontSize:       12.0,
		FontName:       "Helvetica",
		Bold:           false,
		Italic:         false,
		Strikethrough:  false,
		Underline:      false,
		FirstIndent:    0.0,
		LeftIndent:     0.0,
		RightIndent:    0.0,
		TextInset:      1.0,
		TextWrap:       true,
	}
}

// BorderType represents the type of border
type BorderType int

const (
	BorderSolid BorderType = iota
	BorderDashes
	BorderDots
	BorderNone
)

// Border represents a cell border
type Border struct {
	Width float64
	Color RGB
	Style BorderType
}

// NewBorder creates a new border
func NewBorder(width float64, color RGB, style BorderType) *Border {
	return &Border{
		Width: width,
		Color: color,
		Style: style,
	}
}

// CellBorder represents borders for all sides of a cell
type CellBorder struct {
	Top    *Border
	Right  *Border
	Bottom *Border
	Left   *Border
}

// NewCellBorder creates a new cell border
func NewCellBorder() *CellBorder {
	return &CellBorder{}
}

// SetAllBorders sets the same border for all sides
func (cb *CellBorder) SetAllBorders(border *Border) {
	cb.Top = border
	cb.Right = border
	cb.Bottom = border
	cb.Left = border
}