var FileUtils = Java.type("oracle.dbtools.common.utils.FileUtils")
var Collectors = Java.type("java.util.stream.Collectors")
var Files = Java.type("java.nio.file.Files")
var Paths = Java.type("java.nio.file.Paths")
var BufferedReader = Java.type("java.io.BufferedReader")
var StandardOpenOption = Java.type("java.nio.file.StandardOpenOption")
var String = Java.type("java.lang.String")

var cwd = FileUtils.getCWD(ctx);
var outputPath = Paths.get(cwd, "/","featureprop.json")
Files.writeString(outputPath, "", StandardOpenOption.CREATE, StandardOpenOption.TRUNCATE_EXISTING)

var stmt = `
SELECT 
    featureprop_id,
    value
FROM featureprop 
WHERE value IS NOT NULL
`

ret = util.executeReturnList(stmt, {})
ret.forEach(function(row){
	reader = row.VALUE.getCharacterStream()
	content = new BufferedReader(reader)
        	 .lines()
        	 .collect(Collectors.joining("\n"))
	Files.writeString(
		outputPath,
		String.join(",",row.FEATUREPROP_ID,content), 
		StandardOpenOption.APPEND
	)
	Files.writeString(
		outputPath,
		"\n",
		StandardOpenOption.APPEND
	)
	reader.close()
})
