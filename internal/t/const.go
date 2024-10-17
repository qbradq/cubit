package t

// AtlasTextureDims are the X and Y dimensions of the face atlas texture in
// pixels.
const AtlasTextureDims = 2048

// FaceDims are the X and Y dimensions of a face in pixels.
const FaceDims = 16

// atlasDims are the X and Y dimensions of the face atlas in faces.
const AtlasDims = AtlasTextureDims / FaceDims

// Stepping value for face U and V offsets.
const PageStep float32 = float32(1) / float32(AtlasDims)

// VoxelScale is the scale of a single voxel.
const VoxelScale float32 = 1.0 / 16.0

// VirtualScreenWidth is the width of the virtual 2D screen in pixels.
const VirtualScreenWidth int = 320

// VirtualScreenHeight is the width of the virtual 2D screen in pixels.
const VirtualScreenHeight int = 180

// VirtualScreenGlyphSize is the dimensions of a glyph.
const VirtualScreenGlyphSize int = 4

// VirtualScreenGlyphsWide is the width of the screen in glyphs.
const VirtualScreenGlyphsWide int = VirtualScreenWidth / VirtualScreenGlyphSize

// VirtualScreenGlyphsHigh is the width of the screen in glyphs.
const VirtualScreenGlyphsHigh int = VirtualScreenHeight / VirtualScreenGlyphSize

// vsGlyphWidth is the width of one glyph in virtual screen units, minus boarder
// width.
const VSGlyphWidth int = VirtualScreenWidth / VirtualScreenGlyphsWide

// CellDimsVS is the dimensions of a cell in screen units.
const CellDimsVS int = VSGlyphWidth

// vsCellHeight is the height of one font atlas cell in virtual screen units.
const VSCellHeight int = (VSGlyphWidth * FACellHeight) / FACellWidth

// vsBaseline is the Y offset for baseline in virtual screen units.
const VSBaseline int = VSCellHeight / 4

// LineSpacingVS is the line spacing used in print commands in virtual screen
// units.
const LineSpacingVS int = (VSGlyphWidth / 2) * 3 // 1.5

// FADims is the font atlas dimensions.
const FADims int = 2048

// FAGlyphSize is the dimensions of a glyph in pixels.
const FAGlyphSize int = 32

// FACellWidth is the width of one cell in the atlas in pixels.
const FACellWidth int = 32

// FACellHeight is the height of one cell in the atlas in pixels.
const FACellHeight int = 64

// FACellXOfs is the X offset to use when rendering a glyph in pixels.
const FACellXOfs int = (FACellWidth - FAGlyphSize) / 2

// FACellYOfs is the Y offset to use when rendering a glyph in pixels.
const FACellYOfs int = (FACellHeight - FAGlyphSize) / 2

// FACellsWide is the width of the font atlas in glyphs.
const FACellsWide int = FADims / FACellWidth
