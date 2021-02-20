package escseq

import (
	c "image/color"
)

type RGBA c.RGBA

var (
	RGBABlack        = RGBA{0, 0, 0, 255}
	RGBARed          = RGBA{255, 0, 0, 255}
	RGBAGreen        = RGBA{0, 255, 0, 255}
	RGBAYellow       = RGBA{255, 255, 0, 255}
	RGBABlue         = RGBA{0, 0, 255, 255}
	RGBAMagenta      = RGBA{255, 0, 255, 255}
	RGBACyan         = RGBA{0, 255, 255, 255}
	RGBALightGray    = RGBA{211, 211, 211, 255}
	RGBADarkGray     = RGBA{169, 169, 169, 255}
	RGBALightRed     = RGBA{255, 144, 144, 255}
	RGBALightGreen   = RGBA{144, 238, 144, 255}
	RGBALightYellow  = RGBA{255, 255, 224, 255}
	RGBALightBlue    = RGBA{173, 216, 230, 255}
	RGBALightMagenta = RGBA{255, 224, 255, 255}
	RGBALightCyan    = RGBA{224, 255, 255, 255}
	RGBAWhite        = RGBA{255, 255, 255, 255}

	StringMap = map[string]RGBA{
		"black":   RGBABlack,
		"red":     RGBARed,
		"green":   RGBAGreen,
		"yellow":  RGBAYellow,
		"blue":    RGBABlue,
		"magenta": RGBAMagenta,
		"cyan":    RGBACyan,
		"white":   RGBAWhite,
	}

	// \x1b[NNm とかの NN に紐づくRGBA色
	// 例: \x1b[30m
	ANSIMap = map[int]RGBA{
		// 文字色
		30: RGBABlack,
		31: RGBARed,
		32: RGBAGreen,
		33: RGBAYellow,
		34: RGBABlue,
		35: RGBAMagenta,
		36: RGBACyan,
		37: RGBALightGray,
		90: RGBADarkGray,
		91: RGBALightRed,
		92: RGBALightGreen,
		93: RGBALightYellow,
		94: RGBALightBlue,
		95: RGBALightMagenta,
		96: RGBALightCyan,
		97: RGBAWhite,
		// 背景色
		40:  RGBABlack,
		41:  RGBARed,
		42:  RGBAGreen,
		43:  RGBAYellow,
		44:  RGBABlue,
		45:  RGBAMagenta,
		46:  RGBACyan,
		47:  RGBALightGray,
		100: RGBADarkGray,
		101: RGBALightRed,
		102: RGBALightGreen,
		103: RGBALightYellow,
		104: RGBALightBlue,
		105: RGBALightMagenta,
		106: RGBALightCyan,
		107: RGBAWhite,
	}

	// \x1b[38;5;NNNm とかの NNN に紐づくRGBA色
	// 例: \x1b[38;5;114m
	Map256 = map[int]RGBA{
		0:   {0, 0, 0, 255},
		1:   {128, 0, 0, 255},
		2:   {0, 128, 0, 255},
		3:   {128, 128, 0, 255},
		4:   {0, 0, 128, 255},
		5:   {128, 0, 128, 255},
		6:   {0, 128, 128, 255},
		7:   {192, 192, 192, 255},
		8:   {128, 128, 128, 255},
		9:   {255, 0, 0, 255},
		10:  {0, 255, 0, 255},
		11:  {255, 255, 0, 255},
		12:  {0, 0, 255, 255},
		13:  {255, 0, 255, 255},
		14:  {0, 255, 255, 255},
		15:  {255, 255, 255, 255},
		16:  {0, 0, 0, 255},
		17:  {0, 0, 95, 255},
		18:  {0, 0, 135, 255},
		19:  {0, 0, 175, 255},
		20:  {0, 0, 215, 255},
		21:  {0, 0, 255, 255},
		22:  {0, 95, 0, 255},
		23:  {0, 95, 95, 255},
		24:  {0, 95, 135, 255},
		25:  {0, 95, 175, 255},
		26:  {0, 95, 215, 255},
		27:  {0, 95, 255, 255},
		28:  {0, 135, 0, 255},
		29:  {0, 135, 95, 255},
		30:  {0, 135, 135, 255},
		31:  {0, 135, 175, 255},
		32:  {0, 135, 215, 255},
		33:  {0, 135, 255, 255},
		34:  {0, 175, 0, 255},
		35:  {0, 175, 95, 255},
		36:  {0, 175, 135, 255},
		37:  {0, 175, 175, 255},
		38:  {0, 175, 215, 255},
		39:  {0, 175, 255, 255},
		40:  {0, 215, 0, 255},
		41:  {0, 215, 95, 255},
		42:  {0, 215, 135, 255},
		43:  {0, 215, 175, 255},
		44:  {0, 215, 215, 255},
		45:  {0, 215, 255, 255},
		46:  {0, 255, 0, 255},
		47:  {0, 255, 95, 255},
		48:  {0, 255, 135, 255},
		49:  {0, 255, 175, 255},
		50:  {0, 255, 215, 255},
		51:  {0, 255, 255, 255},
		52:  {95, 0, 0, 255},
		53:  {95, 0, 95, 255},
		54:  {95, 0, 135, 255},
		55:  {95, 0, 175, 255},
		56:  {95, 0, 215, 255},
		57:  {95, 0, 255, 255},
		58:  {95, 95, 0, 255},
		59:  {95, 95, 95, 255},
		60:  {95, 95, 135, 255},
		61:  {95, 95, 175, 255},
		62:  {95, 95, 215, 255},
		63:  {95, 95, 255, 255},
		64:  {95, 135, 0, 255},
		65:  {95, 135, 95, 255},
		66:  {95, 135, 135, 255},
		67:  {95, 135, 175, 255},
		68:  {95, 135, 215, 255},
		69:  {95, 135, 255, 255},
		70:  {95, 175, 0, 255},
		71:  {95, 175, 95, 255},
		72:  {95, 175, 135, 255},
		73:  {95, 175, 175, 255},
		74:  {95, 175, 215, 255},
		75:  {95, 175, 255, 255},
		76:  {95, 215, 0, 255},
		77:  {95, 215, 95, 255},
		78:  {95, 215, 135, 255},
		79:  {95, 215, 175, 255},
		80:  {95, 215, 215, 255},
		81:  {95, 215, 255, 255},
		82:  {95, 255, 0, 255},
		83:  {95, 255, 95, 255},
		84:  {95, 255, 135, 255},
		85:  {95, 255, 175, 255},
		86:  {95, 255, 215, 255},
		87:  {95, 255, 255, 255},
		88:  {135, 0, 0, 255},
		89:  {135, 0, 95, 255},
		90:  {135, 0, 135, 255},
		91:  {135, 0, 175, 255},
		92:  {135, 0, 215, 255},
		93:  {135, 0, 255, 255},
		94:  {135, 95, 0, 255},
		95:  {135, 95, 95, 255},
		96:  {135, 95, 135, 255},
		97:  {135, 95, 175, 255},
		98:  {135, 95, 215, 255},
		99:  {135, 95, 255, 255},
		100: {135, 135, 0, 255},
		101: {135, 135, 95, 255},
		102: {135, 135, 135, 255},
		103: {135, 135, 175, 255},
		104: {135, 135, 215, 255},
		105: {135, 135, 255, 255},
		106: {135, 175, 0, 255},
		107: {135, 175, 95, 255},
		108: {135, 175, 135, 255},
		109: {135, 175, 175, 255},
		110: {135, 175, 215, 255},
		111: {135, 175, 255, 255},
		112: {135, 215, 0, 255},
		113: {135, 215, 95, 255},
		114: {135, 215, 135, 255},
		115: {135, 215, 175, 255},
		116: {135, 215, 215, 255},
		117: {135, 215, 255, 255},
		118: {135, 255, 0, 255},
		119: {135, 255, 95, 255},
		120: {135, 255, 135, 255},
		121: {135, 255, 175, 255},
		122: {135, 255, 215, 255},
		123: {135, 255, 255, 255},
		124: {175, 0, 0, 255},
		125: {175, 0, 95, 255},
		126: {175, 0, 135, 255},
		127: {175, 0, 175, 255},
		128: {175, 0, 215, 255},
		129: {175, 0, 255, 255},
		130: {175, 95, 0, 255},
		131: {175, 95, 95, 255},
		132: {175, 95, 135, 255},
		133: {175, 95, 175, 255},
		134: {175, 95, 215, 255},
		135: {175, 95, 255, 255},
		136: {175, 135, 0, 255},
		137: {175, 135, 95, 255},
		138: {175, 135, 135, 255},
		139: {175, 135, 175, 255},
		140: {175, 135, 215, 255},
		141: {175, 135, 255, 255},
		142: {175, 175, 0, 255},
		143: {175, 175, 95, 255},
		144: {175, 175, 135, 255},
		145: {175, 175, 175, 255},
		146: {175, 175, 215, 255},
		147: {175, 175, 255, 255},
		148: {175, 215, 0, 255},
		149: {175, 215, 95, 255},
		150: {175, 215, 135, 255},
		151: {175, 215, 175, 255},
		152: {175, 215, 215, 255},
		153: {175, 215, 255, 255},
		154: {175, 255, 0, 255},
		155: {175, 255, 95, 255},
		156: {175, 255, 135, 255},
		157: {175, 255, 175, 255},
		158: {175, 255, 215, 255},
		159: {175, 255, 255, 255},
		160: {215, 0, 0, 255},
		161: {215, 0, 95, 255},
		162: {215, 0, 135, 255},
		163: {215, 0, 175, 255},
		164: {215, 0, 215, 255},
		165: {215, 0, 255, 255},
		166: {215, 95, 0, 255},
		167: {215, 95, 95, 255},
		168: {215, 95, 135, 255},
		169: {215, 95, 175, 255},
		170: {215, 95, 215, 255},
		171: {215, 95, 255, 255},
		172: {215, 135, 0, 255},
		173: {215, 135, 95, 255},
		174: {215, 135, 135, 255},
		175: {215, 135, 175, 255},
		176: {215, 135, 215, 255},
		177: {215, 135, 255, 255},
		178: {215, 175, 0, 255},
		179: {215, 175, 95, 255},
		180: {215, 175, 135, 255},
		181: {215, 175, 175, 255},
		182: {215, 175, 215, 255},
		183: {215, 175, 255, 255},
		184: {215, 215, 0, 255},
		185: {215, 215, 95, 255},
		186: {215, 215, 135, 255},
		187: {215, 215, 175, 255},
		188: {215, 215, 215, 255},
		189: {215, 215, 255, 255},
		190: {215, 255, 0, 255},
		191: {215, 255, 95, 255},
		192: {215, 255, 135, 255},
		193: {215, 255, 175, 255},
		194: {215, 255, 215, 255},
		195: {215, 255, 255, 255},
		196: {255, 0, 0, 255},
		197: {255, 0, 95, 255},
		198: {255, 0, 135, 255},
		199: {255, 0, 175, 255},
		200: {255, 0, 215, 255},
		201: {255, 0, 255, 255},
		202: {255, 95, 0, 255},
		203: {255, 95, 95, 255},
		204: {255, 95, 135, 255},
		205: {255, 95, 175, 255},
		206: {255, 95, 215, 255},
		207: {255, 95, 255, 255},
		208: {255, 135, 0, 255},
		209: {255, 135, 95, 255},
		210: {255, 135, 135, 255},
		211: {255, 135, 175, 255},
		212: {255, 135, 215, 255},
		213: {255, 135, 255, 255},
		214: {255, 175, 0, 255},
		215: {255, 175, 95, 255},
		216: {255, 175, 135, 255},
		217: {255, 175, 175, 255},
		218: {255, 175, 215, 255},
		219: {255, 175, 255, 255},
		220: {255, 215, 0, 255},
		221: {255, 215, 95, 255},
		222: {255, 215, 135, 255},
		223: {255, 215, 175, 255},
		224: {255, 215, 215, 255},
		225: {255, 215, 255, 255},
		226: {255, 255, 0, 255},
		227: {255, 255, 95, 255},
		228: {255, 255, 135, 255},
		229: {255, 255, 175, 255},
		230: {255, 255, 215, 255},
		231: {255, 255, 255, 255},
		232: {8, 8, 8, 255},
		233: {18, 18, 18, 255},
		234: {28, 28, 28, 255},
		235: {38, 38, 38, 255},
		236: {48, 48, 48, 255},
		237: {58, 58, 58, 255},
		238: {68, 68, 68, 255},
		239: {78, 78, 78, 255},
		240: {88, 88, 88, 255},
		241: {98, 98, 98, 255},
		242: {108, 108, 108, 255},
		243: {118, 118, 118, 255},
		244: {128, 128, 128, 255},
		245: {138, 138, 138, 255},
		246: {148, 148, 148, 255},
		247: {158, 158, 158, 255},
		248: {168, 168, 168, 255},
		249: {178, 178, 178, 255},
		250: {188, 188, 188, 255},
		251: {198, 198, 198, 255},
		252: {208, 208, 208, 255},
		253: {218, 218, 218, 255},
		254: {228, 228, 228, 255},
		255: {238, 238, 238, 255},
	}
)
