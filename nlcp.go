package cryptoprofile

import (
	"fmt"
	"math"

	"github.com/xuri/excelize/v2"
)

type NLCPProfile struct {
	Complex    []int
	JumpHeight []int
	Feedback   []string
	JumpSeq    []int
	EV         []int
	EVJumpSeq  []int
}

type MeanVarience struct {
	ComplexMean       []float64
	ComplexVarience   []float64
	JumpSeqMean       []float64
	JumpSeqVarience   []float64
	EVMean            []float64
	EVVarience        []float64
	EVJumpSeqMean     []float64
	EVJumpSeqVarience []float64
}

func reverse(s string) string {
	rns := []rune(s) // convert to rune
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {

		// swap the letters of the string,
		// like first with last and so on.
		rns[i], rns[j] = rns[j], rns[i]
	}

	// return the reversed string.
	return string(rns)
}

func dist(ct string, h []string) int {
	count := 0
	for i := 0; i < len(h); i++ {
		le := len(h[i])
		if ct[0:le] == h[i] {
			count = count + 1
		}
	}
	return count % 2
}

func LCExpValues(n int) ([]float64, []float64) {
	avg := make([]float64, 0)
	varience := make([]float64, 0)
	for i := 1; i < n+1; i++ {
		p := float64(i) / 2
		avg = append(avg, p)
		varience = append(varience, 1)
	}
	return avg, varience
}

func EVCExpValues(n int) ([]float64, []float64) {
	avg := make([]float64, 0)
	varience := make([]float64, 0)
	for i := float64(1); i < float64(n+1); i++ {
		sum := float64(0)
		for j := float64(1); j < i; j++ {
			if i == 2 {
				sum = 0.5
				break
			} else {
				sum = sum + math.Pow(1-(1/math.Pow(2, i-j)), j)
			}
		}
		avg = append(avg, sum+1)
	}
	for i := float64(1); i < float64(n+1); i++ {
		sum := float64(0)
		for j := float64(1); j < i+1; j++ {
			sum = sum + math.Pow(j-avg[int(i-1)], 2)*(math.Pow(1-(1/math.Pow(2, i-j+1)), j-1)-math.Pow(1-(1/math.Pow(2, i-j)), j))
		}
		varience = append(varience, sum)
	}
	return avg, varience
}

func JLCExpValues(n int) ([]float64, []float64) {
	avg := make([]float64, 0)
	varience := make([]float64, 0)
	for i := float64(1); i < float64(n+1); i++ {
		p := math.Pow(2, i)
		var a, v float64
		if int(i)%2 == 0 {
			a = (i / 4) + (1.0 / 3) - 1/(3*p)
			v = (i / 8) - (2.0 / 9) + (i / (6 * p)) + (1 / (3 * p)) - (1 / (9 * p))
		} else {
			a = (i / 4) + (5.0 / 12) - 1/(3*math.Pow(2, i))
			v = ((i - 1) / 8) + (i / (6 * p)) + (7 / (18 * p)) - (1 / (9 * p))
		}
		avg = append(avg, a)
		varience = append(varience, v)
	}
	return avg, varience
}

func NCR(n int64, r int64) int64 {
	if r > n-r {
		r = n - r
	}
	ans := int64(1)

	for i := int64(1); i <= r; i++ {
		ans *= n - r + i
		ans /= i
	}
	return ans
}

func Probability(j float64, i float64) float64 {
	val := math.Pow(1-(1/math.Pow(2, j)), float64(NCR(int64(i-j+1), int64(2))))
	return val
}

/*
	sum1=0
	averg.append(0.5)
	varnce.append(0.5)
	for i in range(2,N+1):
		sum1=0
		for j in range(1,i+1):
			sum1 += j * ( pow((1-(1/pow(2,j))),ncr(i-j+1,2))- pow((1-(1/pow(2,j-1))),ncr(i-j+2,2)) )
		averg.append(sum1)
	#print(averg)

	for i in range(2,N):
		sum2=0
		varnce.append(0)
		for j in range(1,i+1):
				sum2 +=  pow(j-(2*math.log(i,2)),2) * (probability(j,i))
				#print(sum2)
		varnce.append(sum2)
*/

