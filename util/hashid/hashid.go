package hashid

import (
	"bytes"
	"fmt"
)

const (
	DefaultAlphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	DefaultSeps     string = "cfhistuCFHISTU"
)

type Hasher struct {
	alphabet  []byte
	seps      []byte
	sepMap    map[byte]bool
	salt      []byte
	guards    []byte
	guardMap  map[byte]bool
	minLength int
	//maxLengthPerNumber int
}

func NewDefaultHasher(salt []byte) (*Hasher, error) {
	return NewHasher([]byte(DefaultAlphabet), []byte(DefaultSeps), salt)
}

func NewHasher(alphabet, seps, salt []byte) (*Hasher, error) {
	// seps should not in alphabet
	sepMap := map[byte]bool{}
	for _, sepLetter := range seps {
		sepMap[sepLetter] = true
		//TODO: sep unique
	}
	hashAlphabet := make([]byte, 0, len(alphabet))
	//TODO : alphabet unique check
	for _, letter := range alphabet {
		if _, exists := sepMap[letter]; !exists {
			hashAlphabet = append(hashAlphabet, letter)
		}
	}
	consistentShuffle(hashAlphabet, salt)
	consistentShuffle(seps, salt)

	guardCount := 1
	guards := hashAlphabet[:guardCount]
	hashAlphabet = hashAlphabet[guardCount:]

	guardMap := map[byte]bool{}
	for _, guardLetter := range guards {
		guardMap[guardLetter] = true
	}

	return &Hasher{
		alphabet:  hashAlphabet,
		seps:      seps,
		salt:      salt,
		sepMap:    sepMap,
		guards:    guards,
		guardMap:  guardMap,
		minLength: 10,
	}, nil
}

func (h *Hasher) Encode(numbers []int) (string, error) {
	var numbersHash int
	var result bytes.Buffer
	for i, n := range numbers {
		//TODO: Why add 100 ?
		numbersHash += n % (i + 100)
	}
	alphabetLen := len(h.alphabet)

	iterrateBuffer := make([]byte, alphabetLen+len(h.salt)+1)

	copy(iterrateBuffer[1:len(h.salt)+1], h.salt)

	alphabet := iterrateBuffer[len(h.salt)+1:]
	copy(alphabet, h.alphabet)

	lottery := alphabet[numbersHash%alphabetLen]

	iterrateBuffer[0] = lottery
	result.WriteByte(lottery)

	for i, n := range numbers {
		consistentShuffle(alphabet, iterrateBuffer[:alphabetLen])
		nHash, err := hash(n, alphabet)
		if err != nil {
			return "", err
		}
		result.Write(nHash)
		// Add seps
		n %= int(nHash[0]) + i
		result.WriteByte(h.seps[n%len(h.seps)])

	}
	resultB := result.Bytes()
	// remove last sep
	resultB = resultB[:len(resultB)-1]
	//TODO : minLength guard
	lenWantage := h.minLength - len(resultB)
	if lenWantage > 0 {
		guardIndex := (numbersHash + int(resultB[0])) % len(h.guards)
		resultB = append([]byte{h.guards[guardIndex]}, resultB...)
		// add guard to the head
		if lenWantage > 1 {
			//add guard to the tail
			guardIndex := (numbersHash + int(resultB[2])) % len(h.guards)
			resultB = append(resultB, h.guards[guardIndex])
		}
		if lenWantage > 2 {
			// add last guards
			// excess := alphabetLen - (lenWantage % alphabetLen)
			// ((lenWantage/alphabetLen)+1)*alphabetLen/2 - excess/2
			lenWantage -= 2
			iterationCount := lenWantage/alphabetLen + 1

			prefixLen := (lenWantage + lenWantage%alphabetLen) / 2

			suffixLen := lenWantage - prefixLen

			resultBuf := make([]byte, h.minLength)

			prefixBuf := resultBuf[:prefixLen]
			suffixBuf := resultBuf[h.minLength-suffixLen:]
			copy(resultBuf[prefixLen:], resultB)

			halfLength := len(alphabet) / 2
			for i := 0; i < iterationCount; i++ {
				consistentShuffle(alphabet, alphabet)
				pValveIndex := prefixLen - (i+1)*halfLength
				alphabetOffset := 0
				if pValveIndex < 0 {
					alphabetOffset = -pValveIndex
					pValveIndex = 0
				}

				copy(prefixBuf[pValveIndex:], alphabet[alphabetOffset:halfLength])

				copy(suffixBuf[i*halfLength:], alphabet[halfLength:])

			}

			resultB = resultBuf
		}
	}
	return string(resultB), nil
}

