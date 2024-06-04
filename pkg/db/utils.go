package db

import (
  "crypto/sha256"
  "encoding/csv"
  "fmt"
  "io"
  "log"
  "os"
  "path/filepath"
  "reflect"
)

// readCSV reads records from a csv file
func readCSV(csvPath string) ([][]string, error) {
  csvFile, err := os.Open(csvPath)
  if err != nil {
    return nil, fmt.Errorf("error opening csv file %s: %v", csvPath, err)
  }
  defer csvFile.Close()

  reader := csv.NewReader(csvFile)

  records, err := reader.ReadAll()
  if err != nil {
    return nil, fmt.Errorf("error reading csv file %s: %v", csvPath, err)
  }

  return records, nil
}

// sameImport compares if the new import is same as the previous one
func sameImport(newCSV string) bool {
  newHash, err := hashFile(newCSV)
  if err != nil {
    log.Printf("failed to create hash for the new csv %s: %v", newCSV, err)
    return false
  }

  oldHash, err := hashFile(csvImportBackup)
  if err != nil {
    log.Printf("failed to create hash for the old imported csv %s: %v", newCSV, err)
    return false
  }

  return newHash == oldHash
}

// hashFile returns hexadecimal hash for a file
func hashFile(filePath string) (string, error) {
  file, err := os.Open(filePath)
  if err != nil {
    return "", fmt.Errorf("failed to open %s for hashing: %v", filePath, err)
  }
  defer file.Close()

  hasher := sha256.New()
  if _, err := io.Copy(hasher, file); err != nil {
    return "", fmt.Errorf("failed to generate has for %s: %v", filePath, err)
  }

  return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// backupImport to backup for the newly imported csv
func backupImport(filePath, operation string) error {
  var backupPath string
  switch operation {
  case "import":
    backupPath = csvImportBackup
  case "export":
    backupPath = csvExportBackup
  }

  err := copyFile(filePath, backupPath)
  if err != nil {
    return fmt.Errorf("failed to create backup of %s: %v", filePath, err)
  }

  return nil
}

// copyFile copies file from
func copyFile(src, dst string) error {
  // open the source file
  sourceFile, err := os.Open(src)
  if err != nil {
    return err
  }
  defer sourceFile.Close()

  // create the destination file
  if _, err := os.Stat(dst); err != nil {
    log.Printf("backup file %s does not exist", dst)
    log.Printf("creating the backup %s", dst)
    if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
      return fmt.Errorf("error creating backup file at %s: %v", dst, err)
    }
  }
  destinationFile, err := os.Create(dst)
  if err != nil {
    return err
  }
  defer destinationFile.Close()

  // copy the content from source to destination
  _, err = io.Copy(destinationFile, sourceFile)
  if err != nil {
    return err
  }
  log.Printf("copy of file %s created at %s", src, dst)

  // ensure all content is flushed to the destination file
  err = destinationFile.Sync()
  if err != nil {
    return err
  }

  return nil
}
