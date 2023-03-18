package model

import (
	"log"
	"strconv"
)

// Calc 计算全部
func (s *Stock) Calc() {
	s.CalcPB()
	s.CalcPE()
	s.CalcAAGR()
	s.CalcPEG()
	s.CalcROE()
}

// Discount 估值
func (s *Stock) Discount(r float64) {
	s.CalcDCE(r)
	s.CalcDCER()
	s.CalcDPE(r)
	s.CalcDPER()
}

// CalcPB 计算市净率
func (s *Stock) CalcPB() {

	if len(s.Enterprise) > 0 {
		// 每股净资产
		bps, err := strconv.ParseFloat(*s.Enterprise[0].Bps, 64)
		if err != nil || bps == 0 {
			log.Println(err, *s.Code, *s.Name, "bps")
			return
		}
		s.PB = *s.CurrentPrice / bps
	}

}

// CalcPE 计算市盈率
func (s *Stock) CalcPE() {
	if *s.CurrentPrice == 0 {
		log.Println(*s.Code, *s.Name, "*s.CurrentPrice")
		return
	}

	if len(s.Enterprise) > 0 {
		// 每股未分配利润
		mgwfplr, err := strconv.ParseFloat(*s.Enterprise[0].Mgwfplr, 64)
		if err != nil {
			log.Println(err, *s.Code, *s.Name, "mgfplr")
			return
		}
		s.PE = mgwfplr / *s.CurrentPrice
	}

}

// CalcAAGR 计算平均年增长率
func (s *Stock) CalcAAGR() {
	enterpriseList := s.Enterprise
	len := len(enterpriseList)
	var sum float64

	for k, v := range enterpriseList {
		n := k + 1
		if n >= len {
			break
		}
		lastBps, err := strconv.ParseFloat(*enterpriseList[n].Bps, 64)
		if err != nil || lastBps == 0 {
			log.Println(err, *s.Code, *s.Name, "lastBps")
			continue
		}
		Bps, err := strconv.ParseFloat(*v.Bps, 64)

		if err != nil {
			log.Println(err, *s.Code, *s.Name, "Bps")
			continue
		}

		curAAGR := (Bps - lastBps) / lastBps

		sum += curAAGR

	}
	s.AAGR = sum / float64((len - 1))

}

// CalcPEG 计算市盈增长比
func (s *Stock) CalcPEG() {
	s.PEG = s.PE / s.AAGR

}

// CalcROE 计算净资产收益率
func (s *Stock) CalcROE() {
	// 每股净值
	var mgwfplr, bps float64
	var err error
	if len(s.Enterprise) > 0 {
		mgwfplr, err = strconv.ParseFloat(*s.Enterprise[0].Mgwfplr, 64)
		if err != nil {
			log.Println(err, *s.Code, *s.Name, "mgwfplr")
			return
		}
		// 每股未分配利润
		bps, err = strconv.ParseFloat(*s.Enterprise[0].Bps, 64)
		if err != nil || bps == 0 {
			log.Println(err, *s.Code, *s.Name, "bps")
			return
		}
		s.ROE = mgwfplr / bps
	}

}

// CalcDPE 计算动态利润估值
func (s *Stock) CalcDPE(r float64) {
	if len(s.Enterprise) > 0 {

		bps, err := strconv.ParseFloat(*s.Enterprise[0].Bps, 64)
		if err != nil {
			log.Println(err, *s.Code, *s.Name, "bps")
			return
		}
		s.DPE = bps / (r - s.AAGR)
	}
}

// CalcDPER 估值 现值比
func (s *Stock) CalcDPER() {

	if *s.CurrentPrice == 0 {
		log.Println(*s.Code, *s.Name, "dper")
		return
	}

	s.DPER = s.DPE / *s.CurrentPrice
}

// CalcDCE 计算动态现金估值
func (s *Stock) CalcDCE(r float64) {
	if len(s.Enterprise) > 0 {
		// 每股经营现金流(元)
		mgjyxjje, err := strconv.ParseFloat(*s.Enterprise[0].Mgjyxjje, 64)
		if err != nil {
			log.Println(err, *s.Code, *s.Name, "mgjyxjje")
			return
		}
		s.DCE = mgjyxjje / (r - s.AAGR)
	}

}

// CalcDCER 估值 现值比
func (s *Stock) CalcDCER() {
	if *s.CurrentPrice == 0 {
		log.Println(*s.Code, *s.Name, "dcer")
		return
	}

	s.DCER = s.DCE / *s.CurrentPrice
}
