// Package asciiart_blueprint takes a blueprint schema and draws ASCII art for
// it. The characters chosen are based on the entity type / prototype, using
// a hash mapped to A-Z/a-z/0-9.
//
// The public interface is unstable.
//
// BUG: The code currently has a major misunderstanding: just because
// coordinates can be at position 0.5 does not mean we will in practice observe
// such an offset between items. The granularity really seems to be at 1.0 for
// blueprints. This means any use of *2 and any advice to do it is probably
// misguided. The code is a decent start, but should not be used or seen as
// correct. We can keep it ONLY as an internal piece of code for purposes of
// tilemap (to use integers), but with sparse tilemap maybe even that is not
// necessary.
package asciiart_blueprint // badc0de.net/pkg/factorioblueprint/asciiart_blueprint

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"

	"badc0de.net/pkg/factorioblueprint/schema/blueprint_schema"
)

// Reader is a struct that holds the blueprint schema and reads ASCII art for
// it. Precision more than 0.5 in x and y coordinates is not supported here.
type Reader struct {
	// Note: updating TileWidth and TileHeight is invalid after the initial
	// creation, without invalidating other cache, and these fields will become
	// private in the future.

	Blueprint  *blueprint_schema.Blueprint
	TileWidth  int // TileWidth is the size of the tiles in characters.
	TileHeight int // TileHeight is the size of the tiles in characters.

	displayRune               map[string]rune   // displayRune maps entity type / tile prototype to a character.
	legend                    map[rune][]string // legend maps a character to a list of entity types / tile prototypes.
	cachedWidth, cachedHeight float64           // cachedWidth and cachedHeight are the size of the blueprint in native factorio size.
	cachedSparseTilemap       *SparseTilemap    // cachedSparseTilemap is the sparse tilemap for the blueprint.

	// buffer holds the generated ASCII art data
	buffer []byte
	pos    int
}

// Implement the io.Reader interface
//
// TOOD: instead of using a buffer and a strings.Builder, just write directly
// into p []byte, and *only* record whatever does not fit. We won't allow
// seeking back anyway. If we wanted to build the string in one go, we could
// just directly use a strings.Reader in Read.
func (r *Reader) Read(p []byte) (n int, err error) {
	if r.buffer == nil {
		// Generate the ASCII art and store it in r.buffer
		asciiArt, err := r.generateASCIIArt()
		if err != nil {
			return 0, err
		}
		legend, err := r.generateLegend()
		if err != nil {
			return 0, err
		}
		r.buffer = []byte(asciiArt + "---\n" + legend)
	}

	if r.pos >= len(r.buffer) {
		return 0, io.EOF
	}

	n = copy(p, r.buffer[r.pos:])
	r.pos += n
	return n, nil
}