func NLCExpValues(n int) ([]float64, []float64) {
	avg := make([]float64, 0)
	varience := make([]float64, 0)
	avg = append(avg, 0.5)
	varience = append(varience, 0.5)
	for i := float64(2); i < float64(n+1); i++ {
		sum := 0.0
		for j := float64(1); j < i+1; j++ {
			sum = sum + j*(math.Pow((1-(1/math.Pow(2, j))), float64(NCR(int64(i-j+1), 2)))-math.Pow((1-(1/math.Pow(2, j-1))), float64(NCR(int64(i-j+2), 2))))
		}
		avg = append(avg, sum)
	}
	for i := float64(2); i < float64(n); i++ {
		sum := 0.0
		//varience = append(varience, 0.0)
		for j := float64(1); j < i+1; j++ {
			sum = sum + math.Pow(j-(2*math.Log2(i)), 2)*(Probability(j, i))
		}
		varience = append(varience, sum)
	}
	return avg, varience
}
func (bs BitStream) CreateNLCProfile(seqLength int, numSamples int) []NLCPProfile {
	if seqLength*numSamples > bs.Length {
		return nil
	}
	ncps := make([]NLCPProfile, 0)
	for i := 0; i < numSamples; i++ {
		sample := bs.Value[i*seqLength : i*seqLength+seqLength]
		k := 0
		m := 0
		mm := make([]int, 0)
		js := make([]int, 0)
		jh := make([]int, 0)
		h := make([]string, 0)
		mm = append(mm, 0)
		js = append(js, 0)
		jh = append(jh, 0)
		y0 := 0
		if sample[0:1] == "1" {
			y0 = 1
		}

		for n := 1; n < seqLength; n++ {
			cjs := js[len(js)-1]
			ct := reverse(sample[n-m : n])
			yn := 0
			if sample[n:n+1] == "1" {
				yn = 1
			}
			var d int
			if len(ct) == 0 {
				d = yn - y0
			} else {
				d = yn - dist(ct, h)
			}
			if d == -1 {
				d = 1
			}
			if d != 0 {
				if m == 0 {
					k = n
					m = n
				} else if k <= 0 {
					nbs := BitStream{
						Value:  sample[0:n],
						Length: len(sample[0:n]),
					}
					t := nbs.EigenVaule()
					if int(t) < (n + 1 - m) {
						k = n + 1 - int(t) - m
						m = n + 1 - int(t)
					}
				} else {
					k = k - 1
				}
				f := reverse(sample[n-m : n])
				h = append(h, f)
			} else {
				k = k - 1
			}
			if len(mm) != 0 {
				if mm[len(mm)-1] != m {
					js = append(js, cjs+1)
					jh = append(jh, m-mm[len(mm)-1])
				} else {
					js = append(js, cjs)
					jh = append(jh, 0)
				}
			} else {
				js = append(js, cjs)
				jh = append(jh, 0)
			}
			mm = append(mm, m)
		}

		var evv []int
		var evs []int
		ebs := BitStream{
			Value:  sample,
			Length: seqLength,
		}
		evv, evs = ebs.EigenProfileWitJump()

		prof := NLCPProfile{
			Complex:    mm,
			JumpHeight: jh,
			JumpSeq:    js,
			Feedback:   h,
			EV:         evv,
			EVJumpSeq:  evs,
		}
		ncps = append(ncps, prof)
	}
	return ncps
}

