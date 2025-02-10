// Check if a filename was provided as a command-line argument
if (args.length === 0) {
    ctx.write('Error: No filename provided.\nUsage: @dump_table_to_json.js <filename>\n');
    exit(1);
}

// Get the filename from the command-line arguments
var filename = args[1];

// Import necessary Java classes for file operations
var Files = Java.type('java.nio.file.Files');
var Paths = Java.type('java.nio.file.Paths');
var StandardCharsets = Java.type('java.nio.charset.StandardCharsets');

// Sanitize the tableName to prevent SQL injection
function sanitizeTableName(name) {
    return name.replace(/[^a-zA-Z0-9_]/g, '');
}

// Define the path to the provided filename
var path = Paths.get(filename);

// Check if the file exists
if (!Files.exists(path)) {
    ctx.write('Error: ' + filename + ' file not found.\n');
    exit(1);
}

// Read all lines (table names) from the file
var lines;
try {
    lines = Files.readAllLines(path, StandardCharsets.UTF_8);
    ctx.write('Successfully read ' + lines.length + ' lines from ' + filename + '\n');
    sqlcl.setStmt("SET LOADFORMAT JSON-FORMATTED;")
    sqlcl.run()
} catch (e) {
    ctx.write('Error reading ' + filename + ': ' + e + '\n');
    exit(1);
}

// Loop over each table name
lines.forEach(function(tableName, index) {
    // Trim whitespace from the table name
    tableName = tableName.trim();
    // Skip empty lines
    if (tableName === '') return;
    // Sanitize the table name to prevent SQL injection
    tableName = sanitizeTableName(tableName);
    // Ensure the sanitized table name is not empty
    if (tableName === '') {
        ctx.write('Skipping invalid table name.\n');
        return;
    }
    try {
        // Debugging output
        ctx.write('Processing table: ' + tableName + '\n');
	sqlcl.setStmt("UNLOAD TABLE ".concat(tableName))
	sqlcl.run()
        ctx.write('Data from table ' + tableName + ' successfully written\n');
    } catch (e) {
        ctx.write('Error processing table ' + tableName + ': ' + e + '\n');
        if (e.javaException) {
            ctx.write('Java Exception: ' + e.javaException + '\n');
            ctx.write('Stack trace: ' + e.javaException.getStackTrace().join('\n') + '\n');
        }
    }
});

ctx.write('\nScript execution completed.\n');