func (h *Hasher) Decode(raw string) ([]int, error) {
	//FIXME:  strip guards
	result := make([]int, 0, 10)
	lottery := raw[0]

	hashSlices := make([][]byte, 0, 10)
	// can not skip the first char ,maybe it a guard
	Pos := 0
	numStart := 1

	hasGuard := false

	for Pos < len(raw) {
		//FIXME：各种溢出情况
		runeByte := byte(raw[Pos])
		if _, exists := h.guardMap[runeByte]; exists {
			if !hasGuard {
				hasGuard = true
				Pos++
				// ToDO: 这里可能会溢出
				lottery = raw[Pos]
				Pos++
				numStart = Pos
				// reset hashSlices
				hashSlices = hashSlices[:0]
				continue
			} else {
				// the second guard skip follow chars
				break
			}
		}
		if _, exists := h.sepMap[runeByte]; exists {
			hashSlices = append(hashSlices, []byte(raw[numStart:Pos]))
			numStart = Pos + 1
		}
		Pos++
	}
	if numStart < Pos {
		hashSlices = append(hashSlices, []byte(raw[numStart:Pos]))
	}

	iterrateBuffer := make([]byte, len(h.alphabet)+len(h.salt)+1)
	iterrateBuffer[0] = lottery
	copy(iterrateBuffer[1:len(h.salt)+1], h.salt)
	alphabet := iterrateBuffer[len(h.salt)+1:]
	copy(alphabet, h.alphabet)

	for _, nHash := range hashSlices {
		consistentShuffle(alphabet, iterrateBuffer[:len(alphabet)])
		number, err := unhash(nHash, alphabet)
		if err != nil {
			return nil, err
		}
		result = append(result, number)
	}

	return result, nil
}

// Do a Sanity check after decode for insurance
func (h *Hasher) SanityDecode(raw string) ([]int, error) {
	result, err := h.Decode(raw)
	// if no errer occur do a sanity check
	if err == nil {
		sanityCheck, _ := h.Encode(result)
		if sanityCheck != raw {
			return result, fmt.Errorf("Decode result(%v) reEncode hash(%s) is not raw(%s)", result, sanityCheck, raw)
		}
	}
	return result, err
}

// https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
func consistentShuffle(alphabet, salt []byte) {
	saltLen := len(salt)
	if saltLen == 0 {
		return
	}
	//take a snapShoot or salt may change while Shuffle
	saltSnapshoot := make([]byte, saltLen)
	copy(saltSnapshoot, salt)

	alphabetLen := len(alphabet)
	for i, v, p := alphabetLen-1, 0, 0; i > 0; i-- {
		p += int(saltSnapshoot[v])
		j := (int(saltSnapshoot[v]) + v + p) % i
		alphabet[i], alphabet[j] = alphabet[j], alphabet[i]
		v = (v + 1) % saltLen
	}
	return
}

func hash(number int, alphabet []byte) ([]byte, error) {
	if number < 0 {
		return nil, fmt.Errorf("can not hash negative value %d", number)
	}

	alphabetLen := len(alphabet)
	var buf bytes.Buffer
	for {
		// loop excute at least one time
		buf.WriteByte(alphabet[number%alphabetLen])
		number /= alphabetLen
		if number <= 0 {
			break
		}
	}
	return buf.Bytes(), nil
}

func unhash(input, alphabet []byte) (int, error) {
	var result int
	alphabetLen := len(alphabet)

	shiftStep := 1
	for _, inputRune := range input {
		posInAlphabet := -1
		for pos, alphabetRune := range alphabet {
			if inputRune == alphabetRune {
				posInAlphabet = pos
				break
			}
		}
		if posInAlphabet < 0 {
			return 0, fmt.Errorf("can not found rune [%c] in alphabet", rune(inputRune))
		}
		result = result + posInAlphabet*shiftStep
		shiftStep *= alphabetLen
	}
	return result, nil

}
