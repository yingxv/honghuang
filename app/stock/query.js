query {
  discount(
    discountRate: 0.1665
    createDate: "2020-04-25T00:24:05.842+08:00"
    weights: [
      { name: "PB", weight: 1, gt: false }
      { name: "PE", weight: 2, gt: false }
      { name: "PEG", weight: 3, gt: false }
      { name: "AAGR", weight: 4, gt: true }
      { name: "DCER", weight: 5, gt: true }
      { name: "ROE", weight: 6, gt: true }
      { name: "DPER", weight: 7, gt: true }
    ]
    limit: 20
  ) {
    grade
    classify
    name
    code
    PB
    PE
    PEG
    ROE
    DPE
    DPER
    DCE
    DCER
    AAGR
  }
}
