package sqlite

import (
	"strconv"

	pragmas "github.com/Nevoral/sqlofi/internal/sqlite/Pragmas"
)

func newPragma(schemaName, name string) *Pragma {
	return &Pragma{
		Pragma: pragmas.NewPragma(schemaName, name),
	}
}

type Pragma struct {
	*pragmas.Pragma
}

func (p *Pragma) FuncType(value string) *Pragma {
	p.Pragma.FuncType(value)
	return p
}

func (p *Pragma) ValueType(value string) *Pragma {
	p.Pragma.ValueType(value)
	return p
}

// Analysis limit PRAGMA functions
func AnalysisLimit() *Pragma {
	return newPragma("", "analysis_limit")
}

// Application ID PRAGMA functions
func ApplicationID(schemaName string) *Pragma {
	return newPragma(schemaName, "application_id")
}

// Auto vacuum PRAGMA functions
func AutoVacuum(schemaName string) *Pragma {
	return newPragma(schemaName, "auto_vacuum")
}

func AutoVacuumNone(schemaName string) *Pragma {
	return newPragma(schemaName, "auto_vacuum").ValueType("NONE")
}

func AutoVacuumFull(schemaName string) *Pragma {
	return newPragma(schemaName, "auto_vacuum").ValueType("FULL")
}

func AutoVacuumIncremental(schemaName string) *Pragma {
	return newPragma(schemaName, "auto_vacuum").ValueType("INCREMENTAL")
}

// Automatic index PRAGMA functions
func AutomaticIndex() *Pragma {
	return newPragma("", "automatic_index")
}

// Busy timeout PRAGMA functions
func BusyTimeout() *Pragma {
	return newPragma("", "busy_timeout")
}

// Cache size PRAGMA functions
func CacheSize(schemaName string) *Pragma {
	return newPragma(schemaName, "cache_size")
}

// Cache spill PRAGMA functions
func CacheSpill() *Pragma {
	return newPragma("", "cache_spill")
}

func CacheSpillInSchema(schemaName string) *Pragma {
	return newPragma(schemaName, "cache_spill")
}

// Cell size check PRAGMA functions
func CellSizeCheck() *Pragma {
	return newPragma("", "cell_size_check")
}

// Checkpoint fullfsync PRAGMA functions
func CheckpointFullfsync() *Pragma {
	return newPragma("", "checkpoint_fullfsync")
}

// Collation list PRAGMA functions
func CollationList() *Pragma {
	return newPragma("", "collation_list")
}

// Compile options PRAGMA functions
func CompileOptions() *Pragma {
	return newPragma("", "compile_options")
}

// Data version PRAGMA functions
func DataVersion(schemaName string) *Pragma {
	return newPragma(schemaName, "data_version")
}

// Database list PRAGMA functions
func DatabaseList() *Pragma {
	return newPragma("", "database_list")
}

// Defer foreign keys PRAGMA functions
func DeferForeignKeys() *Pragma {
	return newPragma("", "defer_foreign_keys")
}

// Encoding PRAGMA functions
func Encoding() *Pragma {
	return newPragma("", "encoding")
}

func EncodingUTF8() *Pragma {
	return newPragma("", "encoding").ValueType("UTF-8")
}

func EncodingUTF16() *Pragma {
	return newPragma("", "encoding").ValueType("UTF-16")
}

func EncodingUTF16le() *Pragma {
	return newPragma("", "encoding").ValueType("UTF-16le")
}

func EncodingUTF16be() *Pragma {
	return newPragma("", "encoding").ValueType("UTF-16be")
}

// Foreign key check PRAGMA functions
func ForeignKeyCheck(schemaName string) *Pragma {
	return newPragma(schemaName, "foreign_key_check")
}

func ForeignKeyCheckTable(schemaName, tableName string) *Pragma {
	return newPragma(schemaName, "foreign_key_check").FuncType(tableName)
}

// Foreign key list PRAGMA functions
func ForeignKeyList(tableName string) *Pragma {
	return newPragma("", "foreign_key_list").FuncType(tableName)
}

// Foreign keys PRAGMA functions
func ForeignKeys() *Pragma {
	return newPragma("", "foreign_keys")
}

// Freelist count PRAGMA functions
func FreelistCount(schemaName string) *Pragma {
	return newPragma(schemaName, "freelist_count")
}

// Fullfsync PRAGMA functions
func Fullfsync() *Pragma {
	return newPragma("", "fullfsync")
}

// Function list PRAGMA functions
func FunctionList() *Pragma {
	return newPragma("", "function_list")
}

// Hard heap limit PRAGMA functions
func HardHeapLimit() *Pragma {
	return newPragma("", "hard_heap_limit")
}

