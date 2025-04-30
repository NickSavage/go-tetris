package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	Square int = 1
	Line   int = 2
)

type Piece struct {
	Type        int
	X           int
	Y           int
	Orientation int
}

type BoardSquare struct {
	X        int
	Y        int
	Occupied bool
}

type Board struct {
	Squares     [][]BoardSquare
	Pieces      []*Piece
	ActivePiece *Piece
	IsBust      bool
}

const (
	WIN_HEIGHT int32 = 800
	WIN_WIDTH  int32 = 500
)

const (
	MAX_DELAY int = 5
)

const (
	SQUARE_HEIGHT int32 = 25
	SQUARE_WIDTH  int32 = 25
	SQUARE_ROWS   int32 = 20
	SQUARE_COLS   int32 = 10
	SQUARE_BORDER int32 = 2
)

func createSquareRect(x, y int) sdl.Rect {

	x32 := int32(x)
	y32 := -1 * (int32(y) - SQUARE_ROWS)
	rect := sdl.Rect{
		X: x32*SQUARE_HEIGHT + x32*SQUARE_BORDER,
		Y: y32*SQUARE_HEIGHT + y32*SQUARE_BORDER,
		W: SQUARE_WIDTH,
		H: SQUARE_HEIGHT,
	}
	return rect
}

func (b *Board) rotatePiece(piece *Piece) {
	if piece.Type == Line {
		if piece.Orientation == 1 {
			piece.Orientation = 0
		} else {
			piece.Orientation = 1
		}
	}
}

func (b *Board) newActivePiece() {

	piece := Piece{
		X:    0,
		Y:    int(SQUARE_ROWS),
		Type: Line,
	}

	b.Pieces = append(b.Pieces, &piece)
	b.ActivePiece = &piece

}

func (b *Board) occupySpaces(piece *Piece) {
	if piece.Type == Line {
		if piece.Orientation == 1 {
			b.Squares[piece.Y][piece.X].Occupied = true
			b.Squares[piece.Y-1][piece.X].Occupied = true
			b.Squares[piece.Y-2][piece.X].Occupied = true
			b.Squares[piece.Y-3][piece.X].Occupied = true
		}
	}
}

func (b *Board) checkCollision(piece *Piece) bool {
	if b.Squares[piece.Y-3][piece.X].Occupied {
		return true
	}
	return false

}

func (b *Board) movePieceLeft(piece *Piece) {
	piece.X -= 1
}
func (b *Board) movePieceRight(piece *Piece) {
	piece.X += 1
}
func (b *Board) movePieceDown(piece *Piece) {
	piece.Y -= 1

	if b.checkCollision(piece) {
		b.occupySpaces(piece)
		b.newActivePiece()

	}
}

func drawLinePiece(surface *sdl.Surface, piece *Piece) {

	white_colour := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	white_pixel := sdl.MapRGBA(surface.Format, white_colour.R, white_colour.G, white_colour.B, white_colour.A)

	if piece.Orientation == 0 {

		rect := createSquareRect(piece.X, piece.Y)
		surface.FillRect(&rect, white_pixel)

		rect = createSquareRect(piece.X, piece.Y-1)
		surface.FillRect(&rect, white_pixel)

		rect = createSquareRect(piece.X, piece.Y-2)
		surface.FillRect(&rect, white_pixel)

		rect = createSquareRect(piece.X, piece.Y-3)
		surface.FillRect(&rect, white_pixel)
	} else {

		rect := createSquareRect(piece.X, piece.Y)
		surface.FillRect(&rect, white_pixel)

		rect = createSquareRect(piece.X+1, piece.Y)
		surface.FillRect(&rect, white_pixel)

		rect = createSquareRect(piece.X+2, piece.Y)
		surface.FillRect(&rect, white_pixel)

		rect = createSquareRect(piece.X+3, piece.Y)
		surface.FillRect(&rect, white_pixel)
	}

}

func (b *Board) Draw(surface *sdl.Surface, window *sdl.Window) {

	//	rect := sdl.Rect{0, 0, SQUARE_HEIGHT, SQUARE_WIDTH}
	colour := sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple
	pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(nil, 0)
	for y := len(b.Squares) - 1; y >= 0; y-- {
		col := b.Squares[y]
		for x := range col {
			x32 := int32(x)
			y32 := int32(y)
			rect := sdl.Rect{
				X: x32*SQUARE_HEIGHT + x32*SQUARE_BORDER,
				Y: y32*SQUARE_HEIGHT + y32*SQUARE_BORDER,
				W: SQUARE_WIDTH,
				H: SQUARE_HEIGHT,
			}
			surface.FillRect(&rect, pixel)
		}
	}
	for _, piece := range b.Pieces {

		if piece.Type == Line {
			drawLinePiece(surface, piece)
		}
	}

	window.UpdateSurface()
}

func main() {
	log.Printf("hello world")

	squares := make([][]BoardSquare, SQUARE_ROWS)
	for i := range squares {
		squares[i] = make([]BoardSquare, SQUARE_COLS)
	}
	board := &Board{
		Squares: squares,
		Pieces:  make([]*Piece, 0),
	}

	piece := Piece{
		X:    0,
		Y:    int(SQUARE_ROWS),
		Type: Line,
	}

	board.Pieces = append(board.Pieces, &piece)
	board.ActivePiece = &piece
	b := board
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WIN_WIDTH, WIN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	delay := MAX_DELAY
	running := true
	for running {
		b.Draw(surface, window)

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.KeyboardEvent:
				t := event.(*sdl.KeyboardEvent)

				if t.Type == sdl.KEYDOWN && t.Keysym.Sym == sdl.K_UP {
					b.rotatePiece(board.ActivePiece)
				} else if t.Type == sdl.KEYDOWN && t.Keysym.Sym == sdl.K_LEFT {
					b.movePieceLeft(board.ActivePiece)
				} else if t.Type == sdl.KEYDOWN && t.Keysym.Sym == sdl.K_RIGHT {
					b.movePieceRight(board.ActivePiece)
				} else if t.Type == sdl.KEYDOWN && t.Keysym.Sym == sdl.K_DOWN {
					log.Printf("down")
					b.movePieceDown(board.ActivePiece)
				}
			case *sdl.QuitEvent: // NOTE: Please use `*sdl.QuitEvent` for `v0.4.x` (current version).
				println("Quit")
				running = false
			}
		}

		delay -= 1

		if delay == 0 {
			log.Printf("piece %v", board.ActivePiece)
			b.movePieceDown(board.ActivePiece)
			delay = MAX_DELAY
		}
		sdl.Delay(33)

		if board.IsBust {
			running = false
		}

	}
}