// generateASCIIArt generates the ASCII art representation of the blueprint
func (r *Reader) generateASCIIArt() (string, error) {
	var sb strings.Builder

	// r.Size() returns the width and height, but we need to multiply by two
	// to get integer size.
	widthF, heightF := r.Size()
	width := int(widthF) * 2
	height := int(heightF) * 2

	// Loop over y and x to build the ASCII art.
	// Note that we have to start with the smallest x and smallest y to get the
	// correct order.
	minX, minY, _, _ := r.box()
	startX := int(minX*2) * r.TileWidth
	startY := int(minY*2) * r.TileHeight
	//for y := startY + int(height)*r.TileHeight - 1; y >= startY; y-- {
	for y := startY; y < startY+int(height)*r.TileHeight; y++ { // every line
		for x := startX; x < startX+int(width)*r.TileWidth; {
			str, err := r.StringAtScreenPosition(x, y)
			if err != nil {
				return "", err
			}
			sb.WriteString(str)
			x += len(str)
			if len(str) == 0 {
				break
			} // avoid infinite loop
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// generateLegend generates the legend for the ASCII art representation of the
// blueprint.
func (r *Reader) generateLegend() (string, error) {
	r.computeLegend()

	// Add [x]: entity-name or tile-name for each entry we hold.
	//
	// These would be ideally sorted by the name of the entity or tile, but
	// we'll accept sorting by the rune (i.e. full line) -- as long as it's
	// stable for tests etc. So we won't use a strings.Builder.

	var entries []string

	for rn, names := range r.legend {
		entries = append(entries, fmt.Sprintf("[%c]: %s", rn, strings.Join(names, ", ")))
	}
	sort.Strings(entries)
	return strings.Join(entries, "\n"), nil
}

// Size computes width an height of the coordinates. Factorio holds the
// coordinates at 0.5 precision.
//
// Caller needs to multiply by two to get the integer size, and then also use
// TileWidth and TileHeight to get the actual size in characters.
//
// The blueprint is allowed to contain negative coordinates.
func (r *Reader) Size() (float64, float64) {

	if r.cachedWidth > 0 && r.cachedHeight > 0 {
		return r.cachedWidth, r.cachedHeight
	}

	// Compute the size of the blueprint.
	minX, minY, maxX, maxY := r.box()

	// Compute the size of the blueprint in factorio-native coordinates.
	r.cachedWidth = maxX - minX
	r.cachedHeight = maxY - minY
	return r.cachedWidth, r.cachedHeight
}

// box computes minxy and maxxy of the blueprint. Factorio holds the coordinates
// at 0.5 precision. To get the box, we add 1.0 to the max x and y, since then
// we don't return 0x0 for a blueprint of size of a single tile.
func (r *Reader) box() (float64, float64, float64, float64) {
	minX := math.MaxFloat64
	minY := math.MaxFloat64
	maxX := math.MaxFloat64 * -1
	maxY := math.MaxFloat64 * -1

	for _, entity := range r.Blueprint.Entities {
		if entity.Position.X < minX {
			minX = entity.Position.X
		}
		if entity.Position.Y < minY {
			minY = entity.Position.Y
		}
		if entity.Position.X > maxX {
			maxX = entity.Position.X
		}
		if entity.Position.Y > maxY {
			maxY = entity.Position.Y
		}
	}
	for _, tile := range r.Blueprint.Tiles {
		if tile.Position.X < minX {
			minX = tile.Position.X
		}
		if tile.Position.Y < minY {
			minY = tile.Position.Y
		}
		if tile.Position.X > maxX {
			maxX = tile.Position.X
		}
		if tile.Position.Y > maxY {
			maxY = tile.Position.Y
		}
	}

	return minX, minY, maxX + 1, maxY + 1
}

// computeLegend computes the legend and the displayRune cache for the ASCII art.
// It maps a character to a list of entity types / tile prototypes, and a
// entity type / tile prototype to a character.
//
// This function does NOT care about positioning, it just computes all
// characters you might ever see based on whether they are in the blueprint.
//
// TODO: It might be better to compute the legend later, based only on what
// appears on the drawn tiles, ignoring entities and tiles that are not drawn
// in ASCII art.
func (r *Reader) computeLegend() {
	if len(r.displayRune) > 0 && len(r.legend) > 0 {
		return
	}

	r.displayRune = make(map[string]rune)
	r.legend = make(map[rune][]string)

	for _, entity := range r.Blueprint.Entities {
		// n.b. it might be better to use the icon name here.
		nameForDisplay := entity.Name

		if _, ok := r.displayRune[nameForDisplay]; !ok {
			rn := stringToRune(nameForDisplay)
			r.displayRune[nameForDisplay] = rn
			r.legend[rn] = append(r.legend[rn], nameForDisplay)
		}
	}

	for _, tile := range r.Blueprint.Tiles {
		// n.b. it might be better to use the icon name here.
		nameForDisplay := tile.Name // prototype name of the tile

		if _, ok := r.displayRune[nameForDisplay]; !ok {
			rn := stringToRune(nameForDisplay)
			r.displayRune[nameForDisplay] = rn
			r.legend[rn] = append(r.legend[rn], nameForDisplay)
		}
	}
}

const (
	// validAlphabet is the alphabet used to map a hash to a character.
	// The characters are A-Z, a-z, 0-9.
	validAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// stringToRune converts a string to a rune. If the string is empty, it returns
// a space. The character is computed by cheaply hashing the string and mapping
// it to a character in the range A-Z/a-z/0-9.
func stringToRune(s string) rune {
	if s == "" {
		return ' '
	}
	hash := 0
	for _, c := range s {
		hash += int(c)
	}
	return rune(validAlphabet[hash%len(validAlphabet)])
}

// RuneForEntityOrTile returns a rune for the entity or tile struct.
// If there is no entity or tile, it returns a space.
func (r *Reader) RuneForEntityOrTile(et *entityOrTile) rune {
	if et.entity != nil {
		return r.EntityChar(et.entity)
	}
	if et.tile != nil {
		return r.TileChar(et.tile)
	}
	return ' '
}

// EntityChar returns a character for the entity type.
//
// TODO: we may want to invert responsibility and compute the hash from the
// name here, and then reuse this in computeLegend.
func (r *Reader) EntityChar(entity *blueprint_schema.Entity) rune {
	r.computeLegend() // precompute hashes if needed
	nameForDisplay := entity.Name
	return r.displayRune[nameForDisplay]
}

// TileChar returns a character for the tile prototype.
//
// TODO: we may want to invert responsibility and compute the hash from the
// name here, and then reuse this in computeLegend.
func (r *Reader) TileChar(tile *blueprint_schema.Tile) rune {
	r.computeLegend() // precompute hashes if needed
	nameForDisplay := tile.Name
	return r.displayRune[nameForDisplay]
}

// SparseTilemap is a sparse tilemap for the blueprint. Each tile in the tilemap
// is a struct with slices of entities and tiles at that position. Since these
// are integers, the position in blueprint is multiplied by two to get the
// actual position we will use in the tilemap.
type SparseTilemap struct {
	st map[int]map[int]*sparseTilemapTile
}

// newSparseTilemap returns a new empty SparseTilemap.
func newSparseTilemap() *SparseTilemap {
	return &SparseTilemap{st: make(map[int]map[int]*sparseTilemapTile)}
}

type sparseTilemapTile struct {
	Entities []*blueprint_schema.Entity
	Tiles    []*blueprint_schema.Tile
}

// top returns the top entity or tile at this sparse tilemap tile.
func (s *sparseTilemapTile) top() *entityOrTile {
	if s == nil {
		return nil
	}
	if len(s.Entities) > 0 {
		return &entityOrTile{entity: s.Entities[0]}
	}
	if len(s.Tiles) > 0 {
		return &entityOrTile{tile: s.Tiles[0]}
	}
	return nil
}

// At returns the entities and tiles at a position in the sparse tilemap. If
// there are no entities or tiles at the position, it returns nil. The
// coordinates are the same ones as in the blueprint, not multiplied by two.
// They will be multiplied accordingly.
func (s *SparseTilemap) At(x, y float64) *sparseTilemapTile {
	x *= 2
	y *= 2
	return s.AtInt(int(x), int(y))
}

// AtInt returns the entities and tiles at a position in the sparse tilemap. If
// there are no entities or tiles at the position, it returns nil. The
// coordinates are the same ones as in the blueprint but multiplied by two.
func (s *SparseTilemap) AtInt(x, y int) *sparseTilemapTile {
	if s == nil {
		return nil
	}
	if s.st == nil {
		return nil
	}
	if _, ok := s.st[x]; !ok {
		return nil
	}
	if _, ok := s.st[x][y]; !ok {
		return nil
	}
	return s.st[x][y]
}

// EntityOrTileAtPosition returns the top entity or tile at a position. If there
// is no entity or tile at the position, it returns nil.
//
// The x and y coordinates are the blueprint coordinates multiplied by two.
func (s *SparseTilemap) EntityOrTileAtIntPosition(x, y int) *entityOrTile {
	if s == nil {
		return nil
	}
	et := s.AtInt(x, y)
	return et.top()
}

// EntityOrTileAtPosition returns the top entity or tile at a position. If there
// is no entity or tile at the position, it returns nil.
//
// The x and y coordinates are the blueprint coordinates.
func (s *SparseTilemap) EntityOrTileAtPosition(x, y float64) *entityOrTile {
	x *= 2
	y *= 2
	return s.EntityOrTileAtIntPosition(int(x), int(y))
}

type entityOrTile struct {
	entity *blueprint_schema.Entity
	tile   *blueprint_schema.Tile
}

// AsSparseTilemap returns a sparse tilemap for the blueprint. Each tile in the
// tilemap is a struct with slices of entities and tiles at that position. Since
// these are integers, the position in blueprint is multiplied by two to get the
// actual position we will use in the tilemap.
func (r *Reader) AsSparseTilemap() *SparseTilemap {
	if r.cachedSparseTilemap != nil {
		return r.cachedSparseTilemap
	}

	tilemap := newSparseTilemap()

	for _, entity := range r.Blueprint.Entities {
		x := int(entity.Position.X * 2)
		y := int(entity.Position.Y * 2)
		if _, ok := tilemap.st[x]; !ok {
			tilemap.st[x] = make(map[int]*sparseTilemapTile)
		}
		if _, ok := tilemap.st[x][y]; !ok {
			tilemap.st[x][y] = &sparseTilemapTile{}
		}
		// fmt.Printf("adding entity %d=%q to %d, %d\n", entity.EntityNumber, entity.Name, x, y)
		// entity is a local var and would be overwritten, so copy it
		entity := entity
		tilemap.st[x][y].Entities = append(tilemap.st[x][y].Entities, &entity)
	}

	for _, tile := range r.Blueprint.Tiles {
		x := int(tile.Position.X * 2)
		y := int(tile.Position.Y * 2)
		if _, ok := tilemap.st[x]; !ok {
			tilemap.st[x] = make(map[int]*sparseTilemapTile)
		}
		if _, ok := tilemap.st[x][y]; !ok {
			tilemap.st[x][y] = &sparseTilemapTile{}
		}
		fmt.Printf("adding tile %q to %d, %d\n", tile.Name, x, y)
		// tile is a local var and would be overwritten, so copy it
		tile := tile
		tilemap.st[x][y].Tiles = append(tilemap.st[x][y].Tiles, &tile)
	}

	r.cachedSparseTilemap = tilemap
	return tilemap
}

// StringAtScreenPosition returns a string with one or more runes at a position
// in the returned string. It returns at minimum the rune representing the top
// entity or tile, preferring entities over tiles.
//
// If TileHeight is 2, the bottom line is used to indicate the direction with
// characters ^V<> for up, down, left, right.
//
// If TileHeight is at least 3, the top line is used to indicate the up
// direction, bottom to indicate the down direction, and the TileWidth must be
// at least 3. If TileWidth is at least 3, then the character left of the entity
// will indicate the direction of the entity is left, or if the direciton is
// right, the character right of the entity will indicate the direction of the
// entity. The string is padded to the width of the tile. Only heights 1, 2 and
// 3 are supported. x and y are guaranteed to be between 0 and Size()[0] and
// Size()[1] respectively.
//
// Example with TileWidth 1:
// * TileHeight 1: if we would draw an entity which is at 4.5, 6.5, and the
//   entity points up, and we are correctly invoked with x=9, y=13, then the
//   string returned would be "A".
// * TileHeight 2: if we would draw an entity which is at 4.5, 6.5, and the
//   entity points up, and we are correctly invoked with x=9, y=13*2=26, then
//   the string returned would be "A". If we are then invoked for the line below
//   (x=9, y=13*2+1=27), then the string returned would be "^". (If the entity
//   had no direction indicated, then the string would be "A" and " ").
// * TileHeight 3: this is an unsupported configuration and we return an error.
//
// Example with TileWidth 5:
// * TileHeight 1: if we would draw an entity which is at 4.5, 6.5, and we are
//   correctly invoked with x=9*5 for the leftmost char, y=13 for the only
//   line, then the string returned would be " ^A  ".
// * TileHeight 2: if we would draw an entity which is at 4.5, 6.5, and the
//   entity points up, and we are correctly invoked with x=9*5 for the leftmost
//   char, y=13*2=26 for the top line, then the string returned would be
//   "  A  ".
//   If we are then invoked for the line below (x=9*5, y=13*2+1=27), then the
//   string returned would be "  ^  ".

//   - TileHeight 3: if we would draw an entity which is at 4.5, 6.5, and the
//     entity points up, and we are correctly invoked with x=9*5 for the leftmost
//     char, y=13*3+1=40 for the mid line, then the string returned would be
//     "  A  ".
//     If we are then invoked for the line below (x=9*5, y=13*3+2=41),
//     then the string returned would be "     " (no direction indicator since we
//     point up).
//     If we are then invoked for the line above (x=9*5, y=13*3+0=39), then the
//     string returned would be "  ^  ".
//   - TileHeight 3 (left): if we would draw an entity which is at 4.5, 6.5, and
//     the entity points left, and we are correctly invoked with x=9*5=45 for the
//     leftmost char, y=13*3+1=40 for the mid line, then the string returned would
//     be " <A  ".
//     If we are then invoked for the line below (x=9*5, y=13*3+2=41),
//     then the string returned would be "     " (no direction indicator since we
//     point left).
//     If we are then invoked for the line above (x=9*5, y=13*3+0=39),
//     then the string returned would be "     " (no direction indicator since we
//     point left).
//   - TileHeight 3 (right): if we would draw an entity which is at 4.5, 6.5, and
//     the entity points right, and we are correctly invoked with x=9*5=45 for the
//     leftmost char, y=13*3+1=40 for the mid line, then the string returned would
//     be "  A> ".
//     If we are then invoked for the line below (x=9*5, y=13*3+2=41),
//     then the string returned would be "     " (no direction indicator since we
//     point right).
//     If we are then invoked for the line above (x=9*5, y=13*3+0=39),
//     then the string returned would be "     " (no direction indicator since we
//     point right).
//   - TileHeight 3 (up): if we would draw an entity which is at 4.5, 6.5, and
//     the entity points up, and we are correctly invoked with x=9*5=45 for the
//     leftmost char, y=13*3+1=40 for the mid line, then the string returned would
//     be "  A  ".
//     If we are then invoked for the line below (x=9*5, y=13*3+2=41),
//     then the string returned would be "     " (no direction indicator since we
//     point up).
//     If we are then invoked for the line above (x=9*5, y=13*3+0=39),
//     then the string returned would be "  ^  ".
//   - TileHeight 3 (down): if we would draw an entity which is at 4.5, 6.5, and
//     the entity points down, and we are correctly invoked with x=9*5=45 for the
//     leftmost char, y=13*3+1=40 for the mid line, then the string returned would
//     be "  V  ".
//     If we are then invoked for the line below (x=9*5, y=13*3+2=41),
//     then the string returned would be "     " (no direction indicator since we
//     point down).
//     If we are then invoked for the line above (x=9*5, y=13*3+0=39),
//     then the string returned would be "     " (no direction indicator since we
//     point down).
//   - TileHeight 4: this is an unsupported configuration and we return an error.
//
// The y coordinate will be taken verbatim. The x coordinate will be taken
// divided by TileWidth ignoring any modulo to get the leftmost position in the
// blueprint.
func (r *Reader) StringAtScreenPosition(x, y int) (string, error) {
	switch r.TileHeight {
	case 1:
		return r.stringAtScreenPositionHeight1(x, y, true, true, true)
	case 2:
		return r.stringAtScreenPositionHeight2(x, y)
	case 3:
		return r.stringAtScreenPositionHeight3(x, y)
	default:
		return "", fmt.Errorf("unsupported TileHeight %d", r.TileHeight)
	}
}

// stringAtScreenPositionHeight1 deals with case of TileHeight 1.
//
// Accepted coordinates are factorioX*2*TileWidth, factorioY*2. Note: this
// means the caller is expected to have corrected for TileHeight != 1, usually
// by dividing.
//
// If there is enough space, direction is preferred to be on the left.
// If adaptiveLeftRight, the direction is preferred to be whatever the
// side the direction is pointing at.
func (r *Reader) stringAtScreenPositionHeight1(x, y int, includeUpDown bool, adaptiveLeftRight bool, includePosition bool) (string, error) {
	r.computeLegend()
	x = x / r.TileWidth
	tilemap := r.AsSparseTilemap()
	et := tilemap.EntityOrTileAtIntPosition(x, y)

	if et == nil {
		return " ", nil
	}

	etRune := r.RuneForEntityOrTile(et)

	// if et.entity != nil {
	// 	fmt.Printf("found entity %d=%q at %d, %d\n", et.entity.EntityNumber, et.entity.Name, x, y)
	// }

	// First, the easy case: TileWidth == 1.
	if r.TileWidth == 1 {
		return string(etRune), nil
	}

	// If the width is 2, we can afford to have the direction indicator.
	directionRune := ' '
	if includePosition {
		directionRune = entityToDirectionRune(et.entity)
		if !includeUpDown && (directionRune == '^' || directionRune == 'V') {
			directionRune = ' '
		}
	}
	if r.TileWidth == 2 {
		return r.stringAtScreenPositionHeight1CenterBit(includeUpDown, adaptiveLeftRight, etRune, directionRune)
	}

	// Otherwise we must pad the string. We can still reuse the centerbit
	// function, which always has length of 2.
	centerBit, err := r.stringAtScreenPositionHeight1CenterBit(includeUpDown, adaptiveLeftRight, etRune, directionRune)
	if err != nil {
		return "", err
	}
	out, err := paddedString(centerBit, r.TileWidth)
	if err != nil {
		return "", fmt.Errorf("failed to pad string %q to width %d: %v", centerBit, r.TileWidth, err)
	}
	if len(out) != r.TileWidth {
		return "", fmt.Errorf("unexpected length %d for %q, want %d", len(out), out, r.TileWidth)
	}
	return out, nil
}

// paddedString returns a string with the entity or tile rune, possibly with
// the direction rune, padded to the TileWidth. (For example, it can have both
// the direction and the tile -- e.g. "^A", " B" -- or just the tile -- e.g.
// "A" -- or just the direction -- e.g. "<", which is used for the bottom or
// top line when they need centering of the up or down arrows).
//
// We do want to support well the case of single line with direction indicator.
//
// Based on TileWidth, we must center the string to have a total width
// of TileWidth, with the etRune in the center with equal number of
// characters to the left and right (leaning to be on the right). We accept
// TileWidth in the form of wantedSz, and we accept rn as a string rather than
// a rune to allow the caller to indicate whether they want one character
// centered or both the entity and the direction rune centered.
//
// Total width of the string, called wantedSz:
// must match padLeft + 1 + (directionRune > 0 ? 1 : 0) + padRight
func paddedString(rn string, wantedSz int) (string, error) {
	if wantedSz == len(rn) {
		return rn, nil
	}
	var paddedLeftSz, paddedRightSz int

	switch len(rn) {
	case 1:
		paddedLeftSz = (wantedSz - 1) / 2
		paddedRightSz = wantedSz - 1 - paddedLeftSz
	case 2:
		paddedLeftSz = (wantedSz - 2) / 2
		paddedRightSz = wantedSz - 2 - paddedLeftSz
	default:
		return "", fmt.Errorf("unsupported length %d for %q", len(rn), rn)
	}
	paddedLeft := strings.Repeat(" ", paddedLeftSz)
	paddedRight := strings.Repeat(" ", paddedRightSz)
	return paddedLeft + rn + paddedRight, nil
}

// stringAtScreenPositionHeight1CenterBit deals with the case of TileWidth 2.
//
// The centerbit string is always 2 characters long.
func (r *Reader) stringAtScreenPositionHeight1CenterBit(includeUpDown, adaptiveLeftRight bool, etRune rune, directionRune rune) (string, error) {
	if includeUpDown {
		switch directionRune {
		case '^':
			return "^" + string(etRune), nil
		case 'V':
			if adaptiveLeftRight {
				return string(etRune) + "V", nil
			} else {
				return "V" + string(etRune), nil
			}
		case '<':
			return "<" + string(etRune), nil
		case '>':
			if adaptiveLeftRight {
				return string(etRune) + ">", nil
			} else {
				return ">" + string(etRune), nil
			}
		default:
			return " " + string(etRune), nil
		}
	} else {
		switch directionRune {
		case '<':
			return "<" + string(etRune), nil
		case '>':
			if adaptiveLeftRight {
				return string(etRune) + ">", nil
			}
		default:
			return " " + string(etRune), nil
		}
	}
	return "", fmt.Errorf("unsupported directionRune %c", directionRune)
}

// stringAtScreenPositionHeight2 deals with case of TileHeight 2.
//
// Accepted coordinates are factorioX*2*TileWidth, factorioY*2*TileHeight.
func (r *Reader) stringAtScreenPositionHeight2(x, y int) (string, error) {
	r.computeLegend()
	if y%2 == 0 {
		// We are on the top line.
		y = y / r.TileHeight
		return r.stringAtScreenPositionHeight1(x, y, false, false, false)
	}
	// Next case is the bottom line which draws only the direction, centered.
	// All directions are covered.
	return r.centeredDirectionLine(x, y/r.TileHeight, "^V<>")
}

// centeredDirectionLine draws the direction indicator centered on the line.
//
// It is almost identical to drawing the first line of the height-2 style, but
// we must add the direction indicator, no matter what. Padding is the same.
//
// y is expected to be already divided by TileHeight.
//
// A good default for permittedRunes is "^V<>".
func (r *Reader) centeredDirectionLine(x, y int, permittedRunes string) (string, error) {
	x = x / r.TileWidth
	tilemap := r.AsSparseTilemap()
	et := tilemap.EntityOrTileAtIntPosition(x, y)
	if et == nil {
		return " ", nil
	}

	directionRune := entityToDirectionRune(et.entity)
	if !strings.ContainsRune(permittedRunes, directionRune) {
		// substitute with space
		directionRune = ' '
	}

	if r.TileWidth == 1 {
		return string(directionRune), nil
	}

	return paddedString(string(directionRune), r.TileWidth)
}

// stringAtScreenPositionHeight3 deals with case of TileHeight 3.
//
// Accepted coordinates are factorioX*2*TileWidth, factorioY*2*TileHeight.
//
// The difference from height 2 is that the top line gets the indicator for
// up, the bottom line gets the indicator for down, and the middle line gets
// the entity and the direction indicator. The direction indicator is centered
// on top and bottom line, and on the middle line it is on the left or right
// depending on the direction of the entity.
func (r *Reader) stringAtScreenPositionHeight3(x, y int) (string, error) {
	// Minimum requirement: TileWidth >= 3.
	if r.TileWidth < 3 {
		return "", fmt.Errorf("unsupported TileWidth %d for TileHeight 3", r.TileWidth)
	}

	// Determine: is it top or bottom line?
	switch y % 3 {
	case 0:
		// Top line.
		//
		// As long as the direction rune is ^, treat like height 2's bottom
		// line.
		y /= r.TileHeight
		return r.centeredDirectionLine(x, y, "^")
	case 1:
		// Middle line.
		//
		// Reuse 'single line' code, but only support left and right.
		y /= r.TileHeight
		return r.stringAtScreenPositionHeight1(x, y, false, true, true)
	case 2:
		// Bottom line.
		//
		// As long as the direction rune is V, treat like height 2's top line.
		y /= r.TileHeight
		return r.centeredDirectionLine(x, y, "V")
	default:
		return "", fmt.Errorf("TODO: not fully implemented")
	}
}

// entityToDirection returns a string with a direction indicator for the entity.
// The direction indicator is a character from the set ^V<> for up, down, left,
// right. If the entity has no direction, it returns a space.
func entityToDirectionRune(entity *blueprint_schema.Entity) rune {
	if entity.Direction == nil {
		return ' '
	}
	switch *entity.Direction {
	case 0:
		return '^'
	case 2:
		return '>'
	case 4:
		return 'V'
	case 6:
		return '<'
	default:
		return ' '
	}
}

// NewReader returns a new Reader for the blueprint schema. The tileWidth and
// tileHeight are the size of the tiles, expressed in characters.
//
// Members TileWidth and TileHeight are not allowed to be updated after the
// initial creation.
//
// Use the io.Reader interface to read the ASCII art. Use io.Copy to copy the
// ASCII art to an io.Writer (such as os.Stdout).
func NewReader(blueprint *blueprint_schema.Blueprint, tileWidth, tileHeight int) *Reader {
	return &Reader{
		Blueprint:  blueprint,
		TileWidth:  tileWidth,
		TileHeight: tileHeight,
	}
}
