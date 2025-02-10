load("https://raw.githubusercontent.com/lodash/lodash/4.17.15-npm/core.min.js")
var FileUtils = Java.type("oracle.dbtools.common.utils.FileUtils")
var Files = Java.type("java.nio.file.Files")
var Paths = Java.type("java.nio.file.Paths")
var cwd = FileUtils.getCWD(ctx);
var FileWriter = Java.type("java.io.FileWriter")

// Suppress all output
sqlcl.setStmt("SET ECHO OFF FEEDBACK OFF VERIFY OFF HEADING OFF PAGESIZE 0")
sqlcl.run()
sqlcl.setStmt("SET TERMOUT OFF")
sqlcl.run()
sqlcl.setStmt("SET SERVEROUTPUT OFF")
sqlcl.run()

// Redirect SQLcl output to a null device (on Unix-like systems)
// For Windows, you might use 'NUL' instead of '/dev/null'
sqlcl.setStmt("SPOOL /dev/null")
sqlcl.run()

function exportToJson(table) {
    var outfw = new FileWriter(cwd.concat("/").concat(table).concat(".json"))
    /* var spoolStmt = "SPOOL '" + cwd + table + ".json' REPLACE;"
    sqlcl.setStmt(spoolStmt)
    sqlcl.run() */
    var stmt = "SELECT * FROM ".concat(table).concat(";")
    var ret = util.executeReturnList(stmt,{})
    _.forEach(ret, function(row) { outfw.write(row) })
    // sqlcl.setStmt("SELECT * FROM " + table + ";")
    // sqlcl.run()
    /* sqlcl.setStmt("SPOOL OFF;")
    sqlcl.run() */
    // Redirect any potential output to /dev/null
    sqlcl.setStmt("SPOOL /dev/null")
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

// Ensure all output is off before exiting
sqlcl.setStmt("SET TERMOUT OFF FEEDBACK OFF VERIFY OFF")
sqlcl.run()

// If you need to log anything, use ctx.write instead of console.log
// ctx.write("Script completed successfully\n")