func CalculateMeanVarience(nlcp []NLCPProfile) MeanVarience {
	complexMean := make([]float64, len(nlcp[0].Complex))
	for i := 0; i < len(nlcp[0].Complex); i++ {
		complexMean[i] = 0.0
		for j := 0; j < len(nlcp); j++ {
			complexMean[i] = complexMean[i] + float64(nlcp[j].Complex[i])
		}
		complexMean[i] = complexMean[i] / float64(len(nlcp))
	}
	complexVarience := make([]float64, len(nlcp[0].Complex))
	for i := 0; i < len(nlcp[0].Complex); i++ {
		ss := 0.0
		for j := 0; j < len(nlcp); j++ {
			ss = ss + math.Pow((float64(nlcp[j].Complex[i])-complexMean[i]), 2.0)
		}
		complexVarience[i] = ss / float64((len(nlcp) - 1))
	}

	jumpSeqMean := make([]float64, len(nlcp[0].JumpSeq))
	for i := 0; i < len(nlcp[0].JumpSeq); i++ {
		jumpSeqMean[i] = 0.0
		for j := 0; j < len(nlcp); j++ {
			jumpSeqMean[i] = jumpSeqMean[i] + float64(nlcp[j].JumpSeq[i])
		}
		jumpSeqMean[i] = jumpSeqMean[i] / float64(len(nlcp))
	}
	jumSeqVarience := make([]float64, len(nlcp[0].JumpSeq))
	for i := 0; i < len(nlcp[0].JumpSeq); i++ {
		ss := 0.0
		for j := 0; j < len(nlcp); j++ {
			ss = ss + math.Pow((float64(nlcp[j].JumpSeq[i])-jumpSeqMean[i]), 2.0)
		}
		jumSeqVarience[i] = ss / float64((len(nlcp) - 1))
	}

	evMean := make([]float64, len(nlcp[0].EV))
	for i := 0; i < len(nlcp[0].EV); i++ {
		evMean[i] = 0.0
		for j := 0; j < len(nlcp); j++ {
			evMean[i] = evMean[i] + float64(nlcp[j].EV[i])
		}
		evMean[i] = evMean[i] / float64(len(nlcp))
	}
	evVarience := make([]float64, len(nlcp[0].EV))
	for i := 0; i < len(nlcp[0].EV); i++ {
		ss := 0.0
		for j := 0; j < len(nlcp); j++ {
			ss = ss + math.Pow((float64(nlcp[j].EV[i])-evMean[i]), 2.0)
		}
		evVarience[i] = ss / float64((len(nlcp) - 1))
	}

	evJumSeqMean := make([]float64, len(nlcp[0].EVJumpSeq))
	for i := 0; i < len(nlcp[0].EVJumpSeq); i++ {
		evJumSeqMean[i] = 0.0
		for j := 0; j < len(nlcp); j++ {
			evJumSeqMean[i] = evJumSeqMean[i] + float64(nlcp[j].EVJumpSeq[i])
		}
		evJumSeqMean[i] = evJumSeqMean[i] / float64(len(nlcp))
	}
	evJumpSeqVarience := make([]float64, len(nlcp[0].EVJumpSeq))
	for i := 0; i < len(nlcp[0].EVJumpSeq); i++ {
		ss := 0.0
		for j := 0; j < len(nlcp); j++ {
			ss = ss + math.Pow((float64(nlcp[j].EVJumpSeq[i])-evJumSeqMean[i]), 2.0)
		}
		evJumpSeqVarience[i] = ss / float64((len(nlcp) - 1))
	}

	mv := MeanVarience{
		ComplexMean:       complexMean,
		ComplexVarience:   complexVarience,
		JumpSeqMean:       jumpSeqMean,
		JumpSeqVarience:   jumSeqVarience,
		EVMean:            evMean,
		EVVarience:        evVarience,
		EVJumpSeqMean:     evJumSeqMean,
		EVJumpSeqVarience: evJumpSeqVarience,
	}
	return mv
}

func ExpectedMeanVarience(seqLength int) MeanVarience {
	nlm, nlv := NLCExpValues(seqLength)
	jm, jv := JLCExpValues(seqLength)
	evm, evv := EVCExpValues(seqLength)
	mv := MeanVarience{
		ComplexMean:     nlm,
		ComplexVarience: nlv,
		JumpSeqMean:     jm,
		JumpSeqVarience: jv,
		EVMean:          evm,
		EVVarience:      evv,
	}
	return mv
}