// Ignore check constraints PRAGMA functions
func IgnoreCheckConstraints() *Pragma {
	return newPragma("", "ignore_check_constraints")
}

// Incremental vacuum PRAGMA functions
func IncrementalVacuum(schemaName string) *Pragma {
	return newPragma(schemaName, "incremental_vacuum")
}

func IncrementalVacuumPages(schemaName string, pages int) *Pragma {
	return newPragma(schemaName, "incremental_vacuum").FuncType(strconv.Itoa(pages))
}

// Index info PRAGMA functions
func IndexInfo(indexName string) *Pragma {
	return newPragma("", "index_info").FuncType(indexName)
}

// Index list PRAGMA functions
func IndexList(tableName string) *Pragma {
	return newPragma("", "index_list").FuncType(tableName)
}

// Index xinfo PRAGMA functions
func IndexXInfo(indexName string) *Pragma {
	return newPragma("", "index_xinfo").FuncType(indexName)
}

// Integrity check PRAGMA functions
func IntegrityCheck(schemaName string) *Pragma {
	return newPragma(schemaName, "integrity_check")
}

func IntegrityCheckLimit(schemaName string, limit int) *Pragma {
	return newPragma(schemaName, "integrity_check").FuncType(strconv.Itoa(limit))
}

func IntegrityCheckTable(schemaName, tableName string) *Pragma {
	return newPragma(schemaName, "integrity_check").FuncType(tableName)
}

// Journal mode PRAGMA functions
func JournalMode(schemaName string) *Pragma {
	return newPragma(schemaName, "journal_mode")
}

func JournalModeDelete(schemaName string) *Pragma {
	return newPragma(schemaName, "journal_mode").ValueType("DELETE")
}

func JournalModeTruncate(schemaName string) *Pragma {
	return newPragma(schemaName, "journal_mode").ValueType("TRUNCATE")
}

func JournalModePersist(schemaName string) *Pragma {
	return newPragma(schemaName, "journal_mode").ValueType("PERSIST")
}

func JournalModeMemory(schemaName string) *Pragma {
	return newPragma(schemaName, "journal_mode").ValueType("MEMORY")
}

func JournalModeWAL(schemaName string) *Pragma {
	return newPragma(schemaName, "journal_mode").ValueType("WAL")
}

func JournalModeOff(schemaName string) *Pragma {
	return newPragma(schemaName, "journal_mode").ValueType("OFF")
}

// Journal size limit PRAGMA functions
func JournalSizeLimit(schemaName string) *Pragma {
	return newPragma(schemaName, "journal_size_limit")
}

// Legacy alter table PRAGMA functions
func LegacyAlterTable() *Pragma {
	return newPragma("", "legacy_alter_table")
}

// Locking mode PRAGMA functions
func LockingMode(schemaName string) *Pragma {
	return newPragma(schemaName, "locking_mode")
}

func LockingModeNormal(schemaName string) *Pragma {
	return newPragma(schemaName, "locking_mode").ValueType("NORMAL")
}

func LockingModeExclusive(schemaName string) *Pragma {
	return newPragma(schemaName, "locking_mode").ValueType("EXCLUSIVE")
}

// Max page count PRAGMA functions
func MaxPageCount(schemaName string) *Pragma {
	return newPragma(schemaName, "max_page_count")
}

// Mmap size PRAGMA functions
func MmapSize(schemaName string) *Pragma {
	return newPragma(schemaName, "mmap_size")
}

// Module list PRAGMA functions
func ModuleList() *Pragma {
	return newPragma("", "module_list")
}

// Optimize PRAGMA functions
func Optimize() *Pragma {
	return newPragma("", "optimize")
}

func OptimizeWithMask(mask string) *Pragma {
	return newPragma("", "optimize").FuncType(mask)
}

func OptimizeSchema(schemaName string) *Pragma {
	return newPragma(schemaName, "optimize")
}

func OptimizeSchemaWithMask(schemaName, mask string) *Pragma {
	return newPragma(schemaName, "optimize").FuncType(mask)
}

// Page count PRAGMA functions
func PageCount(schemaName string) *Pragma {
	return newPragma(schemaName, "page_count")
}

// Page size PRAGMA functions
func PageSize(schemaName string) *Pragma {
	return newPragma(schemaName, "page_size")
}

// Pragma list PRAGMA functions
func PragmaList() *Pragma {
	return newPragma("", "pragma_list")
}

// Query only PRAGMA functions
func QueryOnly() *Pragma {
	return newPragma("", "query_only")
}

// Quick check PRAGMA functions
func QuickCheck(schemaName string) *Pragma {
	return newPragma(schemaName, "quick_check")
}

