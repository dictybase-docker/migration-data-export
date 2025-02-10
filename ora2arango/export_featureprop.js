FileUtils = Java.type("oracle.dbtools.common.utils.FileUtils")
Collectors = Java.type("java.util.stream.Collectors")
Files = Java.type("java.nio.file.Files")
Paths = Java.type("java.nio.file.Paths")
BufferedReader = Java.type("java.io.BufferedReader")
StandardOpenOption = Java.type("java.nio.file.StandardOpenOption")
String = Java.type("java.lang.String")

cwd = FileUtils.getCWD(ctx);
outputPath = Paths.get(cwd, "/","featureprop_clob.csv")
Files.writeString(
	outputPath, "", 
	StandardOpenOption.CREATE, StandardOpenOption.TRUNCATE_EXISTING
)

stmt = `
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