const (
	nlp_sheet   string = "Non-Linear-Profile"
	jnlp_sheet  string = "Jumpin-Non-Linear-Profile"
	evp_sheet   string = "Eigenvalue-Profile"
	evjsp_sheet string = "Jumpin-Eigenvalue-Profile"
)

func LetterInc(str string) string {
	var new_str string
	var sub_str string
	length := len(str)
	if str == "" {
		return "A"
	}

	if length == 1 {
		sub_str = ""
	} else {
		sub_str = str[:length-1]
	}

	if str[length-1] == 'Z' {
		new_str = LetterInc(sub_str) + "A"
	} else {
		byte_str := byte(str[length-1])
		byte_str++
		new_str = sub_str + string(byte_str)
	}
	return new_str
}

func SaveNLCPProfile(fileName string, nlcp []NLCPProfile, mv MeanVarience, emv MeanVarience) error {
	f := excelize.NewFile()
	f.SaveAs(fileName)

	_, err := f.Rows("Sheet1")
	if err != nil {
		f.NewSheet(nlp_sheet)
	} else {
		f.SetSheetName("Sheet1", nlp_sheet)
	}

	style, err := f.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"wrap_text":true},"font":{"bold":true}}`)
	if err != nil {
		return err
	}
	endRow := fmt.Sprintf("A%d", len(nlcp)+4)
	err = f.SetCellStyle(nlp_sheet, "A1", endRow, style)
	if err != nil {
		return err
	}
	f.SetColWidth(nlp_sheet, "A", "A", 20)
	row := 1
	for i := 0; i < len(nlcp); i++ {
		rowStr := fmt.Sprintf("A%d", row)
		row++
		f.SetCellValue(nlp_sheet, rowStr, "Sample-"+fmt.Sprintf("%d", i+1))
		colStr := "B"
		for j := 0; j < len(nlcp[i].Complex); j++ {
			f.SetCellValue(nlp_sheet, colStr+fmt.Sprintf("%d", i+1), nlcp[i].Complex[j])
			colStr = LetterInc(colStr)
		}
	}
	f.SetCellValue(nlp_sheet, "A"+fmt.Sprintf("%d", row), "Observed Mean")
	colStr := "B"
	for j := 0; j < len(mv.ComplexMean); j++ {
		f.SetCellValue(nlp_sheet, colStr+fmt.Sprintf("%d", row), mv.ComplexMean[j])
		colStr = LetterInc(colStr)
	}
	row++
	f.SetCellValue(nlp_sheet, "A"+fmt.Sprintf("%d", row), "Expected Mean")
	colStr = "B"
	for j := 0; j < len(emv.ComplexMean); j++ {
		f.SetCellValue(nlp_sheet, colStr+fmt.Sprintf("%d", row), emv.ComplexMean[j])
		colStr = LetterInc(colStr)
	}
	row++

	f.SetCellValue(nlp_sheet, "A"+fmt.Sprintf("%d", row), "Observed Varience")
	colStr = "B"
	for j := 0; j < len(mv.ComplexVarience); j++ {
		f.SetCellValue(nlp_sheet, colStr+fmt.Sprintf("%d", row), mv.ComplexVarience[j])
		colStr = LetterInc(colStr)
	}
	row++
	f.SetCellValue(nlp_sheet, "A"+fmt.Sprintf("%d", row), "Expected Varience")
	colStr = "B"
	for j := 0; j < len(emv.ComplexVarience); j++ {
		f.SetCellValue(nlp_sheet, colStr+fmt.Sprintf("%d", row), emv.ComplexVarience[j])
		colStr = LetterInc(colStr)
	}
	row++

	_, err = f.Rows("Sheet1")
	if err != nil {
		f.NewSheet(jnlp_sheet)
	} else {
		f.SetSheetName("Sheet1", jnlp_sheet)
	}

	style, err = f.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"wrap_text":true},"font":{"bold":true}}`)
	if err != nil {
		return err
	}
	endRow = fmt.Sprintf("A%d", len(nlcp)+4)
	err = f.SetCellStyle(jnlp_sheet, "A1", endRow, style)
	if err != nil {
		return err
	}
	f.SetColWidth(jnlp_sheet, "A", "A", 20)
	row = 1
	for i := 0; i < len(nlcp); i++ {
		rowStr := fmt.Sprintf("A%d", row)
		row++
		f.SetCellValue(jnlp_sheet, rowStr, "Sample-"+fmt.Sprintf("%d", i+1))
		colStr := "B"
		for j := 0; j < len(nlcp[i].JumpSeq); j++ {
			f.SetCellValue(jnlp_sheet, colStr+fmt.Sprintf("%d", i+1), nlcp[i].JumpSeq[j])
			colStr = LetterInc(colStr)
		}
	}
	f.SetCellValue(jnlp_sheet, "A"+fmt.Sprintf("%d", row), "Observed Mean")
	colStr = "B"
	for j := 0; j < len(mv.JumpSeqMean); j++ {
		f.SetCellValue(jnlp_sheet, colStr+fmt.Sprintf("%d", row), mv.JumpSeqMean[j])
		colStr = LetterInc(colStr)
	}
	row++
	f.SetCellValue(jnlp_sheet, "A"+fmt.Sprintf("%d", row), "Expected Mean")
	colStr = "B"
	for j := 0; j < len(emv.JumpSeqMean); j++ {
		f.SetCellValue(jnlp_sheet, colStr+fmt.Sprintf("%d", row), emv.JumpSeqMean[j])
		colStr = LetterInc(colStr)
	}
	row++

	f.SetCellValue(jnlp_sheet, "A"+fmt.Sprintf("%d", row), "Observed Varience")
	colStr = "B"
	for j := 0; j < len(mv.JumpSeqVarience); j++ {
		f.SetCellValue(jnlp_sheet, colStr+fmt.Sprintf("%d", row), mv.JumpSeqVarience[j])
		colStr = LetterInc(colStr)
	}
	row++
	f.SetCellValue(jnlp_sheet, "A"+fmt.Sprintf("%d", row), "Expected Varience")
	colStr = "B"
	for j := 0; j < len(emv.JumpSeqVarience); j++ {
		f.SetCellValue(jnlp_sheet, colStr+fmt.Sprintf("%d", row), emv.JumpSeqVarience[j])
		colStr = LetterInc(colStr)
	}
	row++

	_, err = f.Rows("Sheet1")
	if err != nil {
		f.NewSheet(evp_sheet)
	} else {
		f.SetSheetName("Sheet1", evp_sheet)
	}

	style, err = f.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"wrap_text":true},"font":{"bold":true}}`)
	if err != nil {
		return err
	}
	endRow = fmt.Sprintf("A%d", len(nlcp)+4)
	err = f.SetCellStyle(evp_sheet, "A1", endRow, style)
	if err != nil {
		return err
	}
	f.SetColWidth(evp_sheet, "A", "A", 20)
	row = 1
	for i := 0; i < len(nlcp); i++ {
		rowStr := fmt.Sprintf("A%d", row)
		row++
		f.SetCellValue(evp_sheet, rowStr, "Sample-"+fmt.Sprintf("%d", i+1))
		colStr := "B"
		for j := 0; j < len(nlcp[i].EV); j++ {
			f.SetCellValue(evp_sheet, colStr+fmt.Sprintf("%d", i+1), nlcp[i].EV[j])
			colStr = LetterInc(colStr)
		}
	}
	f.SetCellValue(evp_sheet, "A"+fmt.Sprintf("%d", row), "Observed Mean")
	colStr = "B"
	for j := 0; j < len(mv.EVMean); j++ {
		f.SetCellValue(evp_sheet, colStr+fmt.Sprintf("%d", row), mv.EVMean[j])
		colStr = LetterInc(colStr)
	}
	row++
	f.SetCellValue(evp_sheet, "A"+fmt.Sprintf("%d", row), "Expected Mean")
	colStr = "B"
	for j := 0; j < len(emv.EVMean); j++ {
		f.SetCellValue(evp_sheet, colStr+fmt.Sprintf("%d", row), emv.EVMean[j])
		colStr = LetterInc(colStr)
	}
	row++

	f.SetCellValue(evp_sheet, "A"+fmt.Sprintf("%d", row), "Observed Varience")
	colStr = "B"
	for j := 0; j < len(mv.EVVarience); j++ {
		f.SetCellValue(evp_sheet, colStr+fmt.Sprintf("%d", row), mv.EVVarience[j])
		colStr = LetterInc(colStr)
	}
	row++
	f.SetCellValue(evp_sheet, "A"+fmt.Sprintf("%d", row), "Expected Varience")
	colStr = "B"
	for j := 0; j < len(emv.EVVarience); j++ {
		f.SetCellValue(evp_sheet, colStr+fmt.Sprintf("%d", row), emv.EVVarience[j])
		colStr = LetterInc(colStr)
	}
	row++

	_, err = f.Rows("Sheet1")
	if err != nil {
		f.NewSheet(evjsp_sheet)
	} else {
		f.SetSheetName("Sheet1", evjsp_sheet)
	}

	style, err = f.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"wrap_text":true},"font":{"bold":true}}`)
	if err != nil {
		return err
	}
	endRow = fmt.Sprintf("A%d", len(nlcp)+4)
	err = f.SetCellStyle(evjsp_sheet, "A1", endRow, style)
	if err != nil {
		return err
	}
	f.SetColWidth(evjsp_sheet, "A", "A", 20)
	row = 1
	for i := 0; i < len(nlcp); i++ {
		rowStr := fmt.Sprintf("A%d", row)
		row++
		f.SetCellValue(evjsp_sheet, rowStr, "Sample-"+fmt.Sprintf("%d", i+1))
		colStr := "B"
		for j := 0; j < len(nlcp[i].EVJumpSeq); j++ {
			f.SetCellValue(evjsp_sheet, colStr+fmt.Sprintf("%d", i+1), nlcp[i].EVJumpSeq[j])
			colStr = LetterInc(colStr)
		}
	}
	f.SetCellValue(evjsp_sheet, "A"+fmt.Sprintf("%d", row), "Observed Mean")
	colStr = "B"
	for j := 0; j < len(mv.EVJumpSeqMean); j++ {
		f.SetCellValue(evjsp_sheet, colStr+fmt.Sprintf("%d", row), mv.EVJumpSeqMean[j])
		colStr = LetterInc(colStr)
	}
	row++
	f.SetCellValue(evjsp_sheet, "A"+fmt.Sprintf("%d", row), "Expected Mean")
	colStr = "B"
	for j := 0; j < len(emv.EVJumpSeqMean); j++ {
		f.SetCellValue(evjsp_sheet, colStr+fmt.Sprintf("%d", row), emv.EVJumpSeqMean[j])
		colStr = LetterInc(colStr)
	}
	row++

	f.SetCellValue(evjsp_sheet, "A"+fmt.Sprintf("%d", row), "Observed Varience")
	colStr = "B"
	for j := 0; j < len(mv.EVJumpSeqVarience); j++ {
		f.SetCellValue(evjsp_sheet, colStr+fmt.Sprintf("%d", row), mv.EVJumpSeqVarience[j])
		colStr = LetterInc(colStr)
	}
	row++
	f.SetCellValue(evjsp_sheet, "A"+fmt.Sprintf("%d", row), "Expected Varience")
	colStr = "B"
	for j := 0; j < len(emv.EVJumpSeqVarience); j++ {
		f.SetCellValue(evjsp_sheet, colStr+fmt.Sprintf("%d", row), emv.EVJumpSeqVarience[j])
		colStr = LetterInc(colStr)
	}
	row++

	f.Save()
	return nil
}