func QuickCheckLimit(schemaName string, limit int) *Pragma {
	return newPragma(schemaName, "quick_check").FuncType(strconv.Itoa(limit))
}

func QuickCheckTable(schemaName, tableName string) *Pragma {
	return newPragma(schemaName, "quick_check").FuncType(tableName)
}

// Read uncommitted PRAGMA functions
func ReadUncommitted() *Pragma {
	return newPragma("", "read_uncommitted")
}

// Recursive triggers PRAGMA functions
func RecursiveTriggers() *Pragma {
	return newPragma("", "recursive_triggers")
}

// Reverse unordered selects PRAGMA functions
func ReverseUnorderedSelects() *Pragma {
	return newPragma("", "reverse_unordered_selects")
}

// Schema version PRAGMA functions
func SchemaVersion(schemaName string) *Pragma {
	return newPragma(schemaName, "schema_version")
}

// Secure delete PRAGMA functions
func SecureDelete(schemaName string) *Pragma {
	return newPragma(schemaName, "secure_delete")
}

func SecureDeleteFast(schemaName string) *Pragma {
	return newPragma(schemaName, "secure_delete").ValueType("FAST")
}

// Shrink memory PRAGMA functions
func ShrinkMemory() *Pragma {
	return newPragma("", "shrink_memory")
}

// Soft heap limit PRAGMA functions
func SoftHeapLimit() *Pragma {
	return newPragma("", "soft_heap_limit")
}

// Synchronous PRAGMA functions
func Synchronous(schemaName string) *Pragma {
	return newPragma(schemaName, "synchronous")
}

func SynchronousOff(schemaName string) *Pragma {
	return newPragma(schemaName, "synchronous").ValueType("0")
}

func SynchronousNormal(schemaName string) *Pragma {
	return newPragma(schemaName, "synchronous").ValueType("NORMAL")
}

func SynchronousFull(schemaName string) *Pragma {
	return newPragma(schemaName, "synchronous").ValueType("FULL")
}

func SynchronousExtra(schemaName string) *Pragma {
	return newPragma(schemaName, "synchronous").ValueType("EXTRA")
}

// Table info PRAGMA functions
func TableInfo(tableName string) *Pragma {
	return newPragma("", "table_info").FuncType(tableName)
}

// Table list PRAGMA functions
func TableList() *Pragma {
	return newPragma("", "table_list")
}

func TableListInSchema(schemaName string) *Pragma {
	return newPragma(schemaName, "table_list")
}

func TableListForTable(tableName string) *Pragma {
	return newPragma("", "table_list").FuncType(tableName)
}

// Table xinfo PRAGMA functions
func TableXInfo(tableName string) *Pragma {
	return newPragma("", "table_xinfo").FuncType(tableName)
}

// Temp store PRAGMA functions
func TempStore() *Pragma {
	return newPragma("", "temp_store")
}

func TempStoreDefault() *Pragma {
	return newPragma("", "temp_store").ValueType("DEFAULT")
}

func TempStoreFile() *Pragma {
	return newPragma("", "temp_store").ValueType("FILE")
}

func TempStoreMemory() *Pragma {
	return newPragma("", "temp_store").ValueType("MEMORY")
}

// Threads PRAGMA functions
func Threads() *Pragma {
	return newPragma("", "threads")
}

// Trusted schema PRAGMA functions
func TrustedSchema() *Pragma {
	return newPragma("", "trusted_schema")
}

// User version PRAGMA functions
func UserVersion(schemaName string) *Pragma {
	return newPragma(schemaName, "user_version")
}

// WAL autocheckpoint PRAGMA functions
func WalAutocheckpoint() *Pragma {
	return newPragma("", "wal_autocheckpoint")
}

// WAL checkpoint PRAGMA functions
func WalCheckpoint(schemaName string) *Pragma {
	return newPragma(schemaName, "wal_checkpoint")
}

func WalCheckpointPassive(schemaName string) *Pragma {
	return newPragma(schemaName, "wal_checkpoint").FuncType("PASSIVE")
}

func WalCheckpointFull(schemaName string) *Pragma {
	return newPragma(schemaName, "wal_checkpoint").FuncType("FULL")
}

func WalCheckpointRestart(schemaName string) *Pragma {
	return newPragma(schemaName, "wal_checkpoint").FuncType("RESTART")
}

func WalCheckpointTruncate(schemaName string) *Pragma {
	return newPragma(schemaName, "wal_checkpoint").FuncType("TRUNCATE")
}

// Writable schema PRAGMA functions
func WritableSchema() *Pragma {
	return newPragma("", "writable_schema")
}

func WritableSchemaReset() *Pragma {
	return newPragma("", "writable_schema").ValueType("RESET")
}
