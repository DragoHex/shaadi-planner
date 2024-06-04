package invitee

import (
  "database/sql"
  "fmt"
  "strconv"
)

type Invitee struct {
  Name         string
  NameOnCard   string
  NumOfMembers string
  Address      string
  PhoneNum     string
  Events       string
  Tags         string
  Gifts        string
  Notes        string
}

func (i *Invitee) Scan(rows *sql.Rows) ([][]string, error) {
  var records [][]string
  for rows.Next() {
    var id int
    err := rows.Scan(&id, &i.Name, &i.NameOnCard, &i.NumOfMembers,
      &i.Address, &i.PhoneNum, &i.Events, &i.Tags, &i.Gifts, &i.Notes)
    if err != nil {
      return [][]string{}, fmt.Errorf("error scanning the fields: %v", err)
    }

    record := make([]string, 0, 10)
    record = append(record, strconv.Itoa(id))
    record = append(record, i.Name)
    record = append(record, i.NameOnCard)
    record = append(record, i.NumOfMembers)
    record = append(record, i.Address)
    record = append(record, i.PhoneNum)
    record = append(record, i.Events)
    record = append(record, i.Tags)
    record = append(record, i.Gifts)
    record = append(record, i.Notes)

    records = append(records, record)
  }
  return records, nil
}
