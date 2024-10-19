load("https://raw.githubusercontent.com/lodash/lodash/4.17.15-npm/core.min.js")
var FileUtils = Java.type("oracle.dbtools.common.utils.FileUtils")
var Files = Java.type("java.nio.file.Files")
var Paths = Java.type("java.nio.file.Paths")
var cwd = FileUtils.getCWD(ctx);

// Disable terminal output
sqlcl.setStmt("SET TERMOUT OFF")
sqlcl.run()

// Redirect SQLcl output to a null device (on Unix-like systems)
// For Windows, you might use 'NUL' instead of '/dev/null'
sqlcl.setStmt("SPOOL /dev/null")
sqlcl.run()

function exportToJson(table) {
    var spoolStmt = "SPOOL "+"'"+ cwd + table +".json' " + "REPLACE;"
    sqlcl.setStmt(spoolStmt)
    sqlcl.run()
    sqlcl.setStmt("SELECT * FROM " + table + ";")
    sqlcl.run()
    sqlcl.setStmt("SPOOL OFF;")
    sqlcl.run()
}

sqlcl.setStmt("SET SQLFORMAT json-formatted;")
sqlcl.run()
sqlcl.setStmt("SET FEEDBACK off;")
sqlcl.run()

_.forEach(Files.readAllLines(Paths.get(args[1])), exportToJson)

// Restore default settings
sqlcl.setStmt("SPOOL OFF;")
sqlcl.run()
sqlcl.setStmt("SET TERMOUT ON;")
sqlcl.run()
sqlcl.setStmt("SET FEEDBACK on;")
sqlcl.run()
sqlcl.setStmt("SET SQLFORMAT ansiconsole;")
sqlcl.run()

// If you need to log anything, use ctx.write instead of console.log
// ctx.write("Script completed successfully\n")

